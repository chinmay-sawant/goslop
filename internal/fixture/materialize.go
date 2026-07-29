package fixture

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// MaterializeFixture writes fixture source under root/<lang>/<filename>.
// Returns the absolute path of the written file.
func MaterializeFixture(root string, fixture TextFixture) (string, error) {
	return writeFixtureAt(root, fixture)
}

// MaterializeFixtureFile reads a .txt fixture from disk and materializes it under root.
func MaterializeFixtureFile(txtPath, root string) (string, error) {
	data, err := os.ReadFile(txtPath)
	if err != nil {
		return "", fmt.Errorf("reading fixture %s: %w", txtPath, err)
	}
	fixture, err := ParseFixture(string(data), txtPath)
	if err != nil {
		return "", fmt.Errorf("%s: %w", txtPath, err)
	}
	return writeFixtureAt(root, fixture)
}

// MaterializeTree walks all *.txt fixtures under fixturesRoot and writes
// materialized sources under outRoot/<lang>/<filename>.
//
// Python fixtures are materialized into outRoot/python/ (tagged by lang dir);
// detectors may still skip them. Go fixtures land under outRoot/go/.
func MaterializeTree(fixturesRoot, outRoot string) error {
	if err := os.MkdirAll(outRoot, 0o755); err != nil {
		return fmt.Errorf("creating materialize root %s: %w", outRoot, err)
	}
	return filepath.WalkDir(fixturesRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("fixture traversal failed: %w", err)
		}
		if d.IsDir() {
			return nil
		}
		if !strings.EqualFold(filepath.Ext(path), "."+FIXTURE_EXTENSION) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading fixture %s: %w", path, err)
		}
		fixture, err := ParseFixture(string(data), path)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if _, err := writeFixtureAt(outRoot, fixture); err != nil {
			return err
		}
		return nil
	})
}

// SanitizeFilename rejects absolute paths and parent-directory segments.
// Returns a cleaned relative path suitable for joining under the materialize root.
func SanitizeFilename(filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("fixture filename must not be empty")
	}
	// Normalize to slash form for component checks, but preserve OS path on return.
	if filepath.IsAbs(filename) {
		return "", fmt.Errorf("fixture filename must be relative, got absolute path %q", filename)
	}
	// Also reject Windows-style absolute paths on non-Windows (path with volume).
	if vol := filepath.VolumeName(filename); vol != "" {
		return "", fmt.Errorf("fixture filename must be relative, got absolute path %q", filename)
	}

	cleaned := filepath.Clean(filename)
	// Clean(".") or empty after clean is not a usable relative file path.
	if cleaned == "." || cleaned == "" {
		return "", fmt.Errorf("fixture filename must be relative, got absolute path %q", filename)
	}
	// After Clean, ".." at the start means escape attempt (e.g. "../x" or "a/../../x").
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("fixture filename must not contain '..': %q", filename)
	}
	for _, comp := range strings.Split(filepath.ToSlash(filename), "/") {
		if comp == ".." {
			return "", fmt.Errorf("fixture filename must not contain '..': %q", filename)
		}
	}
	return cleaned, nil
}

func writeFixtureAt(root string, fixture TextFixture) (string, error) {
	filename, err := SanitizeFilename(fixture.Filename)
	if err != nil {
		return "", err
	}

	langDir := filepath.Join(root, fixture.Language.String())
	if err := os.MkdirAll(langDir, 0o755); err != nil {
		return "", fmt.Errorf("creating fixture directory %s: %w", langDir, err)
	}

	outPath := filepath.Join(langDir, filename)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return "", fmt.Errorf("creating fixture directory %s: %w", filepath.Dir(outPath), err)
	}

	// Defense in depth: resolved parent must stay under langDir.
	langAbs, err := filepath.Abs(langDir)
	if err != nil {
		return "", fmt.Errorf("resolving fixture directory %s: %w", langDir, err)
	}
	parentAbs, err := filepath.Abs(filepath.Dir(outPath))
	if err != nil {
		return "", fmt.Errorf("resolving fixture directory %s: %w", filepath.Dir(outPath), err)
	}
	// Ensure parentAbs is langAbs or a descendant.
	rel, err := filepath.Rel(langAbs, parentAbs)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("fixture path escapes materialize root: %s", outPath)
	}

	if err := os.WriteFile(outPath, []byte(fixture.Source), 0o644); err != nil {
		return "", fmt.Errorf("writing materialized fixture %s: %w", outPath, err)
	}
	abs, err := filepath.Abs(outPath)
	if err != nil {
		return outPath, nil
	}
	return abs, nil
}
