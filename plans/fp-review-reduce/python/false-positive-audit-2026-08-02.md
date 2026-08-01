# False-positive audit report

## Run metadata

```yaml
timestamp: 2026-08-01T19:55:00Z
repository: goslop
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop
branch: feat/python-perf-ruleset-plan
commit: 72cb97a4cf37aabae80a666db3dd9d56457480a6
scan_target: /home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/pythoncoreengine
chunk_path: ./scripts/chunks
function_context_path: ./scripts/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop`
- Scan command: `make run-python` (profile all, pythoncoreengine; exports under `./scripts/chunks` and `./scripts/findings/functions`)
- Findings: `102`
- Marker: `./scripts/chunks/.python_fp_audit_marker` (`PYTHON_FP_AUDIT_MARKER … findings=102 target=pythoncoreengine`)
- Chunks reviewed: `./scripts/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_75.txt`, `Chunk_76_100.txt`, `Chunk_101_102.txt`
- Function contexts reviewed: `./scripts/findings/functions/1.txt` … `102.txt`

## Audit checklist

- [x] Read every assigned chunk under `./scripts/chunks`.
- [x] Read `./scripts/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient.
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews ([1–50](2c77607e-55bc-47f5-9b59-98a160234061), [51–102](652b2893-c31c-4595-9fed-003b87ced71b)); spot-checked TPs/uncertain against source; no remaining disagreements.
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 96 | 1–6, 8–38, 41–45, 47–48, 51–102 |
| True positive | 5 | 7, 39, 46, 49, 50 |
| Uncertain | 1 | 40 |

## Classification table (all 102)

| ID | Rule | Classification | One-line rationale |
| --- | --- | --- | --- |
| 1 | PERF-PY-27 | False positive | `analyze_pdf` reads each path once |
| 2–6 | BP-PY-46 | False positive | argparse/`__main__` CLI status output |
| 7 | CWE-88 | True positive | CLI `pdf`/`flavour` in subprocess argv without `--` / allowlist |
| 8–17 | BP-PY-46 | False positive | CLI presentation helpers used by argparse cmds |
| 18 | CWE-328 | False positive | bench PDF md5 fingerprint |
| 19 | CWE-22 | False positive | `_repo_root()` resolves `__file__` only |
| 20–21 | BP-PY-46 | False positive | bench argparse `main()` CLI output |
| 22 | CWE-22 | False positive | same `_repo_root()` / `__file__` pattern |
| 23 | CWE-90 | False positive | `re.search` on PDF bytes, not LDAP |
| 24 | CWE-328 | False positive | bench PDF md5 fingerprint |
| 25–28 | BP-PY-46 | False positive | financial bench argparse CLI output |
| 29 | CWE-22 | False positive | same `_repo_root()` / `__file__` pattern |
| 30 | CWE-328 | False positive | bench md5 fingerprint per tier |
| 31 | CWE-90 | False positive | `re.search` on PDF bytes, not LDAP |
| 32–35 | BP-PY-46 | False positive | zerodha bench argparse CLI output |
| 36 | CWE-1046 | False positive | extends mutable `bytearray`, not immutable text |
| 37 | CWE-328 | False positive | ICC profile-ID MD5 required by format |
| 38 | CWE-328 | False positive | PDF trailer `/ID` document fingerprint |
| 39 | BP-PY-46 | True positive | `ENGINE_DEBUG_BUFFERS` debug print in library render |
| 40 | CWE-328 | Uncertain | PDF encrypt Algorithm 2 MD5 (spec-required key material) |
| 41 | CWE-90 | False positive | `_TRAILER_RE.search`, not LDAP |
| 42 | CWE-256 | False positive | fixture-only `"fixture-password"` |
| 43 | CWE-798 | False positive | same fixture-only password literal |
| 44 | CWE-73 | False positive | intentional `__main__` fixture out-dir |
| 45 | CWE-328 | False positive | font subset-tag fingerprint |
| 46 | PERF-PY-24 | True positive | `wrap_text` in `_row_height` and again in `_draw_row` |
| 47 | CWE-328 | False positive | deterministic XMP UUID seed fingerprint |
| 48 | CWE-22 | False positive | `Path(path)` open; no confined-root join |
| 49 | BP-PY-1 | True positive | bare `except Exception:` swallows decode failures |
| 50 | CWE-396 | True positive | generic `Exception` catch hides distinct failures |
| 51 | CWE-90 | False positive | `re.search` on PDF stream bytes, not LDAP |
| 52 | CWE-328 | False positive | test MD5 equality checksum |
| 53 | CWE-88 | False positive | fixed temp path + literal `-f 4` |
| 54–74 | BP-PY-41 | False positive | `_assert_verapdf_compliant` uses `self.assert*` |
| 75 | BP-PY-41 | False positive | helper raises `AssertionError` on miss |
| 76 | CWE-90 | False positive | `re.search` on PDF encrypt dict, not LDAP |
| 77–78 | CWE-328 | False positive | test MD5 determinism checksum |
| 79 | CWE-90 | False positive | `re.search` on PDF ParentTree bytes, not LDAP |
| 80–81 | CWE-328 | False positive | ICC profile-ID MD5 required by format |
| 82–84 | CWE-90 | False positive | `re.search` / regex `.search` on PDF bytes |
| 85 | CWE-328 | False positive | test MD5 equality checksum |
| 86 | CWE-90 | False positive | `re.search` on PDF `/Pages` count, not LDAP |
| 87 | CWE-88 | False positive | fixed temp `note.pdf` + hardcoded flavour |
| 88–89 | BP-PY-41 | False positive | `_assert_verapdf` uses `self.assert*` |
| 90 | CWE-90 | False positive | `re.search` on PDF ParentTree, not LDAP |
| 91 | CWE-88 | False positive | fixed fixture path + loop flavours `"4"`/`"ua2"` |
| 92 | CWE-90 | False positive | `re.search` on PDF catalog StructTreeRoot |
| 93 | BP-PY-46 | False positive | CLI usage `print` in `main()` / `__main__` |
| 94 | CWE-22 | False positive | intentional CLI `Path(argv[0])` root |
| 95 | PERF-PY-27 | False positive | reads each entry of `unique_paths` once |
| 96–98 | BP-PY-46 | False positive | CLI summary/usage `print` in `main()` / `__main__` |
| 99 | CWE-22 | False positive | intentional CLI `Path(argv[0])` root |
| 100 | PERF-PY-27 | False positive | reads each entry of `unique_paths` once |
| 101–102 | BP-PY-46 | False positive | CLI summary `print` in `main()` / `__main__` |

## False positives

### [ ] Finding 1 — PERF-PY-27

- Function context: `./scripts/findings/functions/1.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/pythoncoreengine/compliance/structure_tree_check.py:83:1`
- Checklist pattern: once-per-distinct-path / no repeated same-path load

Source excerpt:

```
def analyze_pdf(pdf_path: str) -> StructureTreeReport:
    report = StructureTreeReport(pdf=pdf_path, ok=True)
    data = Path(pdf_path).read_bytes()
    objs = parse_objects(data)
