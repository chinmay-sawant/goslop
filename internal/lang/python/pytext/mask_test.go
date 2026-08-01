package pytext

import "testing"

func TestMaskBlanksTripleQuotedAndComments(t *testing.T) {
	t.Parallel()
	src := "x = 1\n\"\"\"\npickle.loads(body)\n\"\"\"\n# pickle.loads(body)\ny = pickle.loads(body)\n"
	masked := Mask(src)
	if !containsCode(masked, "pickle.loads") {
		t.Fatalf("expected live call to remain visible in mask:\n%s", masked)
	}
	// The docstring and comment occurrences should be blanked.
	live := 0
	for i := 0; i+len("pickle.loads") <= len(masked); i++ {
		if masked[i:i+len("pickle.loads")] == "pickle.loads" {
			live++
		}
	}
	if live != 1 {
		t.Fatalf("want exactly 1 visible pickle.loads, got %d\n%s", live, masked)
	}
}

func containsCode(s, needle string) bool {
	return len(s) >= len(needle) && (s == needle || len(needle) == 0 ||
		(func() bool {
			for i := 0; i+len(needle) <= len(s); i++ {
				if s[i:i+len(needle)] == needle {
					return true
				}
			}
			return false
		})())
}
