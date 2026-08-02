package badpractices

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/lang/python/pytext"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

const (
	bareExceptClause  = "except:"
	exceptionName     = "Exception"
	baseExceptionName = "BaseException"
	maxSignatureLines = 30
)

func init() {
	RegisterRule("BP-PY-1", detectBPPY1)
	RegisterRule("BP-PY-2", detectBPPY2)
	RegisterRule("BP-PY-4", detectBPPY4)
	RegisterRule("BP-PY-6", detectBPPY6)
	RegisterRule("BP-PY-7", detectBPPY7)
}

// BP-PY-1: bare `except:` or broad `except Exception` / `except BaseException`
// with weak handling (pass / bare continue).
func detectBPPY1(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-1")
	// Offline tooling / release scripts: broad-except noise in batch jobs
	// (Project_Parva tools/, scripts/release) is audited residual FP.
	if isPythonOfflineScriptPath(unit) {
		return
	}
	if !facts.has("except") {
		return
	}
	// Mask string literals/docstrings so prose examples cannot be handlers
	// (safer __init__.py docstring FPs). Byte offsets stay stable, so raw
	// lines remain usable for exception-variable presence checks.
	maskedLines := buildCodeLines(pytext.Mask(unit.Source))
	rawLines := codeLinesFacts(facts, unit.Source)
	isTest := isPythonTestFile(unit)
	for i, line := range maskedLines {
		t := strings.TrimSpace(line.text)
		if t == "" {
			continue
		}
		// Bare except:
		if t == bareExceptClause || strings.HasPrefix(t, bareExceptClause) {
			// "except:" only (no type)
			rest := strings.TrimSpace(strings.TrimPrefix(t, "except"))
			if rest != ":" && !strings.HasPrefix(rest, ":") {
				continue
			}
		} else if !isBroadExcept(t) {
			continue
		}
		caughtVar := exceptClauseCaughtVar(t)
		if isTest && broadExceptCollectsTestEvidence(maskedLines, i, caughtVar) {
			continue
		}
		// Flag broad Exception/BaseException unless the suite surfaces the
		// failure (re-raise, set_exception, error payload transport, stderr
		// reporting with an exit/continue, or a quiet fallback with no
		// side effects). Log-only handlers, print-only handlers, and
		// record-into-field handlers stay reportable: they still hide the
		// distinct failure conditions the rule targets.
		if suiteSurfacesFailure(maskedLines, rawLines, i, caughtVar) {
			continue
		}
		pushAt(unit, meta, line.byte, "broad except Exception/BaseException hides failures; catch specific types or re-raise", out)
	}
}

// exceptClauseCaughtVar returns the name bound by `except ... as X`, or "".
func exceptClauseCaughtVar(t string) string {
	if i := strings.Index(t, " as "); i >= 0 {
		name := strings.TrimSpace(t[i+len(" as "):])
		if j := strings.Index(name, ":"); j >= 0 {
			name = strings.TrimSpace(name[:j])
		}
		if isSimpleIdent(name) {
			return name
		}
	}
	return ""
}

func broadExceptCollectsTestEvidence(lines []codeLine, exceptIdx int, caughtVar string) bool {
	if exceptIdx < 0 || exceptIdx >= len(lines) {
		return false
	}
	exceptIndent := indentWidth(lines[exceptIdx].raw)
	collects := false
	for j := exceptIdx + 1; j < len(lines); j++ {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		if indentWidth(lines[j].raw) <= exceptIndent {
			break
		}
		raw := lines[j].raw
		collects = collects || strings.Contains(t, ".append(")
		if caughtVar != "" {
			collects = collects || assignsRawVar(raw, caughtVar)
		}
	}
	if !collects {
		return false
	}
	for _, line := range lines[exceptIdx+1:] {
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "assert ") || strings.Contains(t, ".assert") {
			return true
		}
	}
	return false
}

func isBroadExcept(t string) bool {
	// except Exception: / except Exception as e: / except BaseException...
	t = strings.TrimSpace(t)
	if !strings.HasPrefix(t, "except ") {
		return false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(t, "except "))
	// Strip trailing :
	if i := strings.Index(rest, ":"); i >= 0 {
		rest = strings.TrimSpace(rest[:i])
	}
	// as name
	if i := strings.Index(rest, " as "); i >= 0 {
		rest = strings.TrimSpace(rest[:i])
	}
	return rest == exceptionName || rest == baseExceptionName
}

func suiteReraises(lines []codeLine, exceptIdx int) bool {
	return false
}

