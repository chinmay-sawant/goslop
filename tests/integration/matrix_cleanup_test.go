package integration

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
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

func TestMaterializeAndScanOptsCleansTemporaryDirectoryAfterAnalyzerFailure(t *testing.T) {
	tempRoot := filepath.Join(t.TempDir(), "materialized")
	originalMkdirTemp := osMkdirTempFixture
	osMkdirTempFixture = func() (string, error) {
		if err := os.Mkdir(tempRoot, 0o755); err != nil {
			return "", err
		}
		return tempRoot, nil
	}
	t.Cleanup(func() { osMkdirTempFixture = originalMkdirTemp })

	wantErr := errors.New("injected analyzer failure")
	_, err := materializeAndScanOpts(
		"go/perf/PERF-006-vulnerable.txt",
		DefaultMatrixOptions(),
		func(*core.ScanContext, string) (*engine.AnalysisResult, error) {
			return nil, wantErr
		},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected analyzer error to be wrapped, got %v", err)
	}
	if _, statErr := os.Stat(tempRoot); !os.IsNotExist(statErr) {
		t.Fatalf("temporary materialization directory remains: %v", statErr)
	}
}
