package integration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMaterializeAndScanOptsCleansTemporaryDirectoryAfterMaterializeFailure(t *testing.T) {
	tempRoot := filepath.Join(t.TempDir(), "materialized")
	originalMkdirTemp := osMkdirTempFixture
	osMkdirTempFixture = func() (string, error) {
		if err := os.Mkdir(tempRoot, 0o755); err != nil {
			return "", err
		}
		return tempRoot, nil
	}
	t.Cleanup(func() { osMkdirTempFixture = originalMkdirTemp })

	_, err := MaterializeAndScanOpts("go/does-not-exist.txt", DefaultMatrixOptions())
	if err == nil {
		t.Fatal("expected materialization failure")
	}
	if _, statErr := os.Stat(tempRoot); !os.IsNotExist(statErr) {
		t.Fatalf("temporary materialization directory remains: %v", statErr)
	}
}
