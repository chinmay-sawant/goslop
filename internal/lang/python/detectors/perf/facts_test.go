package perf

import "testing"

func TestBuildCodeLinesBlanksTripleQuotedBodies(t *testing.T) {
	t.Parallel()
	src := "def f():\n    \"\"\"\n    Job.objects.filter(status=\"pending\").first()\n    \"\"\"\n    return 1\n"
	lines := buildCodeLines(src)
	joined := ""
	for _, line := range lines {
		joined += line.text
	}
	if containsFold(joined, ".objects.filter") {
		t.Fatalf("docstring ORM call should be blanked, got:\n%q", joined)
	}
}

func TestBuildCodeLinesKeepsOrdinaryStringKeywords(t *testing.T) {
	t.Parallel()
	src := "job = Job.objects.filter(status=\"pending\").first()\n"
	lines := buildCodeLines(src)
	if len(lines) == 0 || !containsFold(lines[0].text, "pending") {
		t.Fatalf("ordinary string keyword should remain for heuristics: %#v", lines)
	}
}

func containsFold(s, needle string) bool {
	return len(s) >= len(needle) && (func() bool {
		// simple contains; case-sensitive is fine for these needles
		for i := 0; i+len(needle) <= len(s); i++ {
			if s[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}