// suiteSurfacesFailure reports whether an except suite propagates or records
// the failure instead of swallowing it. This is the single reconciliation
// guard for the audited FP/TP corpus:
//
//   - fall-through to an explicit raise after the handler (calgebra
//     gcal.py:221/245 — message built, then RuntimeError raised) is safe;
//   - suite-level raise: a bare raise-only handler stays reportable
//     (WeThePeople tracing.py:76); a multi-statement cleanup-then-propagate
//     handler with no logging/printing is safe (niquests utils.py:256,
//     safer __init__.py:312, wasi/_adapter.py:107);
//   - set_exception / _error_result feed the failure forward (niquests
//     _ws.py:44, _sse.py:153; calgebra gcsa.py:440);
//   - a two-field error record (x.error = str(e); x.note = ...) is the
//     httptap analyzer.py:262 audited FP, while a single .error= assignment
//     stays reportable (WeThePeople query_analysis.py:294);
//   - a stderr print followed by continue/return/exit reports the failure to
//     the user (calgebra ical.py:384, sync-with-uv cli.py:161); any other
//     print-of-exception handler stays reportable (pyauto-desktop, Cronboard,
//     httpmorph, logxide, WeThePeople prints);
//   - a single bare print nested inside another except suite is the final
//     fallback of an error-handling ladder (requestSpeedTest
//     increase_limits.py:38);
//   - type-extraction feeding a counter is error accounting, not surfacing
//     (requestSpeedTest rnet_test.py:37);
//   - error=<var> / append(<var>) transport the raw exception in a result
//     payload (calgebra gcsa.py:1117; niquests sgi/_async/__init__.py:661);
//   - a single-statement raw capture x = e is test evidence capture
//     (onlymaps test_pool.py:63, test_query.py:78; niquests
//     sgi/_async/__init__.py:355);
//   - a single-statement pass-through return of the input value is a parse
//     fallback that loses nothing (onlymaps _types.py:189);
//   - a multi-statement suite that never references the caught variable and
//     neither logs nor prints is a documented fallback path (niquests
//     wasi/_async/_adapter.py:99).
func suiteSurfacesFailure(lines, rawLines []codeLine, exceptIdx int, caughtVar string) bool {
	suiteIdx, after, _ := exceptSuiteLineIdx(lines, exceptIdx)
	if len(suiteIdx) == 0 {
		return false
	}
	exceptIndent := indentWidth(lines[exceptIdx].raw)
	// Fall-through to an explicit raise outside the handler propagates
	// (calgebra gcal.py:221/245 — message built, then RuntimeError raised).
	// A pass-only suite followed by a distant raise stays reportable
	// (WeThePeople services/auth.py:396 — the 401 exit is unrelated flow).
	var body []string
	var bodyRaw []string
	allText := ""
	allRaw := ""
	for j := exceptIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		if indentWidth(lines[j].raw) <= exceptIndent {
			break
		}
		allText += mt + "\n"
		allRaw += strings.TrimSpace(rawLines[j].raw) + "\n"
	}
	// Merge multi-line statements (parenthesized returns, list comprehensions)
	// so continuation lines do not count as separate statements. Continuation
	// lines may sit at deeper indent than the direct body (calgebra
	// gcsa.py:1117 list comprehension), so the merge scans every line until
	// the brackets balance.
	merged := make([]bool, len(suiteIdx))
	for k := 0; k < len(suiteIdx); k++ {
		if merged[k] {
			continue
		}
		i := suiteIdx[k]
		mt := strings.TrimSpace(lines[i].text)
		rt := strings.TrimSpace(rawLines[i].raw)
		last := i
		if pyBracketsUnbalanced(mt) {
			for j := i + 1; j < len(lines); j++ {
				txt := strings.TrimSpace(lines[j].text)
				if txt == "" || indentWidth(lines[j].raw) <= exceptIndent {
					break
				}
				mt += " " + txt
				rt += " " + strings.TrimSpace(rawLines[j].raw)
				last = j
				if !pyBracketsUnbalanced(mt) {
					break
				}
			}
		}
		// Mark suite lines consumed by this statement's continuation.
		for m := k + 1; m < len(suiteIdx); m++ {
			if suiteIdx[m] > i && suiteIdx[m] <= last {
				merged[m] = true
			}
		}
		body = append(body, mt)
		bodyRaw = append(bodyRaw, rt)
	}
	stmtCount := len(body)
	joinedRaw := strings.Join(bodyRaw, "\n")

	// Fall-through to an explicit raise outside the handler propagates
	// (calgebra gcal.py:221/245 — message built, then RuntimeError raised).
	// A pass-only suite followed by a distant raise stays reportable
	// (WeThePeople services/auth.py:396 — the 401 exit is unrelated flow).
	if suiteNotPassOnly(body) && nextDedentedStartsWith(lines, after, exceptIndent, "raise") {
		return true
	}

	hasRaise := suiteLineRaiseRE.MatchString(allText)
	hasLog := strings.Contains(allText, "exc_info") || strings.Contains(allText, ".exception(") ||
		strings.Contains(allText, "log.") || strings.Contains(allText, "logger.") ||
		strings.Contains(allText, "logging.") || strings.Contains(allText, "traceback.print_exc")
	hasPrint := strings.Contains(allText, "print(") || strings.Contains(allText, "console.print(")
	hasExitFlow := false
	for _, rt := range bodyRaw {
		if strings.HasPrefix(rt, "continue") || strings.HasPrefix(rt, "break") ||
			strings.HasPrefix(rt, "return") || strings.Contains(rt, "exit(") ||
			strings.Contains(rt, "sys.exit(") {
			hasExitFlow = true
		}
	}
	// Nested returns under if/for inside the suite also exit the enclosing
	// function (WeThePeople orchestrator veritas_strict path).
	if !hasExitFlow && suiteLineReturnRE.MatchString(allText) {
		hasExitFlow = true
	}

	// Suite-level raise: raise-only handlers stay reportable; cleanup /
	// translation handlers that neither log nor print are the audited-FP
	// propagate shape. A single wrapped raise (raise ConnectionError(...) /
	// raise X from e) translates the failure (niquests pyodide bridge FPs);
	// a bare re-raise stays reportable (WeThePeople tracing.py:76).
	// log+re-raise is best practice and fully surfaces the failure
	// (Project_Parva middleware: logger.exception(...); raise).
	if hasRaise {
		if hasLog {
			return true
		}
		if stmtCount >= 2 && !hasPrint {
			return true
		}
		if stmtCount == 1 && !hasPrint && wrappedRaiseStatement(body[0]) {
			return true
		}
		return false
	}

	// Failure fed into a future / error payload helper.
	if strings.Contains(allText, "set_exception(") || strings.Contains(allText, "_error_result(") {
		return true
	}

	// Two-field error record (step.error = ...; step.note = ...) — the
	// httptap analyzer.py:262 audited FP. A lone .error= assignment stays
	// reportable (WeThePeople query_analysis.py:294).
	errorAttr := false
	for _, mt := range body {
		if k := strings.Index(mt, ".error"); k >= 0 {
			if strings.HasPrefix(strings.TrimSpace(mt[k+len(".error"):]), "=") {
				errorAttr = true
			}
		}
	}
	if errorAttr && stmtCount >= 2 {
		return true
	}

	// logging.Handler error contract: sole suite call to handleError /
	// _handle_error (logxide handlers.py / compat_handlers / sentry
	// integration audited FPs). Does not match log-and-return TPs.
	if suiteIsLoggingErrorContract(body) {
		return true
	}

	// Framework HTTP error response after logging models the failure
	// (logxide django/flask JsonResponse/jsonify audited FPs). Bare
	// return {"error": str(e)} stays reportable (WeThePeople connectors).
	if hasLog && suiteReturnsHTTPErrorResponse(joinedRaw) {
		return true
	}

	// Health-check catch-all that records unhealthy/degraded status after
	// logging (logxide examples/* health endpoints). Generic fallback
	// assigns after log stay reportable (WeThePeople rate-limit defaults).
	if hasLog && suiteModelsUnhealthyStatus(joinedRaw) {
		return true
	}

	// Print of the exception. Only stderr reporting followed by an exit
	// flow, a lone print inside an except ladder, or logger.exception +
	// user-facing print + exit (httptap cli.py:589) is safe; every other
	// print handler stays reportable. logger.exception + return None
	// without print stays reportable (WeThePeople veritas_bridge).
	printStmts := 0
	stderrPrint := false
	for i := range body {
		if strings.Contains(body[i], "print(") || strings.Contains(body[i], "console.print(") {
			printStmts++
			if strings.Contains(bodyRaw[i], "file=sys.stderr") {
				stderrPrint = true
			}
		}
	}
	if printStmts > 0 {
		if stderrPrint && hasExitFlow {
			return true
		}
		if stmtCount == 1 && printStmts == 1 && exceptNestedInExceptSuite(lines, exceptIdx) {
			return true
		}
		if strings.Contains(allText, ".exception(") && hasExitFlow {
			return true
		}
		return false
	}

	// Soft best-effort log after overriding a pre-initialized local via a
	// path-gated optional-import enrichment try (Project_Parva middleware:794
	// — if path: import + override default, logger.warning, flow continues).
	// Unconditional lazy-import + warning (WeThePeople token signing) and
	// warning after non-import work stay reportable.
	tryIdxSoft := enclosingTryIdx(lines, exceptIdx)
	if softBestEffortLogContinue(allText) && !hasRaise && !hasExitFlow &&
		suiteHasContinuationAfter(lines, after, exceptIndent) &&
		tryIdxSoft >= 0 && tryNestedUnderIf(lines, tryIdxSoft) &&
		tryBlockHasImportStmt(lines, tryIdxSoft) &&
		tryBlockOverridesPreinitializedLocal(lines, tryIdxSoft) {
		return true
	}

	// Exception-type extraction feeding a failure counter (rnet_test.py:37):
	// counted and logged, not surfaced as a value.
	if caughtVar != "" && strings.Contains(allText, "type("+caughtVar) && strings.Contains(allText, "+=") {
		return true
	}

	// Raw exception transported in a result payload.
	if caughtVar != "" {
		if strings.Contains(joinedRaw, "error="+caughtVar) {
			return true
		}
		// error=str(exc) only when transported via append payload — bare
		// return ApiResult(error=str(e)) stays reportable (WeThePeople
		// regulationsgov / connectors).
		if strings.Contains(joinedRaw, "error=str("+caughtVar+")") &&
			strings.Contains(allRaw, ".append(") {
			return true
		}
		if strings.Contains(allRaw, "append("+caughtVar) {
			return true
		}
		// Structured result recording "error": str(exc):
		//   - results.append({name, passed, error}) (rulelang run_rule_tests)
		//   - multi-statement outcome build without return/exit (rulelang
		//     validate_date: value/warnings/confidence then fall through)
		// Health-check appends without name=, return {"error": str(e)}, and
		// rejection_detail={"error":...}; return stay reportable (WeThePeople).
		if (strings.Contains(allRaw, `"error":`) || strings.Contains(allRaw, `'error':`)) &&
			lineReferencesVar(allRaw, caughtVar) {
			namedAppend := strings.Contains(allRaw, ".append(") &&
				(strings.Contains(allRaw, `"name":`) || strings.Contains(allRaw, `'name':`) ||
					strings.Contains(allRaw, `"name" :`) || strings.Contains(allRaw, `'name' :`))
			if namedAppend || (stmtCount >= 2 && !hasExitFlow) {
				return true
			}
		}
		// Single-statement raw capture x = e (evidence capture / transport).
		if stmtCount == 1 && assignsRawVar(bodyRaw[0], caughtVar) {
			return true
		}
	}
	// Single-statement pass-through return of a non-exception value
	// (onlymaps _types.py:189 — parse fallback returns the input).
	if stmtCount == 1 && strings.HasPrefix(body[0], "return ") {
		rest := strings.TrimSpace(body[0][len("return "):])
		if isSimpleIdent(rest) && rest != caughtVar && rest != "None" && rest != "True" && rest != "False" {
			return true
		}
	}
	// Optional-import fallback: ImportError/ModuleNotFoundError + import-only
	// try + return fallback. Broad Exception around import-only with a
	// constant/list fallback is the langchain/llamaindex optional-deps shape;
	// Exception around import + non-fallback body stays reportable
	// (pdf_oxide / violit / WeThePeople optional-import TPs).
	if stmtCount == 1 && strings.HasPrefix(strings.TrimSpace(bodyRaw[0]), "return") {
		tryIdxImp := enclosingTryIdx(lines, exceptIdx)
		exceptLine := strings.TrimSpace(lines[exceptIdx].text)
		catchTypesImp := exceptCatchTypes(exceptLine)
		if tryIdxImp >= 0 && tryBlockIsImportOnly(lines, tryIdxImp) {
			if importErrorCatch(catchTypesImp) {
				return true
			}
			if broadExceptionCatch(catchTypesImp) && isImportOptionalFallbackReturn(bodyRaw[0]) {
				return true
			}
		}
	}

	// Soft recording into a warnings list is the Project_Parva rulelang
	// evidence-packet fallback (warnings.append(...exc...)) — failure is
	// retained on the result, not swallowed. Use raw text: Mask blanks strings.
	if strings.Contains(allRaw, "warnings.append(") || strings.Contains(allRaw, ".warnings.append(") {
		return true
	}
	// Documented defensive constant fallback after a pure probe try
	// (httptap: cert_info = None / return None with "defensive"/"best effort"
	// marker). Without a marker, Exception: x = None stays reportable.
	if stmtCount == 1 && (isDefinedConstantAssign(bodyRaw[0]) || isDefinedConstantReturn(bodyRaw[0])) {
		tryIdxProbe := enclosingTryIdx(lines, exceptIdx)
		exceptIndentProbe := indentWidth(lines[exceptIdx].raw)
		if tryIdxProbe >= 0 && tryBlockIsProbe(lines, tryIdxProbe) &&
			exceptLineHasFallbackMarker(rawLines, exceptIdx, exceptIndentProbe) {
			return true
		}
	}

	// Documented defensive multi-statement cleanup that never logs/prints/
	// re-raises and does not bind the exception (niquests wasi upload
	// cleanup: failed = True; notify callback). Only the narrow "defensive"
	// / "no cover" markers — not generic "fallback" comments (WeThePeople).
	if !hasLog && !hasPrint && !hasRaise &&
		exceptLineHasDefensiveMarker(rawLines, exceptIdx, exceptIndent) {
		if caughtVar == "" || !lineReferencesVar(allRaw, caughtVar) {
			return true
		}
	}

	// Best-effort JS/pyodide bridge probe with a defined constant fallback or
	// pass swallow (niquests extensions/pyodide getReader / to_py FPs).
	// Requires an explicit JS-bridge signal so ordinary Exception:pass handlers
	// (WeThePeople TPs) stay reportable. CWE-396 is intentionally not mirrored.
	// Constant return is accepted even when the try body nests async with
	// timeout (pyodide _get_next_chunk) — the bridge signal is the gate.
	tryIdx := enclosingTryIdx(lines, exceptIdx)
	if tryIdx >= 0 && tryBlockHasJSBridgeSignal(lines, tryIdx) {
		if stmtCount == 1 && isDefinedConstantReturn(bodyRaw[0]) {
			return true
		}
		if tryBlockIsProbe(lines, tryIdx) {
			if !suiteNotPassOnly(body) {
				return true
			}
			// Use raw suite text: pytext.Mask blanks string literals, so
			// `response.reason = ""` / `body = b""` look non-constant on masked lines.
			if stmtCount == 1 && isDefinedConstantAssign(bodyRaw[0]) {
				return true
			}
		}
	}

	// Stream/websocket disconnect: single break/continue after a lone await of
	// receive/recv/read (niquests test_asgi ws-echo). Exception:break after
	// await client.get/post stays reportable.
	if stmtCount == 1 && (body[0] == "break" || body[0] == "continue") &&
		tryIdx >= 0 && tryBlockIsStreamReceiveAwait(lines, tryIdx) {
		return true
	}
	return false
}

