package perf_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/lang/go/detectors/perf"
)

func TestBatch3RuleIDsRegistered(t *testing.T) {
	d := perf.NewGoPerfScan()
	ids := d.RuleIDs()
	want := []string{
		"PERF-112", "PERF-113", "PERF-114", "PERF-115", "PERF-117",
		"PERF-118", "PERF-119", "PERF-120", "PERF-121", "PERF-122",
		"PERF-123", "PERF-124", "PERF-125", "PERF-126", "PERF-127",
		"PERF-128", "PERF-129", "PERF-130", "PERF-131", "PERF-132",
		"PERF-133", "PERF-134", "PERF-135", "PERF-137", "PERF-138",
		"PERF-139", "PERF-140", "PERF-141", "PERF-142", "PERF-143",
		"PERF-144", "PERF-145", "PERF-146", "PERF-147", "PERF-148",
		"PERF-149", "PERF-150", "PERF-151", "PERF-152", "PERF-153",
		"PERF-154", "PERF-155", "PERF-156", "PERF-157", "PERF-158",
		"PERF-159", "PERF-160", "PERF-161", "PERF-162", "PERF-163",
	}
	have := map[string]bool{}
	for _, id := range ids {
		have[id] = true
	}
	for _, id := range want {
		if !have[id] {
			t.Errorf("missing registered rule %s", id)
		}
		if d.MetadataFor(id) == nil {
			t.Errorf("missing metadata for %s", id)
		}
	}
}

func TestPERF112EqualFold(t *testing.T) {
	src := `package sample
import "strings"
func SameIgnoreCase(a, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}
`
	findings := runPerf(t, src, "PERF-112-vulnerable.go")
	if !hasRule(findings, "PERF-112") {
		t.Fatalf("expected PERF-112, got %#v", findings)
	}
	safe := `package sample
import "strings"
func SameIgnoreCase(a, b string) bool {
	return strings.EqualFold(a, b)
}
`
	if hasRule(runPerf(t, safe, "PERF-112-safe.go"), "PERF-112") {
		t.Fatalf("safe EqualFold should not fire PERF-112")
	}
}

func TestPERF114ManualCopy(t *testing.T) {
	src := `package sample
func Copy(src []int) []int {
	dst := make([]int, len(src))
	for i, v := range src {
		dst[i] = v
	}
	return dst
}
`
	if !hasRule(runPerf(t, src, "PERF-114-vulnerable.go"), "PERF-114") {
		t.Fatalf("expected PERF-114")
	}
}

func TestPERF119And128Appends(t *testing.T) {
	src := `package sample
func Build() []string {
	var out []string
	out = append(out, "a")
	out = append(out, "b")
	out = append(out, "c")
	return out
}
`
	findings := runPerf(t, src, "PERF-119-vulnerable.go")
	if !hasRule(findings, "PERF-119") {
		t.Fatalf("expected PERF-119, got %#v", findings)
	}
	if !hasRule(findings, "PERF-128") {
		t.Fatalf("expected PERF-128, got %#v", findings)
	}
	safe := `package sample
func Build() []string {
	out := make([]string, 0, 3)
	out = append(out, "a", "b", "c")
	return out
}
`
	sf := runPerf(t, safe, "PERF-119-safe.go")
	if hasRule(sf, "PERF-119") || hasRule(sf, "PERF-128") {
		t.Fatalf("merged append should be silent, got %#v", sf)
	}
}

func TestPERF133SortSliceInLoop(t *testing.T) {
	src := `package sample
import "sort"
func TopK(buckets [][]int, k int) [][]int {
	for i := range buckets {
		sort.Slice(buckets[i], func(a, b int) bool {
			return buckets[i][a] > buckets[i][b]
		})
	}
	return buckets
}
`
	if !hasRule(runPerf(t, src, "PERF-133-vulnerable.go"), "PERF-133") {
		t.Fatalf("expected PERF-133")
	}
}

func TestPERF141URLQueryRepeated(t *testing.T) {
	src := `package sample
import "net/http"
func Handle(w http.ResponseWriter, r *http.Request) {
	_ = w
	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	_ = page
	_ = size
}
`
	if !hasRule(runPerf(t, src, "PERF-141-vulnerable.go"), "PERF-141") {
		t.Fatalf("expected PERF-141")
	}
	safe := `package sample
import "net/http"
func Handle(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	_ = q.Get("page")
	_ = q.Get("size")
}
`
	if hasRule(runPerf(t, safe, "PERF-141-safe.go"), "PERF-141") {
		t.Fatalf("cached Query should be silent")
	}
}

func TestPERF142MaxBytesReader(t *testing.T) {
	src := `package sample
import (
	"io"
	"net/http"
)
func Handle(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	_, _ = w.Write(body)
}
`
	if !hasRule(runPerf(t, src, "PERF-142-vulnerable.go"), "PERF-142") {
		t.Fatalf("expected PERF-142")
	}
	safe := `package sample
import (
	"io"
	"net/http"
)
func Handle(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	body, _ := io.ReadAll(r.Body)
	_, _ = w.Write(body)
}
`
	if hasRule(runPerf(t, safe, "PERF-142-safe.go"), "PERF-142") {
		t.Fatalf("MaxBytesReader should be silent")
	}
}

func TestPERF158SortSliceBasic(t *testing.T) {
	src := `package sample
import "sort"
func Sort(xs []int) {
	sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
}
`
	if !hasRule(runPerf(t, src, "PERF-158-vulnerable.go"), "PERF-158") {
		t.Fatalf("expected PERF-158")
	}
}

func TestPERF161RowsErr(t *testing.T) {
	src := `package sample
import "database/sql"
func LoadUsers(db *sql.DB) error {
	rows, err := db.Query("SELECT id FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
	}
	return nil
}
`
	if !hasRule(runPerf(t, src, "PERF-161-vulnerable.go"), "PERF-161") {
		t.Fatalf("expected PERF-161")
	}
	safe := `package sample
import "database/sql"
func LoadUsers(db *sql.DB) error {
	rows, err := db.Query("SELECT id FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
	}
	return rows.Err()
}
`
	if hasRule(runPerf(t, safe, "PERF-161-safe.go"), "PERF-161") {
		t.Fatalf("rows.Err present should be silent")
	}
}

func TestPERF163QueryRow(t *testing.T) {
	src := `package sample
import "database/sql"
func FindUser(db *sql.DB, id int) (string, error) {
	rows, err := db.Query("SELECT name FROM users WHERE id = ?", id)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return "", err
		}
		return name, nil
	}
	return "", sql.ErrNoRows
}
`
	if !hasRule(runPerf(t, src, "PERF-163-vulnerable.go"), "PERF-163") {
		t.Fatalf("expected PERF-163")
	}
}

func TestPERF125NilAppend(t *testing.T) {
	src := `package sample
func Add(s []int, v int) []int {
	if s != nil {
		s = append(s, v)
	}
	return s
}
`
	if !hasRule(runPerf(t, src, "PERF-125-vulnerable.go"), "PERF-125") {
		t.Fatalf("expected PERF-125")
	}
}

func TestPERF120TimeSince(t *testing.T) {
	src := `package sample
import "time"
func Since(t time.Time) time.Duration {
	return time.Now().Sub(t)
}
`
	if !hasRule(runPerf(t, src, "PERF-120-vulnerable.go"), "PERF-120") {
		t.Fatalf("expected PERF-120")
	}
}
