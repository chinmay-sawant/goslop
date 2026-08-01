package perf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type catalogueEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestMetadataMatchesPERFPYCatalogue(t *testing.T) {
	entries := make(map[string]catalogueEntry)
	for _, name := range []string{"perf-py-001-014.json", "perf-py-015-022.json", "perf-py-023-030.json"} {
		path := filepath.Join("..", "..", "..", "..", "..", "ruleset", "python", "chunks", name)
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		var chunk map[string]catalogueEntry
		if err := json.Unmarshal(raw, &chunk); err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		for id, entry := range chunk {
			entries[id] = entry
		}
	}
	if len(entries) != 30 {
		t.Fatalf("catalogue entries = %d, want 30", len(entries))
	}
	for id, entry := range entries {
		meta := MetadataForID(id)
		if meta == nil {
			t.Errorf("missing metadata for %s", id)
			continue
		}
		if meta.ID != entry.ID || meta.Title != entry.Name {
			t.Errorf("metadata %s = (%q, %q), want (%q, %q)", id, meta.ID, meta.Title, entry.ID, entry.Name)
		}
		if meta.EffectiveMaturity() != "experimental" {
			t.Errorf("%s maturity = %q, want experimental", id, meta.EffectiveMaturity())
		}
	}
	for i := 1; i <= 30; i++ {
		id := fmt.Sprintf("PERF-PY-%d", i)
		if _, ok := entries[id]; !ok {
			t.Errorf("catalogue missing %s", id)
		}
	}
}
