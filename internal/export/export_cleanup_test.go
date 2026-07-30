package export

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportContextReturnsOwnedFileRemovalFailure(t *testing.T) {
	dir := t.TempDir()
	stalePath := filepath.Join(dir, "1.txt")
	if err := os.WriteFile(stalePath, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}

	originalRemove := removeOwnedFile
	removeOwnedFile = func(path string) error {
		if path == stalePath {
			return errors.New("injected remove failure")
		}
		return originalRemove(path)
	}
	t.Cleanup(func() { removeOwnedFile = originalRemove })

	_, err := ExportFindings(nil, Options{
		ExportContext:    true,
		ContextOutputDir: dir,
	}, nil)
	if err == nil || !strings.Contains(err.Error(), "clean context output") ||
		!strings.Contains(err.Error(), "injected remove failure") {
		t.Fatalf("expected cleanup failure, got %v", err)
	}
	if _, err := os.Stat(stalePath); err != nil {
		t.Fatalf("owned file should remain when removal fails: %v", err)
	}
}
