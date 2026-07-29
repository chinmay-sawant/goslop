package cwe_test

import (
	"testing"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/lang/go/detectors/cwe"
	"github.com/chinmay/codehound/internal/rules"
)

func TestCWE89VulnerableIdentifier(t *testing.T) {
	src := `package sample

import (
	"database/sql"
	"net/http"
)

func LookupUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	_ = db.Query(name)
}
`
	findings := runCWE89(t, src, "CWE-89-taint-vulnerable.go")
	if !hasRule(findings, "CWE-89") {
		t.Fatalf("expected CWE-89 finding, got %#v", findings)
	}
}

func TestCWE89VulnerableConcat(t *testing.T) {
	src := `package sample

import (
	"database/sql"
	"net/http"
)

func LookupUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	_ = db.Query("SELECT * FROM users WHERE name = '" + name + "'")
}
`
	findings := runCWE89(t, src, "CWE-89-concat.go")
	if !hasRule(findings, "CWE-89") {
		t.Fatalf("expected CWE-89 on concat SQL, got %#v", findings)
	}
}

func TestCWE89VulnerableSprintf(t *testing.T) {
	src := `package sample

import (
	"database/sql"
	"fmt"
	"net/http"
)

func LookupUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	q := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)
	_ = db.Query(q)
}
`
	findings := runCWE89(t, src, "CWE-89-sprintf-var.go")
	if !hasRule(findings, "CWE-89") {
		t.Fatalf("expected CWE-89 via intermediate sprintf var, got %#v", findings)
	}

	src2 := `package sample
import ("database/sql"; "fmt"; "net/http")
func LookupUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	_ = db.Query(fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name))
}
`
	findings2 := runCWE89(t, src2, "CWE-89-sprintf.go")
	if !hasRule(findings2, "CWE-89") {
		t.Fatalf("expected CWE-89 on fmt.Sprintf SQL, got %#v", findings2)
	}
}

func TestCWE89Safe(t *testing.T) {
	src := `package sample

import (
	"database/sql"
	"net/http"
)

func LookupUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	_ = r.URL.Query().Get("name")
	_ = db.Query("SELECT * FROM users WHERE name = 'alice'")
}
`
	findings := runCWE89(t, src, "CWE-89-taint-safe.go")
	if hasRule(findings, "CWE-89") {
		t.Fatalf("safe fixture should not emit CWE-89, got %#v", findings)
	}
}

func runCWE89(t *testing.T, src, path string) []rules.Finding {
	t.Helper()
	d := cwe.NewCWE89()
	unit := core.NewParsedUnit(core.LangGo, path, src)
	ctx := core.DefaultScanContext()
	var out []rules.Finding
	d.Run(ctx, unit, &out)
	return out
}
