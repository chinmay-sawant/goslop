# False-positive audit — calgebra

## Run metadata

```yaml
timestamp: 2026-08-02T07:27:48Z
repository: calgebra
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra
branch: main
commit: 476c3e6
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `bin/goslop` prebuilt (`make build`)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir real-repos/calgebra/scripts/chunks -context-dir real-repos/calgebra/scripts/findings/functions real-repos/calgebra`
- Findings: `42`
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt`, `Chunk_26_42.txt`
- Function contexts reviewed: `./scripts/findings/functions/1.txt` … `42.txt` (all 42 exist; the enclosing source was re-read directly for every finding)

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 23 | 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 22, 23, 24, 30, 32, 33, 34, 35, 37, 38, 39, 42 |
| True positive | 19 | 1, 2, 3, 4, 5, 6, 18, 19, 20, 21, 25, 26, 27, 28, 29, 31, 36, 40, 41 |
| Uncertain | 0 | — |

## False positives

### [ ] Finding `7` — `BP-PY-13`

- Function context: `./scripts/findings/functions/7.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcal.py:14:5`
- Checklist pattern: secret-like assignment token inside a module docstring example, not executable code

Source excerpt:

```
10: Example:
11:     >>> from calgebra.gcal import calendars, Event, Calendar
12:     >>> from calgebra import at_tz
13:     >>>
14:     >>> access_token = "ya29...."  # From Google OAuth
15:     >>> cals = calendars(access_token)
```

Why this is a false positive: the flagged line is doctest prose inside the module docstring (lines 1–27), and the literal `"ya29...."` is a truncated placeholder token; no executable secret assignment exists in the source. BP-PY-13's condition ("a secret-like name is assigned a non-empty string literal in source") cannot be satisfied by documentation text.

Checklist evidence: the flagged line is docstring example text (`>>>` prompt), not an assignment executed in source code.

### [ ] Finding `8` — `BP-PY-7`

- Function context: `./scripts/findings/functions/8.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcal.py:210:12`
- Checklist pattern: `open(` token is an XMLHttpRequest method call, not the file-opening builtin

Source excerpt:

```
209:         xhr = XMLHttpRequest.new()
210:         xhr.open(method, url, False)  # synchronous
211:         xhr.setRequestHeader("Authorization", f"Bearer {access_token}")
```

Why this is a false positive: `xhr.open(...)` is a method call on an `XMLHttpRequest` object (HTTP request setup), not the builtin file-opening `open()`; no file handle is created, so there is no resource that `with` could manage.

Checklist evidence: BP-PY-7's condition (a file opened without a `with` statement, risking resource leaks) is unmet — the flagged call opens an HTTP connection, not a file.

### [ ] Finding `9` — `BP-PY-1`

- Function context: `./scripts/findings/functions/9.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcal.py:221:1`
- Checklist pattern: handler's failure is immediately re-raised, so it is not hidden

Source excerpt:

```
217:         if xhr.status >= 400:
218:             try:
219:                 err = json.loads(xhr.responseText)
220:                 msg = err.get("error", {}).get("message", xhr.responseText)
221:             except Exception:
222:                 msg = f"HTTP {xhr.status}"
223:             raise RuntimeError(f"Google Calendar API error ({xhr.status}): {msg}")
```

Why this is a false positive: the handler only formats an error message; the next statement re-raises `RuntimeError` carrying that message, so the failure propagates to the caller. BP-PY-1's condition is "swallows failures and hides bugs" — nothing is swallowed here.

Checklist evidence: the shown source's immediate flow re-raises the failure (line 223), so the "broad except hides failures" predicate is not satisfied.

### [ ] Finding `10` — `CWE-396`

- Function context: `./scripts/findings/functions/10.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcal.py:221:1`
- Checklist pattern: generic handler cannot hide failures because the flow re-raises immediately after

Source excerpt:

```
218:             try:
219:                 err = json.loads(xhr.responseText)
220:                 msg = err.get("error", {}).get("message", xhr.responseText)
221:             except Exception:
222:                 msg = f"HTTP {xhr.status}"
223:             raise RuntimeError(f"Google Calendar API error ({xhr.status}): {msg}")
```

Why this is a false positive: CWE-396's weakness is that a generic handler "can hide failures that require distinct handling"; here the handler feeds a message into an immediate `RuntimeError` re-raise, so no failure condition is hidden. The sibling rule BP-PY-1 explicitly skips broad excepts whose suite re-raises (`suiteReraises`), confirming this shape is not a weakness.

Checklist evidence: the generic catch at line 221 is followed by `raise RuntimeError(...)` at line 223 using the handler's message — the failure is propagated, not hidden.

### [ ] Finding `11` — `BP-PY-1`

- Function context: `./scripts/findings/functions/11.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcal.py:245:1`
- Checklist pattern: handler's failure is re-raised (`from e`) immediately after the handler

Source excerpt:

```
241:         except urllib.error.HTTPError as e:
242:             try:
243:                 err = json.loads(e.read().decode())
244:                 msg = err.get("error", {}).get("message", str(err))
245:             except Exception:
246:                 msg = f"HTTP {e.code}"
247:             raise RuntimeError(f"Google Calendar API error ({e.code}): {msg}") from e
```

Why this is a false positive: same shape as finding 9 — the handler only builds a message and the enclosing flow re-raises `RuntimeError(...) from e`, propagating the failure with its detail.

Checklist evidence: the shown source re-raises the failure at line 247 with `from e`, so the "hides failures" condition of BP-PY-1 is not satisfied.

### [ ] Finding `12` — `BP-PY-1`

- Function context: `./scripts/findings/functions/12.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcal.py:533:1`
- Checklist pattern: exception object is returned to the caller as an error result, not discarded

Source excerpt:

```
526: def _handle_write_errors(func: _F) -> _F:
527:     """Decorator to wrap write operations with error handling."""
528:
529:     @wraps(func)
530:     def wrapper(*args: Any, **kwargs: Any) -> list[WriteResult]:
531:         try:
532:             return func(*args, **kwargs)
533:         except Exception as e:
534:             return _error_result(e)
```

Why this is a false positive: the handler does not swallow the failure — it packages the exception object into a `WriteResult(success=False, error=e)` (`_error_result(e)`) returned to the caller, which is the documented error channel of this API. The failure is surfaced to the caller, not hidden.

Checklist evidence: the exception is preserved and returned as the error payload of a `WriteResult`, so the "swallows failures" predicate of the rule is unmet.

### [ ] Finding `13` — `BP-PY-1`

- Function context: `./scripts/findings/functions/13.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcal.py:616:1`
- Checklist pattern: handler reports the failure to the logging system with `exc_info=True`

Source excerpt:

```
608:             try:
609:                 url = f"{_API_BASE}/calendars/{self.id}"
610:                 data = _xhr_request("GET", url, self._access_token)
611:                 if data:
612:                     tz_str = data.get("timeZone")
613:                     if tz_str:
614:                         self.timezone = tz_str
615:                         self.__calendar_timezone = ZoneInfo(tz_str)
616:             except Exception:
617:                 logging.getLogger("calgebra.gcal").warning(
618:                     "Failed to fetch timezone for calendar %s, falling back to UTC",
619:                     self.id,
620:                     exc_info=True,
621:                 )
```

Why this is a false positive: the handler reports the failure through the logging module with `exc_info=True` (full traceback) and the UTC fallback is explicit in the message; the failure is surfaced in the logs, not hidden.

Checklist evidence: the suite calls `logging.getLogger(...).warning(..., exc_info=True)` — the failure is reported, so the "broad except hides failures" condition is not satisfied.

### [ ] Finding `14` — `BP-PY-1`

- Function context: `./scripts/findings/functions/14.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcal.py:838:1`
- Checklist pattern: failure detail is returned to the caller as an error result

Source excerpt:

```
835:         try:
836:             url = f"{_API_BASE}/calendars/{self.id}/events/{master_event_id}"
837:             master_data = _xhr_request("GET", url, self._access_token)
838:         except Exception as e:
839:             return _error_result(
840:                 ValueError(f"Failed to fetch master event {master_event_id}: {e}")
841:             )
```

Why this is a false positive: the exception detail is wrapped in a `ValueError` and returned to the caller as an error `WriteResult` — the failure is communicated with its context, not discarded.

Checklist evidence: the handler returns `_error_result(ValueError(...))` carrying `{e}`, surfacing the failure to the caller.

### [ ] Finding `15` — `BP-PY-1`

- Function context: `./scripts/findings/functions/15.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcal.py:870:1`
- Checklist pattern: failure detail is returned to the caller as an error result

Source excerpt:

```
867:             try:
868:                 url = f"{_API_BASE}/calendars/{self.id}/events/{master_event_id}"
869:                 _xhr_request("PUT", url, self._access_token, master_data)
870:             except Exception as e:
871:                 return _error_result(ValueError(f"Failed to update master event: {e}"))
```

Why this is a false positive: identical shape to finding 14 — the exception is converted into an error result containing `{e}` and returned to the caller, so the failure is surfaced, not hidden.

Checklist evidence: the handler returns `_error_result(ValueError(...))` with the exception detail embedded.

### [ ] Finding `16` — `BP-PY-1`

- Function context: `./scripts/findings/functions/16.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcsa.py:440:1`
- Checklist pattern: exception object is returned to the caller as an error result, not discarded

Source excerpt:

```
436:     @wraps(func)
437:     def wrapper(*args: Any, **kwargs: Any) -> list[WriteResult]:
438:         try:
439:             return func(*args, **kwargs)
440:         except Exception as e:
441:             return _error_result(e)
```

Why this is a false positive: same as finding 12 — the wrapper converts the exception into a `WriteResult` whose `error` field carries the exception, and the docstring above documents this as deliberate API design ("Catches all exceptions and converts them to WriteResult lists"). The failure is returned to the caller.

Checklist evidence: the exception object is preserved and returned as the error payload of a `WriteResult`, so the "swallows failures" predicate is unmet.

### [ ] Finding `17` — `CWE-396`

- Function context: `./scripts/findings/functions/17.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcsa.py:440:1`
- Checklist pattern: generic catch's failure is returned to the caller as an error result, so it cannot hide failure conditions

Source excerpt:

```
437:     def wrapper(*args: Any, **kwargs: Any) -> list[WriteResult]:
438:         try:
439:             return func(*args, **kwargs)
440:         except Exception as e:
441:             return _error_result(e)
```

Why this is a false positive: CWE-396's weakness (generic handler hiding distinct failure conditions) is defeated because the exception object itself is returned to the caller in the `WriteResult.error` field; no failure is hidden.

Checklist evidence: the handler returns the exception `e` to the caller, so the "hide failures" condition of CWE-396 is not satisfied.

### [ ] Finding `22` — `BP-PY-1`

- Function context: `./scripts/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcsa.py:1010:1`
- Checklist pattern: failure detail is returned to the caller as an error result

Source excerpt:

```
1005:         try:
1006:             # Fetch the master event
1007:             master_event = self.calendar.get_event(
1008:                 master_event_id, calendar_id=self.calendar_id
1009:             )
1010:         except Exception as e:
1011:             return _error_result(
1012:                 ValueError(f"Failed to fetch master event {master_event_id}: {e}")
1013:             )
```

Why this is a false positive: identical shape to finding 14 — the exception detail is wrapped in a `ValueError` and returned to the caller as an error result; the failure is communicated, not swallowed.

Checklist evidence: the handler returns `_error_result(ValueError(...))` embedding `{e}`.

### [ ] Finding `23` — `BP-PY-1`

- Function context: `./scripts/findings/functions/23.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcsa.py:1040:1`
- Checklist pattern: failure detail is returned to the caller as an error result

Source excerpt:

```
1038:             try:
1039:                 self.calendar.update_event(master_event, calendar_id=self.calendar_id)
1040:             except Exception as e:
1041:                 return _error_result(ValueError(f"Failed to update master event: {e}"))
```

Why this is a false positive: identical shape to finding 15 — the exception is converted into an error result containing `{e}` and returned to the caller.

Checklist evidence: the handler returns `_error_result(ValueError(...))` embedding `{e}`.

### [ ] Finding `24` — `BP-PY-1`

- Function context: `./scripts/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/gcsa.py:1117:1`
- Checklist pattern: failure detail is returned to the caller as per-event error results

Source excerpt:

```
1115:         try:
1116:             batch.execute()
1117:         except Exception as e:
1118:             # If batch execution fails entirely, return error for all
1119:             return [
1120:                 WriteResult(success=False, event=None, error=e) for _ in events_list
1121:             ]
```

Why this is a false positive: the exception is turned into one `WriteResult(success=False, error=e)` per event and returned to the caller (comment line 1118 documents the intent); the failure is surfaced, not hidden.

Checklist evidence: the handler returns the exception to the caller in every result item, so the "swallows failures" predicate is unmet.

### [ ] Finding `30` — `BP-PY-1`

- Function context: `./scripts/findings/functions/30.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/ical.py:384:1`
- Checklist pattern: handler reports the failure to stderr with the exception detail and deliberately skips the item

Source excerpt:

```
380:     for component in cal.walk("VEVENT"):
381:         try:
382:             item = _parse_vevent(component, calendar_name=calendar_name)
383:             timeline.add(item)
384:         except Exception as e:
385:             # We might want to log this or optionally fail
386:             print(f"Warning: Failed to parse VEVENT: {e}", file=sys.stderr)
387:             continue
```

Why this is a false positive: the handler reports the failure to stderr including the exception detail (`{e}`) and skips the malformed component by design (comment line 385); the failure is surfaced to the user, not hidden.

Checklist evidence: the suite prints the warning with the exception text to stderr — the "broad except hides failures" condition is not satisfied.

### [ ] Finding `32` — `BP-PY-46`

- Function context: `./scripts/findings/functions/32.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/interval.py:134:9`
- Checklist pattern: `print` produces the function's documented output, not operational logging

Source excerpt:

```
121: def pprint(
122:     intervals: Iterable["Interval"], tz: str = "UTC", fmt: str = "%Y-%m-%d %H:%M:%S"
123: ) -> None:
124:     """Pretty-print an iterable of Intervals.
125:
126:     Consumes the iterable and prints formatted datetime strings to stdout.
127:
128:     Args:
129:         intervals: Iterable of Interval objects
130:         tz: Target timezone name (default "UTC")
131:         fmt: strftime format string (default "%Y-%m-%d %H:%M:%S")
132:     """
133:     for ivl in intervals:
134:         print(ivl.format(tz=tz, fmt=fmt))
```

Why this is a false positive: the print is the function's intended, documented output ("prints formatted datetime strings to stdout") — a user-facing presentation helper, not operational logging. BP-PY-46's condition is "print used for operational logging in non-script modules"; the "operational logging" clause is unmet.

Checklist evidence: the flagged `print` implements the documented purpose of `pprint()` (output, not logging), so the rule's "operational logging" predicate is not satisfied.

### [ ] Finding `33` — `BP-PY-2`

- Function context: `./scripts/findings/functions/33.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/recurrence.py:437:1`
- Checklist pattern: `pass` is a control-flow fall-through to an explicit raise of the same error

Source excerpt:

```
427:                 if len(s) > 2:
428:                     code = s[-2:]
429:                     prefix = s[:-2]
430:                     # Check code in map (must use lower for map key)
431:                     if code.lower() in _DAY_MAP:
432:                         wd_const = _DAY_MAP[code.lower()]
433:                         try:
434:                             n = int(prefix)
435:                             weekdays.append(wd_const(n))
436:                             continue
437:                         except ValueError:
438:                             pass
439:
440:                 # If we get here, invalid
441:                 valid = ", ".join(sorted(_DAY_MAP.keys()))
442:                 raise ValueError(
443:                     f"Invalid day name: '{d}'\n"
444:                     f"Valid days: {valid} or numbered (e.g. 1MO)\n"
445:                 )
```

Why this is a false positive: the `pass` is a deliberate fall-through — when `int(prefix)`/`wd_const(n)` raises `ValueError`, execution proceeds to the explicit `raise ValueError("Invalid day name: ...")` at line 442. The failure is surfaced with a detailed message, not "discarded silently", which is the condition BP-PY-2 describes.

Checklist evidence: the pass immediately precedes an unconditional `raise ValueError` that reports the invalid day, so the "failures are discarded silently" predicate is not satisfied.

### [ ] Finding `34` — `CWE-390`

- Function context: `./scripts/findings/functions/34.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/recurrence.py:437:1`
- Checklist pattern: the pass path falls through to an explicit raise, so the error condition is acted on

Source excerpt:

```
433:                         try:
434:                             n = int(prefix)
435:                             weekdays.append(wd_const(n))
436:                             continue
437:                         except ValueError:
438:                             pass
439:
440:                 # If we get here, invalid
441:                 valid = ", ".join(sorted(_DAY_MAP.keys()))
442:                 raise ValueError(
443:                     f"Invalid day name: '{d}'\n"
444:                     f"Valid days: {valid} or numbered (e.g. 1MO)\n"
445:                 )
```

Why this is a false positive: CWE-390's condition is "exception is detected but the handler takes no action"; here the handler's `pass` routes control to the explicit `raise ValueError(...)` at line 442, which is the action taken on the error condition.

Checklist evidence: the pass is a fall-through to a raise that reports the invalid input, so the "takes no action" predicate of CWE-390 is not satisfied.

### [ ] Finding `35` — `CWE-1071`

- Function context: `./scripts/findings/functions/35.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/calgebra/recurrence.py:437:25`
- Checklist pattern: the pass is a fall-through to an explicit raise; the handler is not an empty dead block

Source excerpt:

```
433:                         try:
434:                             n = int(prefix)
435:                             weekdays.append(wd_const(n))
436:                             continue
437:                         except ValueError:
438:                             pass
439:
440:                 # If we get here, invalid
441:                 valid = ", ".join(sorted(_DAY_MAP.keys()))
442:                 raise ValueError(
443:                     f"Invalid day name: '{d}'\n"
444:                     f"Valid days: {valid} or numbered (e.g. 1MO)\n"
445:                 )
```

Why this is a false positive: the handler is not an empty dead block — its `pass` is the designated path into the "invalid day name" raise at line 442; the error condition is handled by the fall-through, so the "silently contains only pass" (empty block) condition is not genuinely met.

Checklist evidence: the pass is the control-flow route to the explicit raise, so the handler is not a silent empty code block.

### [ ] Finding `37` — `BP-PY-2`

- Function context: `./scripts/findings/functions/37.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/tests/test_gcsa.py:1062:1`
- Checklist pattern: the pass is the test's expected-exception assertion, not a discarded failure

Source excerpt:

```
1056:     # Verify master event was deleted (not the instance)
1057:     assert len(stub._events) == 0
1058:     # Verify master was deleted, not instance
1059:     try:
1060:         stub.get_event("master-event-id")
1061:         assert False, "Master event should have been deleted"
1062:     except ValueError:
1063:         pass  # Expected
```

Why this is a false positive: the try/except is the test's expected-exception assertion — `assert False` fails the test if the event is still present, and the `pass` catches the `ValueError` that proves the deletion. The failure is verified, not "discarded silently", so BP-PY-2's condition is not met.

Checklist evidence: the pass is the verification branch of a pytest test (`assert False` precedes it), so the "failures are discarded silently" predicate is unmet.

### [ ] Finding `38` — `CWE-390`

- Function context: `./scripts/findings/functions/38.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/tests/test_gcsa.py:1062:1`
- Checklist pattern: the handler's "no action" is the test's expected-exception assertion

Source excerpt:

```
1059:     try:
1060:         stub.get_event("master-event-id")
1061:         assert False, "Master event should have been deleted"
1062:     except ValueError:
1063:         pass  # Expected
```

Why this is a false positive: the handler's `pass` is the test's verification of the expected `ValueError` (the test fails at `assert False` if no exception occurs); the error condition is the tested outcome, so "the handler takes no action" is not satisfied.

Checklist evidence: the pass is the expected-exception assertion of the test, so the "takes no action" predicate of CWE-390 is unmet.

### [ ] Finding `39` — `CWE-1071`

- Function context: `./scripts/findings/functions/39.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/tests/test_gcsa.py:1062:5`
- Checklist pattern: the pass is the test's expected-exception assertion; the handler is not an empty dead block

Source excerpt:

```
1059:     try:
1060:         stub.get_event("master-event-id")
1061:         assert False, "Master event should have been deleted"
1062:     except ValueError:
1063:         pass  # Expected
```

Why this is a false positive: the handler is the expected-exception branch of the test — it is not a silently empty block; the `pass` with the `# Expected` comment is the verification that the delete succeeded.

Checklist evidence: the pass is the test's assertion mechanism (guarded by `assert False`), so the "silently contains only pass" condition of CWE-1071 is not met.

### [ ] Finding `42` — `CWE-1121`

- Function context: `./scripts/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/calgebra/tests/test_recurrence_fuzz.py:17:70`
- Checklist pattern: the 12-branch threshold is reached only by counting `if `/`for ` tokens inside comments

Source excerpt:

```
17: def get_ground_truth(freq_str, interval, start_ts, end_ts, duration):
18:     """
19:     Generate ground truth intervals using rrule from Epoch.
20:     """
21:     if freq_str == "daily":
22:         freq = DAILY
23:         dtstart = EPOCH
24:     elif freq_str == "weekly":
25:         freq = WEEKLY
26:         dtstart = EPOCH_MONDAY
27:     elif freq_str == "monthly":
28:         freq = MONTHLY
29:         dtstart = EPOCH
30:     elif freq_str == "yearly":
31:         freq = YEARLY
32:         dtstart = EPOCH
33:     else:
34:         raise ValueError(f"Unknown freq: {freq_str}")
```

Why this is a false positive: the function's real control flow is 10 branch headers — 4×`if`/`elif` dispatch, `if before:`, `for dt in candidates:`, `if overlap_start < overlap_end:`, `if not results:`, `for next_ivl in results[1:]:`, `if current.end >= next_ivl.start:` — below the 12-branch threshold. The detector's `strings.Count(code, "if "/"for "/...)` also matches comment text (`rrule.between if we are careful about overlaps`, `for the basic fuzz test`, `Convert query bounds ... for filtering`, `in case (for long durations)`, `Overlap logic for exclusive ends`, `Merge if overlap or adjacent`), inflating the count to ≥12.

Checklist evidence: the rule condition is "at least twelve visible control-flow branches"; the shown function has 10 real branch headers, so the threshold is only met via comment-text matches, and the condition is not genuinely satisfied.

## True positives

### BP-PY-1 — broad except without re-raise (Findings 5, 18, 25)

| Finding id | Source | Reason |
| --- | --- | --- |
| 5 | `calgebra/dataframe.py:49:1` | `except Exception: return ZoneInfo("UTC")` — a fixed fallback return that silently substitutes the timezone; the error information is discarded and the caller cannot observe the failure (no re-raise, no report). |
| 18 | `calgebra/gcsa.py:677:1` | `except Exception:` whose suite is only `pass` (after comment stripping) — broad except, no re-raise, no report; the timezone-fetch failure is silently discarded. |
| 25 | `calgebra/ical.py:248:1` | `except Exception:` whose suite is only `pass` — broad except that silently discards the `str(start_dt.tzinfo)` failure, leaving `tz` implicitly `None`. |

### BP-PY-2 — Except Pass (Findings 1, 4, 19, 26)

| Finding id | Source | Reason |
| --- | --- | --- |
| 1 | `calgebra/__init__.py:24:1` | `except ImportError: pass` after `from .ical import (...)` — the suite is solely `pass`; the rule has no exemption for optional-import guards, so the condition (`suite[0] == "pass"`) is met. |
| 4 | `calgebra/__init__.py:29:1` | `except ImportError: pass` after `from .dataframe import to_dataframe` — same structural condition on a distinct line. |
| 19 | `calgebra/gcsa.py:677:1` | `except Exception:` suite is solely `pass` after comment stripping. |
| 26 | `calgebra/ical.py:248:1` | `except Exception:` suite is solely `pass`. |

### CWE-390 — Detection of Error Condition Without Action (Findings 3, 20, 27)

| Finding id | Source | Reason |
| --- | --- | --- |
| 3 | `calgebra/__init__.py:24:1` | `except ImportError:` whose only direct body statement is `pass` — matches `exceptPassStart` exactly. |
| 20 | `calgebra/gcsa.py:677:1` | `except Exception:` whose only direct body statement is `pass` — the rule deliberately does not judge whether the documented UTC fallback is adequate handling. |
| 27 | `calgebra/ical.py:248:1` | `except Exception:` whose only direct body statement is `pass`. |

### CWE-1071 — Empty Code Block (Findings 2, 21, 29)

| Finding id | Source | Reason |
| --- | --- | --- |
| 2 | `calgebra/__init__.py:24:9` | `except ImportError:` immediately followed by a `pass` body — exception handler contains only `pass`. |
| 21 | `calgebra/gcsa.py:677:13` | `except Exception:` handler containing only `pass` (comments are stripped). |
| 29 | `calgebra/ical.py:248:9` | `except Exception:` handler containing only `pass`. |

### CWE-396 — Declaration of Catch for Generic Exception (Findings 6, 28)

| Finding id | Source | Reason |
| --- | --- | --- |
| 6 | `calgebra/dataframe.py:49:1` | `except Exception:` matches `pyGenericExceptRE` as the first match in the non-test module; the handler does not re-raise and silently substitutes a fixed UTC fallback. |
| 28 | `calgebra/ical.py:248:1` | `except Exception:` first match in the non-test module; handler only `pass`. |

### BP-PY-46 — print Debugging In Library Code (Finding 31)

| Finding id | Source | Reason |
| --- | --- | --- |
| 31 | `calgebra/ical.py:386:13` | `print(f"Warning: Failed to parse VEVENT: {e}", file=sys.stderr)` — operational-logging print in library module `calgebra/ical.py` (imported by the package, no `__main__` guard, no CLI function), matching the rule's condition. |

### BP-PY-41 — pytest assert With Side Effects Only (Findings 36, 40, 41)

| Finding id | Source | Reason |
| --- | --- | --- |
| 36 | `tests/test_cache.py:509:1` | `test_cover_interval_hashable` only constructs a `CoverInterval` and calls `hash(cover)`; no `assert`/`pytest.raises`/`self.assert` anywhere in the body. |
| 40 | `tests/test_metrics_groupby.py:14:1` | `test_valid_hour_hour_of_day` only calls `total_duration(...)`; no assertion statement exists in the body. |
| 41 | `tests/test_metrics_groupby.py:19:1` | `test_valid_day_day_of_week` only calls `total_duration(...)`; no assertion statement exists in the body. |

## Uncertain findings

None.

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — pass
