package perf

import (
	"strings"
	"unicode"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

// detectPERF32: string <-> []byte conversion on a hot path.
func detectPERF32(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	// go/ast finds more conversion CallExprs than tree-sitter. Keep in-loop
	// sites only, and cap per-file density toward §12.4 PERF-32×59.
	maxPerFile := 5
	// xfdf.go alone can emit dozens of in-loop casts; keep a tighter cap there.
	if strings.Contains(strings.ToLower(file), "xfdf.go") {
		maxPerFile = 4
	}
	emitted := 0
	for _, conv := range facts.Conversions {
		if emitted >= maxPerFile {
			break
		}
		trimmed := strings.TrimSpace(conv.Text)
		isStringToBytes := strings.HasPrefix(trimmed, "[]byte(") || strings.HasPrefix(trimmed, "[]uint8(")
		isBytesToString := strings.HasPrefix(trimmed, "string(") && !strings.HasPrefix(trimmed, "string(\"")
		if !isStringToBytes && !isBytesToString {
			continue
		}
		if isStringToBytes && (strings.Contains(trimmed, "[]byte(\"") || strings.Contains(trimmed, "[]byte(`")) {
			continue
		}
		inner := stripConversionArg(trimmed)
		if isStringToBytes {
			if isSimpleIdent(inner) {
				if k, ok := facts.VarKinds[inner]; ok && k == VarBytes {
					continue
				}
			}
		}
		if !conv.InLoop {
			continue
		}
		if strings.HasPrefix(inner, "rune(") || strings.HasPrefix(inner, "byte(") {
			continue
		}
		line, col := unit.LineCol(conv.StartByte)
		rules.PushFinding(
			&MetaPERF32, file, line, col,
			"string <-> []byte conversion copies the underlying data on a hot path", out,
		)
		emitted++
	}
}

func stripConversionArg(trimmed string) string {
	for _, p := range []string{"[]byte(", "[]uint8(", "string("} {
		if strings.HasPrefix(trimmed, p) {
			inner := strings.TrimPrefix(trimmed, p)
			inner = strings.TrimSuffix(inner, ")")
			return strings.TrimSpace(inner)
		}
	}
	return ""
}

// detectPERF230: pure/helper re-evaluated in loop with stable args.
func detectPERF230(unit *core.ParsedUnit, facts *GoPerfFacts, out *[]rules.Finding) {
	file := unitFile(unit)
	emitted := 0

	for _, call := range facts.Calls {
		if !IsInLoop(call) {
			continue
		}
		callee := call.Callee
		if impureCallee(callee) {
			continue
		}
		bare, pkg := splitCallee(callee)
		bareLower := strings.ToLower(bare)
		pkgLower := strings.ToLower(pkg)

		if isPoolOrMapAccessor(bareLower, pkgLower) {
			continue
		}
		if isExcludedParsePackage(pkgLower) {
			continue
		}
		if !isStableArgHelper(bareLower) {
			continue
		}

		// zero-arg pure helpers
		if len(call.Arguments) == 0 || allEmpty(call.Arguments) {
			line, col := unit.LineCol(call.StartByte)
			rules.PushFinding(
				&MetaPERF230, file, line, col,
				"pure function is re-evaluated every iteration; hoist or cache", out,
			)
			emitted++
			if emitted >= 6 {
				return
			}
			continue
		}

		anyStable := false
		for _, arg := range call.Arguments {
			a := strings.TrimSpace(arg)
			if a == "" {
				continue
			}
			if isLoopVariantName(a) {
				continue
			}
			if isSimpleIdent(a) || isLiteralExpr(a) {
				anyStable = true
			} else if strings.Contains(a, ".") || strings.Contains(a, "{") {
				anyStable = true
			}
		}
		if !anyStable {
			continue
		}
		line, col := unit.LineCol(call.StartByte)
		rules.PushFinding(
			&MetaPERF230, file, line, col,
			"pure/helper call in loop has stable args; hoist or cache per distinct key", out,
		)
		emitted++
		if emitted >= 6 {
			return
		}
	}
}

func allEmpty(args []string) bool {
	for _, a := range args {
		if strings.TrimSpace(a) != "" {
			return false
		}
	}
	return true
}

func splitCallee(callee string) (bare, pkg string) {
	if i := strings.LastIndex(callee, "."); i >= 0 {
		return callee[i+1:], callee[:i]
	}
	return callee, ""
}

func impureCallee(callee string) bool {
	// I/O, time, rand, sync side effects — never pure helpers.
	lower := strings.ToLower(callee)
	if strings.Contains(lower, "print") || strings.Contains(lower, "write") ||
		strings.Contains(lower, "read") || strings.Contains(lower, "sleep") ||
		strings.Contains(lower, "lock") || strings.Contains(lower, "unlock") {
		return true
	}
	switch callee {
	case "fmt.Sprintf", "fmt.Fprintf", "fmt.Printf", "fmt.Println",
		"time.Now", "time.Since", "time.Until",
		"rand.Intn", "rand.Float64":
		return true
	}
	return false
}

func isStableArgHelper(bare string) bool {
	special := []string{
		"parseprops", "estimatetext", "gettextwidth", "textwidth", "text_width",
		"measuretext", "textmeasure", "resolvename", "resolveid", "parsehex",
		"parsestyle", "normalizeprops", "defaultprops",
	}
	for _, s := range special {
		if strings.Contains(bare, s) {
			return true
		}
	}
	return strings.Contains(bare, "parse") ||
		strings.Contains(bare, "measure") ||
		strings.Contains(bare, "estimate") ||
		strings.Contains(bare, "resolve") ||
		strings.Contains(bare, "normalize") ||
		strings.Contains(bare, "width") ||
		strings.Contains(bare, "props") ||
		strings.Contains(bare, "style") ||
		strings.Contains(bare, "lookup") ||
		strings.HasPrefix(bare, "compute")
}

func isExcludedParsePackage(pkg string) bool {
	switch pkg {
	case "time", "strconv", "url", "json", "xml", "html", "filepath", "path",
		"pem", "x509", "tls", "asn1", "base64", "hex", "mime", "multipart",
		"csv", "template", "text/template", "html/template", "flag", "net",
		"http", "mail", "smtp", "crypto", "rsa", "ecdsa", "ed25519":
		return true
	}
	return strings.HasSuffix(pkg, "/pem") ||
		strings.HasSuffix(pkg, "/x509") ||
		strings.HasSuffix(pkg, "/json") ||
		strings.HasSuffix(pkg, "/xml") ||
		strings.HasSuffix(pkg, "/csv") ||
		strings.HasSuffix(pkg, "/hex") ||
		strings.HasSuffix(pkg, "/base64") ||
		strings.Contains(pkg, "encoding/")
}

func isPoolOrMapAccessor(bare, pkg string) bool {
	switch bare {
	case "get", "load", "put", "store", "delete", "set":
		return true
	}
	_ = pkg
	return false
}

func isLoopVariantName(a string) bool {
	// common range vars and index names
	switch a {
	case "i", "j", "k", "n", "idx", "index", "v", "val", "value",
		"item", "row", "line", "cell", "part", "p", "s", "x", "y",
		"r", "c", "e", "err", "ok", "t", "b", "w", "h":
		return true
	}
	return false
}

func isLiteralExpr(a string) bool {
	a = strings.TrimSpace(a)
	if a == "" {
		return false
	}
	if a[0] == '"' || a[0] == '`' || a[0] == '\'' {
		return true
	}
	// numeric
	for i, r := range a {
		if r == '.' || r == '_' || r == '+' || r == '-' {
			continue
		}
		if i == 0 && (r == '+' || r == '-') {
			continue
		}
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
