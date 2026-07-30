package golang_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/engine"
	golang "github.com/chinmay-sawant/goslop/internal/lang/go"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestGoPluginRegistration(t *testing.T) {
	reg := engine.DefaultRegistry()
	p, ok := reg.Plugin(core.LangGo)
	if !ok {
		t.Fatal("Go plugin not registered")
	}
	exts := p.Extensions()
	if len(exts) != 1 || exts[0] != "go" {
		t.Fatalf("extensions = %v", exts)
	}
	ids := reg.RuleIDs()
	want := map[string]bool{
		"CWE-78": false, "CWE-89": false,
		"PERF-116": false, "PERF-6": false, "PERF-32": false,
	}
	for _, id := range ids {
		if _, ok := want[id]; ok {
			want[id] = true
		}
	}
	for id, seen := range want {
		if !seen {
			t.Errorf("missing rule %s in registry", id)
		}
	}
}

func TestSeedDetectorsEndToEnd(t *testing.T) {
	cases := []struct {
		rule string
		src  string
		fire bool
	}{
		{
			rule: "CWE-78",
			src: `package sample
import ("net/http"; "os/exec")
func H(w http.ResponseWriter, r *http.Request) {
	c := r.URL.Query().Get("cmd")
	_ = exec.Command("sh", "-c", c)
}`,
			fire: true,
		},
		{
			rule: "CWE-89",
			src: `package sample
import ("database/sql"; "net/http")
func H(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	n := r.FormValue("n")
	_ = db.Query(n)
}`,
			fire: true,
		},
		{
			rule: "PERF-116",
			src: `package sample
import "strings"
func Has(s, sub string) bool { return strings.Index(s, sub) != -1 }`,
			fire: true,
		},
		{
			rule: "PERF-116",
			src: `package sample
import "strings"
func Has(s, sub string) bool { return strings.Contains(s, sub) }`,
			fire: false,
		},
	}

	p := golang.NewPlugin()
	for _, tc := range cases {
		unit := core.NewParsedUnit(core.LangGo, "sample.go", tc.src)
		ctx := core.DefaultScanContext()
		var out []rules.Finding
		for _, d := range p.Detectors() {
			for _, id := range d.RuleIDs() {
				if id != tc.rule {
					continue
				}
				d.Run(ctx, unit, &out)
			}
		}
		got := false
		for _, f := range out {
			if f.RuleID == tc.rule {
				got = true
			}
		}
		if got != tc.fire {
			t.Errorf("rule %s fire=%v got=%v findings=%v", tc.rule, tc.fire, got, out)
		}
	}
}
