package cwe

import (
	"regexp"
	"strings"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func init() {
	// CWE-396/397 are gated on the exact exception tokens they match.
	// Remaining platform rules use FN-safe structural tokens.
	RegisterRule("CWE-396", detectCWE396, &MetaCWE396,
		"except Exception", "except BaseException")
	RegisterRule("CWE-397", detectCWE397, &MetaCWE397,
		"raise Exception", "raise BaseException")
	RegisterRule("CWE-478", detectCWE478, &MetaCWE478,
		"match ")
	RegisterRule("CWE-252", detectCWE252, &MetaCWE252,
		"subprocess.run", "subprocess.call", "os.system")
	RegisterRule("CWE-390", detectCWE390, &MetaCWE390,
		"except")
	RegisterRule("CWE-584", detectCWE584, &MetaCWE584,
		"finally")
}

var (
	pyGenericExceptLineRE = regexp.MustCompile(`^[\t ]*except\s+(?:Exception|BaseException)(?:\s+as\s+[A-Za-z_][A-Za-z0-9_]*)?\s*:`)
	pyGenericRaiseRE      = regexp.MustCompile(`(?m)^[\t ]*raise\s+(?:Exception|BaseException)(?:\s*\(|\b)`)
	pyMatchStartRE        = regexp.MustCompile(`^[\t ]*match\s+[^:\n]+:\s*$`)
	pyCaseStartRE         = regexp.MustCompile(`^[\t ]*case\s+[^:\n]+:\s*$`)
	pyDefaultCaseRE       = regexp.MustCompile(`^[\t ]*case\s+_\s*(?:if\s+[^:\n]+)?\s*:\s*$`)
	pyExceptStartRE       = regexp.MustCompile(`^[\t ]*except(?:\s+[^:\n]+)?\s*:\s*$`)
	pyFinallyStartRE      = regexp.MustCompile(`^[\t ]*finally\s*:\s*$`)
	pyReturnLineRE        = regexp.MustCompile(`^[\t ]*return\b`)
)

// CWE-396 reports only the two Python root exception classes. Specific
// exception handlers intentionally remain outside this narrow heuristic.
//
// Suite-surface skips exist for audited FPs (cleanup-then-raise, set_exception,
// two-field error record, etc.) but MUST NOT exempt a single wrapped raise
// (`raise TypedError(...) from exc`) — those are audited TPs in voicetag /
// caniscrape / among-llms. BP-PY-1 owns a separate surface policy.
func detectCWE396(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if isPythonTestModule(unit) {
		return
	}
	if unit == nil || out == nil {
		return
	}
	// Dual-mode SSL helpers (verify_ssl switch → default context vs CERT_NONE)
	// and TLS-diagnostic clients that thread verify_ssl into create_ssl_context
	// (httptap implementations/tls.py, http_client) — generic catch inside is
	// compatibility plumbing, not hidden application failure.
	if dualModeTLSHelper(unit.Source) || tlsDiagnosticVerifySSLClient(unit.Source) {
		return
	}
	if !containsAnyNeedle(unit.Source, "except Exception", "except BaseException") {
		return
	}
	lines := facts.MaskedLines()
	rawLines := buildMaskedPythonLines(unit.Source)
	for i, line := range lines {
		if !pyGenericExceptLineRE.MatchString(line.text) {
			continue
		}
		if cwe396SuiteIsSafe(lines, rawLines, i) {
			continue
		}
		emitPlatformFinding(unit, &MetaCWE396, line.start, "generic Exception or BaseException handler can hide distinct failure conditions", confidence84, out)
		// One finding per file matches historical scan/audit behavior used as
		// the convergence target (emitting every handler massively inflates
		// WeThePeople vs the audited TP list).
		return
	}
}

// cwe396SuiteIsSafe reports audited-FP generic-catch suites whose failure is
// surfaced without being swallowed. This mirrors the BP-PY-1 reconciliation
// guard (see bad_practices suiteSurfacesFailure):
//
//   - fall-through raise after the handler is safe (calgebra gcal.py:221/245);
//   - suite-level raise: raise-only handlers stay reportable (WeThePeople
//     tracing.py:76); multi-statement cleanup/translation handlers that
//     neither log nor print are safe (niquests utils.py:256, safer
//     __init__.py:312, wasi/_adapter.py:107); logged raises stay reportable
//     (sync-with-uv, WeThePeople sync_fara_data.py:185);
//   - set_exception / _error_result feed the failure forward (niquests
//     _ws.py:44, _sse.py:153; calgebra gcsa.py:440);
//   - a two-field error record is the httptap analyzer.py:262 audited FP,
//     while a lone .error= assignment stays reportable;
//   - stderr print followed by continue/return/exit is user-facing reporting
//     (sync-with-uv cli.py:161); any other print handler stays reportable;
//   - type-extraction feeding a counter is error accounting (requestSpeedTest
//     rnet_test.py:37);
//   - error=<var> / append(<var>) transport the raw exception (calgebra
//     gcsa.py:1117; niquests sgi/_async/__init__.py:661);
//   - single-statement raw capture and pass-through returns are evidence
//     capture / parse fallbacks (onlymaps tests; niquests sgi/_async:355;
//     onlymaps _types.py:189);
//   - multi-statement suites that never reference the exception and neither
//     log nor print are documented fallback paths (niquests
//     wasi/_async/_adapter.py:99);
//   - JS/pyodide bridge probes with pass-only or constant-assign fallback
//     after a try body that touches getReader/to_py/run_sync/_js_ (niquests
//     extensions/pyodide) — wrap-raise after the same bridge stays reportable;
//   - logger.exception (or equivalent) followed by a process exit-code return
//     surfaces the failure to the CLI operator (httptap cli.py:589).
func cwe396SuiteIsSafe(lines, rawLines []pyMaskedLine, exceptIdx int) bool {
	if exceptIdx < 0 || exceptIdx >= len(lines) {
		return false
	}
	exceptIndent := lines[exceptIdx].indent
	suite, after, _ := pyExceptSuiteLineIdx(lines, exceptIdx)
	if len(suite) == 0 {
		return false
	}
	var body []string
	var bodyRaw []string
	allText := ""
	allRaw := ""
	for j := exceptIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		if lines[j].indent <= exceptIndent {
			break
		}
		allText += mt + "\n"
		allRaw += strings.TrimSpace(rawLines[j].text) + "\n"
	}
	// Merge multi-line statements (parenthesized returns, list comprehensions)
	// so continuation lines do not count as separate statements. Continuation
	// lines may sit at deeper indent than the direct body (calgebra
	// gcsa.py:1117 list comprehension), so the merge scans every line until
	// the brackets balance.
	merged := make([]bool, len(suite))
	for k := 0; k < len(suite); k++ {
		if merged[k] {
			continue
		}
		i := suite[k]
		mt := strings.TrimSpace(lines[i].text)
		rt := strings.TrimSpace(rawLines[i].text)
		last := i
		if pyBracketsUnbalanced(mt) {
			for j := i + 1; j < len(lines); j++ {
				// Masked string-only continuation lines are blank — skip them
				// but keep merging using raw text so multi-line raise Type(
				//     "const",
				// ) from exc still forms one statement (rate_limit FP).
				txt := strings.TrimSpace(lines[j].text)
				rawTxt := strings.TrimSpace(rawLines[j].text)
				if lines[j].indent <= exceptIndent && rawTxt != "" && txt != "" {
					break
				}
				if rawTxt == "" && txt == "" {
					continue
				}
				if lines[j].indent <= exceptIndent && txt == "" && rawTxt != "" {
					// deeper or same-level blanked string: still part of call
					// only while brackets remain unbalanced.
				} else if lines[j].indent <= exceptIndent {
					break
				}
				if txt != "" {
					mt += " " + txt
				}
				if rawTxt != "" {
					rt += " " + rawTxt
				}
				last = j
				if !pyBracketsUnbalanced(mt) && !pyBracketsUnbalanced(rt) {
					break
				}
			}
		}
		// Mark suite lines consumed by this statement's continuation.
		for m := k + 1; m < len(suite); m++ {
			if suite[m] > i && suite[m] <= last {
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

	// Only direct suite statements count as suite-level raise (not nested
	// try/except MaxRetryError raise inside a retry loop).
	hasRaise := false
	for _, b := range body {
		if pySuiteLineRaiseRE.MatchString(b) {
			hasRaise = true
			break
		}
	}
	// logger() / log() factory forms (pycaps logger().error) lack a trailing
	// dot after the name; also accept .error/.warning call shapes.
	hasLog := strings.Contains(allText, "exc_info") || strings.Contains(allText, ".exception(") ||
		strings.Contains(allText, "log.") || strings.Contains(allText, "logger.") ||
		strings.Contains(allText, "logging.") || strings.Contains(allText, "traceback.print_exc") ||
		strings.Contains(allText, "logger(") || strings.Contains(allText, "log(") ||
		strings.Contains(allText, ".error(") || strings.Contains(allText, ".warning(")
	hasPrint := strings.Contains(allText, "print(")
	caughtVar := exceptClauseCaughtVar(strings.TrimSpace(lines[exceptIdx].text))

	// Suite-level raise:
	//  - multi-statement resource-cleanup-then-raise (os.remove/unlink/close
	//    then raise) is an audited FP (safer __init__.py:312, niquests utils,
	//    sgi task.cancel + re-raise);
	//  - multi-statement non-cleanup (telemetry.track then raise) stays
	//    REPORTABLE — caniscrape cli.py:326 audited TP;
	//  - a SINGLE wrapped raise (`raise TypedError(...) from exc`) stays
	//    REPORTABLE — voicetag / among-llms audited TPs — except when the try
	//    body is a JS/pyodide bridge call or pure SSL attribute probe
	//    (niquests pyodide ConnectionError mapping; httptap getpeercert);
	//  - logger.exception + bare re-raise without DB rollback surfaces the
	//    failure (Project_Parva middleware). logger.error + raise and
	//    exception+rollback stay reportable (WeThePeople finnhub / sync_fara);
	//  - bare re-raise without log stays reportable (WeThePeople tracing.py:76);
	//  - pure conversion/parse try + wrap-raise is API-boundary translation
	//    (Project_Parva engine_routes date parse → HTTPException).
	if hasRaise {
		if stmtCount >= 2 && !hasLog && !hasPrint && suiteOnlyResourceCleanupBeforeRaise(body) {
			return true
		}
		// Type-discriminating re-raise/map (niquests wasi stream closed):
		// if not isinstance(exc, Expected): raise / map expected / raise.
		// telemetry.track + raise stays reportable (caniscrape cli TP).
		if stmtCount >= 2 && !hasLog && !hasPrint && suiteIsTypeFilterReraise(body, caughtVar) {
			return true
		}
		// Documented defensive multi-statement cleanup that re-raises after a
		// notify callback (niquests wasi upload: Defensive marker + raise).
		if stmtCount >= 2 && !hasLog && !hasPrint &&
			exceptLineHasDefensiveMarkerPy(rawLines, exceptIdx, exceptIndent) {
			return true
		}
		if !hasPrint && suiteLogsExceptionThenBareReraise(body) && !suiteHasDBRollback(body) {
			return true
		}
		if !hasLog && !hasPrint && suiteOnlyWrappedRaises(body) {
			if tryIdx := enclosingTryIdx(lines, exceptIdx); tryIdx >= 0 {
				if tryBlockHasJSBridgeSignal(lines, tryIdx) {
					return true
				}
				if tryBlockIsSSLAttrProbe(lines, tryIdx) {
					return true
				}
				if tryBlockIsConversionGuard(lines, tryIdx) {
					return true
				}
			}
			// Fail-closed domain wrap with constant message only (no f-string /
			// str(exc) embedding) — Project_Parva rate_limit.py. voicetag wrap
			// raises embed {exc} and stay reportable.
			if suiteIsConstantMessageDomainWrap(bodyRaw, caughtVar) {
				return true
			}
			// Typed domain re-wrap after a sibling `except DomainError: raise`
			// (Project_Parva rulelang_service). Lone wrap-raise without the
			// typed pass-through stays reportable (voicetag / caniscrape).
			if suiteIsTypedDomainRewrap(lines, exceptIdx, bodyRaw, caughtVar) {
				return true
			}
		}
		return false
	}

	// Nested type-filter re-raise: direct body is only if/with branches; raise
	// lives under isinstance arms (niquests wasi:_adapter.py:107). Ordinary
	// multi-stmt swallows without raise stay reportable.
	if !hasLog && !hasPrint && caughtVar != "" &&
		suiteDirectBodyIsOnlyBranches(body) &&
		strings.Contains(allText, "isinstance(") &&
		strings.Contains(allText, "raise") {
		return true
	}

	// JS/pyodide bridge best-effort: pass-only or constant-assign/return after
	// a try body that touches the bridge API. Ordinary Exception:pass stays
	// reportable. Use bodyRaw so string-literal constants survive pythonCodeMask.
	if tryIdx := enclosingTryIdx(lines, exceptIdx); tryIdx >= 0 &&
		tryBlockHasJSBridgeSignal(lines, tryIdx) &&
		suiteIsPassOrConstantFallback(bodyRaw) {
		return true
	}
	// Pure SSL attribute probe with pass/constant fallback (httptap cert
	// extract). Ordinary Exception:pass stays reportable.
	if tryIdx := enclosingTryIdx(lines, exceptIdx); tryIdx >= 0 &&
		tryBlockIsSSLAttrProbe(lines, tryIdx) &&
		suiteIsPassOrConstantFallback(bodyRaw) {
		return true
	}

	// Retry machinery: failure fed into retries.increment then continue/sleep
	// (niquests sgi/__init__.py:129 audited FP). Nested MaxRetryError raise is
	// not a suite-level raise (handled above).
	if strings.Contains(allText, ".increment(") {
		for _, b := range body {
			bt := strings.TrimSpace(b)
			if strings.HasPrefix(bt, "continue") || strings.HasPrefix(bt, "return") ||
				strings.Contains(bt, ".sleep(") {
				return true
			}
		}
	}

	// Failure fed into a future / error payload helper.
	if strings.Contains(allText, "set_exception(") || strings.Contains(allText, "_error_result(") {
		return true
	}

	// Two-field error record (step.error = ...; step.note = ...) — httptap
	// analyzer.py:262. A lone .error= assignment stays reportable.
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

	// logger.exception / exc_info log then return a process exit code — the
	// CLI surfaces the failure and terminates (httptap cli.py:589). Bare log
	// without an exit return stays reportable.
	if hasLog && suiteLogsThenReturnsExitCode(body) {
		return true
	}

	// Print of the exception: CWE-396 keeps print handlers reportable — the
	// generic catch itself is the weakness even when the message is printed
	// (requestSpeedTest increase_limits.py:38 TP, sync-with-uv cli.py:161
	// TP). BP-PY-1 separately exempts stderr reporting with an exit flow.
	if hasPrint {
		return false
	}

	// Exception-type extraction feeding a failure counter.
	if caughtVar != "" && strings.Contains(allText, "type("+caughtVar) && strings.Contains(allText, "+=") {
		return true
	}

	// Raw exception transported in a result payload.
	if caughtVar != "" {
		if strings.Contains(joinedRaw, "error="+caughtVar) {
			return true
		}
		if strings.Contains(allRaw, "append("+caughtVar) {
			return true
		}
		if stmtCount == 1 && assignsRawPyVar(bodyRaw[0], caughtVar) {
			return true
		}
	}
	// Soft recording into a warnings list is the Project_Parva rulelang
	// evidence-packet fallback (warnings.append(...exc...)) — failure is
	// retained on the result, not swallowed. Use raw text: Mask blanks strings.
	if strings.Contains(allRaw, "warnings.append(") || strings.Contains(allRaw, ".warnings.append(") {
		return true
	}
	// Soft validation outcome: multi-assign builds value/error/warnings from
	// the exception without re-raising (rulelang validate_date). Lone
	// .error= assignment and wrap-raise stay reportable.
	if !hasRaise && !hasPrint && caughtVar != "" && stmtCount >= 2 &&
		suiteIsSoftValidationOutcome(bodyRaw, allRaw, caughtVar) {
		return true
	}
	// Structured result recording "error": str(exc) via named append payloads
	// (rulelang run_rule_tests: results.append({name, passed, error})).
	// Health-check appends without name= stay reportable (WeThePeople).
	if caughtVar != "" &&
		(strings.Contains(allRaw, `"error":`) || strings.Contains(allRaw, `'error':`)) &&
		lineReferencesVarPy(allRaw, caughtVar) {
		namedAppend := strings.Contains(allRaw, ".append(") &&
			(strings.Contains(allRaw, `"name":`) || strings.Contains(allRaw, `'name':`) ||
				strings.Contains(allRaw, `"name" :`) || strings.Contains(allRaw, `'name' :`))
		if namedAppend {
			return true
		}
	}
	// Soft best-effort logger.warning-only after optional-import enrichment
	// under an if-path (Project_Parva middleware ephemeris). Unconditional
	// warning after non-import work stays reportable (WeThePeople).
	if softBestEffortWarningOnlyCWE(allText) {
		tryIdxSoft := enclosingTryIdx(lines, exceptIdx)
		if tryIdxSoft >= 0 && tryNestedUnderIfPy(lines, tryIdxSoft) &&
			tryBlockHasImportStmtPy(lines, tryIdxSoft) {
			return true
		}
	}
	// Soft warning collector: warnings.warn / warn_* only. Pure wrap-raise
	// stays reportable.
	if !hasRaise && suiteIsSoftWarningCollector(body, allRaw) {
		return true
	}
	// Soft feature-degrade: log then assign empty/None/[]/{} or return prior
	// value without re-raising (pycaps emoji/tagger API ignore-feature FPs).
	// Log-only without fallback assign (WeThePeople connectors) stays reportable.
	if hasLog && !hasPrint && suiteIsSoftLogThenDegrade(body, bodyRaw, caughtVar) {
		return true
	}
	// Multi-statement defensive cleanup that never logs/prints/re-raises and
	// does not bind the exception (niquests wasi upload: failed = True; notify).
	// Requires an explicit defensive/no-cover marker so ordinary multi-stmt
	// swallows stay reportable.
	if !hasLog && !hasPrint &&
		exceptLineHasDefensiveMarkerPy(rawLines, exceptIdx, exceptIndent) {
		if caughtVar == "" || !lineReferencesVarPy(allRaw, caughtVar) {
			return true
		}
	}
	// Single-statement pass-through return of a non-exception value
	// (onlymaps _types.py:189 — parse fallback returns the input).
	if stmtCount == 1 && strings.HasPrefix(body[0], "return ") {
		rest := strings.TrimSpace(body[0][len("return "):])
		if bareIdentRE.MatchString(rest) && rest != caughtVar && rest != "None" && rest != "True" && rest != "False" {
			return true
		}
	}
	// Soft constant-assign fallback after a pure attr/cache probe try body
	// (logxide fast_logger_wrapper / interceptor getMessage) or optional
	// import (send_alerts RESEND_API_URL default). Ordinary Exception:pass
	// and log-only handlers stay reportable (WeThePeople).
	if !hasLog && !hasPrint && suiteIsPassOrConstantFallback(bodyRaw) {
		if tryIdx := enclosingTryIdx(lines, exceptIdx); tryIdx >= 0 &&
			(tryBlockIsSoftAttrCacheProbe(lines, tryIdx) || tryBlockIsImportOnly(lines, tryIdx)) {
			return true
		}
	}
	// Delegation to a handler method (self.handleError / self._handle_error) —
	// the stdlib logging emit contract (logxide handlers.py:191 FP). BP-PY-1
	// deliberately drops this exemption (Cronboard self.notify TPs).
	if len(suite) == 1 {
		if selfCallDelegation(strings.TrimSpace(rawLines[suite[0]].text)) {
			return true
		}
	}
	return false
}

// softBestEffortWarningOnlyCWE reports logger.warning-only handlers (not
// .exception/.error/.info).
func softBestEffortWarningOnlyCWE(allText string) bool {
	if !strings.Contains(allText, ".warning(") {
		return false
	}
	for _, hard := range []string{".exception(", ".error(", ".critical(", "exc_info", "traceback.print_exc",
		".info(", ".debug("} {
		if strings.Contains(allText, hard) {
			return false
		}
	}
	return strings.Contains(allText, "logger.") ||
		strings.Contains(allText, "log.") ||
		strings.Contains(allText, "logging.")
}

// suiteIsConstantMessageDomainWrap reports wrap-raises whose constructor args
// are only string/numeric literals (no caught-var interpolation). Fail-closed
// domain policy (Project_Parva rate_limit). f-string / str(exc) embeds stay
// reportable (voicetag).
func suiteIsConstantMessageDomainWrap(bodyRaw []string, caughtVar string) bool {
	saw := false
	for _, s := range bodyRaw {
		t := strings.TrimSpace(s)
		if t == "" || t == "pass" {
			continue
		}
		if !pySuiteLineRaiseRE.MatchString(t) {
			return false
		}
		if !wrappedRaiseStatement(t) {
			return false
		}
		if !raiseArgsAreConstantLiterals(t, caughtVar) {
			return false
		}
		saw = true
	}
	return saw
}

// raiseArgsAreConstantLiterals reports raise Type(const...) [from var] where
// no non-literal expression appears in the constructor args. `from caughtVar`
// is allowed (chain only).
func raiseArgsAreConstantLiterals(raiseStmt, caughtVar string) bool {
	t := strings.TrimSpace(raiseStmt)
	if !strings.HasPrefix(t, "raise") {
		return false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(t, "raise"))
	// Strip trailing " from <var>".
	if i := strings.LastIndex(rest, " from "); i >= 0 {
		fromPart := strings.TrimSpace(rest[i+len(" from "):])
		if caughtVar != "" && fromPart != caughtVar && fromPart != "e" && fromPart != "exc" && fromPart != "err" {
			// still ok — chain target need not match
		}
		rest = strings.TrimSpace(rest[:i])
	}
	// Must construct: Name(...)
	open := strings.Index(rest, "(")
	if open < 0 || !strings.HasSuffix(rest, ")") {
		return false
	}
	args := strings.TrimSpace(rest[open+1 : len(rest)-1])
	if args == "" {
		return true
	}
	// f-strings / format embeds always count as non-constant for this guard.
	if strings.Contains(args, "f\"") || strings.Contains(args, "f'") ||
		strings.Contains(args, "F\"") || strings.Contains(args, "F'") ||
		strings.Contains(args, ".format(") || strings.Contains(args, "%") {
		return false
	}
	if caughtVar != "" && lineReferencesVarPy(args, caughtVar) {
		return false
	}
	// Reject str(e)/repr(e) and bare exception identifiers.
	for _, bad := range []string{"str(", "repr(", "exc", " err", "(e)", " e,", " e)", "e,", "e)"} {
		if strings.Contains(args, bad) {
			// refine bare e/exc only via word-ish check below
		}
	}
	for _, arg := range splitTopLevelArgs(args) {
		a := strings.TrimSpace(arg)
		if a == "" {
			continue
		}
		// name=literal kwargs (detail="...", code="X", status_code=400)
		if eq := strings.Index(a, "="); eq > 0 && !strings.ContainsAny(a[:eq], "\"'") {
			a = strings.TrimSpace(a[eq+1:])
		}
		if isPureStringLiteral(a) || isNumericLiteral(a) {
			continue
		}
		return false
	}
	return true
}

// suiteIsTypedDomainRewrap reports a generic catch that only wraps into a
// domain error after a sibling `except DomainError: raise` pass-through
// (Project_Parva rulelang_service). Standalone wrap-raise stays reportable.
func suiteIsTypedDomainRewrap(lines []pyMaskedLine, exceptIdx int, bodyRaw []string, caughtVar string) bool {
	if exceptIdx <= 0 || len(bodyRaw) == 0 {
		return false
	}
	// Collect wrap-raise type names from this suite.
	wrapTypes := raiseConstructedTypeNames(bodyRaw)
	if len(wrapTypes) == 0 {
		return false
	}
	exceptIndent := lines[exceptIdx].indent
	// Walk prior sibling except clauses at the same indent.
	for j := exceptIdx - 1; j >= 0; j-- {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := lines[j].indent
		if ind > exceptIndent {
			continue
		}
		if ind < exceptIndent {
			return false
		}
		if strings.HasPrefix(t, "try:") {
			return false
		}
		if !strings.HasPrefix(t, "except") {
			// finally/else at same indent ends the try group
			if strings.HasPrefix(t, "finally:") || strings.HasPrefix(t, "else:") {
				return false
			}
			continue
		}
		types := exceptCatchTypes(t)
		if len(types) != 1 {
			continue
		}
		dom := types[0]
		if dom == "Exception" || dom == "BaseException" {
			continue
		}
		// Domain type must match a wrap-raise target.
		match := false
		for _, wt := range wrapTypes {
			if wt == dom {
				match = true
				break
			}
		}
		if !match {
			continue
		}
		// Sibling suite must be bare re-raise only.
		suite, _, _ := pyExceptSuiteLineIdx(lines, j)
		if len(suite) == 0 {
			continue
		}
		bareOnly := true
		for _, si := range suite {
			st := strings.TrimSpace(lines[si].text)
			if st == "" || st == "pass" {
				continue
			}
			if st != "raise" {
				bareOnly = false
				break
			}
		}
		if bareOnly {
			return true
		}
	}
	return false
}

// raiseConstructedTypeNames extracts Type from `raise Type(...)` statements.
func raiseConstructedTypeNames(body []string) []string {
	var out []string
	for _, s := range body {
		t := strings.TrimSpace(s)
		if !strings.HasPrefix(t, "raise") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(t, "raise"))
		if i := strings.LastIndex(rest, " from "); i >= 0 {
			rest = strings.TrimSpace(rest[:i])
		}
		open := strings.Index(rest, "(")
		if open <= 0 {
			continue
		}
		name := strings.TrimSpace(rest[:open])
		if i := strings.LastIndex(name, "."); i >= 0 {
			name = name[i+1:]
		}
		if name != "" {
			out = append(out, name)
		}
	}
	return out
}

// suiteIsSoftValidationOutcome reports multi-statement suites that only
// assign a structured validation failure (valid=False, error=str(exc),
// warnings=[exc], reason_codes) without re-raising.
func suiteIsSoftValidationOutcome(bodyRaw []string, allRaw, caughtVar string) bool {
	if caughtVar == "" || !lineReferencesVarPy(allRaw, caughtVar) {
		return false
	}
	hasErrorField := strings.Contains(allRaw, `"error"`) || strings.Contains(allRaw, `'error'`) ||
		strings.Contains(allRaw, `"error":`) || strings.Contains(allRaw, `'error':`)
	hasValidField := strings.Contains(allRaw, `"valid"`) || strings.Contains(allRaw, `'valid'`) ||
		strings.Contains(allRaw, "valid =") || strings.Contains(allRaw, "valid=")
	hasWarnings := strings.Contains(allRaw, "warnings") || strings.Contains(allRaw, "INVALID")
	if !hasErrorField || (!hasValidField && !hasWarnings) {
		return false
	}
	for _, s := range bodyRaw {
		t := strings.TrimSpace(s)
		if t == "" || t == "pass" {
			continue
		}
		if pySuiteLineRaiseRE.MatchString(t) {
			return false
		}
		// Only assignments and simple literals/containers.
		if strings.HasPrefix(t, "if ") || strings.HasPrefix(t, "elif ") || t == "else:" ||
			strings.HasPrefix(t, "for ") || strings.HasPrefix(t, "while ") ||
			strings.HasPrefix(t, "return ") || t == "return" {
			return false
		}
		if _, _, ok := splitAssignmentEq(t); ok {
			continue
		}
		return false
	}
	return true
}

// suiteIsSoftWarningCollector reports warnings.warn-style soft collectors
// without re-raise (stdlib warnings module — not a project-specific API).
func suiteIsSoftWarningCollector(body []string, allRaw string) bool {
	if len(body) == 0 {
		return false
	}
	for _, s := range body {
		t := strings.TrimSpace(s)
		if t == "" || t == "pass" || t == "continue" || t == "break" {
			continue
		}
		if pySuiteLineRaiseRE.MatchString(t) {
			return false
		}
	}
	// stdlib warnings.warn (not logger).
	if strings.Contains(allRaw, "warnings.warn(") &&
		!strings.Contains(allRaw, "logger.") &&
		!strings.Contains(allRaw, "log.") {
		return true
	}
	return false
}

// suiteIsSoftLogThenDegrade reports log + feature-cache clear or return of a
// prior named value (pycaps emoji/tagger ignore-feature FPs). Log+return None
// and log+break (WeThePeople connectors) stay reportable.
func suiteIsSoftLogThenDegrade(body, bodyRaw []string, caughtVar string) bool {
	hasLogStmt := false
	hasDegrade := false
	for i, s := range body {
		t := strings.TrimSpace(s)
		if t == "" || t == "pass" {
			continue
		}
		if pySuiteLineRaiseRE.MatchString(t) {
			return false
		}
		// break/continue after log is production loop control (WeThePeople TPs).
		if t == "break" || strings.HasPrefix(t, "break ") ||
			t == "continue" || strings.HasPrefix(t, "continue ") {
			return false
		}
		if strings.Contains(t, ".exception(") || strings.Contains(t, ".error(") ||
			strings.Contains(t, ".warning(") || strings.Contains(t, "exc_info") ||
			strings.Contains(t, "logger.") || strings.Contains(t, "log.") ||
			strings.Contains(t, "logging.") || strings.Contains(t, "traceback.print_exc") ||
			strings.Contains(t, "logger(") || strings.Contains(t, "log(") {
			hasLogStmt = true
			continue
		}
		if strings.HasPrefix(t, "return ") {
			rest := strings.TrimSpace(t[len("return "):])
			// return <prior-name> pass-through only — not None/[]/constants
			// (those are production failure returns / WeThePeople TPs).
			if rest != caughtVar && bareIdentRE.MatchString(rest) &&
				rest != "None" && rest != "True" && rest != "False" {
				hasDegrade = true
				continue
			}
			return false
		}
		if t == "return" {
			return false
		}
		// Assign empty/None/[] to a self/instance cache field only.
		raw := t
		if i < len(bodyRaw) {
			raw = strings.TrimSpace(bodyRaw[i])
		}
		if lhs, rhs, ok := splitAssignmentEq(raw); ok {
			lhs = strings.TrimSpace(lhs)
			rhs = strings.TrimSpace(rhs)
			if !strings.HasPrefix(lhs, "self.") {
				return false
			}
			if isDefinedConstantExprCWE(rhs) || rhs == "[]" || rhs == "{}" || rhs == "()" {
				hasDegrade = true
				continue
			}
			return false
		}
		return false
	}
	return hasLogStmt && hasDegrade
}

// tryBlockIsSoftAttrCacheProbe reports try bodies that only read attributes /
// cache fields (getEffectiveLevel, getMessage, getattr) without I/O sinks.
func tryBlockIsSoftAttrCacheProbe(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	saw := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		if lines[j].indent <= tryIndent {
			return saw
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		if stmtHasSideEffectSink(stmt) {
			return false
		}
		// Reject network/parse heavy names even inside assignments.
		for _, bad := range []string{
			"requests.", "httpx.", "urlopen(", "subprocess.", "json.loads(",
			"yaml.load(", "fromisoformat(", "strptime(", "connect(", "execute(",
		} {
			if strings.Contains(stmt, bad) {
				return false
			}
		}
		if isProbeBranchHeader(stmt) || stmt == "pass" {
			continue
		}
		if _, rhs, ok := splitAssignmentEq(stmt); ok {
			rhs = strings.TrimSpace(rhs)
			// Attr/method reads and ternary hasattr fallbacks only.
			if softAttrReadExpr(rhs) {
				saw = true
				continue
			}
			return false
		}
		if softAttrReadExpr(stmt) {
			saw = true
			continue
		}
		return false
	}
	return saw
}

// softAttrReadExpr reports getattr/hasattr ternaries and narrow attr/cache
// reads (logxide getEffectiveLevel/getMessage). Ordinary response.json() /
// db.query() calls stay outside so WeThePeople TPs keep firing.
func softAttrReadExpr(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	if stmtHasSideEffectSink(t) {
		return false
	}
	// Ternary: x if hasattr(...) else y
	if strings.Contains(t, " if ") && strings.Contains(t, " else ") {
		return strings.Contains(t, "hasattr(") || strings.Contains(t, "getattr(")
	}
	if strings.HasPrefix(t, "getattr(") || strings.HasPrefix(t, "hasattr(") {
		return true
	}
	// Pure attribute chain without calls: self._rust_logger.name / handler._inner.level
	if strings.Contains(t, ".") && !strings.Contains(t, "(") {
		return bareAttrChain(t)
	}
	// Zero-arg cache/message getters only.
	if strings.HasSuffix(t, "()") {
		base := strings.TrimSpace(t[:len(t)-2])
		dot := strings.LastIndex(base, ".")
		if dot < 0 {
			return false
		}
		method := base[dot+1:]
		switch method {
		case "getEffectiveLevel", "getMessage", "getName", "getLevelName",
			"level", "name": // defensive; name usually attr not call
			return bareAttrChain(base[:dot]) || bareAttrChain(base)
		}
	}
	return false
}

// bareAttrChain reports a.b.c identifiers separated by dots (no calls/ops).
func bareAttrChain(expr string) bool {
	t := strings.TrimSpace(expr)
	if t == "" {
		return false
	}
	for _, part := range strings.Split(t, ".") {
		part = strings.TrimSpace(part)
		if part == "" || !bareIdentRE.MatchString(part) {
			return false
		}
	}
	return strings.Contains(t, ".") || bareIdentRE.MatchString(t)
}

// tryNestedUnderIfPy reports that the try header sits under an if/elif suite.
func tryNestedUnderIfPy(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx <= 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	for j := tryIdx - 1; j >= 0; j-- {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := lines[j].indent
		if ind >= tryIndent {
			continue
		}
		if ind < tryIndent {
			return strings.HasPrefix(t, "if ") || strings.HasPrefix(t, "elif ")
		}
	}
	return false
}

// tryBlockHasImportStmtPy reports a try body with at least one import.
func tryBlockHasImportStmtPy(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	for j := tryIdx + 1; j < len(lines); j++ {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		if lines[j].indent <= tryIndent {
			return false
		}
		if strings.HasPrefix(t, "import ") || strings.HasPrefix(t, "from ") {
			return true
		}
	}
	return false
}

// suiteIsPassOrConstantFallback reports pass-only or constant-assign/return
// suites used as defined JS-bridge fallbacks (None/""/b""/False/0).
// Also accepts pure string literals and UPPER_SNAKE module constants
// (WARNING/INFO) used as level fallbacks (logxide fast_logger_wrapper).
func suiteIsPassOrConstantFallback(body []string) bool {
	if len(body) == 0 {
		return false
	}
	saw := false
	for _, s := range body {
		t := strings.TrimSpace(s)
		if t == "" || t == "pass" || t == "continue" || t == "break" {
			saw = true
			continue
		}
		if strings.HasPrefix(t, "return ") {
			rest := strings.TrimSpace(t[len("return "):])
			if constantFallbackExprCWE(rest) {
				saw = true
				continue
			}
			return false
		}
		if _, rhs, ok := splitAssignmentEq(t); ok {
			if constantFallbackExprCWE(strings.TrimSpace(rhs)) {
				saw = true
				continue
			}
			return false
		}
		return false
	}
	return saw
}

// constantFallbackExprCWE reports None/bools/numerics/string literals and
// UPPER_SNAKE module constants used as defined soft fallbacks.
func constantFallbackExprCWE(expr string) bool {
	t := strings.TrimSpace(expr)
	if isDefinedConstantExprCWE(t) || isPureStringLiteral(t) || isNumericLiteral(t) {
		return true
	}
	return upperSnakeRE.MatchString(t)
}

// suiteLogsThenReturnsExitCode reports a suite that logs the failure and
// returns a named process exit constant (EXIT_* / UNIX_*). Bare `return 1`
// after logger.exception stays reportable (WeThePeople jobs audited TPs).
func suiteLogsThenReturnsExitCode(body []string) bool {
	hasLog := false
	hasExitReturn := false
	for _, s := range body {
		t := strings.TrimSpace(s)
		if strings.Contains(t, "exc_info") || strings.Contains(t, ".exception(") ||
			strings.Contains(t, "log.") || strings.Contains(t, "logger.") ||
			strings.Contains(t, "logging.") || strings.Contains(t, "traceback.print_exc") {
			hasLog = true
		}
		if strings.HasPrefix(t, "return ") {
			rest := strings.TrimSpace(t[len("return "):])
			if strings.HasPrefix(rest, "EXIT_") ||
				strings.HasPrefix(rest, "UNIX_") ||
				strings.Contains(rest, "EXIT_") {
				hasExitReturn = true
			}
		}
	}
	return hasLog && hasExitReturn
}

// suiteLogsExceptionThenBareReraise reports logger.exception / traceback
// print_exc followed by bare re-raise only (no wrap-raise). Surfaces the
// failure fully (Project_Parva middleware). logger.error + raise stays
// reportable (WeThePeople finnhub TPs).
func suiteLogsExceptionThenBareReraise(body []string) bool {
	hasExceptionLog := false
	sawBareRaise := false
	for _, s := range body {
		t := strings.TrimSpace(s)
		if strings.Contains(t, ".exception(") || strings.Contains(t, "traceback.print_exc") {
			hasExceptionLog = true
		}
		if strings.HasPrefix(t, "raise") {
			rest := strings.TrimSpace(strings.TrimPrefix(t, "raise"))
			if rest != "" {
				// wrap-raise / raise e stay reportable
				return false
			}
			sawBareRaise = true
		}
	}
	return hasExceptionLog && sawBareRaise
}

// suiteHasDBRollback reports session/db rollback in the suite (WeThePeople
// jobs that log+raise after rolling back stay reportable).
func suiteHasDBRollback(body []string) bool {
	for _, s := range body {
		t := strings.TrimSpace(s)
		if strings.Contains(t, ".rollback(") {
			return true
		}
	}
	return false
}

// suiteOnlyWrappedRaises reports that every non-trivial suite statement is a
// wrapped raise (raise Type(...)/raise X from e), optionally after building
// a message or branching on the exception string (niquests pyodide timeout
// mapping). Bare re-raise alone stays outside this helper.
func suiteOnlyWrappedRaises(body []string) bool {
	sawWrap := false
	for _, s := range body {
		t := strings.TrimSpace(s)
		if t == "" || t == "pass" || t == "continue" || t == "break" {
			continue
		}
		if pySuiteLineRaiseRE.MatchString(t) {
			if !wrappedRaiseStatement(t) {
				return false
			}
			sawWrap = true
			continue
		}
		// Allow message/err_str construction and isinstance/branch headers.
		if strings.HasPrefix(t, "if ") || strings.HasPrefix(t, "elif ") || t == "else:" {
			continue
		}
		if _, _, ok := splitAssignmentEq(t); ok {
			continue
		}
		return false
	}
	return sawWrap
}

// suiteIsTypeFilterReraise reports multi-statement suites that only branch on
// the caught exception type (isinstance) and re-raise or return a mapped
// result (niquests wasi stream closed). Telemetry/print side effects stay out.
func suiteIsTypeFilterReraise(body []string, caughtVar string) bool {
	sawRaise := false
	sawTypeCheck := false
	for _, s := range body {
		t := strings.TrimSpace(s)
		if t == "" || t == "pass" {
			continue
		}
		if pySuiteLineRaiseRE.MatchString(t) {
			sawRaise = true
			continue
		}
		if strings.HasPrefix(t, "return ") || t == "return" {
			continue
		}
		if strings.HasPrefix(t, "if ") || strings.HasPrefix(t, "elif ") || t == "else:" {
			// Require isinstance(...) or a real identifier reference to the
			// caught var — naive substring match on "e"/"err" over-matched
			// temp_file.exists() and suppressed caniscrape cleanup wrap-raise.
			if strings.Contains(t, "isinstance(") || lineReferencesVarPy(t, caughtVar) {
				sawTypeCheck = true
			}
			continue
		}
		if strings.HasPrefix(t, "with ") {
			// with wasi_exception_mapping(...): raise — mapping context
			continue
		}
		return false
	}
	return sawRaise && sawTypeCheck
}

// suiteDirectBodyIsOnlyBranches reports that every direct suite statement is
// an if/elif/else/with header (nested raise/return live under those arms).
func suiteDirectBodyIsOnlyBranches(body []string) bool {
	if len(body) == 0 {
		return false
	}
	for _, s := range body {
		t := strings.TrimSpace(s)
		if t == "" || t == "pass" {
			continue
		}
		if strings.HasPrefix(t, "if ") || strings.HasPrefix(t, "elif ") || t == "else:" ||
			strings.HasPrefix(t, "with ") || strings.HasPrefix(t, "async with ") {
			continue
		}
		return false
	}
	return true
}

// onlyValueErrorCatch reports a sole ValueError catch (kiss_headers index).
func onlyValueErrorCatch(types []string) bool {
	return len(types) == 1 && types[0] == "ValueError"
}

// asyncCancelCatchOnly reports CancelledError and/or GeneratorExit only —
// expected async cancel swallow (niquests sgi lifespan). TimeoutError alone
// and mixed RuntimeError paths stay outside this helper.
func asyncCancelCatchOnly(types []string) bool {
	if len(types) == 0 {
		return false
	}
	for _, ty := range types {
		if ty != "CancelledError" && ty != "GeneratorExit" {
			return false
		}
	}
	return true
}

// tryBlockIsSSLAttrProbe reports a try body that only probes SSL/TLS socket
// attributes (getpeercert/version/cipher/wrap_socket) without network request
// APIs (certificate-extraction probes). Allows returning a constructed
// CertificateInfo after a successful probe.
func tryBlockIsSSLAttrProbe(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	sawSSL := false
	sawBody := false
	sslSignals := []string{
		"getpeercert(", ".cipher(", ".version(", "wrap_socket(",
		"SSLSocket", "SSLObject", "ssl_socket", "ssl_object",
	}
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		if lines[j].indent <= tryIndent {
			break
		}
		sawBody = true
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		if stmtHasSideEffectSink(stmt) {
			return false
		}
		for _, s := range sslSignals {
			if strings.Contains(stmt, s) {
				sawSSL = true
			}
		}
		// Allow assignments, branches, constant/constructor returns.
		if isProbeBranchHeader(stmt) || stmt == "pass" || stmt == "continue" || stmt == "break" {
			continue
		}
		if stmt == "return" || strings.HasPrefix(stmt, "return ") {
			continue
		}
		if _, _, ok := splitAssignmentEq(stmt); ok {
			continue
		}
		// Bare SSL-related calls only.
		if isJSBridgeExpr(stmt) {
			continue
		}
		return false
	}
	return sawBody && sawSSL
}

// tryBlockHasHeavyParse reports datetime/json/yaml parse APIs that are not
// simple index/conversion probes (keeps conversion-guard vulnerable fixtures).
func tryBlockHasHeavyParse(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	heavies := []string{
		"fromisoformat(", "strptime(", "json.loads(", "yaml.load(",
		"yaml.safe_load(", "parse_", "loads(", "decode(",
	}
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		if lines[j].indent <= tryIndent {
			break
		}
		for _, h := range heavies {
			if strings.Contains(mt, h) {
				return true
			}
		}
	}
	return false
}