// tryBlockIsStreamReceiveAwait reports a try body that is exactly one await of
// a stream/websocket receive-style call (not generic client.get).
func tryBlockIsStreamReceiveAwait(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	stmts := 0
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			break
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		stmts++
		if stmts > 1 {
			return false
		}
		// data = await websocket.receive() or bare await ws.recv()
		rhs := stmt
		if _, r, ok := splitAssignmentEq(stmt); ok {
			rhs = strings.TrimSpace(r)
		}
		if !strings.HasPrefix(rhs, "await ") {
			return false
		}
		rest := strings.TrimSpace(rhs[len("await "):])
		// Reject network client verbs.
		if strings.Contains(rest, ".get(") || strings.Contains(rest, ".post(") ||
			strings.Contains(rest, ".request(") || strings.Contains(rest, ".send(") {
			return false
		}
		matched := false
		for _, m := range []string{".receive(", ".recv(", ".read(", ".readline(", ".__anext__("} {
			if strings.Contains(rest, m) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return stmts == 1
}

// suiteIsLoggingErrorContract reports a sole suite statement that delegates
// to the stdlib logging.Handler error path (handleError) or a private
// _handle_error helper used for the same contract.
func suiteIsLoggingErrorContract(body []string) bool {
	if len(body) != 1 {
		return false
	}
	t := strings.TrimSpace(body[0])
	// Require a call form; avoid matching comments or names alone.
	return strings.Contains(t, "handleError(") || strings.Contains(t, "_handle_error(")
}

// suiteReturnsHTTPErrorResponse reports a suite that returns a framework
// HTTP error constructor (Django JsonResponse, Flask jsonify, Starlette/
// FastAPI JSONResponse/HTTPException). Plain dict returns stay reportable.
func suiteReturnsHTTPErrorResponse(joinedRaw string) bool {
	if !strings.Contains(joinedRaw, "return ") {
		return false
	}
	for _, ctor := range []string{"JsonResponse(", "jsonify(", "JSONResponse(", "HTTPException("} {
		if strings.Contains(joinedRaw, ctor) {
			return true
		}
	}
	return false
}

// suiteModelsUnhealthyStatus reports a suite that assigns the conventional
// health-check failure token "unhealthy" (raw text — Mask blanks string
// literals). "degraded" alone is intentionally omitted: WeThePeople
// routers/common.py health endpoint assigns status = "degraded" and is an
// audited BP-PY-1 TP. Generic status/default assigns stay reportable.
func suiteModelsUnhealthyStatus(joinedRaw string) bool {
	return strings.Contains(joinedRaw, `= "unhealthy"`) ||
		strings.Contains(joinedRaw, `= 'unhealthy'`)
}

// softBestEffortLogContinue reports logger.warning-only handlers (not info/
// debug, not .exception/.error/.critical/exc_info). Optional non-fatal surfaces.
// info/debug alone stays reportable (WeThePeople claims/database lookup TPs).
func softBestEffortLogContinue(allText string) bool {
	if !strings.Contains(allText, ".warning(") {
		return false
	}
	for _, hard := range []string{".exception(", ".error(", ".critical(", "exc_info", "traceback.print_exc",
		".info(", ".debug("} {
		if hard == ".warning(" {
			continue
		}
		if strings.Contains(allText, hard) {
			return false
		}
	}
	return strings.Contains(allText, "logger.") ||
		strings.Contains(allText, "log.") ||
		strings.Contains(allText, "logging.")
}

// suiteHasContinuationAfter reports that more code runs after the except suite
// at the try/except indentation level (or outer), i.e. the handler is not the
// terminal body of the enclosing function.
func suiteHasContinuationAfter(lines []codeLine, after, exceptIndent int) bool {
	for j := after; j < len(lines); j++ {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind > exceptIndent {
			continue
		}
		if ind < exceptIndent {
			if strings.HasPrefix(t, "def ") || strings.HasPrefix(t, "async def ") ||
				strings.HasPrefix(t, "class ") {
				return false
			}
			return true
		}
		// Same indent as except: sibling clause or next statement after try.
		if strings.HasPrefix(t, "except") || strings.HasPrefix(t, "finally:") ||
			t == "else:" || strings.HasPrefix(t, "else:") {
			continue
		}
		if strings.HasPrefix(t, "def ") || strings.HasPrefix(t, "async def ") ||
			strings.HasPrefix(t, "class ") {
			return false
		}
		return true
	}
	return false
}

// tryNestedUnderIf reports that the try header sits inside an if/elif suite
// (conditional optional surface), not at the bare function body level.
func tryNestedUnderIf(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	for j := tryIdx - 1; j >= 0; j-- {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind >= tryIndent {
			continue
		}
		return strings.HasPrefix(t, "if ") || strings.HasPrefix(t, "elif ")
	}
	return false
}

// tryBlockHasImportStmt reports that the try body contains an import / from
// / import_module statement.
func tryBlockHasImportStmt(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			return false
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		if importLikeStmt(stmt) {
			return true
		}
	}
	return false
}

// tryBlockOverridesPreinitializedLocal reports a try body that assigns to a
// simple name which was already assigned before the try (best-effort override
// of a default). Used with soft-log continue for optional surfaces.
func tryBlockOverridesPreinitializedLocal(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	// Names assigned inside the try (direct body).
	assigned := map[string]bool{}
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			break
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		if stmtHasSideEffectSink(stmt) {
			return false
		}
		if lhs, _, ok := splitAssignmentEq(stmt); ok {
			name := strings.TrimSpace(lhs)
			if i := strings.Index(name, ":"); i >= 0 {
				name = strings.TrimSpace(name[:i])
			}
			if isSimpleIdent(name) {
				assigned[name] = true
			}
		}
	}
	if len(assigned) == 0 {
		return false
	}
	// Look upward for a prior assignment of the same simple name at a
	// strictly outer indent (pre-initialized default before this try).
	// Same-indent assignments are sibling try bodies / loop peers and must
	// not count (WeThePeople detect_stories multi-detector loop TPs).
	for j := tryIdx - 1; j >= 0; j-- {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind >= tryIndent {
			continue
		}
		if strings.HasPrefix(mt, "def ") || strings.HasPrefix(mt, "async def ") ||
			strings.HasPrefix(mt, "class ") {
			break
		}
		if lhs, _, ok := splitAssignmentEq(mt); ok {
			name := strings.TrimSpace(lhs)
			if i := strings.Index(name, ":"); i >= 0 {
				name = strings.TrimSpace(name[:i])
			}
			if assigned[name] {
				return true
			}
		}
	}
	return false
}

// exceptLineHasDefensiveMarker is a narrow subset of fallback markers used for
// multi-statement cleanup exemptions. Generic "fallback"/"best effort" prose
// is too common in WeThePeople comments to gate BP-PY-1 alone.
func exceptLineHasDefensiveMarker(rawLines []codeLine, exceptIdx, exceptIndent int) bool {
	// Intentionally omits "fallback"/"best effort"/"non-fatal" — those appear
	// in ordinary WeThePeople comments and log messages.
	markers := []string{"no cover", "defensive"}
	if lineHasFallbackMarkerText(rawLines[exceptIdx].raw, markers) {
		return true
	}
	for j := exceptIdx + 1; j < len(rawLines); j++ {
		if strings.TrimSpace(rawLines[j].raw) == "" {
			continue
		}
		if indentWidth(rawLines[j].raw) <= exceptIndent {
			break
		}
		if lineHasFallbackMarkerText(rawLines[j].raw, markers) {
			return true
		}
	}
	return false
}

// isImportOptionalFallbackReturn reports return None/[]/{}/list-comp fallbacks
// used when an optional framework import is missing (langchain/llamaindex).
func isImportOptionalFallbackReturn(stmt string) bool {
	if isDefinedConstantReturn(stmt) {
		return true
	}
	t := strings.TrimSpace(stmt)
	if !strings.HasPrefix(t, "return ") {
		return t == "return"
	}
	rest := strings.TrimSpace(t[len("return "):])
	// list / dict / tuple display forms and comprehensions
	return strings.HasPrefix(rest, "[") || strings.HasPrefix(rest, "{") || strings.HasPrefix(rest, "(")
}

// isDefinedConstantAssign reports `name = <constant>` (or attr = constant).
func isDefinedConstantAssign(stmt string) bool {
	lhs, rhs, ok := splitAssignmentEq(strings.TrimSpace(stmt))
	if !ok || strings.TrimSpace(lhs) == "" {
		return false
	}
	return isDefinedConstantExpr(strings.TrimSpace(rhs))
}

// isDefinedConstantReturn reports `return None` / `return ""` / `return b""` /
// `return False` / `return 0` / `return []` / `return {}` (and bare `return`).
func isDefinedConstantReturn(stmt string) bool {
	t := strings.TrimSpace(stmt)
	if t == "return" {
		return true
	}
	if !strings.HasPrefix(t, "return ") {
		return false
	}
	return isDefinedConstantExpr(strings.TrimSpace(t[len("return "):]))
}

// isDefinedConstantExpr reports a trivial constant used as a defined fallback.
func isDefinedConstantExpr(expr string) bool {
	switch strings.TrimSpace(expr) {
	case "None", "True", "False", "0", "1", "\"\"", "''", "b\"\"", "b''", "[]", "{}", "()", "0.0":
		return true
	}
	return false
}

// suiteNotPassOnly reports a suite whose statements are not solely pass.
func suiteNotPassOnly(body []string) bool {
	for _, s := range body {
		if strings.TrimSpace(s) != "pass" {
			return true
		}
	}
	return false
}

// suiteLineRaiseRE matches a suite statement that is a raise.
var suiteLineRaiseRE = regexp.MustCompile(`(?m)^raise\b`)

// suiteLineReturnRE matches a suite statement that is a return (any indent in
// the joined suite text which has already been left-trimmed per line).
var suiteLineReturnRE = regexp.MustCompile(`(?m)^return\b`)

// wrappedRaiseStatement reports a raise of a newly constructed exception
// (`raise Type(...)`, `raise X from e`) — a translation, not a bare
// re-propagation of the caught condition.
func wrappedRaiseStatement(t string) bool {
	if !strings.HasPrefix(t, "raise") {
		return false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(t, "raise"))
	if rest == "" {
		return false
	}
	if strings.Contains(rest, "(") || strings.Contains(rest, " from ") {
		return true
	}
	return false
}

// exceptNestedInExceptSuite reports an except clause whose enclosing block is
// itself another except suite — the error-handling ladder shape (requestSpeedTest
// increase_limits.py:38 is the innermost fallback of an except ValueError).
func exceptNestedInExceptSuite(lines []codeLine, exceptIdx int) bool {
	exceptIndent := indentWidth(lines[exceptIdx].raw)
	for j := exceptIdx - 1; j >= 0; j-- {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		if indentWidth(lines[j].raw) >= exceptIndent {
			continue
		}
		return strings.HasPrefix(t, "except ")
	}
	return false
}

// exceptSuiteLineIdx returns the line indices of the direct-body statements of
// the except suite at exceptIdx (indent > exceptIndent, trimmed non-empty),
// the index of the first line after the suite, and the suite's body indent.
func exceptSuiteLineIdx(lines []codeLine, exceptIdx int) (suite []int, after int, bodyIndent int) {
	if exceptIdx < 0 || exceptIdx >= len(lines) {
		return nil, exceptIdx + 1, -1
	}
	exceptIndent := indentWidth(lines[exceptIdx].raw)
	bodyIndent = -1
	j := exceptIdx + 1
	for ; j < len(lines); j++ {
		raw := lines[j].raw
		if strings.TrimSpace(raw) == "" {
			continue
		}
		ind := indentWidth(raw)
		if ind <= exceptIndent {
			break
		}
		if bodyIndent < 0 {
			bodyIndent = ind
		}
		if ind == bodyIndent {
			suite = append(suite, j)
		}
	}
	return suite, j, bodyIndent
}

