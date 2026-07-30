package python_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Ensures Phase 4 Python catalogue seeds parse as JSON objects with the
// expected shape. No detector wiring.
func TestPythonCatalogueJSONParses(t *testing.T) {
	t.Parallel()

	dir := "."
	// Allow running from repo root: go test ./ruleset/python
	if _, err := os.Stat(filepath.Join(dir, "bad-practices.json")); err != nil {
		dir = "ruleset/python"
	}

	bpRaw, err := os.ReadFile(filepath.Join(dir, "bad-practices.json"))
	if err != nil {
		t.Fatalf("read bad-practices.json: %v", err)
	}
	var bp map[string]json.RawMessage
	if err := json.Unmarshal(bpRaw, &bp); err != nil {
		t.Fatalf("bad-practices.json: %v", err)
	}
	if len(bp) != 0 {
		t.Fatalf("bad-practices.json: want empty object for Phase 4 seed, got %d keys", len(bp))
	}

	cweRaw, err := os.ReadFile(filepath.Join(dir, "chunks", "cwe-seed.json"))
	if err != nil {
		t.Fatalf("read cwe-seed.json: %v", err)
	}
	var cwe map[string]map[string]any
	if err := json.Unmarshal(cweRaw, &cwe); err != nil {
		t.Fatalf("cwe-seed.json: %v", err)
	}

	wantIDs := []string{"CWE-22", "CWE-78", "CWE-79", "CWE-89"}
	if len(cwe) != len(wantIDs) {
		t.Fatalf("cwe-seed.json: want %d rules, got %d", len(wantIDs), len(cwe))
	}
	required := []string{
		"id", "name", "description", "original_description", "detection_notes",
		"category", "status", "weakness_abstraction", "go_relevance", "applicable_to",
	}
	for _, id := range wantIDs {
		entry, ok := cwe[id]
		if !ok {
			t.Fatalf("missing %s", id)
			continue
		}
		for _, k := range required {
			if _, ok := entry[k]; !ok {
				t.Errorf("%s: missing field %q", id, k)
			}
		}
		apps, ok := entry["applicable_to"].([]any)
		if !ok || len(apps) == 0 {
			t.Errorf("%s: applicable_to must be a non-empty array", id)
			continue
		}
		foundPython := false
		for _, a := range apps {
			if s, ok := a.(string); ok && s == "python" {
				foundPython = true
				break
			}
		}
		if !foundPython {
			t.Errorf("%s: applicable_to must include \"python\", got %v", id, apps)
		}
	}
}
