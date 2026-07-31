# Python ruleset catalogues (WIP)

> **Status:** catalogues + priority **CWE**, **BP-PY**, and experimental **PERF-PY** detectors
> **Issues:** #52 CWE · #53 BP · #54 PERF 22/22 fixture-backed detectors · epic #51 · foundation #39
> **Plans:** `plans/v0.0.2/heuristics/python-heuristics-*.md`  
> **CWE mapping:** `plans/v0.0.2/init-python/python-cwe-from-699-mapping.md` (from `699.csv`)  
> **BP audit:** `plans/v0.0.2/init-python/ruleset-reuse-audit.md`

**Runnable detectors (priority batches):**
- CWE-22, CWE-78, CWE-79, CWE-89, CWE-502 — `internal/lang/python/detectors/cwe`
- BP-PY priority subset (core + security + sample Flask/Django) — `internal/lang/python/detectors/bad_practices`
- PERF-PY-1 through PERF-PY-22 — `internal/lang/python/detectors/perf` (experimental)

Pure-Go source patterns; catalogue JSON remains source of truth. Full CWE-344 / all 50 BP / full PERF detector parity not claimed.

## Layout

```text
ruleset/python/
  README.md
  bad-practices.json          # BP-PY-* (core + Flask/Django/FastAPI/SQLAlchemy/Jinja2)
  chunks/
    cwe-001-050.json          # ≤50 rules per file; ID-range names like golang
    cwe-051-100.json
    …
    cwe-1351-1400.json
    perf-py-001-014.json       # initial Python PERF seed
    perf-py-015-022.json       # remaining static evaluation gaps
```

Sibling Go catalogues remain the default product source of truth:

- `ruleset/golang/bad-practices.json`
- `ruleset/golang/chunks/*.json`

## CWE catalogue (from `699.csv`)

Source export: **CWE view 699** (Software Development), **399** rows.

| Class | Count | Meaning |
|-------|------:|---------|
| **Included · generic** | ~339 | `Not Language-Specific` (portable weakness) |
| **Included · python-specific** | 5 | Platforms list `LANGUAGE NAME:Python` |
| **Excluded** | ~55 | C/C++/Java/… only, or C-style memory corruption |

**~344** CWE entries are written under `chunks/cwe-*.json`.

### File naming and 50-rule limit

Same idea as golang (`cwe-001-050.json`, …):

- Filename is a **CWE-ID range** of width 50 (`001-050`, `051-100`, …, then `1001-1050`, …).
- Each file has **at most 50** JSON rule objects.
- Sparse ranges still get a file when at least one included rule falls in that band.

Regenerate from repo root:

```sh
# (generator was applied for this commit; see mapping doc for filter rules)
python3 -c 'print("see plans/v0.0.2/python-cwe-from-699-mapping.md")'
```

### JSON shape (per rule)

Same fields as golang CWE chunks:

`id`, `name`, `original_description`, `description`, `detection_notes`,  
`category`, `status`, `weakness_abstraction`, `python_relevance`, `applicable_to`

- `applicable_to` always includes `"python"`.
- `description` / `detection_notes` are Python-oriented (framework/stdlib sinks).
- `python_relevance` is the Python-catalogue counterpart of golang’s `go_relevance` (`High` / `Medium` / `Low`).

### Python-specific platform CWEs (from CSV filter)

| CWE | Name |
|-----|------|
| CWE-396 | Declaration of Catch for Generic Exception |
| CWE-397 | Declaration of Throws for Generic Exception |
| CWE-478 | Missing Default Case in Multiple Condition Expression |
| CWE-502 | Deserialization of Untrusted Data |
| CWE-915 | Improperly Controlled Modification of Dynamically-Determined Object Attributes |

## Bad practices

| File | Contents |
|------|----------|
| `bad-practices.json` | **50** `BP-PY-*` entries (catalogue metadata) |

IDs use **`BP-PY-*`** so they never collide with Go `BP-*` when both catalogues are listed.

### Detector status (issue #53)

Source-pattern heuristics live under `internal/lang/python/detectors/bad_practices/`.
They run only when the Python plugin is enabled (`languages = ["python"]` / multi-language registry).

**Shipped (first land):**

| Batch | IDs |
|-------|-----|
| A — Core | `BP-PY-1`, `2`, `4`, `6`, `7` |
| B — Security hygiene | `BP-PY-8` … `13` |
| C — Framework (high-signal) | `BP-PY-16`, `17`, `21` |

Remaining catalogue IDs (`BP-PY-3`, `5`, rest of C/D/E) are deferred — see the BP ledger.
The Python plugin contains the priority CWE/BP batches and all 22 seeded PERF-PY rules.

## Performance

`chunks/perf-py-001-014.json` and `chunks/perf-py-015-022.json` contain the Python-only PERF catalogue:
twenty-two source-verifiable cost signals extracted from the FastAPI, Django, and Flask performance
evaluation corpus. IDs use the
`PERF-PY-*` namespace, matching `BP-PY-*`, so Python never consumes or collides with Go's numeric PERF
sequence. The seed deliberately leaves async blocking and HTTP timeouts with `BP-PY-30` and
`BP-PY-14` respectively, avoiding duplicate findings.

Every PERF entry is marked `Seeded`, has `applicable_to` including `python`, and uses
Python-specific `detection_notes`. All 22 are emitted by the opt-in Python plugin and have paired
vulnerable/safe fixtures under `tests/fixtures/python/perf/`. They remain experimental until a
reference-corpus canary supports a maturity decision.
FA-8 remains a correctness/operability issue, while XC-1/XC-3/XC-4/XC-5 are benchmark, deployment,
observability, or security work rather than source PERF rules; XC-2 is represented by `PERF-PY-17`.

## Reuse policy

1. **Same JSON shape** as golang CWE/PERF chunks and BP maps.
2. **Do not bulk-copy** Go PERF or Go BP into Python. Python PERF uses the distinct `PERF-PY-*` namespace.
3. **Do not** point Go generators (`metadata_gen.go`) at this directory.
4. Prefer rewriting detection notes for Python sinks rather than cloning gin/gorm text.

## What this is not

- Not a claim that experimental PERF-PY findings belong in the recommended/perf profile
- Not a substitute for reference-corpus canary and false-positive review

## Validation

```sh
go test ./ruleset/python/
```
