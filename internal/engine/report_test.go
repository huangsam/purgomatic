package engine

import (
	"testing"
)

func TestGetStrategy(t *testing.T) {
	tests := []struct {
		name    string
		ext     string
		year    int
		wantCat string
		wantDir string
	}{
		{"Video", ".mp4", 2023, "Family Video", "Synology/Videos/"},
		{"RAW", ".dng", 2023, "Master RAW", "ExternalSSD/Archive/RAWs/"},
		{"Photo", ".jpg", 2023, "Photo Library", "Synology/Photos/2023/"},
		{"Unknown", ".txt", 2023, "Uncategorized", "Misc/Other/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStrategy(tt.ext, tt.year)
			if got.Category != tt.wantCat || got.TargetDir != tt.wantDir {
				t.Errorf("getStrategy() = %v, %v; want %v, %v", got.Category, got.TargetDir, tt.wantCat, tt.wantDir)
			}
		})
	}
}

func TestGetHomeWinners(t *testing.T) {
	files := []dbFile{
		{Hash: "abc", Path: "Synology/Photos/2023/a.jpg", Ext: ".jpg", Year: 2023},
		{Hash: "abc", Path: "Local/a.jpg", Ext: ".jpg", Year: 2023},
		{Hash: "def", Path: "Local/b.jpg", Ext: ".jpg", Year: 2023},
	}

	winners := getHomeWinners(files)
	if len(winners) != 1 {
		t.Errorf("expected 1 winner, got %d", len(winners))
	}
	if winners["abc"] != "Synology/Photos/2023/a.jpg" {
		t.Errorf("expected winner to be Synology path, got %s", winners["abc"])
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
	}

	for _, tt := range tests {
		got := formatBytes(tt.bytes)
		if got != tt.want {
			t.Errorf("formatBytes(%d) = %s; want %s", tt.bytes, got, tt.want)
		}
	}
}

func TestNameToExt(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Family Video", ".mp4"},
		{"Master RAW", ".dng"},
		{"Photo Library", ".jpg"},
		{"Unknown", ""},
	}

	for _, tt := range tests {
		got := nameToExt(tt.name)
		if got != tt.want {
			t.Errorf("nameToExt(%s) = %s; want %s", tt.name, got, tt.want)
		}
	}
}
