package cwe

import "testing"

func TestNewAndFormat(t *testing.T) {
	r := New(89, "SQL Injection", "")
	if r.ID != "CWE-89" {
		t.Fatalf("ID = %q", r.ID)
	}
	if r.NumericID() != 89 {
		t.Fatalf("NumericID = %d", r.NumericID())
	}
	if r.String() != "CWE-89" {
		t.Fatalf("String = %q", r.String())
	}
	r2 := NewFromID("cwe-22")
	if r2.ID != "CWE-22" {
		t.Fatalf("NewFromID = %q", r2.ID)
	}
	r3 := NewFromID("78")
	if r3.ID != "CWE-78" {
		t.Fatalf("NewFromID numeric = %q", r3.ID)
	}
}

func TestFormatList(t *testing.T) {
	if FormatList(nil) != "" {
		t.Fatal("empty")
	}
	got := FormatList([]CweRef{
		New(89, "SQL Injection", ""),
		New(22, "Path Traversal", ""),
	})
	want := "CWE-22 (Path Traversal), CWE-89 (SQL Injection)"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRefsFromIDs(t *testing.T) {
	refs := RefsFromIDs([]string{"CWE-78", "89"})
	if len(refs) != 2 || refs[0].ID != "CWE-78" || refs[1].ID != "CWE-89" {
		t.Fatalf("%+v", refs)
	}
}
