package python_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const maxRulesPerChunk = 50

func catalogueDir(t *testing.T) string {
	t.Helper()
	dir := "."
	if _, err := os.Stat(filepath.Join(dir, "bad-practices.json")); err != nil {
		dir = "ruleset/python"
	}
	if _, err := os.Stat(filepath.Join(dir, "bad-practices.json")); err != nil {
		t.Fatalf("cannot locate ruleset/python (cwd bad-practices.json): %v", err)
	}
	return dir
}

// Ensures Python catalogue seeds parse as JSON objects with the expected shape.
// No detector wiring.
func TestPythonBadPracticesCatalogue(t *testing.T) {
	t.Parallel()
	dir := catalogueDir(t)

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
					t.Logf("%s: applicable_to tag %q not in seed target set", key, s)
				}
			}
			if !foundPython {
				t.Errorf("%s: applicable_to must include \"python\"", key)
			}
		}
	}

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
}

func TestPythonCWEChunksFrom699(t *testing.T) {
	t.Parallel()
	dir := catalogueDir(t)
	chunkDir := filepath.Join(dir, "chunks")

	entries, err := os.ReadDir(chunkDir)
	if err != nil {
		t.Fatalf("read chunks: %v", err)
	}

	required := []string{
		"id", "name", "description", "original_description", "detection_notes",
		"category", "status", "weakness_abstraction", "go_relevance", "applicable_to",
	}
	total := 0
	seenIDs := map[int]string{}
	mustHave := map[int]bool{22: false, 78: false, 79: false, 89: false, 502: false}
	var cweFiles int

	for _, ent := range entries {
		name := ent.Name()
		if ent.IsDir() || !strings.HasPrefix(name, "cwe-") || !strings.HasSuffix(name, ".json") {
			continue
		}
		cweFiles++
		raw, err := os.ReadFile(filepath.Join(chunkDir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		var doc map[string]map[string]any
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if len(doc) == 0 {
			t.Errorf("%s: empty object", name)
			continue
		}
		if len(doc) > maxRulesPerChunk {
			t.Errorf("%s: %d rules exceeds max %d", name, len(doc), maxRulesPerChunk)
		}
		total += len(doc)

		for key, entry := range doc {
			if !strings.HasPrefix(key, "CWE-") {
				t.Errorf("%s: key %q want CWE-*", name, key)
				continue
			}
			for _, k := range required {
				if _, ok := entry[k]; !ok {
					t.Errorf("%s %s: missing field %q", name, key, k)
				}
			}
			idFloat, ok := entry["id"].(float64) // JSON numbers
			if !ok {
				t.Errorf("%s %s: id must be number", name, key)
				continue
			}
			id := int(idFloat)
			if key != "CWE-"+itoa(id) {
				t.Errorf("%s: key %s does not match id %d", name, key, id)
			}
			if prev, dup := seenIDs[id]; dup {
				t.Errorf("duplicate CWE-%d in %s and %s", id, prev, name)
			}
			seenIDs[id] = name
			if _, want := mustHave[id]; want {
				mustHave[id] = true
			}

			apps, ok := entry["applicable_to"].([]any)
			if !ok || len(apps) == 0 {
				t.Errorf("%s %s: applicable_to must be non-empty array", name, key)
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
				t.Errorf("%s %s: applicable_to must include python", name, key)
			}
		}
	}

	if cweFiles == 0 {
		t.Fatal("no cwe-*.json chunk files found")
	}
	if total < 100 {
		t.Fatalf("expected a large Python CWE set from 699.csv, got %d rules in %d files", total, cweFiles)
	}
	for id, ok := range mustHave {
		if !ok {
			t.Errorf("missing required CWE-%d in catalogue", id)
		}
	}
	t.Logf("python CWE catalogue: %d rules across %d files", total, cweFiles)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [16]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