```

Why this is a false positive: One `read_bytes` per call; callers iterate distinct paths without re-loading the same path in a hot loop.

Checklist evidence: No repeated load of the same filesystem path without a cache.

---

### [ ] Findings 2–6 — BP-PY-46

- Function contexts: `./scripts/findings/functions/2.txt` … `6.txt`
- Source: `…/compliance/structure_tree_check.py:162–191`
- Checklist pattern: CLI / `__main__` user-facing output

Source excerpt:

```
def main() -> int:
    ...
            print(f"FAIL structure-tree: file not found: {pdf}", file=sys.stderr)
    ...
                print(f"FAIL structure-tree: {rel}")
...
if __name__ == "__main__":
    raise SystemExit(main())
```

Why this is a false positive: Prints are intentional CLI status/report output inside argparse `main()`, not library debug logging.

Checklist evidence: argparse entrypoint + `__name__ == "__main__"` guard.

---

### [ ] Findings 8–17 — BP-PY-46

- Function contexts: `./scripts/findings/functions/8.txt` … `17.txt`
- Source: `…/compliance/verapdf_report.py:224–295`
- Checklist pattern: CLI / intentional user messaging via argparse helpers

Source excerpt:

```
def print_result_lines(...):
    print(f"{status} {profile}: {result.summary}")
    ...
            print(line)

def print_table(...):
    print("")
    print(colors.bold("veraPDF compliance summary"))
    ...