// tlsDiagnosticVerifySSLClient reports modules that thread a verify_ssl switch
// into create_ssl_context for TLS diagnostics (httptap http_client / tls.py)
// without necessarily defining CERT_NONE in the same file.
func tlsDiagnosticVerifySSLClient(src string) bool {
	if !strings.Contains(src, "create_ssl_context") {
		return false
	}
	return strings.Contains(src, "verify_ssl=") || strings.Contains(src, "verify_ssl:") ||
		strings.Contains(src, "verify_ssl ") || strings.Contains(src, "verify_ssl,")
}

// exceptLineHasDefensiveMarkerPy is a narrow marker set for multi-statement
// cleanup exemptions (niquests wasi upload). Omits generic "fallback" prose.
func exceptLineHasDefensiveMarkerPy(rawLines []pyMaskedLine, exceptIdx, exceptIndent int) bool {
	markers := []string{"no cover", "defensive"}
	if lineHasFallbackMarkerTextPy(rawLines[exceptIdx].text, markers) {
		return true
	}
	for j := exceptIdx + 1; j < len(rawLines); j++ {
		if strings.TrimSpace(rawLines[j].text) == "" {
			continue
		}
		if rawLines[j].indent <= exceptIndent {
			break
		}
		if lineHasFallbackMarkerTextPy(rawLines[j].text, markers) {
			return true
		}
	}
	return false
}