// assignsRawVar reports whether line is an assignment whose RHS is exactly the
// exception variable (the exception object itself is captured/transported).
// `x = str(e)` does not match — only the raw object.
func assignsRawVar(line, varName string) bool {
	if varName == "" || !strings.Contains(line, "=") {
		return false
	}
	i := strings.Index(line, "=")
	if i == 0 {
		return false
	}
	prev := line[i-1]
	if prev == '=' || prev == '!' || prev == '<' || prev == '>' || prev == ':' {
		return false
	}
	if i+1 < len(line) && line[i+1] == '=' {
		return false
	}
	rhs := strings.TrimSpace(line[i+1:])
	return rhs == varName
}

// lineReferencesVar reports whether the raw line mentions the caught variable
// as a whole identifier (e.g. in an f-string interpolation or call argument).
func lineReferencesVar(line, varName string) bool {
	if varName == "" {
		return false
	}
	start := 0
	for {
		i := strings.Index(line[start:], varName)
		if i < 0 {
			return false
		}
		abs := start + i
		ok := true
		if abs > 0 {
			prev := line[abs-1]
			ok = ok && !(isIdentByte(prev) || prev == '.')
		}
		end := abs + len(varName)
		if end < len(line) {
			next := line[end]
			ok = ok && !isIdentByte(next)
		}
		if ok {
			return true
		}
		start = abs + len(varName)
	}
}

// nextDedentedStartsWith reports whether a statement starting with prefix
// (raise/return) follows the except suite at the try's indentation level,
// within the next few statements, before any control-flow header intervenes.
// This captures fall-through handlers whose dedented continuation propagates
// the failure (calgebra gcal.py:221/245, recurrence.py:437) while leaving
// raise-after-if shapes reportable (voicetag encoder.py:79).
func nextDedentedStartsWith(lines []codeLine, start, exceptIndent int, prefix string) bool {
	count := 0
	for j := start; j < len(lines) && count < 3; j++ {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind > exceptIndent {
			continue
		}
		count++
		if strings.HasPrefix(t, prefix) {
			return true
		}
		if isControlFlowHeader(t) {
			return false
		}
	}
	return false
}

// isControlFlowHeader reports statements that open a new block or end the
// current flow, i.e. anything between the fall-through statement and a raise
// that is not itself the propagation.
func isControlFlowHeader(t string) bool {
	for _, kw := range []string{"if ", "elif ", "else:", "for ", "while ", "with ", "async ", "def ", "class ", "try:", "except", "finally:", "match ", "return", "break", "continue", "pass", "import ", "from "} {
		if strings.HasPrefix(t, kw) {
			return true
		}
	}
	return false
}

// tryBlockDeliberatelyRaises reports whether the try block preceding a
// pass-only except contains an explicit raise or pytest.fail — the
// expected-exception test idiom (onlymaps tests/test_database.py:339,516;
// niquests test_requests.py ReadTimeout/ConnectTimeout). Assertions alone
// are not enough: httpmorph's test_proxy.py:40-style assert-in-try handlers
// are audited true positives and stay reportable.
func tryBlockDeliberatelyRaises(lines []codeLine, exceptIdx int) bool {
	exceptIndent := indentWidth(lines[exceptIdx].raw)
	for j := exceptIdx - 1; j >= 0; j-- {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind < exceptIndent {
			return false
		}
		if ind > exceptIndent {
			if suiteLineRaiseRE.MatchString(t) || strings.Contains(t, "pytest.fail(") {
				return true
			}
			continue
		}
		if strings.HasPrefix(t, "try:") {
			return false
		}
		if strings.HasPrefix(t, "except") || strings.HasPrefix(t, "finally:") || strings.HasPrefix(t, "else:") {
			continue
		}
	}
	return false
}

// exceptPassIsSafeBP reports pass-only except suites that are audited false
// positives:
//
//   - optional-dependency ImportError/ModuleNotFoundError fallbacks whose try
//     block is solely imports (module or function scope) — the try block must
//     not execute side effects, so httpmorph's dotenv bootstraps (import +
//     load_dotenv call) stay reportable;
//   - optional load_extension probes that only load and purely use the loaded
//     name (niquests async_session.py:620);
//   - parsing/termination probes (ValueError/TypeError/OSError/AttributeError/
//     UnsupportedOperation/StopIteration/StopAsyncIteration/TimeoutError/
//     CancelledError/GeneratorExit/RequestException) whose try block is a pure
//     assignment/return probe and whose enclosing flow returns a defined
//     default, raises, or continues with a fallback — while wse's pass probes
//     that terminate their method (connection.py:459) stay reportable;
//   - pure conversion guards (int/float/bool) whose invalid input is non-fatal
//     and the enclosing flow continues (parva-mcp content-length parse),
//     including pure follow-on uses of the converted value (Project_Parva
//     backtest hour counter; timegraph fact-id enrichment appends);
//   - sole ValueError after a pure .index()/.find() lookup probe (kiss_headers
//     header_strip); fromisoformat/ipaddress ValueError:pass stays reportable
//     (WeThePeople TPs);
//   - try/except/else recovery where pass is the failure arm and else is the
//     success arm (re-validate-then-redownload FPs);
//   - best-effort Exception/BaseException probes whose try body is a pure
//     attribute/assignment/method-chain probe without network/file/subprocess
//     side effects (niquests pyodide getReader); bare cleanup calls
//     (__aexit__/close) stay reportable (wse connection.py);
//   - best-effort RequestException fetches (single-statement try block);
//   - documented optional/defensive fallbacks (httptap / "Fall back" markers);
//   - expected-exception tests whose try block deliberately raises;
//   - expected-exception async tests whose try is a single await call with a
//     specific (non-broad) catch (niquests wasi_guest edge_cases); bare /
//     Exception:pass in tests stay reportable (httpmorph test_proxy);
//   - expected-exception sync tests under tests/ with a single bare call and
//     a specific catch (niquests wasi_guest unwrap Err);
//   - dual-path tests with a specific catch and an "expected"/"may raise"
//     marker (Project_Parva test_calculator year-outside-range).
//
// Generic ValueError:pass after fromisoformat/strptime stays reportable
// (WeThePeople TPs). Offline tooling paths are skipped entirely in
// detectBPPY2 (Project_Parva backend/tools/load_test).
func exceptPassIsSafeBP(lines, rawLines []codeLine, exceptIdx int, isTest, isTestDir bool) bool {
	exceptLine := strings.TrimSpace(lines[exceptIdx].text)
	if isTest && tryBlockDeliberatelyRaises(lines, exceptIdx) {
		return true
	}
	suite, after, _ := exceptSuiteLineIdx(lines, exceptIdx)
	if len(suite) != 1 || strings.TrimSpace(lines[suite[0]].text) != "pass" {
		return false
	}
	exceptIndent := indentWidth(lines[exceptIdx].raw)
	if exceptLineHasFallbackMarker(rawLines, exceptIdx, exceptIndent) {
		return true
	}
	// try/except/else: pass is the failure arm; else is the success path.
	if exceptFollowedByElse(lines, after, exceptIndent) {
		return true
	}
	tryIdx := enclosingTryIdx(lines, exceptIdx)
	catchTypes := exceptCatchTypes(exceptLine)
	if importErrorCatch(catchTypes) && tryIdx >= 0 &&
		(tryBlockIsImportOnly(lines, tryIdx) ||
			tryBlockIsOptionalExtensionLoad(lines, tryIdx) ||
			tryBlockIsOptionalFeatureRegistration(lines, tryIdx)) {
		return true
	}
	// Designed iterator/async termination: only when the try body is a pure
	// next()/__anext__() probe (niquests async_session resolve_redirects FP).
	// Generator-teardown `next(g, None)` / `aclose()` in finally (PyDepends
	// audited TPs) stay reportable.
	if tryIdx >= 0 && terminationCatchOnly(catchTypes) && tryBlockIsIteratorNextProbe(lines, tryIdx) {
		return true
	}
	// Pure cancel/generator-exit wait (sgi lifespan await Future forever).
	if tryIdx >= 0 && terminationCatchOnly(catchTypes) && tryBlockIsAsyncCancelWait(lines, tryIdx) {
		return true
	}
	if asyncShutdownCatch(catchTypes) {
		return true
	}
	// AttributeError after pure attr probes (super().__getattribute__/hasattr).
	if attributeErrorCatch(catchTypes) && tryIdx >= 0 && tryBlockIsAttrProbe(lines, tryIdx) {
		return true
	}
	// Best-effort close()/aclose() that ignores OSError (testserver teardown).
	if osErrorCatch(catchTypes) && tryIdx >= 0 && tryBlockIsCloseOnly(lines, tryIdx) {
		return true
	}
	// Expected-exception async test paths: single await call, specific catch,
	// only under tests/ directories (niquests wasi_guest edge_cases).
	// Scripts named *_test.py (Project_Parva load_test) and bare/Exception:pass
	// tests (httpmorph test_proxy) stay reportable; sync chaos swallows stay.
	if isTestDir && tryIdx >= 0 && tryBlockIsSingleAwaitCall(lines, tryIdx) &&
		!broadExceptionCatch(catchTypes) && !isBareExceptClause(exceptLine) {
		return true
	}
	// Expected-exception sync test: single bare call + specific catch under
	// tests/ (niquests wasi_guest edge_cases unwrap Err). Exception:pass and
	// bare except stay reportable; dual-path assert/except uses markers below.
	if isTestDir && tryIdx >= 0 && tryBlockIsSingleBareCall(lines, tryIdx) &&
		!broadExceptionCatch(catchTypes) && !isBareExceptClause(exceptLine) {
		return true
	}
	// Dual-path / documented expected-exception tests (Project_Parva
	// test_calculator year-outside-range): specific catch + "expected" marker.
	// Exception:pass stays reportable even with a comment (httpmorph TPs).
	if isTest && !broadExceptionCatch(catchTypes) && !isBareExceptClause(exceptLine) &&
		exceptLineHasExpectedExceptionMarker(rawLines, exceptIdx, exceptIndent) {
		return true
	}
	if probeCatch(catchTypes) {
		// Alternate return/raise after a try that itself returned/raised is a
		// defined fallback (niquests kiss_headers builder: ValueError:pass then
		// return str(...)). A plain function epilogue `return` after try/except
		// that never returns in the try body is NOT a fallback (pingram
		// conftest monkeypatch TP — return sleep always runs).
		if tryIdx >= 0 && tryBlockHasReturnOrRaise(lines, tryIdx) {
			if nextDedentedStartsWith(lines, after, exceptIndent, "return") ||
				nextDedentedStartsWith(lines, after, exceptIndent, "raise") {
				return true
			}
		}
		if nextDedentedStartsWith(lines, after, exceptIndent, "raise") {
			return true
		}
		// Pure int/float/bool conversion guards may continue with a streaming
		// or re-measure fallback without an immediate return (parva-mcp client).
		if tryIdx >= 0 && tryBlockIsConversionGuard(lines, tryIdx) {
			return true
		}
		// Sole TypeError or sole OSError after a pure attr/method probe
		// (niquests kiss_headers TypeError:pass; body.tell() OSError).
		// Mixed (ValueError, TypeError) stays reportable (WeThePeople TPs).
		if tryIdx >= 0 && tryBlockIsProbe(lines, tryIdx) &&
			(onlyTypeErrorCatch(catchTypes) || onlyOSErrorCatch(catchTypes)) {
			return true
		}
		// Sole ValueError after a pure .index()/.find() lookup (kiss_headers
		// header_strip). ipaddress/fromisoformat ValueError:pass stay reportable.
		if tryIdx >= 0 && onlyValueErrorCatch(catchTypes) && tryBlockIsIndexLookupProbe(lines, tryIdx) {
			return true
		}
		// Best-effort network/JSON enrichment (single-stmt RequestException, or
		// multi-stmt RequestException+JSONDecodeError+HTTPError update checks).
		if bestEffortNetworkCatch(catchTypes) && tryIdx >= 0 &&
			(tryBlockSingleStatement(lines, tryIdx) || requestExceptionCatch(catchTypes)) {
			return true
		}
	}
	// Broad Exception/BaseException: only best-effort JS/pyodide bridge probes
	// (niquests extensions/pyodide). Ordinary Exception:pass after attribute
	// probes or optional imports (WeThePeople TPs) stay reportable. Bare
	// cleanup calls (wse __aexit__/close) are not probes and stay reportable.
	if broadExceptionCatch(catchTypes) && tryIdx >= 0 &&
		tryBlockIsProbe(lines, tryIdx) && tryBlockHasJSBridgeSignal(lines, tryIdx) {
		return true
	}
	return false
}