```

Why this is a false positive: Dedicated CLI presentation helpers invoked from argparse `cmd_check` / `cmd_table`, not operational library debugging.

Checklist evidence: intentional CLI user output under argparse commands.

---

### [ ] Findings 18, 24, 30 — CWE-328

- Function contexts: `./scripts/findings/functions/{18,24,30}.txt`
- Source: `…/engine/bench.py:146`, `bench_financial.py:137`, `bench_zerodha.py:189`
- Checklist pattern: non-security checksum/fingerprint

Source excerpt:

```
results["md5"] = hashlib.md5(data).hexdigest()
# …
"md5": hashlib.md5(last_pdf).hexdigest() if last_pdf else "",
# …
"md5": {tier: hashlib.md5(data).hexdigest() for tier, data in outputs.items()},
```

Why this is a false positive: MD5 fingerprints bench PDF bytes for determinism/reporting, not passwords or trust integrity.

Checklist evidence: non-security fingerprint in benchmark harnesses.

---

### [ ] Findings 19, 22, 29 — CWE-22

- Function contexts: `./scripts/findings/functions/{19,22,29}.txt`
- Source: `…/engine/bench.py:207`, `bench_financial.py:68`, `bench_zerodha.py:80`
- Checklist pattern: no confined-directory join of untrusted segment

Source excerpt:

```
def _repo_root() -> Path:
    return Path(__file__).resolve().parent.parent
```

Why this is a false positive: Resolves `__file__` only; no dynamic/user segment join into a restricted root.

Checklist evidence: CWE-22 restricted-directory condition not present.

---

### [ ] Findings 20–21, 25–28, 32–35 — BP-PY-46

- Function contexts: `./scripts/findings/functions/{20,21,25,26,27,28,32,33,34,35}.txt`
- Source: `…/engine/bench.py`, `bench_financial.py`, `bench_zerodha.py` argparse `main()` prints
- Checklist pattern: CLI / `__main__` user-facing output

Source excerpt:

```
def main(...):
    ...
    for line in lines:
        print(line)
    ...
    print("\nappended to %s" % report)
```

Why this is a false positive: Benchmark CLIs report results/paths to stdout from argparse `main()`.

Checklist evidence: intentional CLI/user output under argparse main.

---

### [ ] Findings 23, 31, 41, 51, 76, 79, 82, 83, 84, 86, 90, 92 — CWE-90

- Function contexts: `./scripts/findings/functions/{23,31,41,51,76,79,82,83,84,86,90,92}.txt`
- Checklist pattern: wrong sink (not LDAP)

Source excerpt:

```
match = _PAGE_COUNT_RE.search(data)
# …
match = _TRAILER_RE.search(data)
# …
match = _STREAM_HEADER_RE.search(data, offset)
```

Why this is a false positive: Hits are Python `re` searches over PDF object/stream bytes, not LDAP search-filter construction or escaping.

Checklist evidence: sink is `re.Pattern.search` / `re.search`, not an LDAP filter API.

---

### [ ] Finding 36 — CWE-1046

- Function context: `./scripts/findings/functions/36.txt`
- Source: `…/engine/cipher.py:304:1`
- Checklist pattern: not immutable text concatenation

Source excerpt:

```
output = bytearray()
...
            output += bytes(byte ^ prev for byte, prev in zip(plain, previous))