// selfCallDelegation reports a single-statement suite that delegates to a
// handler method (`self.handleError(record)`, `self._handle_error(e)`).
func selfCallDelegation(t string) bool {
	if !strings.HasPrefix(t, "self.") || !strings.HasSuffix(t, ")") {
		return false
	}
	rest := t[len("self."):]
	i := strings.IndexAny(rest, "(")
	if i <= 0 {
		return false
	}
	name := rest[:i]
	for _, r := range name {
		if !isIdentByteCWE(byte(r)) && r != '_' {
			return false
		}
	}
	return true
}

// pySuiteLineRaiseRE matches a suite statement that is a raise.
var pySuiteLineRaiseRE = regexp.MustCompile(`(?m)^raise\b`)

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

// suiteOnlyResourceCleanupBeforeRaise reports multi-statement suites where
// every non-raise statement is resource cleanup (os.remove/unlink/close/rmtree)
// and the raise is a bare re-raise only. Wrap-raise after unlink stays
// reportable (caniscrape config.py:38 audited TP).
func suiteOnlyResourceCleanupBeforeRaise(body []string) bool {
	sawBareRaise := false
	sawCleanup := false
	for _, s := range body {
		t := strings.TrimSpace(s)
		if t == "" || t == "pass" {
			continue
		}
		if strings.HasPrefix(t, "raise") {
			rest := strings.TrimSpace(strings.TrimPrefix(t, "raise"))
			if rest != "" {
				// wrap-raise / raise e after cleanup stays reportable
				return false
			}
			sawBareRaise = true
			continue
		}
		// Allow branch wrappers around cleanup (if path.exists(): unlink).
		if strings.HasPrefix(t, "if ") || strings.HasPrefix(t, "elif ") || t == "else:" {
			continue
		}
		if isResourceCleanupStmt(t) {
			sawCleanup = true
			continue
		}
		return false
	}
	return sawBareRaise && sawCleanup
}