// isBareExceptClause reports `except:` / `except :` with no type list.
func isBareExceptClause(t string) bool {
	t = strings.TrimSpace(t)
	if !strings.HasPrefix(t, "except") {
		return false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(t, "except"))
	return rest == ":" || strings.HasPrefix(rest, ":")
}

// onlyValueErrorCatch reports a catch list that is exactly ValueError.
func onlyValueErrorCatch(types []string) bool {
	return len(types) == 1 && types[0] == "ValueError"
}

// tryBlockIsIndexLookupProbe reports a single-statement try body that only
// probes str/list membership via .index()/.find()/.rindex()/.rfind().
func tryBlockIsIndexLookupProbe(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	stmts := 0
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			break
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		stmts++
		if stmts > 1 {
			return false
		}
		if stmtHasSideEffectSink(stmt) {
			return false
		}
		// Assignment preferred: x = s.index(...)
		_, rhs, ok := splitAssignmentEq(stmt)
		if !ok {
			rhs = stmt
		}
		if !indexLookupCallIn(rhs) {
			return false
		}
	}
	return stmts == 1
}

func indexLookupCallIn(expr string) bool {
	for _, m := range []string{".index(", ".rindex(", ".find(", ".rfind("} {
		if strings.Contains(expr, m) {
			return true
		}
	}
	return false
}

// tryBlockIsAsyncCancelWait reports a single-await try body that waits on a
// Future/Event/sleep until cancelled (sgi lifespan shutdown path).
func tryBlockIsAsyncCancelWait(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	stmts := 0
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			break
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		stmts++
		if stmts > 1 {
			return false
		}
		if !strings.HasPrefix(stmt, "await ") {
			return false
		}
		rest := strings.TrimSpace(stmt[len("await "):])
		// asyncio.Future() / loop.create_future() / event.wait() / asyncio.sleep
		if strings.Contains(rest, "Future(") || strings.Contains(rest, "create_future(") ||
			strings.HasSuffix(rest, ".wait()") || strings.Contains(rest, ".wait(") ||
			strings.Contains(rest, "asyncio.sleep(") || strings.Contains(rest, "sleep(") {
			continue
		}
		return false
	}
	return stmts == 1
}

// tryBlockIsSingleAwaitCall reports a try body that is exactly one await
// expression (no assert) — expected-exception async test shape.
func tryBlockIsSingleAwaitCall(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	stmts := 0
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			break
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		stmts++
		if stmts > 1 {
			return false
		}
		if strings.HasPrefix(stmt, "assert ") || strings.Contains(stmt, ".assert") {
			return false
		}
		if strings.HasPrefix(stmt, "raise ") || stmt == "raise" {
			return false
		}
		if !strings.HasPrefix(stmt, "await ") {
			return false
		}
	}
	return stmts == 1
}

// exceptFollowedByElse reports a try/except/else where the except suite is
// only pass and the next same-indent clause is else: (success arm).
func exceptFollowedByElse(lines []codeLine, after, exceptIndent int) bool {
	for j := after; j < len(lines); j++ {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind > exceptIndent {
			continue
		}
		if ind < exceptIndent {
			return false
		}
		return t == "else:" || strings.HasPrefix(t, "else:")
	}
	return false
}

// tryBlockIsConversionGuard reports a try body that only guards pure numeric/
// boolean conversions via assignment (x = int(...)), not return. Assignment
// leaves the prior/unset state on failure (niquests content_length / SSE
// retry). After a conversion, pure follow-on uses of the converted value are
// still non-fatal (Project_Parva backtest hour bucket counter; timegraph
// fact-id enrichment appends). `return int(value); except: pass` with no
// fallback stays reportable (parsing-fallback-vulnerable fixture).
// fromisoformat/strptime stay out.
func tryBlockIsConversionGuard(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	sawConversion := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			return sawConversion
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		if stmtHasSideEffectSink(stmt) {
			return false
		}
		// Non-pure parsers must stay reportable (WeThePeople fromisoformat TPs).
		if conversionHasBadParser(stmt) {
			return false
		}
		// return int(...) is not a non-fatal guard by itself — needs the
		// return/raise-after-pass path instead of this guard.
		if stmt == "return" || strings.HasPrefix(stmt, "return ") {
			return false
		}
		if isProbeBranchHeader(stmt) || stmt == "pass" || stmt == "continue" || stmt == "break" ||
			strings.HasPrefix(stmt, "raise ") || stmt == "raise" {
			if conversionCallIn(stmt) {
				sawConversion = true
			}
			continue
		}
		if _, rhs, ok := splitAssignmentEq(stmt); ok {
			if conversionCallIn(rhs) || conversionCallIn(stmt) {
				sawConversion = true
				continue
			}
			// Pure follow-on after conversion: counter updates via += / rebinding
			// of derived values (backtest mismatch_by_ingress_hour[hour] += 1).
			// splitAssignmentEq treats `+=` as assignment with a non-conversion RHS.
			if sawConversion {
				continue
			}
			return false
		}
		// Bare conversion expression is not typical.
		if conversionCallIn(stmt) {
			sawConversion = true
			continue
		}
		// Pure enrichment after conversion: list/set append of derived values
		// (timegraph fact_ids.append(bs_ad_fact_id(...))).
		if sawConversion && conversionFollowOnStmt(stmt) {
			continue
		}
		return false
	}
	return sawConversion
}

// conversionHasBadParser reports non-pure parsers that must stay reportable
// under conversion-guard (WeThePeople fromisoformat / json.loads TPs).
func conversionHasBadParser(stmt string) bool {
	for _, bad := range []string{
		"fromisoformat(", "strptime(", "json.loads(", "yaml.load(",
		"yaml.safe_load(", "parse_", "loads(", "decode(",
	} {
		if strings.Contains(stmt, bad) {
			return true
		}
	}
	return false
}

// conversionFollowOnStmt reports pure enrichment statements allowed after a
// conversion assignment (append/extend/add of derived values).
func conversionFollowOnStmt(stmt string) bool {
	for _, m := range []string{".append(", ".extend(", ".add("} {
		if strings.Contains(stmt, m) {
			return true
		}
	}
	return false
}

// tryBlockIsSingleBareCall reports a try body that is exactly one bare call
// expression (no assignment, no assert, no raise) — expected-exception sync
// test shape (niquests wasi_guest edge_cases unwrap Err).
func tryBlockIsSingleBareCall(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	stmts := 0
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			break
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		stmts++
		if stmts > 1 {
			return false
		}
		if strings.HasPrefix(stmt, "assert ") || strings.Contains(stmt, ".assert") {
			return false
		}
		if strings.HasPrefix(stmt, "raise ") || stmt == "raise" {
			return false
		}
		if strings.HasPrefix(stmt, "await ") {
			return false // covered by tryBlockIsSingleAwaitCall
		}
		if _, _, ok := splitAssignmentEq(stmt); ok {
			return false
		}
		if !strings.Contains(stmt, "(") {
			return false
		}
	}
	return stmts == 1
}

// exceptLineHasExpectedExceptionMarker reports a comment documenting that the
// pass handler is an expected-exception arm (Project_Parva test_calculator).
// Distinct from generic fallback markers so production KeyError:pass TPs stay.
func exceptLineHasExpectedExceptionMarker(rawLines []codeLine, exceptIdx, exceptIndent int) bool {
	markers := []string{"expected", "may raise", "intentionally"}
	if lineHasFallbackMarkerText(rawLines[exceptIdx].raw, markers) {
		return true
	}
	for j := exceptIdx + 1; j < len(rawLines); j++ {
		if strings.TrimSpace(rawLines[j].raw) == "" {
			continue
		}
		if indentWidth(rawLines[j].raw) <= exceptIndent {
			break
		}
		if lineHasFallbackMarkerText(rawLines[j].raw, markers) {
			return true
		}
	}
	return false
}

func conversionCallIn(stmt string) bool {
	for _, fn := range []string{"int(", "float(", "bool(", "complex(", "ord("} {
		if strings.Contains(stmt, fn) {
			return true
		}
	}
	return false
}

// broadExceptionCatch reports Exception or BaseException in the catch list.
func broadExceptionCatch(types []string) bool {
	for _, ty := range types {
		if ty == exceptionName || ty == baseExceptionName {
			return true
		}
	}
	return false
}

// enclosingTryIdx returns the try header line index whose suite contains the
// except clause at exceptIdx, or -1.
func enclosingTryIdx(lines []codeLine, exceptIdx int) int {
	exceptIndent := indentWidth(lines[exceptIdx].raw)
	for j := exceptIdx - 1; j >= 0; j-- {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind < exceptIndent {
			return -1
		}
		if ind > exceptIndent {
			continue
		}
		if strings.HasPrefix(t, "try:") {
			return j
		}
		if strings.HasPrefix(t, "except") || strings.HasPrefix(t, "finally:") || strings.HasPrefix(t, "else:") {
			continue
		}
		return -1
	}
	return -1
}

