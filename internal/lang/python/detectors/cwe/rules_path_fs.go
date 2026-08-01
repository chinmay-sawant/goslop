package cwe

import (
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// No SourceIndex gates: all rule evidence is intentionally assessed from
	// parsed call arguments and a gate would risk omitting a valid spelling.
	RegisterRule("CWE-73", detectCWE73, &MetaCWE73)
	RegisterRule("CWE-59", detectCWE59, &MetaCWE59)
	RegisterRule("CWE-41", detectCWE41, &MetaCWE41)
	RegisterRule("CWE-276", detectCWE276, &MetaCWE276)
	RegisterRule("CWE-378", detectCWE378, &MetaCWE378)
	RegisterRule("CWE-426", detectCWE426, &MetaCWE426)
	RegisterRule("CWE-250", detectCWE250, &MetaCWE250)
	RegisterRule("CWE-494", detectCWE494, &MetaCWE494)
}

// CWE-73 reports only direct framework request input at a filesystem sink.
// A variable that might have been validated elsewhere is deliberately outside
// this source-only rule, as are basename/secure_filename-safe expressions.
func detectCWE73(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"open", "os.remove", "os.unlink", "os.rename", "shutil.move", "shutil.copy", "Path", "pathlib.Path"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			args := splitTopLevelArgs(call.ArgsText)
			if len(args) == 0 || !isDirectFilePathSource(args[0]) {
				continue
			}
			emitPathFSFinding(unit, &MetaCWE73, call.Start, "directly request-controlled path reaches a filesystem API", confidence82, out)
			return
		}
	}
}

// CWE-59 recognizes the explicit but race-prone islink/lexists check followed
// by name-based file access. It avoids treating ordinary file access as proof
// of a link-following weakness.
func detectCWE59(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	functions := facts.Functions()
	if len(functions) == 0 {
		if start := linkCheckThenUse(unit.Source, facts.Masked); start >= 0 {
			emitPathFSFinding(unit, &MetaCWE59, start, "symbolic-link check is followed by access to the same path name (check-then-use race)", confidence80, out)
		}
		return
	}
	for _, fn := range functions {
		masked := facts.codeMask(fn.body, fn.bodyStart)
		if start := linkCheckThenUse(fn.body, masked); start >= 0 {
			emitPathFSFinding(unit, &MetaCWE59, fn.bodyStart+start, "symbolic-link check is followed by access to the same path name (check-then-use race)", confidence80, out)
			return
		}
	}
}

func linkCheckThenUse(source, masked string) int {
	for _, check := range findCallsMasked(source, masked, "os.path.islink", "os.path.lexists") {
		checked := strings.TrimSpace(check.ArgsText)
		if checked == "" {
			continue
		}
		for _, use := range findCallsMasked(source, masked, "open", "os.remove", "os.unlink", "os.rename", "shutil.move", "shutil.copy") {
			if use.Start <= check.Start {
				continue
			}
			args := splitTopLevelArgs(use.ArgsText)
			if len(args) > 0 && strings.TrimSpace(args[0]) == checked {
				return check.Start
			}
		}
	}
	return -1
}

// CWE-41 is intentionally narrower than CWE-22: it needs both a direct
// request source and normpath at the same file-access expression. Canonical
// resolution (realpath/resolve) in that expression is a safe suppression.
func detectCWE41(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "open", "os.remove", "os.unlink", "Path", "pathlib.Path") {
		args := splitTopLevelArgs(call.ArgsText)
		if len(args) == 0 {
			continue
		}
		path := args[0]
		if strings.Contains(path, "os.path.normpath(") && isDirectFilePathSource(path) &&
			!strings.Contains(path, "os.path.realpath(") && !strings.Contains(path, ".resolve(") {
			emitPathFSFinding(unit, &MetaCWE41, call.Start, "request-controlled path is normalized with normpath before file access without canonical resolution", confidence78, out)
			return
		}
	}
}

func detectCWE276(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "os.umask") {
		if strings.TrimSpace(call.ArgsText) == "0" {
			emitPathFSFinding(unit, &MetaCWE276, call.Start, "process umask is disabled, allowing insecure default file permissions", confidence84, out)
			return
		}
	}
	for _, name := range []string{"os.chmod", ".chmod", "os.makedirs"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if hasWorldWritableMode(call.ArgsText) {
				emitPathFSFinding(unit, &MetaCWE276, call.Start, "filesystem permissions are explicitly world-writable", confidence86, out)
				return
			}
		}
	}
}

func detectCWE378(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if calls := findCalls(facts, unit.Source, "tempfile.mktemp"); len(calls) > 0 {
		emitPathFSFinding(unit, &MetaCWE378, calls[0].Start, "tempfile.mktemp creates an unreserved temporary pathname", confidence90, out)
	}
}

func detectCWE426(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, call := range findCalls(facts, unit.Source, "sys.path.insert", "sys.path.append") {
		args := splitTopLevelArgs(call.ArgsText)
		if isUntrustedSearchPath(call.Name, args) {
			emitPathFSFinding(unit, &MetaCWE426, call.Start, "untrusted path is added to Python's import search path", confidence84, out)
			return
		}
	}
}

func detectCWE250(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"os.setuid", "os.seteuid", "os.setgid", "os.setegid"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if strings.TrimSpace(call.ArgsText) == "0" {
				emitPathFSFinding(unit, &MetaCWE250, call.Start, "process explicitly switches to the root user or group", confidence88, out)
				return
			}
		}
	}
}

// CWE-494 reports a direct, same-expression download-and-execute sequence.
// Downloading data alone and executing locally authored code remain safe.
func detectCWE494(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"exec", "eval"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if isDownloadedCode(call.ArgsText) {
				emitPathFSFinding(unit, &MetaCWE494, call.Start, "downloaded HTTP response is executed without an integrity verification step", confidence92, out)
				return
			}
		}
	}
}

func isDirectFilePathSource(expr string) bool {
	compact := compactWhitespace(expr)
	if strings.Contains(compact, "secure_filename(") || strings.Contains(compact, "os.path.basename(") {
		return false
	}
	return isDirectRequestExpr(expr) || strings.Contains(compact, "input(") || strings.Contains(compact, "sys.argv[") || strings.Contains(compact, "os.environ[")
}

func hasWorldWritableMode(args string) bool {
	compact := compactWhitespace(args)
	return strings.Contains(compact, "0o777") || strings.Contains(compact, "0o666") || strings.Contains(compact, "stat.S_IWOTH")
}

func isUntrustedSearchPath(name string, args []string) bool {
	if name == "sys.path.insert" {
		if len(args) < 2 || strings.TrimSpace(args[0]) != "0" {
			return false
		}
		path := compactWhitespace(args[1])
		return path == "'.'" || path == `"."` || path == "os.getcwd()" || isDirectFilePathSource(args[1])
	}
	return len(args) > 0 && isDirectFilePathSource(args[0])
}

func isDownloadedCode(args string) bool {
	compact := strings.ToLower(compactWhitespace(args))
	if strings.Contains(compact, "urlopen(") && strings.Contains(compact, ".read(") {
		return true
	}
	return (strings.Contains(compact, "requests.get(") || strings.Contains(compact, "httpx.get(")) &&
		(strings.Contains(compact, ".text") || strings.Contains(compact, ".content"))
}

func emitPathFSFinding(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}