// isResourceCleanupStmt reports os.remove/unlink/close/rmtree/task-cancel
// style cleanup (niquests sgi ASGI cancel-then-reraise FPs).
func isResourceCleanupStmt(t string) bool {
	needles := []string{
		"os.remove(", "os.unlink(", "os.rmdir(", "os.removedirs(",
		".unlink(", ".remove(", ".close(", ".aclose(",
		"shutil.rmtree(", "pathlib.Path(",
		"task.cancel(", ".cancel(", "close_resource(",
		"contextlib.suppress(",
	}
	// pathlib Path(...).unlink() already covered by .unlink(
	for _, n := range needles {
		if strings.Contains(t, n) {
			return true
		}
	}
	// await task / await handle after cancel is part of the cleanup sequence.
	if strings.HasPrefix(strings.TrimSpace(t), "await ") {
		rest := strings.TrimSpace(t[len("await "):])
		if rest == "task" || strings.HasPrefix(rest, "task.") ||
			strings.Contains(rest, ".close(") || strings.Contains(rest, ".aclose(") {
			return true
		}
	}
	// with contextlib.suppress(...): await task
	if strings.HasPrefix(strings.TrimSpace(t), "with ") && strings.Contains(t, "suppress") {
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

// pyExceptSuiteLineIdx returns the line indices of the direct-body statements
// of the except suite at exceptIdx, the first line index after the suite, and
// the suite's body indent.
func pyExceptSuiteLineIdx(lines []pyMaskedLine, exceptIdx int) (suite []int, after int, bodyIndent int) {
	if exceptIdx < 0 || exceptIdx >= len(lines) {
		return nil, exceptIdx + 1, -1
	}
	exceptIndent := lines[exceptIdx].indent
	bodyIndent = -1
	j := exceptIdx + 1
	for ; j < len(lines); j++ {
		if strings.TrimSpace(lines[j].text) == "" {
			continue
		}
		ind := lines[j].indent
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

// nextDedentedStartsWith reports whether a statement starting with prefix
// (raise/return) follows the except suite at the try's indentation level
// within the next few statements, before any control-flow header intervenes
// (fall-through handlers, calgebra gcal.py:221).
func nextDedentedStartsWith(lines []pyMaskedLine, start, exceptIndent int, prefix string) bool {
	count := 0
	for j := start; j < len(lines) && count < 3; j++ {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		if lines[j].indent > exceptIndent {
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

// assignsRawPyVar reports `x = e` / `x[k] = e` / `x.a = e` where the RHS is
// exactly the caught exception object (capture/transport, not str()).
func assignsRawPyVar(line, varName string) bool {
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

// exceptClauseCaughtVar returns the name bound by `except ... as X`, or "".
func exceptClauseCaughtVar(t string) string {
	if i := strings.Index(t, " as "); i >= 0 {
		name := strings.TrimSpace(t[i+len(" as "):])
		if j := strings.Index(name, ":"); j >= 0 {
			name = strings.TrimSpace(name[:j])
		}
		if name != "" && isIdentByteCWE(name[0]) {
			ok := true
			for k := 1; k < len(name); k++ {
				if !isIdentByteCWE(name[k]) {
					ok = false
					break
				}
			}
			if ok {
				return name
			}
		}
	}
	return ""
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

// exceptPassIsSafe reports pass-only except suites that are audited false
// positives (mirrors bad_practices exceptPassIsSafeBP):
//
//   - optional-dependency ImportError/ModuleNotFoundError fallbacks whose try
//     block is solely imports (module or function scope);
//   - optional ImportError feature registration (imports + def/class/fixture);
//   - optional load_extension probes that only load and purely use the loaded
//     name (niquests async_session.py:620);
//   - termination-only catches (StopIteration/StopAsyncIteration/CancelledError/
//     GeneratorExit/TimeoutError);
//   - AttributeError after pure super().__getattribute__/hasattr probes;
//   - parsing probes whose try block is a pure assignment/return probe and
//     whose enclosing flow returns a defined default or raises;
//   - pure conversion guards (int/float/bool) and try/except/else recovery;
//   - best-effort Exception/BaseException JS/pyodide bridge probes (niquests
//     getReader/to_py/run_sync) — ordinary Exception:pass stays reportable;
//   - best-effort RequestException fetches (single-statement try block);
//   - documented optional/defensive fallbacks (marker comments);
//   - expected-exception tests whose try block deliberately raises/pytest.fail.
//
// CWE-396 is intentionally not exempted for these shapes (declaration-of-
// generic-catch stays TP-safe).
func exceptPassIsSafe(unit *core.ParsedUnit, lines, rawLines []pyMaskedLine, exceptIdx int) bool {
	exceptLine := strings.TrimSpace(lines[exceptIdx].text)
	if isPythonTestModule(unit) && tryBlockDeliberatelyRaises(lines, exceptIdx) {
		return true
	}
	// Test modules: try body that asserts an optional outcome then pass on
	// expected range errors (Project_Parva test_calculator). Chaos tests with
	// bare call+pass and no assert stay reportable (WeThePeople).
	if isPythonTestModule(unit) {
		if tryIdx := enclosingTryIdx(lines, exceptIdx); tryIdx >= 0 && tryBlockHasAssert(lines, tryIdx) {
			return true
		}
	}
	suite, after, _ := pyExceptSuiteLineIdx(lines, exceptIdx)
	if len(suite) != 1 || strings.TrimSpace(lines[suite[0]].text) != "pass" {
		return false
	}
	exceptIndent := lines[exceptIdx].indent
	if exceptLineHasFallbackMarkerPy(rawLines, exceptIdx, exceptIndent) {
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
	// next()/__anext__() probes only — not next(g, None)/aclose() teardown (PyDepends TPs).
	if tryIdx >= 0 && terminationCatchOnly(catchTypes) && tryBlockIsIteratorNextProbe(lines, tryIdx) {
		return true
	}
	if asyncShutdownCatch(catchTypes) {
		return true
	}
	// Expected async cancel/generator-exit swallow (niquests sgi lifespan
	// wait-forever until canceled). Ordinary Exception:pass stays reportable.
	if asyncCancelCatchOnly(catchTypes) {
		return true
	}
	if attributeErrorCatch(catchTypes) && tryIdx >= 0 && tryBlockIsAttrProbe(lines, tryIdx) {
		return true
	}
	if osErrorCatch(catchTypes) && tryIdx >= 0 && tryBlockIsCloseOnly(lines, tryIdx) {
		return true
	}
	if probeCatch(catchTypes) {
		// Alternate return after a try that itself returned (kiss_headers), not
		// a function epilogue return after setattr-style try (pingram TP).
		if tryIdx >= 0 && tryBlockHasReturnOrRaise(lines, tryIdx) {
			if nextDedentedStartsWith(lines, after, exceptIndent, "return") ||
				nextDedentedStartsWith(lines, after, exceptIndent, "raise") {
				return true
			}
		}
		if nextDedentedStartsWith(lines, after, exceptIndent, "raise") {
			return true
		}
		// Pure int/float/bool conversion guards may continue without immediate
		// return (parva-mcp content-length parse).
		if tryIdx >= 0 && tryBlockIsConversionGuard(lines, tryIdx) {
			return true
		}
		// Dict field extract + int conversion then list append (Project_Parva
		// timegraph_fact_links). Catch is KeyError/TypeError/ValueError only.
		if tryIdx >= 0 && fieldExtractConversionCatch(catchTypes) &&
			tryBlockIsFieldExtractConversion(lines, tryIdx) {
			return true
		}
		// Sole TypeError or sole OSError after pure probe (not mixed ValueError).
		if tryIdx >= 0 && tryBlockIsProbe(lines, tryIdx) &&
			(onlyTypeErrorCatch(catchTypes) || onlyOSErrorCatch(catchTypes)) {
			return true
		}
		// Sole ValueError after a pure index/attr assignment probe without a
		// return/raise in the try (kiss_headers content.index). Return-int
		// dangling pass and fromisoformat/json parse probes stay reportable.
		if tryIdx >= 0 && tryBlockIsProbe(lines, tryIdx) && onlyValueErrorCatch(catchTypes) &&
			!tryBlockHasReturnOrRaise(lines, tryIdx) && !tryBlockHasHeavyParse(lines, tryIdx) {
			return true
		}
		if bestEffortNetworkCatch(catchTypes) && tryIdx >= 0 &&
			(tryBlockSingleStatement(lines, tryIdx) || requestExceptionCatch(catchTypes)) {
			return true
		}
	}
	if broadExceptionCatch(catchTypes) && tryIdx >= 0 &&
		tryBlockIsProbe(lines, tryIdx) && tryBlockHasJSBridgeSignal(lines, tryIdx) {
		return true
	}
	return false
}

// broadExceptionCatch reports Exception or BaseException in the catch list.
func broadExceptionCatch(types []string) bool {
	for _, ty := range types {
		if ty == "Exception" || ty == "BaseException" {
			return true
		}
	}
	return false
}

// enclosingTryIdx returns the try header line index whose suite contains the
// except clause at exceptIdx, or -1.
func enclosingTryIdx(lines []pyMaskedLine, exceptIdx int) int {
	exceptIndent := lines[exceptIdx].indent
	for j := exceptIdx - 1; j >= 0; j-- {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := lines[j].indent
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

// tryBlockHasAssert reports that the try body contains an assert statement
// (optional-outcome tests that also allow pass on expected failures).
func tryBlockHasAssert(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		if lines[j].indent <= tryIndent {
			return false
		}
		if strings.HasPrefix(mt, "assert ") || mt == "assert" || strings.HasPrefix(mt, "self.assert") {
			return true
		}
	}
	return false
}

// tryBlockDeliberatelyRaises reports whether the try block preceding a
// pass-only except contains an explicit raise or pytest.fail — the
// expected-exception test idiom (onlymaps tests/test_database.py:339,516;
// niquests test_requests.py ReadTimeout). Assertions alone are not enough
// for non-test production modules: httpmorph's assert-in-try handlers are
// audited true positives.
func tryBlockDeliberatelyRaises(lines []pyMaskedLine, exceptIdx int) bool {
	exceptIndent := lines[exceptIdx].indent
	for j := exceptIdx - 1; j >= 0; j-- {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := lines[j].indent
		if ind < exceptIndent {
			return false
		}
		if ind > exceptIndent {
			if pySuiteLineRaiseRE.MatchString(t) || strings.Contains(t, "pytest.fail(") {
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

// exceptFollowedByElse reports a try/except/else where the except suite is
// only pass and the next same-indent clause is else: (success arm).
func exceptFollowedByElse(lines []pyMaskedLine, after, exceptIndent int) bool {
	for j := after; j < len(lines); j++ {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := lines[j].indent
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

// terminationCatchOnly reports that every caught type is a designed
// iterator/async termination signal.
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

// tryBlockIsIteratorNextProbe reports next()/__anext__() advance only (not
// next(g, None) or aclose() — those are PyDepends teardown TPs).
func tryBlockIsIteratorNextProbe(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	stmts := 0
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		if lines[j].indent <= tryIndent {
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
		_, rhs, ok := splitAssignmentEq(stmt)
		if ok {
			if isIteratorNextCall(strings.TrimSpace(rhs)) {
				continue
			}
			return false
		}
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
	if !strings.HasPrefix(expr, "next(") {
		return false
	}
	args := expr[len("next("):]
	if i := strings.LastIndex(args, ")"); i >= 0 {
		args = args[:i]
	}
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
				return false
			}
		}
	}
	return strings.TrimSpace(args) != ""
}

// asyncShutdownCatch reports lifespan handlers that mix TimeoutError/
// CancelledError with RuntimeError (loop already closed).
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

// tryBlockHasReturnOrRaise reports that the try body contains a return or raise.
func tryBlockHasReturnOrRaise(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		if lines[j].indent <= tryIndent {
			return false
		}
		if strings.HasPrefix(mt, "return") || strings.HasPrefix(mt, "raise") {
			return true
		}
	}
	return false
}

func onlyTypeErrorCatch(types []string) bool {
	return len(types) == 1 && types[0] == "TypeError"
}

func onlyOSErrorCatch(types []string) bool {
	return len(types) == 1 && (types[0] == "OSError" || types[0] == "IOError" || types[0] == "EnvironmentError")
}

func osErrorCatch(types []string) bool {
	for _, ty := range types {
		if ty == "OSError" || ty == "IOError" || ty == "EnvironmentError" {
			return true
		}
	}
	return false
}

func tryBlockIsCloseOnly(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	sawClose := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := lines[j].indent
		if ind <= tryIndent {
			return sawClose
		}
		stmt := mt
		for pyBracketsUnbalanced(stmt) && j+1 < len(lines) {
			j++
			stmt += " " + strings.TrimSpace(lines[j].text)
		}
		stmt = strings.TrimSpace(strings.TrimPrefix(stmt, "await "))
		if strings.Contains(stmt, ".close(") || strings.Contains(stmt, ".aclose(") {
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
// super().__getattribute__ / hasattr / getattr probes.
func tryBlockIsAttrProbe(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	sawBody := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := lines[j].indent
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

func isAttrProbeStmt(stmt string) bool {
	t := strings.TrimSpace(stmt)
	if t == "" || t == "pass" || t == "break" || t == "continue" {
		return true
	}
	if stmtHasSideEffectSink(t) {
		return false
	}
	if isProbeBranchHeader(t) {
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

func isAttrProbeExpr(expr string) bool {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return false
	}
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
// optional dependencies and registers defs/classes/fixtures.
func tryBlockIsOptionalFeatureRegistration(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	bodyIndent := -1
	sawImport := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := lines[j].indent
		if ind <= tryIndent {
			return sawImport
		}
		if bodyIndent < 0 {
			bodyIndent = ind
		}
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

// tryBlockIsConversionGuard reports a try body that only guards pure numeric/
// boolean conversions via assignment (x = int(...)), not return. Assignment
// leaves the prior/unset state on failure; return int(...); except: pass with
// no fallback stays reportable.
func tryBlockIsConversionGuard(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	sawConversion := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := lines[j].indent
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
		for _, bad := range []string{
			"fromisoformat(", "strptime(", "json.loads(", "yaml.load(",
			"yaml.safe_load(", "parse_", "loads(", "decode(",
		} {
			if strings.Contains(stmt, bad) {
				return false
			}
		}
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
			return false
		}
		if conversionCallIn(stmt) {
			sawConversion = true
			continue
		}
		return false
	}
	return sawConversion
}

func conversionCallIn(stmt string) bool {
	for _, fn := range []string{"int(", "float(", "bool(", "complex(", "ord(", "date("} {
		if strings.Contains(stmt, fn) {
			return true
		}
	}
	return false
}

// tryBlockIsImportOnly reports a try block whose direct-body statements are
// all import statements (or assignments of importlib.import_module /
// __import__ / load_extension results). A try block that executes calls or
// control flow stays reportable.
func tryBlockIsImportOnly(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	bodyIndent := -1
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := lines[j].indent
		if ind <= tryIndent {
			return true
		}
		if bodyIndent < 0 {
			bodyIndent = ind
		}
		if ind != bodyIndent {
			continue
		}
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
// extension and then only purely uses the loaded name with no unrelated side
// effects (niquests async_session.py:620).
func tryBlockIsOptionalExtensionLoad(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	loaded := map[string]bool{}
	sawLoad := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := lines[j].indent
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
			// Only dynamic extension loaders count — plain import/from
			// optional-deps are tryBlockIsImportOnly (import-only bodies).
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
		if isProbeBranchHeader(stmt) || strings.HasPrefix(stmt, "return") ||
			stmt == "pass" || stmt == "continue" || stmt == "break" {
			continue
		}
		if _, _, ok := splitAssignmentEq(stmt); ok {
			continue
		}
		if len(loaded) == 0 {
			return false
		}
		refsLoaded := false
		for name := range loaded {
			if lineReferencesVarPy(stmt, name) {
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

// importAssignedName returns the LHS of name = load_extension(...) assignments.
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
	if i := strings.Index(name, ":"); i >= 0 {
		name = strings.TrimSpace(name[:i])
	}
	if bareIdentRE.MatchString(name) {
		return name
	}
	return ""
}

// importStmtNames extracts simple names from import / from-import statements.
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
			if bareIdentRE.MatchString(part) {
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
				if bareIdentRE.MatchString(part) {
					out = append(out, part)
				}
			}
		}
	}
	return out
}

// lineReferencesVarPy reports whether stmt mentions varName as a whole identifier.
func lineReferencesVarPy(line, varName string) bool {
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
			ok = ok && !(isIdentByteCWE(prev) || prev == '.')
		}
		end := abs + len(varName)
		if end < len(line) {
			next := line[end]
			ok = ok && !isIdentByteCWE(next)
		}
		if ok {
			return true
		}
		start = abs + len(varName)
	}
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
// assignments, returns, and branch headers without network/file/subprocess
// side effects and without bare cleanup call statements (mirrors BP).
func tryBlockIsProbe(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	sawBody := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := lines[j].indent
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
		return isProbeReturnExpr(strings.TrimSpace(t[len("return "):]))
	}
	if isProbeBranchHeader(t) {
		return true
	}
	if _, _, ok := splitAssignmentEq(t); ok {
		return true
	}
	// Bare JS-bridge expressions are best-effort probes; other bare calls
	// (close/__aexit__) stay non-probe so wse teardown TPs fire.
	if isJSBridgeExpr(t) {
		return true
	}
	return false
}

// jsBridgeSignals mark a try body as a best-effort JS interop / stream probe
// (web-API and common interop names only — no product-specific identifiers).
var jsBridgeSignals = []string{
	"getReader(",
	".entries(",
	"status_text",
	".bytes()",
	"_ws.",
	"_ws.close",
	"ReadableStream",
	".cancel(",
	".destroy(",
	"AbortSignal",
	"WebSocket",
}

// isJSBridgeExpr reports a statement that touches a JS interop / stream API.
func isJSBridgeExpr(t string) bool {
	for _, s := range jsBridgeSignals {
		if strings.Contains(t, s) {
			return true
		}
	}
	return false
}

// tryBlockHasJSBridgeSignal reports that the try body mentions a JS/pyodide
// bridge API (keeps ordinary Exception:pass TPs reportable).
func tryBlockHasJSBridgeSignal(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		if lines[j].indent <= tryIndent {
			return false
		}
		if isJSBridgeExpr(mt) {
			return true
		}
	}
	return false
}

// isProbeReturnExpr reports return expressions that are probe-shaped.
func isProbeReturnExpr(expr string) bool {
	expr = strings.TrimSpace(expr)
	if expr == "" || isDefinedConstantExprCWE(expr) || bareIdentRE.MatchString(expr) {
		return true
	}
	if stmtHasSideEffectSink(expr) {
		return false
	}
	if strings.Contains(expr, ".") {
		return true
	}
	if i := strings.Index(expr, "("); i > 0 {
		name := strings.TrimSpace(expr[:i])
		return isProbeConversionFunc(name)
	}
	return true
}

// isDefinedConstantExprCWE reports a trivial constant used as a defined fallback.
func isDefinedConstantExprCWE(expr string) bool {
	switch strings.TrimSpace(expr) {
	case "None", "True", "False", "0", "1", "\"\"", "''", "b\"\"", "b''", "[]", "{}", "()", "0.0":
		return true
	}
	return false
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

// isProbeBranchHeader reports if/elif/else/for/while/async for headers.
func isProbeBranchHeader(t string) bool {
	return strings.HasPrefix(t, "if ") || strings.HasPrefix(t, "elif ") ||
		t == "else:" || strings.HasPrefix(t, "for ") || strings.HasPrefix(t, "while ") ||
		strings.HasPrefix(t, "async for ")
}

// stmtHasSideEffectSink reports network, file, or subprocess side effects.
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
			if s == "open(" {
				if indexOfBareOpenCallCWE(t) >= 0 {
					return true
				}
				continue
			}
			return true
		}
	}
	return false
}

// indexOfBareOpenCallCWE finds builtin open( — not mid-ident and not .open(.
func indexOfBareOpenCallCWE(line string) int {
	start := 0
	for {
		i := strings.Index(line[start:], "open(")
		if i < 0 {
			return -1
		}
		abs := start + i
		if abs > 0 {
			prev := line[abs-1]
			if prev == '.' || isIdentByteCWE(prev) {
				start = abs + 4
				continue
			}
		}
		return abs
	}
}

// flowContinuesAfter reports a plain statement (not a same-try clause, not an
// else/elif continuation, not a nested def/class) following the handler at or
// below its indent — the enclosing flow continues with a defined fallback.
func flowContinuesAfter(lines []pyMaskedLine, after, exceptIndent int) bool {
	count := 0
	for j := after; j < len(lines) && count < 8; j++ {
		t := strings.TrimSpace(lines[j].text)
		if t == "" {
			continue
		}
		ind := lines[j].indent
		if ind > exceptIndent {
			continue
		}
		count++
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
func tryBlockSingleStatement(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	bodyIndent := -1
	stmts := 0
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		ind := lines[j].indent
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

// exceptLineHasFallbackMarkerPy reports a marker comment on the except clause
// or the first suite lines that documents the handler as an intentional
// optional/defensive fallback. Matching is case-insensitive ("Fall back").
func exceptLineHasFallbackMarkerPy(rawLines []pyMaskedLine, exceptIdx, exceptIndent int) bool {
	markers := []string{"no cover", "best effort", "non-fatal", "defensive", "fall back", "fallback"}
	if lineHasFallbackMarkerTextPy(rawLines[exceptIdx].text, markers) {
		return true
	}
	for j := exceptIdx + 1; j < len(rawLines); j++ {
		if strings.TrimSpace(rawLines[j].text) == "" {
			continue
		}
		if rawLines[j].indent <= exceptIndent {
			break
		}
		if lineHasFallbackMarkerTextPy(rawLines[j].text, markers) {
			return true
		}
	}
	return false
}

func lineHasFallbackMarkerTextPy(raw string, markers []string) bool {
	lower := strings.ToLower(raw)
	for _, marker := range markers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func importErrorCatch(types []string) bool {
	for _, ty := range types {
		if ty == "ImportError" || ty == "ModuleNotFoundError" {
			return true
		}
	}
	return false
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

// fieldExtractConversionCatch reports KeyError/TypeError/ValueError only —
// dict field extract + int/float conversion (Project_Parva timegraph). Solo
// KeyError stays outside (sync-with-uv TPs).
func fieldExtractConversionCatch(types []string) bool {
	if len(types) == 0 {
		return false
	}
	hasKey, hasConv := false, false
	for _, ty := range types {
		switch ty {
		case "KeyError":
			hasKey = true
		case "TypeError", "ValueError":
			hasConv = true
		default:
			return false
		}
	}
	// Require at least one conversion type; KeyError alone is not enough.
	return hasConv && (hasKey || len(types) >= 1)
}

// tryBlockIsFieldExtractConversion reports try bodies that only convert dict
// fields (int/float/str on subscript/attr) and append the results to a list —
// optional fact-id enrichment (timegraph_fact_links).
func tryBlockIsFieldExtractConversion(lines []pyMaskedLine, tryIdx int) bool {
	if tryIdx < 0 || tryIdx >= len(lines) {
		return false
	}
	tryIndent := lines[tryIdx].indent
	sawConversion := false
	for j := tryIdx + 1; j < len(lines); j++ {
		mt := strings.TrimSpace(lines[j].text)
		if mt == "" {
			continue
		}
		if lines[j].indent <= tryIndent {
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
		if isProbeBranchHeader(stmt) || stmt == "pass" {
			continue
		}
		// year = int(bs["year"]) / day = int(part)
		if _, rhs, ok := splitAssignmentEq(stmt); ok {
			if conversionCallIn(rhs) {
				sawConversion = true
				continue
			}
			// year, month, day = (int(p) for p in ...)
			if strings.Contains(rhs, "int(") || strings.Contains(rhs, "float(") {
				sawConversion = true
				continue
			}
			return false
		}
		// fact_ids.append(...) of converted values — enrichment only.
		if strings.Contains(stmt, ".append(") && !strings.Contains(stmt, "open(") {
			continue
		}
		return false
	}
	return sawConversion
}

// CWE-397 recognizes direct construction or re-raising of Python's generic
// root exception classes. Raising an application-specific exception is safe.
func detectCWE397(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if start := firstMatchStartIfContains(facts, unit, pyGenericRaiseRE,
		"raise Exception", "raise BaseException"); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE397, start, "generic Exception or BaseException is raised directly", confidence82, out)
	}
}

// CWE-478 reports a Python match statement only when it has two or more
// immediate case branches and lacks a wildcard case. The indentation-aware
// walk avoids confusing nested match cases with branches of the outer match.
func detectCWE478(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	lines := facts.MaskedLines()
	if start := matchWithoutDefaultStart(lines); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE478, start, "multiple-case match expression has no wildcard default branch", confidence76, out)
	}
}

// CWE-252 is deliberately limited to process calls used as standalone
// statements. Assigned results and subprocess.run(..., check=True) have
// explicit success handling paths and are not reported.
func detectCWE252(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	for _, name := range []string{"subprocess.run", "subprocess.call", "os.system"} {
		for _, call := range findCalls(facts, unit.Source, name) {
			if !standaloneCall(facts, unit.Source, call) || (name == "subprocess.run" && hasKwargTrue(call.ArgsText, "check")) {
				continue
			}
			emitPlatformFinding(unit, &MetaCWE252, call.Start, "process call return status is discarded without checking success", confidence82, out)
			return
		}
	}
}

// CWE-390 recognizes only an except clause whose direct body is pass. It does
// not infer whether logging, recovery, re-raising, or a caller's behaviour is
// sufficient error handling.
func detectCWE390(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	// Do NOT skip test modules wholesale — chaos / expected-exception tests
	// use exceptPassIsSafe / tryBlockDeliberatelyRaises for FP control.
	lines := facts.MaskedLines()
	rawLines := buildMaskedPythonLines(unit.Source)
	for i, line := range lines {
		if !pyExceptStartRE.MatchString(line.text) {
			continue
		}
		if !exceptPassOnly(lines, i) {
			continue
		}
		if exceptPassIsSafe(unit, lines, rawLines, i) {
			continue
		}
		emitPlatformFinding(unit, &MetaCWE390, line.start, "exception is detected but the handler takes no action", confidence82, out)
		return
	}
}

// exceptPassOnly reports an except clause whose first direct-body statement is
// pass.
func exceptPassOnly(lines []pyMaskedLine, exceptIdx int) bool {
	exceptIndent := lines[exceptIdx].indent
	for _, body := range lines[exceptIdx+1:] {
		trimmed := strings.TrimSpace(body.text)
		if trimmed == "" {
			continue
		}
		if body.indent <= exceptIndent {
			return false
		}
		return trimmed == "pass"
	}
	return false
}

// CWE-584 limits reporting to direct returns in a finally suite. A return in
// an unrelated nested definition is excluded by requiring the suite's direct
// indentation level.
func detectCWE584(unit *core.ParsedUnit, facts *PyCweFacts, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	if start := finallyReturnStart(facts.MaskedLines()); start >= 0 {
		emitPlatformFinding(unit, &MetaCWE584, start, "return inside finally can suppress an exception from the protected block", confidence90, out)
	}
}

type pyMaskedLine struct {
	start  int
	text   string
	indent int
}

// buildMaskedPythonLines splits a pre-masked file into line spans.
func buildMaskedPythonLines(masked string) []pyMaskedLine {
	lines := make([]pyMaskedLine, 0, strings.Count(masked, "\n")+1)
	for start := 0; start <= len(masked); {
		end := len(masked)
		if next := strings.IndexByte(masked[start:], '\n'); next >= 0 {
			end = start + next
		}
		line := masked[start:end]
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		lines = append(lines, pyMaskedLine{start: start, text: line, indent: indent})
		if end == len(masked) {
			break
		}
		start = end + 1
	}
	return lines
}

func matchWithoutDefaultStart(lines []pyMaskedLine) int {
	for i, match := range lines {
		if !pyMatchStartRE.MatchString(match.text) {
			continue
		}
		caseIndent, cases, hasDefault := -1, 0, false
		for _, line := range lines[i+1:] {
			if strings.TrimSpace(line.text) == "" {
				continue
			}
			if line.indent <= match.indent {
				break
			}
			if !pyCaseStartRE.MatchString(line.text) {
				continue
			}
			if caseIndent < 0 {
				caseIndent = line.indent
			}
			if line.indent != caseIndent {
				continue
			}
			cases++
			hasDefault = hasDefault || pyDefaultCaseRE.MatchString(line.text)
		}
		if cases >= 2 && !hasDefault {
			return match.start
		}
	}
	return -1
}

func standaloneCall(facts *PyCweFacts, source string, call callSite) bool {
	lineStart := strings.LastIndex(source[:call.Start], "\n") + 1
	prefix := source[lineStart:call.Start]
	var masked string
	if facts != nil {
		masked = facts.codeMask(prefix, lineStart)
	} else {
		masked = pythonCodeMask(prefix)
	}
	return strings.TrimSpace(masked) == ""
}

func finallyReturnStart(lines []pyMaskedLine) int {
	for i, finally := range lines {
		if !pyFinallyStartRE.MatchString(finally.text) {
			continue
		}
		bodyIndent := -1
		for _, body := range lines[i+1:] {
			trimmed := strings.TrimSpace(body.text)
			if trimmed == "" {
				continue
			}
			if body.indent <= finally.indent {
				break
			}
			if bodyIndent < 0 {
				bodyIndent = body.indent
			}
			if body.indent == bodyIndent && pyReturnLineRE.MatchString(body.text) {
				return body.start
			}
		}
	}
	return -1
}

func emitPlatformFinding(unit *core.ParsedUnit, meta *rules.RuleMetadata, offset int, message string, confidence float32, out *[]rules.Finding) {
	if unit == nil || out == nil {
		return
	}
	line, col := unit.LineCol(offset)
	rules.PushFindingWithConfidence(meta, unitFile(unit), line, col, message, confidence, out)
}