// exceptLineHasFallbackMarker reports a marker comment on the except clause or
// the first suite lines that documents the handler as an intentional
// optional/defensive fallback. Matching is case-insensitive so "# Fall back"
// (Project_Parva calculator_v2) is recognized.
func exceptLineHasFallbackMarker(rawLines []codeLine, exceptIdx, exceptIndent int) bool {
	markers := []string{"no cover", "best effort", "non-fatal", "defensive", "fall back", "fallback"}
	if lineHasFallbackMarkerText(rawLines[exceptIdx].raw, markers) {
		return true
	}
	for j := exceptIdx + 1; j < len(rawLines); j++ {
		if strings.TrimSpace(rawLines[j].raw) == "" {
			continue
		}
		if indentWidth(rawLines[j].raw) <= exceptIndent {
			break
		}
		if lineHasFallbackMarkerText(rawLines[j].raw, markers) {
			return true
		}
	}
	return false
}

func lineHasFallbackMarkerText(raw string, markers []string) bool {
	lower := strings.ToLower(raw)
	for _, marker := range markers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// tryBlockIsImportOnly reports a try block whose direct-body statements are
// all import statements (or assignments of importlib.import_module /
// __import__ / load_extension results). An optional-dependency fallback whose
// try block executes calls or control flow stays reportable.
func tryBlockIsImportOnly(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	bodyIndent := -1
	for j := tryIdx + 1; j < len(lines); j++ {
		raw := lines[j].raw
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(raw)
		if ind <= tryIndent {
			return true
		}
		if bodyIndent < 0 {
			bodyIndent = ind
		}
		if ind != bodyIndent {
			continue
		}
		// Merge multi-line import statements.
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		if !importLikeStmt(stmt) {
			return false
		}
	}
	return true
}

// tryBlockIsOptionalExtensionLoad reports a try block that loads an optional
// extension (load_extension / import_module / __import__) and then only
// purely uses the loaded name (for/if/return/attr calls on it) with no
// unrelated side effects — niquests async_session.py:620.
func tryBlockIsOptionalExtensionLoad(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	loaded := map[string]bool{}
	sawLoad := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			return sawLoad
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		if stmtHasSideEffectSink(stmt) {
			return false
		}
		if importLikeStmt(stmt) {
			// Only dynamic extension loaders count — plain `import`/`from`
			// optional-deps are handled by tryBlockIsImportOnly (import-only
			// bodies). Calling an imported name (notify_db_query()) is not an
			// extension-load probe (WeThePeople database.py TPs).
			isExt := strings.Contains(stmt, "load_extension(") ||
				strings.Contains(stmt, "importlib.import_module(") ||
				strings.Contains(stmt, "__import__(")
			if !isExt {
				return false
			}
			sawLoad = true
			if name := importAssignedName(stmt); name != "" {
				loaded[name] = true
			}
			continue
		}
		// Pure use of loaded names: for/if/return/assignment referencing them.
		if isProbeBranchHeader(stmt) || strings.HasPrefix(stmt, "return") ||
			stmt == "pass" || stmt == "continue" || stmt == "break" {
			continue
		}
		if _, _, ok := splitAssignmentEq(stmt); ok {
			continue
		}
		// Bare expression must reference a loaded name (extension.supported…).
		if len(loaded) == 0 {
			return false
		}
		refsLoaded := false
		for name := range loaded {
			if lineReferencesVar(stmt, name) {
				refsLoaded = true
				break
			}
		}
		if !refsLoaded {
			return false
		}
	}
	return sawLoad
}

// importAssignedName returns the LHS name of `name = load_extension(...)` /
// import_module / __import__ assignments.
func importAssignedName(t string) string {
	lhs, rhs, ok := splitAssignmentEq(t)
	if !ok {
		return ""
	}
	rhs = strings.TrimSpace(rhs)
	if !strings.Contains(rhs, "load_extension(") &&
		!strings.Contains(rhs, "importlib.import_module(") &&
		!strings.Contains(rhs, "__import__(") {
		return ""
	}
	name := strings.TrimSpace(lhs)
	// Strip simple annotation: `ext: Any`
	if i := strings.Index(name, ":"); i >= 0 {
		name = strings.TrimSpace(name[:i])
	}
	if isSimpleIdent(name) {
		return name
	}
	return ""
}

// importStmtNames extracts simple names from `import a, b as c` / `from x import y`.
func importStmtNames(t string) []string {
	var out []string
	if strings.HasPrefix(t, "import ") {
		rest := strings.TrimSpace(strings.TrimPrefix(t, "import "))
		for _, part := range strings.Split(rest, ",") {
			part = strings.TrimSpace(part)
			if i := strings.Index(part, " as "); i >= 0 {
				part = strings.TrimSpace(part[i+len(" as "):])
			} else if i := strings.Index(part, "."); i >= 0 {
				part = part[:i]
			}
			if isSimpleIdent(part) {
				out = append(out, part)
			}
		}
		return out
	}
	if strings.HasPrefix(t, "from ") {
		if i := strings.Index(t, " import "); i >= 0 {
			rest := strings.TrimSpace(t[i+len(" import "):])
			if rest == "*" {
				return out
			}
			for _, part := range strings.Split(rest, ",") {
				part = strings.TrimSpace(part)
				if j := strings.Index(part, " as "); j >= 0 {
					part = strings.TrimSpace(part[j+len(" as "):])
				}
				if isSimpleIdent(part) {
					out = append(out, part)
				}
			}
		}
	}
	return out
}

// importLikeStmt reports a statement that is an import or an assignment of an
// importlib.import_module / __import__ / load_extension call.
func importLikeStmt(t string) bool {
	if strings.HasPrefix(t, "import ") || strings.HasPrefix(t, "from ") ||
		strings.Contains(t, "importlib.import_module(") || strings.Contains(t, "__import__(") ||
		strings.Contains(t, "load_extension(") {
		return true
	}
	lhs, rhs, ok := splitAssignmentEq(t)
	if !ok || strings.TrimSpace(lhs) == "" {
		return false
	}
	rhs = strings.TrimSpace(rhs)
	return strings.Contains(rhs, "importlib.import_module(") || strings.Contains(rhs, "__import__(") ||
		strings.Contains(rhs, "load_extension(")
}

// tryBlockIsProbe reports a try block whose statements are pure probes:
// assignments, returns, and branch headers (if/for/while) over attribute /
// method-chain expressions — without network, file, or subprocess side
// effects and without bare cleanup call statements. Nested lines under
// branch headers are accepted when they share the same shape (niquests
// pyodide getReader / headers.entries FPs). Bare expression statements
// (`await cm.__aexit__(...)`, `self._ws.close()`) keep the try non-probe so
// wse connection.py teardown TPs stay reportable.
func tryBlockIsProbe(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	sawBody := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			return sawBody
		}
		sawBody = true
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		if !isProbeStmt(stmt) {
			return false
		}
	}
	return sawBody
}

// isProbeStmt reports a single try-body statement that is probe-shaped.
func isProbeStmt(stmt string) bool {
	t := strings.TrimSpace(stmt)
	if t == "" || t == "pass" || t == "break" || t == "continue" {
		return true
	}
	if stmtHasSideEffectSink(t) {
		return false
	}
	if t == "return" {
		return true
	}
	if strings.HasPrefix(t, "return ") {
		// `return do_write(payload)` is business logic, not a probe expression;
		// allow constants, names, attr chains, and known pure conversions only.
		return isProbeReturnExpr(strings.TrimSpace(t[len("return "):]))
	}
	if isProbeBranchHeader(t) {
		return true
	}
	if _, _, ok := splitAssignmentEq(t); ok {
		return true
	}
	// Bare JS-bridge expressions (run_sync(reader.cancel()), getReader()) are
	// best-effort probes; other bare calls (close/__aexit__) stay non-probe.
	if isJSBridgeExpr(t) {
		return true
	}
	// Bare expression statements (cleanup calls, sends) are not probes.
	return false
}

// jsBridgeSignals are tokens that mark a try body as a best-effort JS/pyodide
// bridge probe (niquests extensions/pyodide audited FPs).
var jsBridgeSignals = []string{
	"getReader(",
	"to_py(",
	"run_sync(",
	"js_response",
	"js_headers",
	"js_body",
	".entries(",
	"status_text",
	".bytes()",
	"pyodide",
	"_js_",
	"_ws.",
	"_proxies",
	"ReadableStream",
	"JsProxy",
	"JsException",
	".cancel(",
	".destroy(",
}

// isJSBridgeExpr reports a statement that touches the JS/pyodide bridge API.
func isJSBridgeExpr(t string) bool {
	for _, s := range jsBridgeSignals {
		if strings.Contains(t, s) {
			return true
		}
	}
	return false
}

// tryBlockHasJSBridgeSignal reports that the try body mentions a JS/pyodide
// bridge API. Used to keep ordinary Exception:pass TPs reportable.
func tryBlockHasJSBridgeSignal(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			return false
		}
		if isJSBridgeExpr(mt) {
			return true
		}
	}
	return false
}

// isProbeReturnExpr reports return expressions that are probe-shaped (constants,
// names, attribute chains, pure conversions) rather than bare business calls.
func isProbeReturnExpr(expr string) bool {
	expr = strings.TrimSpace(expr)
	if expr == "" || isDefinedConstantExpr(expr) || isSimpleIdent(expr) {
		return true
	}
	if stmtHasSideEffectSink(expr) {
		return false
	}
	// Attribute / method chain: value.to_py(), self._reader.read()
	if strings.Contains(expr, ".") {
		return true
	}
	// Bare call: only known pure conversions (int/bytes/str/…).
	if i := strings.Index(expr, "("); i > 0 {
		name := strings.TrimSpace(expr[:i])
		return isProbeConversionFunc(name)
	}
	return true
}

// isProbeConversionFunc reports builtin/pure conversion callees allowed in
// probe return expressions.
func isProbeConversionFunc(name string) bool {
	switch name {
	case "int", "float", "str", "bytes", "bool", "len", "abs", "ord", "chr",
		"hex", "bin", "oct", "round", "complex", "list", "dict", "set", "tuple",
		"frozenset", "bytearray", "memoryview", "repr", "ascii", "format":
		return true
	}
	return false
}

// isProbeBranchHeader reports if/elif/else/for/while/async for headers used
// to guard attribute probes. with/try/def are intentionally excluded.
func isProbeBranchHeader(t string) bool {
	return strings.HasPrefix(t, "if ") || strings.HasPrefix(t, "elif ") ||
		t == "else:" || strings.HasPrefix(t, "for ") || strings.HasPrefix(t, "while ") ||
		strings.HasPrefix(t, "async for ")
}

// stmtHasSideEffectSink reports network, file, or subprocess side effects in
// a try body — those must not be treated as pure probes.
func stmtHasSideEffectSink(t string) bool {
	sinks := []string{
		"open(",
		"requests.",
		"httpx.",
		"urllib.",
		"subprocess.",
		"os.system(",
		"os.popen(",
		"os.remove(",
		"os.unlink(",
		"os.rmdir(",
		".write(",
		".writelines(",
		"urlopen(",
		"urlretrieve(",
		"load_dotenv(",
	}
	for _, s := range sinks {
		if strings.Contains(t, s) {
			// Avoid matching attribute `.open(` / names containing open(
			if s == "open(" {
				if indexOfBareOpenCall(t) >= 0 {
					return true
				}
				continue
			}
			return true
		}
	}
	return false
}

