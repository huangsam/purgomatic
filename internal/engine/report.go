// Package engine implements the core logic for Purgomatic.
package engine

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// HandleAudit combines scanning and reporting into a single operation.
func HandleAudit(batchFile string) error {
	data, err := os.ReadFile(batchFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", batchFile, err)
	}

	var targets []ScanTarget
	if err := json.Unmarshal(data, &targets); err != nil {
		return err
	}

	fmt.Printf("Starting audit of %d targets...\n", len(targets))
	for _, t := range targets {
		s := t.Source
		if s == "" {
			s = filepath.Base(t.Path)
		}
		if err := HandleScan(s, t.Path); err != nil {
			fmt.Printf("Error scanning %s [%s]: %v\n", s, t.Path, err)
		}
	}

	return HandleReport()
}

// HandleReport generates a global SRE-grade audit dashboard.
func HandleReport() error {
	db, err := sql.Open("sqlite", GetDBPath())
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	fmt.Println("\n--- [ Purgomatic Global Library Audit ] ---")

	var allFiles []dbFile
	rows, err := db.Query("SELECT hostname, source, path, name, ext, size, hash, year FROM files ORDER BY size DESC")
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var f dbFile
		if err := ScanStruct(rows, &f); err == nil {
			allFiles = append(allFiles, f)
		}
	}

	homeWinners := getHomeWinners(allFiles)

	var totalBytes int64
	var toilBytes int64
	var libraryBytes int64
	var junkBytes int64

	uniqueHashes := make(map[string]bool)
	yearStats := make(map[int]int64)
	yearArchived := make(map[int]int64)

	type catSummary struct {
		count int
		size  int64
	}
	summaries := make(map[string]*catSummary)

	var globalPending []dbFile

	for _, f := range allFiles {
		totalBytes += f.Size

		lowName := strings.ToLower(f.FullName)
		if strings.HasPrefix(lowName, ".") || strings.Contains(lowName, "thumb") || f.Ext == ".log" || f.Ext == ".tmp" {
			junkBytes += f.Size
			continue
		}

		strat := getStrategy(f.Ext, f.Year)
		if _, ok := summaries[strat.Category]; !ok {
			summaries[strat.Category] = &catSummary{}
		}
		summaries[strat.Category].count++
		summaries[strat.Category].size += f.Size

		isArchived := strings.Contains(f.Path, "Synology") || strings.Contains(f.Path, "ExternalSSD")

		if _, atHome := homeWinners[f.Hash]; atHome && !isArchived {
			toilBytes += f.Size
		}

		if !uniqueHashes[f.Hash] {
			uniqueHashes[f.Hash] = true
			yearStats[f.Year] += f.Size
			if isArchived {
				libraryBytes += f.Size
				yearArchived[f.Year] += f.Size
			} else if len(globalPending) < 10 {
				globalPending = append(globalPending, f)
			}
		}
	}

	fmt.Printf("\n[ STORAGE OVERVIEW ]\n")
	fmt.Printf("Total Managed Space     : %s (%d files)\n", formatBytes(totalBytes), len(allFiles))
	fmt.Printf("Recoverable Toil        : %s (Redundant copies of Home assets)\n", formatBytes(toilBytes))
	fmt.Printf("Unique Library Assets   : %s (%d unique hashes)\n", formatBytes(libraryBytes), len(uniqueHashes))
	fmt.Printf("System Junk             : %s (Metadata/Logs)\n", formatBytes(junkBytes))

	fmt.Printf("\n[ LIBRARY HEALTH BY YEAR ]\n")
	fmt.Printf("%-6s | %-10s | %-12s | %-8s\n", "YEAR", "STATUS", "UNIQUE SIZE", "ARCHIVED")
	fmt.Printf("-------+------------+--------------+----------\n")
	var years []int
	for y := range yearStats {
		years = append(years, y)
	}
	sort.Ints(years)
	for _, y := range years {
		archived := yearArchived[y]
		total := yearStats[y]
		pct := 0.0
		if total > 0 {
			pct = (float64(archived) / float64(total)) * 100
		}
		status := "PENDING"
		if pct >= 99.0 {
			status = "ARCHIVED"
		}
		fmt.Printf("%-6d | %-10s | %-12s | %6.1f%%\n", y, status, formatBytes(total), pct)
	}

	if len(globalPending) > 0 {
		fmt.Printf("\n[ ACTIONABLE: TOP PENDING ASSETS ]\n")
		for _, f := range globalPending {
			fmt.Printf("! %-40s | %s\n", f.FullName, formatBytes(f.Size))
		}
	}

	fmt.Printf("\n[ STRATEGIC RECOMMENDATIONS ]\n")
	orderedCats := []string{"Family Video", "Master RAW", "Photo Library", "Uncategorized"}
	for _, name := range orderedCats {
		s, ok := summaries[name]
		if !ok || s.count == 0 {
			continue
		}
		strat := getStrategy(nameToExt(name), 0)
		fmt.Printf("- %-15s: %d files (%s) -> Move to [%s]\n", name, s.count, formatBytes(s.size), strat.TargetDir)
		fmt.Printf("  └─ %s\n", strat.Advice)
	}

	fmt.Printf("\n[ HOST FOOTPRINT ]\n")
	srows, _ := db.Query("SELECT hostname, source, count(*), sum(size) FROM files GROUP BY hostname, source ORDER BY sum(size) DESC")
	if srows != nil {
		for srows.Next() {
			var host, src string
			var count int
			var total int64
			if err := srows.Scan(&host, &src, &count, &total); err == nil {
				fmt.Printf("- %s @ %-15s: %d files, %s\n", host, src, count, formatBytes(total))
			}
		}
		_ = srows.Close()
	}

	return nil
}

func getStrategy(ext string, year int) AssetStrategy {
	switch ext {
	case ".mp4", ".mov":
		return AssetStrategy{
			Category:  "Family Video",
			TargetDir: "Synology/Videos/",
			Advice:    "High-bandwidth media. Best hosted on primary NAS for streaming.",
		}
	case ".dng", ".arw":
		return AssetStrategy{
			Category:  "Master RAW",
			TargetDir: "ExternalSSD/Archive/RAWs/",
			Advice:    "High fidelity archival. Keep off-site or on dedicated fast storage.",
		}
	case ".jpg", ".jpeg", ".png", ".heic":
		return AssetStrategy{
			Category:  "Photo Library",
			TargetDir: fmt.Sprintf("Synology/Photos/%d/", year),
			Advice:    "Primary asset history. Consolidate into year-based folders.",
		}
	default:
		return AssetStrategy{
			Category:  "Uncategorized",
			TargetDir: "Misc/Other/",
			Advice:    "Miscellaneous files. Requires manual audit.",
		}
	}
}

func getHomeWinners(files []dbFile) map[string]string {
	winners := make(map[string]string)
	for _, f := range files {
		strat := getStrategy(f.Ext, f.Year)
		if strat.TargetDir != "" && strings.HasPrefix(f.Path, strat.TargetDir) && f.Hash != "" {
			if _, exists := winners[f.Hash]; !exists {
				winners[f.Hash] = f.Path
			}
		}
	}
	return winners
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func nameToExt(name string) string {
	switch name {
	case "Family Video":
		return ".mp4"
	case "Master RAW":
		return ".dng"
	case "Photo Library":
		return ".jpg"
	default:
		return ""
	}
}
