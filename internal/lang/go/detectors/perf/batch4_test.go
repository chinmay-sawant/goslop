package perf_test

import (
	"fmt"
	"testing"

	"github.com/chinmay/goslop/internal/lang/go/detectors/perf"
)

func TestBatch4RuleIDsRegistered(t *testing.T) {
	d := perf.NewGoPerfScan()
	ids := d.RuleIDs()
	set := map[string]bool{}
	for _, id := range ids {
		set[id] = true
	}
	want := []int{
		164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180,
		181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197,
		198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 209, 210, 211, 212, 213, 214,
	}
	for _, n := range want {
		id := fmt.Sprintf("PERF-%d", n)
		if !set[id] {
			t.Errorf("missing registered rule %s", id)
		}
	}
	if len(want) != 50 {
		t.Fatalf("batch4 should cover 50 rules, list has %d", len(want))
	}
}

func TestBatch4SelectedFixtures(t *testing.T) {
	cases := []struct {
		vuln, safe, rule string
	}{
		{"PERF-164-vulnerable.txt", "PERF-164-safe.txt", "PERF-164"},
		{"PERF-165-vulnerable.txt", "PERF-165-safe.txt", "PERF-165"},
		{"PERF-167-vulnerable.txt", "PERF-167-safe.txt", "PERF-167"},
		{"PERF-170-vulnerable.txt", "PERF-170-safe.txt", "PERF-170"},
		{"PERF-176-vulnerable.txt", "PERF-176-safe.txt", "PERF-176"},
		{"PERF-190-vulnerable.txt", "PERF-190-safe.txt", "PERF-190"},
		{"PERF-192-vulnerable.txt", "PERF-192-safe.txt", "PERF-192"},
		{"PERF-195-vulnerable.txt", "PERF-195-safe.txt", "PERF-195"},
		{"PERF-213-vulnerable.txt", "PERF-213-safe.txt", "PERF-213"},
		{"PERF-214-vulnerable.txt", "PERF-214-safe.txt", "PERF-214"},
		{"PERF-181-vulnerable.txt", "PERF-181-safe.txt", "PERF-181"},
		{"PERF-171-vulnerable.txt", "PERF-171-safe.txt", "PERF-171"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.rule, func(t *testing.T) {
			runFixtureRule(t, tc.vuln, tc.rule, true)
			if tc.safe != "" {
				runFixtureRule(t, tc.safe, tc.rule, false)
			}
		})
	}
}
