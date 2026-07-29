# Taint Tracking (Experimental)

CodeHound includes an experimental taint-tracking engine for Go that augments
substring-based CWE detectors for the **taint-core** family:

`CWE-22`, `CWE-78`, `CWE-79`, `CWE-89`, `CWE-90`, `CWE-91`

When enabled, it traces data flow from untrusted sources to dangerous sinks and
suppresses findings where a recognized sanitizer intercepts the flow.

**Honesty bar:** name-string and intra/same-package analysis — useful for triage,
**not** a CodeQL / whole-program security gate. See [ADR 0003](./adr/0003-taint-model.md)
and [SECURITY.md](../SECURITY.md).

## Enabling

| Method | How |
|--------|-----|
| CLI flag | `codehound --taint` |
| Config file | `[codehound.taint]\nenabled = true` |
| Profile | `--profile security` (taint **on** by default) |
| Disable | `--no-taint` or `[codehound.taint]\nenabled = false` |
| Show paths | `--taint-show-paths` or `[codehound.taint]\nshow_paths = true` |
| Inter-proc depth | `--taint-depth N` (1–4, default **1** = direct caller→callee only) |

Taint is disabled by default under `recommended` (CWE IDs stay allow-listed).
The substring-based heuristic still runs as a fallback when taint is off.

### Pipeline (high level)

1. **Extract** source/sink/sanitizer annotations from the Go AST (`taint/extract/`).
2. **Build** an intra-procedural data-flow graph per function.
3. **Query** paths from sinks back to sources; apply sanitizers.
4. **Accumulate** per-file units; on project **finalize**, refine
   same-package call summaries with multi-hop depth 1–4.
5. Emit evidence (`TaintFlow`) when `--taint-show-paths` is set.

Cache hits with taint enabled still re-parse or re-accumulate project units when
`requires_cache_state` is set — warm cache is not free for taint-heavy scans.
See [incremental-cache.md](./incremental-cache.md).

### Intra-proc precision (Phase 8)

- **Versioned last-write:** each assignment is a versioned node; uses resolve to
  the latest declaration with `decl_byte ≤ use_byte` (overwrite with a constant
  kills live taint).
- **Field keys:** LHS/RHS like `user.Path` are tracked as qualified keys (not
  only the base `user`).
- **Map/slice index:** `m[k] = t` conservatively taints base `m` (low precision;
  no per-key model).

## Model

The engine builds an intra-procedural data-flow graph per file. Each node is
a variable or expression; edges represent assignment, arithmetic, or
function-call return. The graph is searched backward from each sink to find
paths to recognized sources.

### Sources

Classifier lives in `taint/extract/classify.rs` (source of truth over this table).

| Kind | Examples |
|------|----------|
| `UserInput` | `r.URL.Query()`, `r.FormValue()`, `r.PostForm`, `r.Header.Get`, `r.PathValue`, framework params (Gin/Echo/Chi/Fiber) |
| `Args` | `os.Args`, `flag.Args()`, `flag.String()`, `flag.Int` |
| `EnvVar` | `os.Getenv()`, `os.LookupEnv()` |
| `File` | `os.ReadFile()`, `ioutil.ReadFile()`, `os.Open()`, `io.ReadAll` on files |
| `Network` | `net.Conn.Read()`, `http.Request.Body`, scanners over network readers |

### Sinks

| Kind | CWE | Examples |
|------|-----|----------|
| `CommandExec` | CWE-78 | `exec.Command()`, `(*Cmd).Run/Start/Output` |
| `SQLQuery` | CWE-89 | `(*sql.DB).Query/Exec`, `(*sql.Tx).Query/Exec` |
| `FileOpen` | CWE-22 | `os.Create()`, `os.OpenFile()`, `os.WriteFile()` |
| `Template` | CWE-79 | `(*template.Template).Execute()` |
| `HTTPWrite` | CWE-79 | `w.Write()`, `w.WriteHeader()` |
| `LDAPQuery` | CWE-90 | LDAP search / filter construction sinks |
| `XMLQuery` | CWE-91 | XML query / decode sinks treated as injection targets |
| `Deserialization` | — | `json.Unmarshal()`, `xml.Unmarshal()` |

### Sanitizers

Source of truth is the runtime **classifier**, not enum comments in `model.rs`.

| Kind | Examples |
|------|----------|
| `Path` | `filepath.Base()` only (strips to final component). **`filepath.Clean` / `path.Clean` alone are not path-safe** and are not treated as sanitizers. |
| `HTML` | `html.EscapeString()`, `template.HTMLEscaper()` |
| `URL` | `url.QueryEscape()`, `url.PathEscape()` |
| `SQL` | **Not** bare `.Prepare` (dynamic Prepare is still injectable). Safe paths: (1) **literal** first arg at Query/Exec; (2) same-function same-var literal `Prepare`/`PrepareContext` → `Stmt.Query`/`Exec` (or `*Context` forms). |
| `LDAP` / `XML` | Escape helpers recognized for CWE-90/91 when classified |
| `Validation` | `strconv.Atoi()`, `strconv.ParseInt()`, allowlist-style `sanitize*` / `escape*` / `validate*` helpers |
| `Bounded` | `len` (weak; prefer explicit validation) |

**Path confinement (CWE-22):** findings are suppressed only when a path
variable is confined with `filepath.Abs` / `filepath.EvalSymlinks` **and**
`strings.HasPrefix` on that same binding — not file-level co-presence of
`Clean`.

