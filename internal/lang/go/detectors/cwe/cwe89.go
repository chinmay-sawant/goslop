package cwe

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/lang/go/detectors/sourceutil"
	"github.com/chinmay/codehound/internal/rules"
)

const ruleCWE89 = "CWE-89"

// SQL sink method names we treat as query/exec entry points.
var sqlSinkMethods = []string{
	"Query",
	"Exec",
	"QueryRow",
	"QueryContext",
	"ExecContext",
	"QueryRowContext",
	"Raw", // GORM
}

// CWE89Detector flags string-concat / dynamic SQL into Query/Exec patterns.
// Seed structural/taint-lite heuristic (same-file).
type CWE89Detector struct {
	core.BaseDetector
}

// NewCWE89 returns a CWE-89 detector.
func NewCWE89() *CWE89Detector {
	return &CWE89Detector{}
}

// Language implements core.Detector.
func (d *CWE89Detector) Language() core.LanguageID { return core.LangGo }

// RuleIDs implements core.Detector.
func (d *CWE89Detector) RuleIDs() []string { return []string{ruleCWE89} }

// MetadataFor implements core.Detector.
func (d *CWE89Detector) MetadataFor(ruleID string) *rules.RuleMetadata {
	if ruleID == ruleCWE89 {
		return &MetaCWE89
	}
	return nil
}

// Run implements core.Detector.
func (d *CWE89Detector) Run(ctx *core.ScanContext, unit *core.ParsedUnit, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if ctx != nil && !ctx.Allows(ruleCWE89) {
		return
	}
	src := unit.Source
	tainted := sourceutil.FindTaintedIdents(src)
	file := unit.DisplayPath
	if file == "" {
		file = unit.Path
	}

	for _, call := range findSQLSinkCalls(src) {
		args := sourceutil.SplitTopLevelArgs(call.ArgsText)
		if len(args) == 0 {
			continue
		}
		// *Context forms take context as first arg; SQL is second.
		sqlArg := args[0]
		if strings.HasSuffix(call.Name, "Context") && len(args) >= 2 {
			sqlArg = args[1]
		}
		if !isDynamicSQLArg(sqlArg, tainted) {
			continue
		}
		line, col := unit.LineCol(call.Start)
		rules.PushFindingWithConfidence(
			&MetaCWE89,
			file,
			line,
			col,
			"user-controlled input reaches an SQL execution sink (heuristic; not full SQLi coverage)",
			0.7,
			out,
		)
	}
}

func isDynamicSQLArg(arg string, tainted map[string]struct{}) bool {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return false
	}
	// Pure string literal (parameterized query template) is safe.
	if sourceutil.IsPureStringLiteral(arg) {
		return false
	}
	// String concat / Sprintf into SQL.
	if strings.Contains(arg, "+") || strings.Contains(arg, "fmt.Sprintf") || strings.Contains(arg, "fmt.Sprintf(") {
		// still allow only if it looks like building SQL rather than unrelated expr
		return true
	}
	// Direct request source in the SQL argument.
	if sourceutil.HasRequestSource(arg) {
		return true
	}
	// Identifier assigned from request input.
	for name := range tainted {
		if sourceutil.ContainsIdent(arg, name) {
			return true
		}
	}
	// Simple bare identifier / expression that is not a literal — only when
	// that identifier is known tainted (already covered) or the whole arg is
	// exactly an identifier that we cannot prove safe. We deliberately do not
	// flag all dynamic args without taint signal (too noisy for seed).
	return false
}

type sqlCall struct {
	Name     string
	Start    int
	ArgsText string
}

// findSQLSinkCalls finds recv.Method( calls for SQL methods without requiring
// a known receiver name (db.Query, stmt.Exec, tx.QueryContext, ...).
func findSQLSinkCalls(source string) []sqlCall {
	var out []sqlCall
	for _, method := range sqlSinkMethods {
		needle := "." + method
		start := 0
		for {
			idx := strings.Index(source[start:], needle)
			if idx < 0 {
				break
			}
			abs := start + idx
			// method name must not be prefix of longer name
			afterName := abs + len(needle)
			if afterName < len(source) {
				r, _ := utf8.DecodeRuneInString(source[afterName:])
				if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
					start = afterName
					continue
				}
			}
			// left of '.' should be identifier/end of selector
			if abs == 0 {
				start = afterName
				continue
			}
			// skip whitespace then require '('
			j := afterName
			for j < len(source) && (source[j] == ' ' || source[j] == '\t' || source[j] == '\n' || source[j] == '\r') {
				j++
			}
			if j >= len(source) || source[j] != '(' {
				start = afterName
				continue
			}
			// expand start left for a simple receiver identifier (for location, use '.')
			closeAt, args := scanParen(source, j)
			if closeAt < 0 {
				start = afterName
				continue
			}
			out = append(out, sqlCall{
				Name:     method,
				Start:    abs + 1, // point at method name
				ArgsText: args,
			})
			start = closeAt + 1
		}
	}
	return out
}

func scanParen(source string, open int) (closeAt int, args string) {
	if open >= len(source) || source[open] != '(' {
		return -1, ""
	}
	depth := 0
	inStr := byte(0)
	escape := false
	for i := open; i < len(source); i++ {
		c := source[i]
		if inStr != 0 {
			if escape {
				escape = false
				continue
			}
			if inStr == '"' || inStr == '\'' {
				if c == '\\' {
					escape = true
					continue
				}
				if c == inStr {
					inStr = 0
				}
				continue
			}
			if inStr == '`' && c == '`' {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '`', '\'':
			inStr = c
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i, source[open+1 : i]
			}
		}
	}
	return -1, ""
}
