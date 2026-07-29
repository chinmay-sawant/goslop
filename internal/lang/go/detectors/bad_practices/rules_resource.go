package badpractices

import (
	goast "go/ast"
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

func init() {
	RegisterRule("BP-95", detectBP95)
	RegisterRule("BP-96", detectBP96)
	RegisterRule("BP-97", detectBP97)
	RegisterRule("BP-98", detectBP98)
	RegisterRule("BP-99", detectBP99)
	RegisterRule("BP-126", detectBP126)
	RegisterRule("BP-128", detectBP128)
	RegisterRule("BP-131", detectBP131)
	RegisterRule("BP-132", detectBP132)
	RegisterRule("BP-133", detectBP133)
	RegisterRule("BP-134", detectBP134)
	RegisterRule("BP-135", detectBP135)
	RegisterRule("BP-136", detectBP136)
	RegisterRule("BP-140", detectBP140)
	RegisterRule("BP-142", detectBP142)
	RegisterRule("BP-143", detectBP143)
	RegisterRule("BP-145", detectBP145)
}

func detectBP95(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-95")
	src := unit.Source
	// HTTP client call without Body.Close
	clientCalls := []string{"http.Get(", "http.Post(", "http.Head(", "http.Do(", ".Do("}
	hasClient := false
	for _, c := range clientCalls {
		if strings.Contains(src, c) {
			hasClient = true
			break
		}
	}
	if !hasClient {
		return
	}
	// need response body usage pattern
	if !strings.Contains(src, "resp") && !strings.Contains(src, "Response") {
		// still may assign
	}
	if strings.Contains(src, "Body.Close") || strings.Contains(src, "defer resp.Body.Close") || strings.Contains(src, ".Body.Close(") {
		return
	}
	// if we got a response variable
	if strings.Contains(src, "http.Get(") || strings.Contains(src, "http.Post(") ||
		strings.Contains(src, "Client.Do(") || strings.Contains(src, "DefaultClient.Do(") ||
		strings.Contains(src, ".Do(req") {
		pos := indexOfIdent(src, "http.Get")
		if pos < 0 {
			pos = indexOfIdent(src, "http.Post")
		}
		if pos < 0 {
			pos = strings.Index(src, ".Do(")
		}
		if pos < 0 {
			pos = 0
		}
		pushAt(unit, meta, pos, "HTTP response body is never closed; defer resp.Body.Close()", out)
	}
}

func detectBP96(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-96")
	if !strings.Contains(unit.Source, ".Query(") && !strings.Contains(unit.Source, ".QueryContext(") {
		return
	}
	if strings.Contains(unit.Source, ".Close()") || strings.Contains(unit.Source, "rows.Close") {
		return
	}
	if strings.Contains(unit.Source, "Query(") || strings.Contains(unit.Source, "QueryContext(") {
		if pos := strings.Index(unit.Source, ".Query"); pos >= 0 {
			pushAt(unit, meta, pos, "sql.Rows result is never closed; defer rows.Close()", out)
		}
	}
}

func detectBP97(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-97")
	// bufio.Writer without Flush before read
	if strings.Contains(unit.Source, "bufio.NewWriter") && strings.Contains(unit.Source, ".Read") && !strings.Contains(unit.Source, ".Flush(") {
		if pos := strings.Index(unit.Source, "bufio.NewWriter"); pos >= 0 {
			pushAt(unit, meta, pos, "writer is read from without an intervening Flush", out)
		}
	}
}

// BP-98: local os.Open/os.OpenFile result neither closed nor returned (Rust parity).
// Per-function: only the assigned identifier is tracked; same-function Close or
// a return that mentions the name is ownership transfer. os.Create is gated for
// short-circuit only (Rust does not flag Create callees).
func detectBP98(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-98")
	src := unit.Source
	if !(strings.Contains(src, "os.Open") || strings.Contains(src, "os.OpenFile") || strings.Contains(src, "os.Create")) {
		return
	}
	msg := "opened file is neither closed nor transferred to the caller"
	if facts != nil && facts.tree != nil && facts.tree.File != nil {
		tree := facts.tree
		goast.Inspect(tree.File, func(n goast.Node) bool {
			var body *goast.BlockStmt
			switch x := n.(type) {
			case *goast.FuncDecl:
				body = x.Body
			case *goast.FuncLit:
				body = x.Body
			default:
				return true
			}
			if body == nil {
				return true
			}
			bodyText := tree.NodeText(body)
			// Collect assign statements that open files within this function body only
			// (not nested function literals — those are visited separately).
			goast.Inspect(body, func(inner goast.Node) bool {
				// Do not descend into nested function bodies; they get their own visit.
				switch inner.(type) {
				case *goast.FuncLit:
					return false
				case *goast.FuncDecl:
					return false
				}
				as, ok := inner.(*goast.AssignStmt)
				if !ok {
					return true
				}
				for _, rhs := range as.Rhs {
					call, ok := rhs.(*goast.CallExpr)
					if !ok {
						continue
					}
					callee := tree.NodeText(call.Fun)
					if callee != "os.Open" && callee != "os.OpenFile" {
						continue
					}
					if len(as.Lhs) == 0 {
						continue
					}
					id, ok := as.Lhs[0].(*goast.Ident)
					if !ok || id.Name == "" || id.Name == "_" {
						continue
					}
					if hasCloseOrTransferBP98(bodyText, id.Name) {
						continue
					}
					pushAt(unit, meta, tree.Offset(call.Pos()), msg, out)
				}
				return true
			})
			return true
		})
		return
	}
	// Text fallback: function-scoped open without same-function close/return.
	for _, fn := range splitTopLevelFuncs(src) {
		if !strings.Contains(fn.body, "os.Open(") && !strings.Contains(fn.body, "os.OpenFile(") {
			continue
		}
		for _, open := range []string{"os.Open(", "os.OpenFile("} {
			pos := strings.Index(fn.body, open)
			if pos < 0 {
				continue
			}
			lineStart := strings.LastIndex(fn.body[:pos], "\n") + 1
			line := fn.body[lineStart:pos]
			name := ""
			if i := strings.Index(line, ","); i >= 0 {
				name = firstIdent(strings.TrimSpace(line[:i]))
			} else if i := strings.Index(line, ":="); i >= 0 {
				name = firstIdent(strings.TrimSpace(line[:i]))
			} else if i := strings.Index(line, "="); i >= 0 {
				name = firstIdent(strings.TrimSpace(line[:i]))
			}
			if name == "" || name == "_" {
				continue
			}
			if hasCloseOrTransferBP98(fn.body, name) {
				continue
			}
			pushAt(unit, meta, fn.start+pos, msg, out)
		}
	}
}

func hasCloseOrTransferBP98(body, name string) bool {
	if strings.Contains(body, name+".Close(") {
		return true
	}
	for _, line := range strings.Split(body, "\n") {
		t := strings.TrimSpace(line)
		if !strings.HasPrefix(t, "return ") {
			continue
		}
		rest := strings.TrimPrefix(t, "return ")
		for _, tok := range strings.FieldsFunc(rest, func(r rune) bool {
			return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_')
		}) {
			if tok == name {
				return true
			}
		}
	}
	return false
}

type funcChunk struct {
	start int
	body  string
}

func splitTopLevelFuncs(src string) []funcChunk {
	var out []funcChunk
	depth := 0
	i := 0
	for i < len(src) {
		if src[i] == '{' {
			depth++
			i++
			continue
		}
		if src[i] == '}' {
			if depth > 0 {
				depth--
			}
			i++
			continue
		}
		if depth == 0 && strings.HasPrefix(src[i:], "func ") {
			brace := strings.Index(src[i:], "{")
			if brace < 0 {
				break
			}
			start := i + brace
			d := 1
			j := start + 1
			for j < len(src) && d > 0 {
				if src[j] == '{' {
					d++
				} else if src[j] == '}' {
					d--
				}
				j++
			}
			out = append(out, funcChunk{start: start, body: src[start:j]})
			i = j
			continue
		}
		i++
	}
	return out
}

func detectBP99(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-99")
	if strings.Contains(unit.Source, ".Wait()") && strings.Contains(unit.Source, "sync.Cond") {
		// Cond.Wait without Locker held — hard; flag Wait without Lock nearby
		if !strings.Contains(unit.Source, ".L.Lock") && !strings.Contains(unit.Source, ".Lock()") {
			if pos := strings.Index(unit.Source, ".Wait()"); pos >= 0 {
				pushAt(unit, meta, pos, "sync.Cond.Wait without holding the associated locker", out)
			}
		}
	}
}

func detectBP126(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-126")
	if !strings.Contains(unit.Source, ".Begin(") && !strings.Contains(unit.Source, ".BeginTx(") {
		return
	}
	if strings.Contains(unit.Source, ".Commit(") || strings.Contains(unit.Source, ".Rollback(") {
		return
	}
	if pos := strings.Index(unit.Source, ".Begin"); pos >= 0 {
		pushAt(unit, meta, pos, "transaction started without Commit or Rollback", out)
	}
}

func detectBP128(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-128")
	if !strings.Contains(unit.Source, "QueryRow") || !strings.Contains(unit.Source, ".Scan(") {
		return
	}
	if strings.Contains(unit.Source, "sql.ErrNoRows") || strings.Contains(unit.Source, "ErrNoRows") {
		return
	}
	if pos := strings.Index(unit.Source, ".Scan("); pos >= 0 {
		pushAt(unit, meta, pos, "QueryRow.Scan error handling does not distinguish sql.ErrNoRows", out)
	}
}

func detectBP131(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-131")
	// Query used for INSERT/UPDATE/DELETE
	for _, line := range codeLines(unit.Source) {
		t := strings.ToUpper(line.text)
		if strings.Contains(line.text, ".Query(") || strings.Contains(line.text, ".QueryContext(") {
			if strings.Contains(t, "INSERT ") || strings.Contains(t, "UPDATE ") || strings.Contains(t, "DELETE ") {
				pushAt(unit, meta, line.byte, "use Exec for DML statements that do not return rows", out)
			}
		}
	}
}

func detectBP132(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-132")
	if !strings.Contains(unit.Source, ".Exec(") && !strings.Contains(unit.Source, "UPDATE") {
		return
	}
	if strings.Contains(strings.ToUpper(unit.Source), "UPDATE") && strings.Contains(unit.Source, ".Exec") {
		if !strings.Contains(unit.Source, "RowsAffected") {
			if pos := strings.Index(unit.Source, ".Exec"); pos >= 0 {
				pushAt(unit, meta, pos, "UPDATE/Exec ignores RowsAffected; optimistic locks may silently no-op", out)
			}
		}
	}
}

func detectBP133(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-133")
	if !strings.Contains(unit.Source, "gorm") && !strings.Contains(unit.Source, ".Find(") && !strings.Contains(unit.Source, ".Create(") {
		return
	}
	// db.Where(...).Find without .Error check
	if strings.Contains(unit.Source, ".Find(") || strings.Contains(unit.Source, ".Create(") || strings.Contains(unit.Source, ".Save(") {
		if !strings.Contains(unit.Source, ".Error") && !strings.Contains(unit.Source, "err :=") {
			pos := strings.Index(unit.Source, ".Find(")
			if pos < 0 {
				pos = strings.Index(unit.Source, ".Create(")
			}
			if pos >= 0 {
				pushAt(unit, meta, pos, "GORM chain result used without checking Error", out)
			}
		}
	}
}

func detectBP134(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-134")
	if !strings.Contains(unit.Source, ".First(") && !strings.Contains(unit.Source, ".Take(") {
		return
	}
	if strings.Contains(unit.Source, "ErrRecordNotFound") {
		return
	}
	pos := strings.Index(unit.Source, ".First(")
	if pos < 0 {
		pos = strings.Index(unit.Source, ".Take(")
	}
	if pos >= 0 {
		pushAt(unit, meta, pos, "GORM First/Take without ErrRecordNotFound handling", out)
	}
}

func detectBP135(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-135")
	if strings.Contains(unit.Source, "gorm.DB") || strings.Contains(unit.Source, "*gorm.DB") {
		// global db without Session
		if strings.Contains(unit.Source, "var db ") || strings.Contains(unit.Source, "var DB ") {
			if !strings.Contains(unit.Source, ".Session(") && !strings.Contains(unit.Source, ".WithContext(") {
				if pos := strings.Index(unit.Source, "gorm"); pos >= 0 {
					pushAt(unit, meta, pos, "global GORM DB used without Session/WithContext isolation", out)
				}
			}
		}
	}
}

func detectBP136(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-136")
	if !strings.Contains(unit.Source, "AutoMigrate") {
		return
	}
	// in handler-shaped file
	if strings.Contains(unit.Source, "http.ResponseWriter") || strings.Contains(unit.Source, "gin.Context") ||
		strings.Contains(unit.Source, "echo.Context") || strings.Contains(unit.Source, "fiber.Ctx") {
		if pos := strings.Index(unit.Source, "AutoMigrate"); pos >= 0 {
			pushAt(unit, meta, pos, "GORM AutoMigrate on the request path; migrate at startup instead", out)
		}
	}
}

func detectBP140(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-140")
	if !strings.Contains(unit.Source, "sqlx") && !strings.Contains(unit.Source, ".Get(") && !strings.Contains(unit.Source, ".Select(") {
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if (strings.Contains(t, ".Get(") || strings.Contains(t, ".Select(") || strings.Contains(t, ".NamedExec(")) &&
			!strings.Contains(t, "err") && !strings.Contains(t, ":=") {
			if strings.HasPrefix(t, "_ =") || (!strings.Contains(t, "=") && strings.Contains(t, "(")) {
				pushAt(unit, meta, line.byte, "sqlx retrieval error is ignored", out)
			}
		}
	}
}

func detectBP142(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-142")
	if strings.Contains(unit.Source, "sqlx.In(") && !strings.Contains(unit.Source, ".Rebind(") {
		if pos := strings.Index(unit.Source, "sqlx.In("); pos >= 0 {
			pushAt(unit, meta, pos, "sqlx.In query without Rebind for the driver bindvar style", out)
		}
	}
}

func detectBP143(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-143")
	// redis cmd without .Err()
	if strings.Contains(unit.Source, "redis.") || strings.Contains(unit.Source, ".Get(") {
		for _, line := range codeLines(unit.Source) {
			t := strings.TrimSpace(line.text)
			if (strings.Contains(t, ".Get(") || strings.Contains(t, ".Set(")) && strings.Contains(unit.Source, "redis") {
				if !strings.Contains(unit.Source, ".Err()") && !strings.Contains(t, "err") {
					// only fire once
					pushAt(unit, meta, line.byte, "Redis command result error is ignored", out)
					return
				}
			}
		}
	}
}

func detectBP145(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-145")
	if !strings.Contains(unit.Source, ".Acquire(") && !strings.Contains(unit.Source, "pool.Acquire") {
		return
	}
	if strings.Contains(unit.Source, ".Release(") || strings.Contains(unit.Source, "defer") && strings.Contains(unit.Source, "Release") {
		return
	}
	if pos := strings.Index(unit.Source, "Acquire"); pos >= 0 {
		pushAt(unit, meta, pos, "pgx pool connection acquired but never released", out)
	}
}
