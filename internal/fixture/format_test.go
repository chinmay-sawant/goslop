package fixture_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/fixture"
)

func TestParseFixture_CWE22Vulnerable(t *testing.T) {
	path := filepath.Join("..", "..", "tests", "fixtures", "go", "taint", "CWE-22-vulnerable.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	fx, err := fixture.ParseFixture(string(data), path)
	if err != nil {
		t.Fatalf("ParseFixture: %v", err)
	}
	if fx.Language != fixture.LangGo {
		t.Errorf("language = %q, want go", fx.Language)
	}
	if fx.Filename != "CWE-22-taint-vulnerable.go" {
		t.Errorf("filename = %q, want CWE-22-taint-vulnerable.go", fx.Filename)
	}
	if !strings.Contains(fx.Source, "package sample") {
		t.Errorf("source missing package sample:\n%s", fx.Source)
	}
	if !strings.Contains(fx.Source, "os.Open") {
		t.Errorf("source missing os.Open")
	}
	if strings.HasPrefix(fx.Source, "\n") {
		t.Errorf("source should not start with leading newline")
	}
}

func TestParseFixture_CWE22Safe(t *testing.T) {
	path := filepath.Join("..", "..", "tests", "fixtures", "go", "taint", "CWE-22-safe.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	fx, err := fixture.ParseFixture(string(data), path)
	if err != nil {
		t.Fatalf("ParseFixture: %v", err)
	}
	if fx.Filename != "CWE-22-taint-safe.go" {
		t.Errorf("filename = %q", fx.Filename)
	}
	if !strings.Contains(fx.Source, "filepath.Base") {
		t.Errorf("safe fixture should use filepath.Base")
	}
}

func TestParseFixture_PERF001(t *testing.T) {
	path := filepath.Join("..", "..", "tests", "fixtures", "go", "perf", "PERF-001-vulnerable.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	fx, err := fixture.ParseFixture(string(data), path)
	if err != nil {
		t.Fatalf("ParseFixture: %v", err)
	}
	if fx.Filename != "PERF-001-vulnerable.go" {
		t.Errorf("filename = %q", fx.Filename)
	}
	if !strings.Contains(fx.Source, "regexp.MustCompile") {
		t.Errorf("source missing regexp.MustCompile")
	}
}

func TestParseFixture_DefaultFilename(t *testing.T) {
	text := "lang: go\n---\npackage main\n"
	fx, err := fixture.ParseFixture(text, "/tmp/my_sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	if fx.Filename != "my_sample.go" {
		t.Errorf("default filename = %q, want my_sample.go", fx.Filename)
	}
}

func TestParseFixture_PythonAlias(t *testing.T) {
	text := "language: py\nfile: x.py\n---\nprint(1)\n"
	fx, err := fixture.ParseFixture(text, "x.txt")
	if err != nil {
		t.Fatal(err)
	}
	if fx.Language != fixture.LangPython {
		t.Errorf("language = %q, want python", fx.Language)
	}
}

func TestParseFixture_IgnoresCommentsAndUnknownKeys(t *testing.T) {
	text := `# comment
lang: go
variant: taint
expect: CWE-22
file: ok.go
---
package p
`
	fx, err := fixture.ParseFixture(text, "ok.txt")
	if err != nil {
		t.Fatal(err)
	}
	if fx.Filename != "ok.go" {
		t.Errorf("filename = %q", fx.Filename)
	}
}

func TestParseFixture_Errors(t *testing.T) {
	t.Run("missing separator", func(t *testing.T) {
		_, err := fixture.ParseFixture("lang: go\npackage main\n", "a.txt")
		if err == nil || !strings.Contains(err.Error(), "separator") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("missing lang", func(t *testing.T) {
		_, err := fixture.ParseFixture("file: a.go\n---\nx\n", "a.txt")
		if err == nil || !strings.Contains(err.Error(), "lang") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("unknown lang", func(t *testing.T) {
		_, err := fixture.ParseFixture("lang: rust\n---\nx\n", "a.txt")
		if err == nil || !strings.Contains(err.Error(), "unknown fixture language") {
			t.Fatalf("err = %v", err)
		}
	})
}
