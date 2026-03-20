// Package engine implements the core logic for Purgomatic.
package engine

import (
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
)

// DBInit handles schema creation.
const DBInit = `
CREATE TABLE IF NOT EXISTS files (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	hostname TEXT,
	source TEXT,
	path TEXT,
	name TEXT,
	ext TEXT,
	size INTEGER,
	hash TEXT,
	year INTEGER,
	mtime INTEGER,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(hostname, source, path)
);
CREATE INDEX IF NOT EXISTS idx_hash ON files(hash);
CREATE INDEX IF NOT EXISTS idx_source ON files(source);
CREATE INDEX IF NOT EXISTS idx_mtime ON files(mtime);
`

// GetDBPath determines the database file path based on environment, CWD, or home directory.
func GetDBPath() string {
	if env := os.Getenv("PURGOMATIC_DB"); env != "" {
		return env
	}
	// Check current directory
	cwd, _ := os.Getwd()
	cwdPath := filepath.Join(cwd, "purgomatic.db")
	if _, err := os.Stat(cwdPath); err == nil {
		return cwdPath
	}
	// Check home directory
	home, _ := os.UserHomeDir()
	homePath := filepath.Join(home, "purgomatic.db")
	if _, err := os.Stat(homePath); err == nil {
		return homePath
	}
	// Default to home directory
	return homePath
}

// ScanStruct maps a *sql.Rows row to a struct tagged with `db`.
func ScanStruct(rows *sql.Rows, dest any) error {
	v := reflect.ValueOf(dest).Elem()
	t := v.Type()
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	ptrs := make([]any, len(cols))
	for i, col := range cols {
		found := false
		for j := 0; j < t.NumField(); j++ {
			f := t.Field(j)
			if f.Tag.Get("db") == col {
				ptrs[i] = v.Field(j).Addr().Interface()
				found = true
				break
			}
		}
		if !found {
			var ignored any
			ptrs[i] = &ignored
		}
	}
	return rows.Scan(ptrs...)
}