// flowContinuesAfter reports a plain statement (not a same-try clause, not an
// else/elif continuation, not a nested def/class) following the handler at or
// below its indent — the enclosing flow continues with a defined fallback.
func flowContinuesAfter(lines []codeLine, after, exceptIndent int) bool {
	count := 0
	for j := after; j < len(lines) && count < 8; j++ {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind > exceptIndent {
			continue
		}
		count++
		// Same-try clauses (else/finally/except) belong to the try statement.
		if strings.HasPrefix(t, "else:") || strings.HasPrefix(t, "elif ") ||
			strings.HasPrefix(t, "except") || strings.HasPrefix(t, "finally:") ||
			strings.HasPrefix(t, "def ") || strings.HasPrefix(t, "async def ") ||
			strings.HasPrefix(t, "class ") {
			continue
		}
		return true
	}
	return false
}

// tryBlockSingleStatement reports a try block with exactly one direct-body
// statement (best-effort fetch shape).
func tryBlockSingleStatement(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	bodyIndent := -1
	stmts := 0
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			return stmts == 1
		}
		if bodyIndent < 0 {
			bodyIndent = ind
		}
		if ind == bodyIndent {
			stmts++
			if stmts > 1 {
				return false
			}
		}
	}
	return stmts == 1
}

func requestExceptionCatch(types []string) bool {
	for _, ty := range types {
		if ty == "RequestException" {
			return true
		}
	}
	return false
}

// exceptCatchTypes returns the exception type names in an except clause
// (unwrapping a single-level tuple).
func exceptCatchTypes(t string) []string {
	t = strings.TrimSpace(t)
	if !strings.HasPrefix(t, "except ") {
		return nil
	}
	rest := strings.TrimSpace(strings.TrimPrefix(t, "except "))
	if i := strings.Index(rest, ":"); i >= 0 {
		rest = strings.TrimSpace(rest[:i])
	}
	if i := strings.Index(rest, " as "); i >= 0 {
		rest = strings.TrimSpace(rest[:i])
	}
	rest = strings.TrimPrefix(rest, "(")
	rest = strings.TrimSuffix(rest, ")")
	types := strings.Split(rest, ",")
	out := make([]string, 0, len(types))
	for _, ty := range types {
		ty = strings.TrimSpace(ty)
		if i := strings.LastIndex(ty, "."); i >= 0 {
			ty = ty[i+1:]
		}
		if ty != "" {
			out = append(out, ty)
		}
	}
	return out
}

func importErrorCatch(types []string) bool {
	for _, ty := range types {
		if ty == "ImportError" || ty == "ModuleNotFoundError" {
			return true
		}
	}
	return false
}

// terminationCatchOnly reports that every caught type is a designed
// iterator/async termination signal. Mixed catches with ValueError/OSError/
// RuntimeError stay reportable (except asyncShutdownCatch).
func terminationCatchOnly(types []string) bool {
	if len(types) == 0 {
		return false
	}
	for _, ty := range types {
		switch ty {
		case "StopIteration", "StopAsyncIteration", "CancelledError", "GeneratorExit", "TimeoutError":
			// ok
		default:
			return false
		}
	}
	return true
}

// tryBlockIsIteratorNextProbe reports a try body that only advances an
// iterator/async generator (`x = next(it)`, `x = await gen.__anext__()`).
// `next(g, None)` (default arg swallows StopIteration) and `aclose()` are
// NOT probes — those are PyDepends-style teardown TPs.
func tryBlockIsIteratorNextProbe(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	stmts := 0
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			break
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		stmts++
		if stmts > 1 {
			return false
		}
		// Assignment form preferred.
		_, rhs, ok := splitAssignmentEq(stmt)
		if ok {
			rhs = strings.TrimSpace(rhs)
			if isIteratorNextCall(rhs) {
				continue
			}
			return false
		}
		// Bare await gen.__anext__() without assignment.
		if isIteratorNextCall(stmt) {
			continue
		}
		return false
	}
	return stmts == 1
}

// isIteratorNextCall reports next(x) / await x.__anext__() without a default
// second argument to next().
func isIteratorNextCall(expr string) bool {
	expr = strings.TrimSpace(expr)
	if strings.HasPrefix(expr, "await ") {
		expr = strings.TrimSpace(expr[len("await "):])
	}
	if strings.Contains(expr, ".__anext__(") {
		return true
	}
	// next(it) — not next(it, default)
	if !strings.HasPrefix(expr, "next(") {
		return false
	}
	args := expr[len("next("):]
	// Strip trailing )
	if i := strings.LastIndex(args, ")"); i >= 0 {
		args = args[:i]
	}
	// A default arg is present if there is a top-level comma.
	depth := 0
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				return false // next(g, None)
			}
		}
	}
	return strings.TrimSpace(args) != ""
}

// asyncShutdownCatch reports lifespan/async-shutdown handlers that mix
// TimeoutError/CancelledError with RuntimeError (loop already closed) —
// niquests sgi/_async/__init__.py lifespan teardown.
func asyncShutdownCatch(types []string) bool {
	if len(types) < 2 {
		return false
	}
	hasTerm, hasRuntime := false, false
	for _, ty := range types {
		switch ty {
		case "TimeoutError", "CancelledError", "GeneratorExit":
			hasTerm = true
		case "RuntimeError":
			hasRuntime = true
		default:
			return false
		}
	}
	return hasTerm && hasRuntime
}

// tryBlockHasReturnOrRaise reports that the try body contains a return or raise
// (used to distinguish alternate-return fallbacks from function epilogues).
func tryBlockHasReturnOrRaise(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			return false
		}
		if strings.HasPrefix(mt, "return") || strings.HasPrefix(mt, "raise") {
			return true
		}
	}
	return false
}

// onlyTypeErrorCatch reports a catch list that is exactly TypeError.
func onlyTypeErrorCatch(types []string) bool {
	return len(types) == 1 && types[0] == "TypeError"
}

// onlyOSErrorCatch reports a catch list that is exactly OSError (or alias).
func onlyOSErrorCatch(types []string) bool {
	return len(types) == 1 && (types[0] == "OSError" || types[0] == "IOError" || types[0] == "EnvironmentError")
}

// osErrorCatch reports OSError (or aliases) among catch types.
func osErrorCatch(types []string) bool {
	for _, ty := range types {
		if ty == "OSError" || ty == "IOError" || ty == "EnvironmentError" {
			return true
		}
	}
	return false
}

// tryBlockIsCloseOnly reports a try body that only closes a resource.
func tryBlockIsCloseOnly(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	sawClose := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			return sawClose
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		stmt = strings.TrimSpace(strings.TrimPrefix(stmt, "await "))
		if strings.HasSuffix(stmt, ".close()") || strings.HasSuffix(stmt, ".aclose()") ||
			strings.HasSuffix(stmt, ".close()") || strings.Contains(stmt, ".close(") ||
			strings.Contains(stmt, ".aclose(") {
			// Pure close call (optionally with args).
			if strings.Contains(stmt, "=") {
				return false
			}
			sawClose = true
			continue
		}
		return false
	}
	return sawClose
}

// bestEffortNetworkCatch reports RequestException and/or JSON/HTTP parse types
// used for optional enrichment (help.py check_update, revocation AIA fetch).
func bestEffortNetworkCatch(types []string) bool {
	if len(types) == 0 {
		return false
	}
	hasNet := false
	for _, ty := range types {
		switch ty {
		case "RequestException", "HTTPError", "ConnectionError", "Timeout",
			"ConnectTimeout", "ReadTimeout", "JSONDecodeError", "URLError":
			hasNet = true
		default:
			return false
		}
	}
	return hasNet
}

// attributeErrorCatch reports AttributeError among the catch types.
func attributeErrorCatch(types []string) bool {
	for _, ty := range types {
		if ty == "AttributeError" {
			return true
		}
	}
	return false
}

// tryBlockIsAttrProbe reports a try body composed only of
// super().__getattribute__ / hasattr / getattr probes (and if-guards over
// them). Bare cleanup calls and network side effects stay non-probe.
func tryBlockIsAttrProbe(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	sawBody := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			return sawBody
		}
		sawBody = true
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		if !isAttrProbeStmt(stmt) {
			return false
		}
	}
	return sawBody
}

// isAttrProbeStmt reports a single try-body statement that only probes
// attribute presence/value via super().__getattribute__, hasattr, or getattr.
func isAttrProbeStmt(stmt string) bool {
	t := strings.TrimSpace(stmt)
	if t == "" || t == "pass" || t == "break" || t == "continue" {
		return true
	}
	if stmtHasSideEffectSink(t) {
		return false
	}
	if isProbeBranchHeader(t) {
		// if/elif conditions must still look like attr probes.
		return strings.Contains(t, "__getattribute__") ||
			strings.Contains(t, "hasattr(") ||
			strings.Contains(t, "getattr(") ||
			t == "else:"
	}
	if _, rhs, ok := splitAssignmentEq(t); ok {
		return isAttrProbeExpr(rhs)
	}
	return isAttrProbeExpr(t)
}

// isAttrProbeExpr reports an expression built around super().__getattribute__,
// hasattr, or getattr (including method chains on the getattribute result).
func isAttrProbeExpr(expr string) bool {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return false
	}
	// Strip a leading "await ".
	expr = strings.TrimSpace(strings.TrimPrefix(expr, "await "))
	if strings.Contains(expr, "super().__getattribute__(") ||
		strings.Contains(expr, "super().__getattr__(") {
		return true
	}
	if strings.HasPrefix(expr, "hasattr(") || strings.HasPrefix(expr, "getattr(") {
		return true
	}
	return false
}

// tryBlockIsOptionalFeatureRegistration reports a try block that only imports
// optional dependencies and registers defs/classes/fixtures — no side-effect
// calls at the try body's direct level (niquests conftest.py pytest_pyodide).
func tryBlockIsOptionalFeatureRegistration(lines []codeLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := indentWidth(lines[tryIdx].raw)
	bodyIndent := -1
	sawImport := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := indentWidth(lines[j].raw)
		if ind <= tryIndent {
			return sawImport
		}
		if bodyIndent < 0 {
			bodyIndent = ind
		}
		// Nested suite content under a def/class/fixture is part of the
		// registration body; only direct-body lines are constrained.
		if ind > bodyIndent {
			continue
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		if importLikeStmt(stmt) {
			sawImport = true
			continue
		}
		if strings.HasPrefix(stmt, "@") ||
			strings.HasPrefix(stmt, "def ") ||
			strings.HasPrefix(stmt, "async def ") ||
			strings.HasPrefix(stmt, "class ") {
			continue
		}
		return false
	}
	return sawImport
}

