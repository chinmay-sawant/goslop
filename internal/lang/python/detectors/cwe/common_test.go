package cwe

import "testing"

func TestFindCallsSkipsCommentsAndStrings(t *testing.T) {
	source := "# exec(user_code)\nexample = \"eval(user_code)\"\ndoc = '''compile(user_code)'''\nexec(user_code)\n"
	calls := findCalls(source, "exec", "eval", "compile")
	if len(calls) != 1 || calls[0].Name != "exec" {
		t.Fatalf("findCalls() = %#v, want only executable exec call", calls)
	}
}
