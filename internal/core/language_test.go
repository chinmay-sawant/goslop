package core_test

import (
	"strings"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
)

func TestLanguageIDString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		id   core.LanguageID
		want string
	}{
		{core.LanguageGo, "go"},
		{core.LangGo, "go"},
		{core.LanguagePython, "python"},
		{core.LangPython, "python"},
		{core.LanguageID(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.id.String(); got != tc.want {
			t.Errorf("%v.String() = %q, want %q", tc.id, got, tc.want)
		}
	}
}

func TestParseLanguage(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want core.LanguageID
		ok   bool
	}{
		{"go", core.LanguageGo, true},
		{"GO", core.LanguageGo, true},
		{" Go ", core.LanguageGo, true},
		{"python", core.LanguagePython, true},
		{"Python", core.LanguagePython, true},
		{"py", core.LanguagePython, true},
		{"PY", core.LanguagePython, true},
		{" py ", core.LanguagePython, true},
		{"", 0, false},
		{"rust", 0, false},
		{"golang", 0, false},
		{"python3", 0, false},
	}
	for _, tc := range cases {
		got, ok := core.ParseLanguage(tc.in)
		if ok != tc.ok || (ok && got != tc.want) {
			t.Errorf("ParseLanguage(%q) = (%v, %v), want (%v, %v)", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestLanguageFromExtension(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want core.LanguageID
		ok   bool
	}{
		{"go", core.LanguageGo, true},
		{"GO", core.LanguageGo, true},
		{".go", core.LanguageGo, true},
		{"py", core.LanguagePython, true},
		{"PY", core.LanguagePython, true},
		{".py", core.LanguagePython, true},
		{" py ", core.LanguagePython, true},
		{"", 0, false},
		{"rs", 0, false},
		{"txt", 0, false},
		{"python", 0, false}, // extension map is not name map
	}
	for _, tc := range cases {
		got, ok := core.LanguageFromExtension(tc.in)
		if ok != tc.ok || (ok && got != tc.want) {
			t.Errorf("LanguageFromExtension(%q) = (%v, %v), want (%v, %v)", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestDefaultEnabledLanguages(t *testing.T) {
	t.Parallel()
	got := core.DefaultEnabledLanguages()
	if len(got) != 1 || got[0] != core.LanguageGo {
		t.Fatalf("DefaultEnabledLanguages() = %v, want [go]", got)
	}
	// Second call must not share a mutated slice side effect (value identity ok,
	// but content must remain Go-only after caller mutates a returned copy).
	got[0] = core.LanguagePython
	again := core.DefaultEnabledLanguages()
	if len(again) != 1 || again[0] != core.LanguageGo {
		t.Fatalf("DefaultEnabledLanguages mutated after prior call: %v", again)
	}
}

func TestParseLanguages(t *testing.T) {
	t.Parallel()

	t.Run("stable order and aliases", func(t *testing.T) {
		t.Parallel()
		got, err := core.ParseLanguages([]string{"python", "go", "py"})
		if err != nil {
			t.Fatal(err)
		}
		want := []core.LanguageID{core.LanguagePython, core.LanguageGo}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("got %v, want %v", got, want)
			}
		}
	})

	t.Run("dedupe same language", func(t *testing.T) {
		t.Parallel()
		got, err := core.ParseLanguages([]string{"go", "GO", "go"})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0] != core.LanguageGo {
			t.Fatalf("got %v, want [go]", got)
		}
	})

	t.Run("skip blanks", func(t *testing.T) {
		t.Parallel()
		got, err := core.ParseLanguages([]string{"", "  ", "go", ""})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0] != core.LanguageGo {
			t.Fatalf("got %v, want [go]", got)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		got, err := core.ParseLanguages(nil)
		if err != nil || got != nil {
			t.Fatalf("nil names: got (%v, %v)", got, err)
		}
		got, err = core.ParseLanguages([]string{})
		if err != nil || got != nil {
			t.Fatalf("empty names: got (%v, %v)", got, err)
		}
		got, err = core.ParseLanguages([]string{"", "  "})
		if err != nil || got != nil {
			t.Fatalf("blank-only: got (%v, %v)", got, err)
		}
	})

	t.Run("unknown token", func(t *testing.T) {
		t.Parallel()
		got, err := core.ParseLanguages([]string{"go", "rust"})
		if err == nil {
			t.Fatalf("expected error, got %v", got)
		}
		if got != nil {
			t.Fatalf("partial result must be nil on error, got %v", got)
		}
		if !strings.Contains(err.Error(), "rust") {
			t.Errorf("error should mention unknown token: %v", err)
		}
	})
}
