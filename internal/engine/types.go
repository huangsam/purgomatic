// Package engine implements the core logic for Purgomatic.
package engine

import (
	"os"
)

// AssetStrategy defines the organizational plan for a file type.
type AssetStrategy struct {
	Category  string
	TargetDir string
	Advice    string
}

// dbFile is the persistent model for a file in the SQLite index.
type dbFile struct {
	Host     string `db:"hostname"`
	Src      string `db:"source"`
	Path     string `db:"path"`
	FullName string `db:"name"`
	Ext      string `db:"ext"`
	Hash     string `db:"hash"`
	Size     int64  `db:"size"`
	Year     int    `db:"year"`
}

// ScanTarget represents a target directory and source name for batch scanning.
type ScanTarget struct {
	Source string `json:"source"`
	Path   string `json:"path"`
}

// DbStat represents the file metadata stored in the database.
type DbStat struct {
	Path  string `db:"path"`
	Size  int64  `db:"size"`
	Mtime int64  `db:"mtime"`
	Hash  string `db:"hash"`
}

// Job represents a file to be processed by a scanner worker.
type Job struct {
	Path string
	Info os.FileInfo
}

// ScanResult represents the calculated metadata for a processed file.
type ScanResult struct {
	Path  string
	Hash  string
	Size  int64
	Year  int
	Mtime int64
}