When a recognized sanitizer is found on every path from the source to the
sink, the finding is suppressed. CWE-78 (command injection) does **not**
accept Path sanitizers.

## Limitations

- **`filepath.Clean` is not a path sanitizer.** Clean-only code should still
  flag CWE-22 (or show low confidence in a future release).
- **CWE-89 is a heuristic, not full SQLi.** Literal-first-arg Query/Exec is
  treated as parameterized/safe; GORM `.Raw` / sqlx helpers with dynamic SQL
  still fire. No claim of complete ORM coverage.
- **CWE-79 is not full XSS.** Template + HTTPWrite sinks only; no DOM model.
- **CWE-90/91** use real LDAP/XML sink classification; quarantine long-tail
  CWE IDs that only match fixture needles (catalog honesty — Phase 1).
- **Bare `.Prepare` is not a sanitizer.** CWE-89 suppresses only when the
  **same function** binds a simple receiver via literal `Prepare`/`PrepareContext`
  and that binding is the latest write before `Stmt.Query`/`Exec` (or `*Context`).
  Dynamic Prepare SQL, rebinding, and cross-function Prepare factories are not
  proven safe. Residual FN: dynamic Prepare with no tainted Query/Exec arg.
- **Inter-procedural tracking is depth-limited and same-package only.**
  Cross-function analysis works for direct chains (A→B→C), return
  propagation, method calls, and recursion within one package (keyed by
  package directory + clause, receiver type, and function name). Unqualified
  callees never resolve into another package; import-path wiring is not
  implemented. Method calls with an unknown receiver type are resolved only
  when a single same-package receiver exposes that method name; if multiple
  receivers share the name, summary resolution is declined (false negative
  preferred over the wrong summary) until local type inference exists.
  Receiver type is preserved when the call is on the enclosing method's own
  receiver parameter. Depth is bounded (not a full fixpoint). Mutual
  recursion and deep chains may miss flows.
- **Limited field keys (not full field-sensitive analysis).** Qualified names
  like `user.Path` are tracked as keys; map/slice index writes taint the base
  (`m[k] = t` → `m`). No full field-sensitive / element-precise model.
- **String concatenation via `+` only.** `fmt.Sprintf` is partially handled
  but not all format-string call graphs are resolved.
- **No type-based aliasing.** Two variables of the same type pointing to the
  same allocation are treated as independent.
- **Interface dispatch.** Methods called on interface types are treated as
  opaque — taint flows through arguments but the return value is not tracked
  because the concrete implementation is unknown.
- **Channel/goroutine.** G5 v0 pilot models **same-function** unbuffered
  handoff when pairing rules hold: exactly one send and one receive on the
  same channel binding in one lexical function/method body → a dedicated
  `ChannelTransfer` edge (not assignment sugar). Path finder follows that
  edge so classic `source → send → recv → sink` can fire.
  Declined shapes (`select`, multi-send/recv, buffered `make(chan T, N)`,
  cross-goroutine / `go` MHP) stay `UnsupportedFlow::{Channel,Goroutine}`
  — honest FN. The old IP-010 source-on-send attribution quirk is
  quarantined (send-value sources no longer bind the channel identifier).
  Do **not** treat residual IP-010 fixtures as channel-support success.
- **Pointer dereference.** `*p = tainted` and `json.Unmarshal(data, &target)`
  are handled for a small set of known functions (`json.Unmarshal`,
  `xml.Unmarshal`). General pointer tracking requires type inference.
- **Name-string sinks.** Callees are matched by identifier text, not types —
  renamed wrappers and interface methods may FN or FP.
- **No SSA / no Go types.** Intra-proc last-write and call facts are AST-level.
- **Depth.** Inter-procedural summary is bounded (not a full fixpoint). Mutual
  recursion and deep chains may miss flows.
- **Product positioning.** Enable with `--profile security` or `--taint` for
  triage. **Do not use as a sole security gate** — pair with govulncheck,
  code review, and stronger SAST where required.

See also [ADR 0003 — taint model honesty](./adr/0003-taint-model.md).

## Output

With `--taint-show-paths`, findings include a `TaintFlow` evidence block with
the source kind, sink kind, and hop count (plus per-hop details for
inter-procedural findings):

```json
{
  "evidence": {
    "kind": "TaintFlow",
    "source": { "kind": "UserInput", "function": "r.URL.Query", "variable": "host" },
    "sink": {
      "kind": "CommandExec",
      "function": "exec.Command",
      "hop_details": [
        { "function": "runCommand", "kind": "CommandExec",
          "variable": "cmd", "file": "handler.go", "line": 42 }
      ]
    },
    "hops": 1,
    "sanitized": false
  }
}
```

The text reporter shows a summary line:

```
taint flow UserInput.r.URL.Query -> CommandExec.exec.Command across 1 hop
  hop: runCommand(cmd) at handler.go:42
```

## Custom Sanitizers

Taint recognizes sanitizers by function name matching. Any function whose name
matches the regex `^(sanitize|clean|escape|validate|purify)` is treated as a
`Validation` sanitizer. Framework bind methods (`c.ShouldBind`,
`c.ShouldBindJSON`) are treated as `Validation` sanitizers when the Gin or
Echo packages are imported.

To extend sources/sinks/sanitizers, edit the classifier and walkers under
`src/lang/go/detectors/cwe/taint/extract/` (`classify.rs`, walkers, call graph)
and graph query helpers under `taint/graph_query/`. Enum kinds live in
`taint/model.rs` — runtime classification is the product source of truth.
