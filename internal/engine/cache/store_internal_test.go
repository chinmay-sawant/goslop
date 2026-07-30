package cache

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClearFilesDirReturnsRemovalFailure(t *testing.T) {
	dir := t.TempDir()
	entry := filepath.Join(dir, "stale.json")
	if err := os.WriteFile(entry, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := clearFilesDir(dir, func(path string) error {
		if path != entry {
			t.Fatalf("remove path = %q, want %q", path, entry)
		}
		return errors.New("injected removal failure")
	})
	if err == nil || !strings.Contains(err.Error(), "injected removal failure") {
		t.Fatalf("expected removal failure, got %v", err)
	}
	if _, err := os.Stat(entry); err != nil {
		t.Fatalf("entry should remain after removal failure: %v", err)
	}
}
