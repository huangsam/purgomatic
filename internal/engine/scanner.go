// Package engine implements the core logic for Purgomatic.
package engine

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// HandleScan scans a directory and indexes metadata into SQLite.
func HandleScan(source, root string) error {
	host, _ := os.Hostname()
	db, err := sql.Open("sqlite", GetDBPath())
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	fmt.Printf("--- Stat-First Scan [%s] on host [%s]: %s ---\n", source, host, root)

	existing := make(map[string]DbStat)
	rows, _ := db.Query("SELECT path, size, mtime, hash FROM files WHERE hostname=? AND source=?", host, source)
	if rows != nil {
		for rows.Next() {
			var s DbStat
			if err := ScanStruct(rows, &s); err == nil {
				existing[s.Path] = s
			}
		}
		_ = rows.Close()
	}

	jobChan := make(chan Job, 1000)
	resultChan := make(chan ScanResult, 1000)
	var wg sync.WaitGroup

	allPhysicallySeen := make(map[string]bool)
	go func() {
		_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			allPhysicallySeen[path] = true
			jobChan <- Job{Path: path, Info: info}
			return nil
		})
		close(jobChan)
	}()

	numWorkers := runtime.NumCPU() * 2
	for range numWorkers {
		wg.Go(func() {
			for j := range jobChan {
				mt := j.Info.ModTime().Unix()
				sz := j.Info.Size()

				if ex, ok := existing[j.Path]; ok && ex.Size == sz && ex.Mtime == mt {
					continue
				}

				h, _ := fastHash(j.Path)
				resultChan <- ScanResult{
					Path:  j.Path,
					Hash:  h,
					Size:  sz,
					Year:  j.Info.ModTime().Year(),
					Mtime: mt,
				}
			}
		})
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	fmt.Println("Indexing files (Transactional Stat-Skip + Sequential Writes)...")
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	count := 0
	for r := range resultChan {
		name := filepath.Base(r.Path)
		ext := strings.ToLower(filepath.Ext(r.Path))

		_, err := tx.Exec(`INSERT INTO files (hostname, source, path, name, ext, size, hash, year, mtime)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(hostname, source, path) DO UPDATE SET
			size=excluded.size, hash=excluded.hash, year=excluded.year, mtime=excluded.mtime`,
			host, source, r.Path, name, ext, r.Size, r.Hash, r.Year, r.Mtime)
		if err != nil {
			fmt.Printf("DB error for %s: %v\n", r.Path, err)
		}
		count++
		if count%5000 == 0 {
			fmt.Printf("Indexed %d files...\n", count)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Println("Pruning ghost entries (files no longer on disk)...")
	pruneCount := 0
	ptx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = ptx.Rollback() }()

	for p := range existing {
		if !allPhysicallySeen[p] && strings.HasPrefix(p, root) {
			_, err := ptx.Exec("DELETE FROM files WHERE hostname=? AND source=? AND path=?", host, source, p)
			if err == nil {
				pruneCount++
			}
		}
	}

	if err := ptx.Commit(); err != nil {
		return err
	}

	fmt.Printf("Scan completed. Indexed %d changes. Pruned %d ghosts. Total files physically found: %d\n", count, pruneCount, len(allPhysicallySeen))
	return nil
}

func fastHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	info, err := f.Stat()
	if err != nil {
		return "", err
	}

	size := info.Size()
	const chunk = 16 * 1024
	h := sha256.New()

	if size <= chunk*3 {
		_, _ = io.Copy(h, f)
	} else {
		buf := make([]byte, chunk)
		_, _ = f.Read(buf)
		_, _ = h.Write(buf)

		_, _ = f.Seek(size/2, io.SeekStart)
		_, _ = f.Read(buf)
		_, _ = h.Write(buf)

		_, _ = f.Seek(size-chunk, io.SeekStart)
		_, _ = f.Read(buf)
		_, _ = h.Write(buf)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
