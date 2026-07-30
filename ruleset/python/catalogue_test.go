package python_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ensures Python catalogue seeds parse as JSON objects with the expected shape.
// No detector wiring.
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
	var bp map[string]map[string]any
	if err := json.Unmarshal(bpRaw, &bp); err != nil {
		t.Fatalf("bad-practices.json: %v", err)
	}
	if len(bp) == 0 {
		t.Fatal("bad-practices.json: expected non-empty Python BP catalogue")
	}

	bpRequired := []string{"id", "name", "description", "detection_notes", "severity", "category"}
	validSeverity := map[string]bool{
		"info": true, "low": true, "medium": true, "high": true, "critical": true,
	}
	// Framework / stack tags we intentionally target in this seed set.
	wantTags := map[string]bool{
		"python": true, "flask": true, "django": true, "fastapi": true,
		"starlette": true, "jinja2": true, "sqlalchemy": true, "requests": true, "httpx": true,
	}
	seenCategories := map[string]int{}
	seenTags := map[string]int{}

	for key, entry := range bp {
		if !strings.HasPrefix(key, "BP-PY-") {
			t.Errorf("key %q: want BP-PY-* prefix (distinct from Go BP-*)", key)
		}
		for _, k := range bpRequired {
			if _, ok := entry[k]; !ok {
				t.Errorf("%s: missing field %q", key, k)
			}
		}
		if id, _ := entry["id"].(string); id != key {
			t.Errorf("%s: id field %q must match map key", key, id)
		}
		sev, _ := entry["severity"].(string)
		if !validSeverity[sev] {
			t.Errorf("%s: invalid severity %q", key, sev)
		}
		cat, _ := entry["category"].(string)
		if cat == "" {
			t.Errorf("%s: empty category", key)
		}
		seenCategories[cat]++

		if apps, ok := entry["applicable_to"].([]any); ok {
			foundPython := false
			for _, a := range apps {
				s, ok := a.(string)
				if !ok {
					continue
				}
				seenTags[s]++
				if s == "python" {
					foundPython = true
				}
				if !wantTags[s] {
					// Allow extra tags, but warn-level only via test log
					t.Logf("%s: applicable_to tag %q not in seed target set", key, s)
				}
			}
			if !foundPython {
				t.Errorf("%s: applicable_to must include \"python\"", key)
			}
		}
	}

	// Ensure framework-oriented coverage is present in the seed set.
	for _, need := range []string{"Flask", "Django", "FastAPI", "Database", "Templates", "Security Hygiene"} {
		if seenCategories[need] == 0 {
			t.Errorf("bad-practices.json: missing category %q in seed set", need)
		}
	}
	for _, tag := range []string{"flask", "django", "fastapi", "sqlalchemy", "jinja2"} {
		if seenTags[tag] == 0 {
			t.Errorf("bad-practices.json: no rule tagged applicable_to %q", tag)
		}
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
