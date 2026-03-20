package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFastHash(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "purgomatic-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	t.Run("SmallFile", func(t *testing.T) {
		path := filepath.Join(tmpDir, "small.txt")
		content := []byte("Hello, Purgomatic!")
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatal(err)
		}

		h1, err := fastHash(path)
		if err != nil {
			t.Fatal(err)
		}
		if h1 == "" {
			t.Error("expected non-empty hash")
		}

		// Verify consistency
		h2, _ := fastHash(path)
		if h1 != h2 {
			t.Error("hashes do not match")
		}
	})

	t.Run("LargeFile", func(t *testing.T) {
		path := filepath.Join(tmpDir, "large.bin")
		// Create a file larger than 48KB (chunk*3) to trigger multi-point sampling
		size := 100 * 1024
		content := make([]byte, size)
		for i := range content {
			content[i] = byte(i % 256)
		}
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatal(err)
		}

		h1, err := fastHash(path)
		if err != nil {
			t.Fatal(err)
		}
		if h1 == "" {
			t.Error("expected non-empty hash")
		}

		// Verify consistency
		h2, _ := fastHash(path)
		if h1 != h2 {
			t.Error("hashes do not match")
		}
	})
}
