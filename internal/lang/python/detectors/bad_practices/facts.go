package badpractices

import (
	"strings"

	"github.com/chinmay-sawant/goslop/internal/ast"
	"github.com/chinmay-sawant/goslop/internal/core"
)

// Shared needle table for BP-PY fast-paths.
var bpNeedles = []string{
	"except", "except:", "pass", "def ", "async def ",
	"assert ", "open(", ".open(",
	"subprocess.", "shell=True",
	"os.system", "os.popen",
	"pickle.", "cloudpickle.", "_pickle.",
	"yaml.load", "yaml.safe_load",
	"eval(", "exec(", "compile(",
	"password", "secret", "api_key", "token", "private_key", "SECRET_KEY",
	"debug=True", "DEBUG", "app.run", "secret_key",
	"app.config",
	// Batch 01-02: core extra, http, flask remainder
	"raise", "Exception", "BaseException",
	"import *",
	"requests.", "session.", "sess.", "timeout",
	"AsyncClient", "aclose",
	"request.form", "request.get_json", "request.json", "request.files",
	".route(", "errorhandler", "register_error_handler", "traceback",
	"send_file", "send_from_directory", "safe_join",
}

// bpFacts is a light fact bag for Python source-pattern BP detectors.
type bpFacts struct {
	Source string
	Index  ast.SourceIndex
	lines  []codeLine
}

func buildFacts(unit *core.ParsedUnit) *bpFacts {
	f := &bpFacts{}
	if unit == nil {
		return f
	}
	f.Source = unit.Source
	f.lines = buildCodeLines(unit.Source)
	f.Index = ast.Build(unit.Source, bpNeedles)
	return f
}

func (f *bpFacts) close() {}

func (f *bpFacts) has(needle string) bool {
	if f == nil {
		return false
	}
	if f.Index.Has(needle) {
		return true
	}
	return strings.Contains(f.Source, needle)
}

func (f *bpFacts) hasAny(needles ...string) bool {
	for _, n := range needles {
		if f.has(n) {
			return true
		}
	}
	return false
}
