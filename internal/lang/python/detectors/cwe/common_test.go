package cwe

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/fixture"
)

func TestFindCallsSkipsCommentsAndStrings(t *testing.T) {
	contents, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "..", "tests", "fixtures", "python", "cwe", "CWE-94-safe.txt"))
	if err != nil {
		t.Fatal(err)
	}
	fx, err := fixture.ParseFixture(string(contents), "CWE-94-safe.txt")
	if err != nil {
		t.Fatal(err)
	}
	calls := findCalls(fx.Source, "eval")
	if len(calls) != 1 || calls[0].Name != "eval" {
		t.Fatalf("findCalls() = %#v, want only executable eval call", calls)
	}
}