// probeCatch lists exception types whose pass handler is the documented
// optional-probe/parsing/termination fallback in the audited FP corpus
// (niquests, httptap, Project_Parva). KeyError is deliberately excluded
// (sync-with-uv TPs) as are Exception/BaseException (among-llms TPs) and
// UnicodeDecodeError / RuntimeError / queue types (wse TPs).
func probeCatch(types []string) bool {
	for _, ty := range types {
		switch ty {
		case "ValueError", "TypeError", "OSError", "AttributeError", "UnsupportedOperation",
			"StopIteration", "StopAsyncIteration", "CancelledError", "GeneratorExit",
			"TimeoutError", "RequestException":
			return true
		}
	}
	return false
}

// BP-PY-2: except suite is solely pass (optional comment already stripped).
func detectBPPY2(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-2")
	// Offline tooling / release scripts: pass-only warmup and batch-job
	// handlers (Project_Parva backend/tools/load_test) are audited residual FP.
	if isPythonOfflineScriptPath(unit) {
		return
	}
	if !facts.has("except") || !facts.has("pass") {
		return
	}
	// Masked lines: docstring prose cannot be a handler.
	lines := buildCodeLines(pytext.Mask(unit.Source))
	rawLines := codeLinesFacts(facts, unit.Source)
	isTest := isPythonTestFile(unit)
	isTestDir := isPythonTestDirPath(unit)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "except") || !strings.Contains(t, ":") {
			continue
		}
		// Must look like an except clause (not a variable named except_foo).
		if t != bareExceptClause && !strings.HasPrefix(t, "except ") && !strings.HasPrefix(t, bareExceptClause) {
			continue
		}
		exceptIndent := indentWidth(line.raw)
		// Collect suite statements at greater indent until dedent.
		var suite []string
		for j := i + 1; j < len(lines); j++ {
			st := strings.TrimSpace(lines[j].text)
			if st == "" {
				continue
			}
			ind := indentWidth(lines[j].raw)
			if ind <= exceptIndent {
				break
			}
			// Nested block headers still count as suite content.
			suite = append(suite, st)
			// Only consider immediate suite lines at the first indent level for "solely pass".
			// If we already have more than one statement, not solely pass.
			if len(suite) > 1 {
				break
			}
		}
		if len(suite) == 1 && suite[0] == "pass" && !exceptPassIsSafeBP(lines, rawLines, i, isTest, isTestDir) {
			pushAt(unit, meta, line.byte, "except handler body is only pass; failures are discarded silently", out)
		}
	}
}

var mutableDefaultRe = regexp.MustCompile(`=\s*(\[\s*\]|\{\s*\}|set\s*\(\s*\)|list\s*\(\s*\)|dict\s*\(\s*\))`)

// BP-PY-4: mutable default arguments [] {} set() list() dict()
func detectBPPY4(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-4")
	if !facts.has("def ") && !facts.has("async def ") {
		return
	}
	src := unit.Source
	// Scan for def / async def signatures (possibly multi-line until ')' before ':').
	lines := codeLinesFacts(facts, unit.Source)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if !strings.HasPrefix(t, "def ") && !strings.HasPrefix(t, "async def ") {
			continue
		}
		// Accumulate signature text until we see ')' that closes params and optional '->' then ':'.
		sig := t
		j := i
		// If already complete on one line
		for !signatureComplete(sig) && j+1 < len(lines) {
			j++
			sig += " " + strings.TrimSpace(lines[j].text)
			// Safety: stop after many lines
			if j-i > maxSignatureLines {
				break
			}
		}
		// Only look inside parentheses of the def.
		open := strings.Index(sig, "(")
		closeParen := strings.LastIndex(sig, ")")
		if open < 0 || closeParen <= open {
			continue
		}
		params := sig[open+1 : closeParen]
		if mutableDefaultRe.MatchString(params) {
			pushAt(unit, meta, line.byte, "mutable default argument is shared across calls; use None and assign inside the body", out)
		}
	}
	_ = src
}

func signatureComplete(sig string) bool {
	// Rough: has '(' and matching ')' and ends with ':' (possibly after return annotation).
	if !strings.Contains(sig, "(") {
		return false
	}
	depth := 0
	inStr := byte(0)
	escape := false
	for i := 0; i < len(sig); i++ {
		c := sig[i]
		if inStr != 0 {
			if escape {
				escape = false
				continue
			}
			if c == '\\' {
				escape = true
				continue
			}
			if c == inStr {
				inStr = 0
			}
			continue
		}
		switch c {
		case '"', '\'':
			inStr = c
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				// rest should contain ':'
				return strings.Contains(sig[i:], ":")
			}
		}
	}
	return false
}

// BP-PY-6: assert used for request/CLI/authz/path validation in non-test modules.
// Internal invariant asserts (no security/input needle) are intentionally missed.
func detectBPPY6(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-6")
	if isPythonTestFile(unit) {
		return
	}
	if !facts.has("assert ") && !strings.Contains(unit.Source, "assert ") {
		if !strings.Contains(unit.Source, "assert") {
			return
		}
	}
	lines := codeLinesFacts(facts, unit.Source)
	for _, line := range lines {
		t := strings.TrimSpace(line.text)
		if !(strings.HasPrefix(t, "assert ") || t == "assert" || strings.HasPrefix(t, "assert(")) {
			continue
		}
		if !assertLooksLikeRuntimeValidation(t) {
			continue
		}
		pushAt(unit, meta, line.byte, "assert is stripped with python -O; use if + raise for runtime validation", out)
	}
}

func assertLooksLikeRuntimeValidation(line string) bool {
	lower := strings.ToLower(line)
	// Intentionally no bare "request." / "request(" / "request " — those match
	// library PreparedRequest invariant checks (assert request.url is not None)
	// that are type/shape guards, not user-input or authz validation.
	needles := []string{
		// Request *input* attributes (Django/Flask/Starlette-style), not internal fields like .url/.method.
		"request.get", "request.post", "request.args", "request.form",
		"request.files", "request.data", "request.json", "request.user",
		"request.cookies", "request.query_params", "request.path_params",
		// Auth via headers is still validation; bare request.headers is not (niquests FP).
		"authorization",
		"args.", "form.", "files.",
		"is_authenticated", "is_anonymous", "has_perm", "permission", "authorize",
		"csrf", "safe_join",
		"filename", "filepath", "dirname",
		"sys.argv", "click.", "argparse",
		"user_input", "untrusted",
	}
	for _, n := range needles {
		if strings.Contains(lower, n) {
			return true
		}
	}
	// Word-ish tokens to avoid matching "author"/"apathy"/etc.
	for _, re := range []*regexp.Regexp{
		regexp.MustCompile(`\bauth\b`),
		regexp.MustCompile(`\brole\b`),
		regexp.MustCompile(`\btoken\b`),
		regexp.MustCompile(`\bpath\b`),
		regexp.MustCompile(`\bcli\b`),
		regexp.MustCompile(`\bargv\b`),
	} {
		if re.MatchString(lower) {
			return true
		}
	}
	return false
}

// BP-PY-7: builtin open( without with.
// Attribute methods (fitz.open, Image.open, self.open, Client.open, os.open)
// and function definitions (def open) are out of scope. Docstring/comment
// matches are blanked via pytext.Mask before the line scan.
func detectBPPY7(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-PY-7")
	if !facts.has("open(") && !facts.has(".open(") {
		return
	}
	masked := pytext.Mask(unit.Source)
	lines := buildCodeLines(masked)
	for _, line := range lines {
		t := line.text
		if !strings.Contains(t, "open(") {
			continue
		}
		if lineContainsWithOpen(t) {
			continue
		}
		trimmed := strings.TrimSpace(t)
		if strings.HasPrefix(trimmed, "from ") || strings.HasPrefix(trimmed, "import ") {
			continue
		}
		if strings.HasPrefix(trimmed, "def open(") || strings.HasPrefix(trimmed, "async def open(") {
			continue
		}
		if i := indexOfBareOpenCall(t); i >= 0 {
			pushAt(unit, meta, line.byte+i, "open without with risks resource leaks; use a context manager", out)
		}
	}
}

// indexOfBareOpenCall finds builtin open( — not mid-ident and not attribute .open(.
func indexOfBareOpenCall(line string) int {
	start := 0
	for {
		i := strings.Index(line[start:], "open(")
		if i < 0 {
			return -1
		}
		abs := start + i
		if abs > 0 {
			prev := line[abs-1]
			if prev == '.' || isIdentByte(prev) {
				start = abs + 4
				continue
			}
		}
		return abs
	}
}

// splitAssignmentEq splits a true assignment (=), rejecting ==, !=, <=, >=, :=.
func splitAssignmentEq(line string) (string, string, bool) {
	inStr := byte(0)
	esc := false
	triple := false
	depth := 0
	for i := 0; i < len(line); i++ {
		c := line[i]
		if inStr != 0 {
			i, inStr, esc, triple = advanceInString(line, i, inStr, esc, triple)
			continue
		}
		switch c {
		case '\'', '"':
			inStr = c
			if i+2 < len(line) && line[i+1] == c && line[i+2] == c {
				triple = true
				i += 2
			}
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			if depth > 0 {
				depth--
			}
		case '=':
			if depth != 0 {
				continue
			}
			if i > 0 {
				prev := line[i-1]
				if prev == '=' || prev == '!' || prev == '<' || prev == '>' || prev == ':' {
					continue
				}
			}
			if i+1 < len(line) && line[i+1] == '=' {
				continue
			}
			return line[:i], line[i+1:], true
		}
	}
	return "", "", false
}

// advanceInString steps one byte inside a string literal scanner state.
func advanceInString(line string, i int, inStr byte, esc, triple bool) (int, byte, bool, bool) {
	c := line[i]
	if triple {
		if c == inStr && i+2 < len(line) && line[i+1] == inStr && line[i+2] == inStr {
			return i + 2, 0, false, false
		}
		return i, inStr, esc, triple
	}
	if esc {
		return i, inStr, false, triple
	}
	if c == '\\' {
		return i, inStr, true, triple
	}
	if c == inStr {
		return i, 0, false, triple
	}
	return i, inStr, esc, triple
}

// pyBracketsUnbalanced reports a text with unclosed ( [ { brackets outside
// string literals — i.e. a statement that continues on the next line.
func pyBracketsUnbalanced(text string) bool {
	depth := 0
	inStr := byte(0)
	esc := false
	triple := false
	for i := 0; i < len(text); i++ {
		if inStr != 0 {
			var t bool
			i, inStr, esc, t = advanceInString(text, i, inStr, esc, triple)
			triple = t
			continue
		}
		c := text[i]
		if c == '"' || c == '\'' {
			inStr = c
			if i+2 < len(text) && text[i+1] == c && text[i+2] == c {
				triple = true
				i += 2
			}
			continue
		}
		switch c {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
		}
	}
	return depth > 0
}
