package core_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
)

func TestParseLanguage(t *testing.T) {
	cases := []struct {
		in   string
		want core.LanguageID
		ok   bool
	}{
		{"go", core.LanguageGo, true},
		{"Go", core.LanguageGo, true},
		{"  python ", core.LanguagePython, true},
		{"py", core.LanguagePython, true},
		{"ruby", 0, false},
		{"", 0, false},
	}
	for _, tc := range cases {
		got, ok := core.ParseLanguage(tc.in)
		if ok != tc.ok || (ok && got != tc.want) {
			t.Fatalf("ParseLanguage(%q)=(%v,%v) want (%v,%v)", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestParseLanguages(t *testing.T) {
	got, err := core.ParseLanguages([]string{"go", "python", "go", "py"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != core.LanguageGo || got[1] != core.LanguagePython {
		t.Fatalf("got %v", got)
	}

	if _, err := core.ParseLanguages(nil); err == nil {
		t.Fatal("expected empty list error")
	}
	if _, err := core.ParseLanguages([]string{"", "  "}); err == nil {
		t.Fatal("expected all-blank list error")
	}
	if _, err := core.ParseLanguages([]string{"ruby"}); err == nil {
		t.Fatal("expected unknown language error")
	}
}

func TestDefaultEnabledLanguages(t *testing.T) {
	def := core.DefaultEnabledLanguages()
	if len(def) != 1 || def[0] != core.LanguageGo {
		t.Fatalf("default=%v want [go]", def)
	}
}

func TestLanguageFromExtension(t *testing.T) {
	if id, ok := core.LanguageFromExtension("go"); !ok || id != core.LanguageGo {
		t.Fatalf("go: %v %v", id, ok)
	}
	if id, ok := core.LanguageFromExtension("py"); !ok || id != core.LanguagePython {
		t.Fatalf("py: %v %v", id, ok)
	}
	if _, ok := core.LanguageFromExtension("rs"); ok {
		t.Fatal("rs should be unknown")
	}
}

func TestLanguageEnabled(t *testing.T) {
	enabled := []core.LanguageID{core.LanguageGo}
	if !core.LanguageEnabled(enabled, core.LanguageGo) {
		t.Fatal("go should be enabled")
	}
	if core.LanguageEnabled(enabled, core.LanguagePython) {
		t.Fatal("python should not be enabled")
	}
}
