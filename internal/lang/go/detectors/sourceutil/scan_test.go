package sourceutil

import (
	"reflect"
	"testing"
)

func TestFindCallsSkipsLongerIdentifiersAndHandlesNestedArguments(t *testing.T) {
	source := `
package sample

func run() {
	exec.Commander("ignore")
	exec.Command("sh", "-c", format(call(1, 2), ` + "`a,b`" + `))
}
`
	calls := FindCalls(source, "exec.Command")
	if len(calls) != 1 {
		t.Fatalf("calls=%#v", calls)
	}
	if got, want := SplitTopLevelArgs(calls[0].ArgsText), []string{"\"sh\"", "\"-c\"", "format(call(1, 2), `a,b`)"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("args=%q want=%q", got, want)
	}
}

func TestFindTaintedIdentsHandlesControlFlowAndPropagation(t *testing.T) {
	source := `
if name := r.URL.Query().Get("name"); name != "" {
	command := prefix + name
	_ = command
}
safe := "r.URL.Query().Get("
`
	tainted := FindTaintedIdents(source)
	for _, name := range []string{"name", "command"} {
		if _, ok := tainted[name]; !ok {
			t.Fatalf("missing tainted identifier %q: %#v", name, tainted)
		}
	}
}

func TestSourceHelpersRespectStringAndIdentifierEdges(t *testing.T) {
	if !IsPureStringLiteral(`"a\\\"b"`) || !IsPureStringLiteral("`a,b`") || IsPureStringLiteral(`"a" + name`) {
		t.Fatal("unexpected string literal classification")
	}
	if !ContainsIdent("run(name, name2)", "name") || ContainsIdent("name2", "name") || ContainsIdent("pkg.name", "na") {
		t.Fatal("unexpected identifier classification")
	}
	if got := SplitTopLevelArgs(`map[string][]string{"a,b": {"c"}}, fn("d,e"), x`); !reflect.DeepEqual(got, []string{`map[string][]string{"a,b": {"c"}}`, `fn("d,e")`, "x"}) {
		t.Fatalf("split=%q", got)
	}
}

func FuzzSplitTopLevelArgs(f *testing.F) {
	for _, seed := range []string{"a,b", `fn("a,b"), c`, "`a,b`, map[string]int{\"x\": 1}"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, args string) {
		_ = SplitTopLevelArgs(args)
		_ = FindCalls("exec.Command("+args+")", "exec.Command")
	})
}