```

Why this is a false positive: `bytearray` is mutable; `+=` extends in place, so the immutable-text concatenation condition is not met.

Checklist evidence: mutable buffer extend, not str/bytes reallocation concat.

---

### [ ] Findings 37, 80, 81 — CWE-328

- Function contexts: `./scripts/findings/functions/{37,80,81}.txt`
- Source: `…/engine/color.py:230` (+ ICC tests)
- Checklist pattern: format-required MD5 (ICC profile ID)

Source excerpt:

```
profile_id = hashlib.md5(profile).digest()
return profile[:84] + profile_id + profile[100:]
```

Why this is a false positive: ICC profile ID field is filled/validated with MD5 as required by the ICC profile format.

Checklist evidence: format-mandated non-security MD5.

---

### [ ] Finding 38 — CWE-328

- Function context: `./scripts/findings/functions/38.txt`
- Source: `…/engine/doc.py:288:18`
- Checklist pattern: PDF format ID fingerprint

Source excerpt:

```
digest = hashlib.md5()
writer.feed_digest(digest)
```

Why this is a false positive: MD5 builds the PDF trailer `/ID` document identifier, not password/auth material.

Checklist evidence: PDF document fingerprint, non-security use.

---

### [ ] Findings 42–43 — CWE-256 / CWE-798

- Function contexts: `./scripts/findings/functions/42.txt`, `43.txt`
- Source: `…/engine/fixtures.py:594:1`
- Checklist pattern: test/fixture literal, not production secret

Source excerpt:

```
        encrypt=EncryptSpec(
            password="fixture-password",
            revision=revision,
```

Why this is a false positive: Explicit fixture builder password for deterministic encrypted sample PDFs, not a deployed credential.

Checklist evidence: fixture-only known password literal.

---

### [ ] Finding 44 — CWE-73

- Function context: `./scripts/findings/functions/44.txt`
- Source: `…/engine/fixtures.py:725:15`
- Checklist pattern: intentional CLI path, no injection into confined tree

Source excerpt:

```
if __name__ == "__main__":
    out_dir = Path(sys.argv[1]) if len(sys.argv) > 1 else DEFAULT_OUTPUT_DIR
    for name, path in generate_fixtures(out_dir).items():
```

Why this is a false positive: Operator-chosen fixture output directory under `__main__`; written basenames come from a fixed builder map.

Checklist evidence: intentional CLI destination, not untrusted path injection into a sandboxed app.

---

### [ ] Findings 45, 47 — CWE-328

- Function contexts: `./scripts/findings/functions/45.txt`, `47.txt`
- Source: `…/engine/font.py:698`, `meta.py:126`
- Checklist pattern: non-security fingerprint

Source excerpt:

```
digest = hashlib.md5(face_name.encode("ascii")).digest()
# …
digest = hashlib.md5("|".join(parts).encode("utf-8")).digest()
```

Why this is a false positive: MD5 derives a deterministic font subset tag / XMP UUID seed, not crypto auth.

Checklist evidence: non-security naming/ID fingerprint.

---

### [ ] Finding 48 — CWE-22

- Function context: `./scripts/findings/functions/48.txt`
- Source: `…/engine/model.py:389:12`
- Checklist pattern: no restricted-directory join/escape

Source excerpt:

```
def load_template(path: Union[str, Path]) -> PDFTemplate:
    path = Path(path)
    with open(path, encoding="utf-8") as handle:
```

Why this is a false positive: Converts and opens the caller-supplied path; there is no join into a confined root to escape.

Checklist evidence: CWE-22 restricted-directory condition not present.

---

### [ ] Findings 52, 77, 78, 85 — CWE-328

- Function contexts: `./scripts/findings/functions/{52,77,78,85}.txt`
- Checklist pattern: format/test MD5 checksum (not security hashing)

Source excerpt:

```
self.assertEqual(
    hashlib.md5(first).hexdigest(), hashlib.md5(second).hexdigest()
)
```

Why this is a false positive: MD5 is used only as a test/determinism checksum of generated PDF bytes, not for authentication or password storage.

Checklist evidence: adjacent fixture determinism tests.

---

### [ ] Findings 53, 87, 91 — CWE-88

- Function contexts: `./scripts/findings/functions/{53,87,91}.txt`
- Checklist pattern: no user-controlled unsanitized option sink

Source excerpt:

```
result = subprocess.run(
    [VERAPDF_BIN, "-f", "4", str(path)],
```

Why this is a false positive: argv entries are a fixed binary, fixed/hardcoded `-f` flavour, and an absolute temp or fixed fixture path from test code — nothing attacker-controlled that can become a new option.

Checklist evidence: rule needs user control + unsanitized option-capable operand; absent here.

---

### [ ] Findings 54–74, 88, 89 — BP-PY-41

- Function contexts: `./scripts/findings/functions/54.txt` … `74.txt`, `88.txt`, `89.txt`
- Checklist pattern: helper `_assert_*` / `self.assert*` inside helper

Source excerpt:

```
def test_minimal_text_passes_pdfa4(self) -> None:
    self._assert_verapdf_compliant("minimal-text")
# helper contains:
self.assertEqual(result.returncode, 0, ...)
self.assertIn('isCompliant="true"', result.stdout)
```

Why this is a false positive: Tests only call `_assert_*` helpers that perform the unittest assertions.

Checklist evidence: `_assert_*` helper embeds `self.assert*`.

---

### [ ] Finding 75 — BP-PY-41

- Function context: `./scripts/findings/functions/75.txt`
- Source: `…/engine/tests/test_doc.py:83`
- Checklist pattern: raises inside helper

Source excerpt:

```
for marker in (b"/Type /Catalog", b"/Type /Pages /", b"/Type /Page /"):
    find_object_with(self.data, marker, self.offsets)
# helpers.py:
raise AssertionError(f"no object contains {marker!r}")
```

Why this is a false positive: Failure is asserted via `AssertionError` raised from `find_object_with`, not a pure side-effect-only test.

Checklist evidence: raises inside helper.

---

### [ ] Findings 93, 96, 97, 98, 101, 102 — BP-PY-46

- Function contexts: `./scripts/findings/functions/{93,96,97,98,101,102}.txt`
- Checklist pattern: CLI / `__main__` user-facing output

Source excerpt:

```
def main(...):
    ...
    print("usage: summarize_runs.py ...", file=sys.stderr)
    ...
    print("\n".join(lines))

if __name__ == "__main__":
    raise SystemExit(main())
```

Why this is a false positive: Prints are CLI usage/summary output from `main()` invoked under `__main__`, not library debug prints.

Checklist evidence: CLI/`__main__` user-facing output.

---

### [ ] Findings 94, 99 — CWE-22

- Function contexts: `./scripts/findings/functions/{94,99}.txt`
- Checklist pattern: no confined-directory join of untrusted segment

Source excerpt:

```
out = Path(argv[0])
stats_path = Path(argv[1])
```

Why this is a false positive: CLI intentionally takes OUT_DIR/STATS_PATH as whole paths; there is no base-directory join of a relative segment that escapes a restricted root.

Checklist evidence: CWE-22 needs user control **into** a restricted pathname; here the argv values *are* the intended roots/targets.

---

### [ ] Findings 95, 100 — PERF-PY-27

- Function contexts: `./scripts/findings/functions/{95,100}.txt`
- Checklist pattern: once-per-distinct-path

Source excerpt:

```
for path in unique_paths:
    text = path.read_text(encoding="utf-8", errors="replace")
```

Why this is a false positive: Loop iterates deduped `unique_paths`, so each distinct path is loaded once.

Checklist evidence: once-per-distinct-path (not repeated same path).

## True positives (summary)

| ID | Rule | Why true positive |
| --- | --- | --- |
| 7 | CWE-88 | `run_verapdf_check` passes CLI `flavour`/`pdf` into `subprocess.run([... pdf])` without `--` before the path operand or an allowlisted flavour set |
| 39 | BP-PY-46 | `print(...)` under `ENGINE_DEBUG_BUFFERS` inside reusable `engine/doc.py` render path |
| 46 | PERF-PY-24 | `wrap_text(...)` runs in `_row_height` and again in `_draw_row` for the same cells |
| 49 | BP-PY-1 | `except Exception:` around `base64.b64decode` swallows all failures |
| 50 | CWE-396 | same generic `Exception` handler hides distinct decode/type failures |

## Uncertain findings

### [ ] Finding 40 — CWE-328

- Function context: `./scripts/findings/functions/40.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/codehound-python-perf-targets/pythoncoreengine/engine/encrypt.py:222:18`

Source excerpt:

```
def _md5_iterate(data: bytes, iterations: int, length: int) -> bytes:
    """Algorithm 2 step h: ``iterations`` rounds of MD5 over ``length`` bytes."""
    digest = data
    for _ in range(iterations):
        digest = hashlib.md5(digest[:length]).digest()
    return digest
```

Why it is uncertain: Dual-use — this is PDF Standard-security password/key derivation (security-relevant) but MD5 is mandated by the PDF encryption algorithm for R4 compatibility, so “weak hash” vs “format-required” cannot be decided from the rule condition alone.

## Final evidence

- Delegated reviewers: [1–50](2c77607e-55bc-47f5-9b59-98a160234061), [51–102](652b2893-c31c-4595-9fed-003b87ced71b); parent reconciled TPs/uncertain against source
- Chunk evidence: `./scripts/chunks`
- Function evidence: `./scripts/findings/functions`
- Validation: `git diff --check` — pass
