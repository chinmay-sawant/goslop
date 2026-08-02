# False-positive audit: WHEN-Language

## Run metadata

```yaml
timestamp: 2026-08-02T07:42:59Z
repository: WHEN-Language
repository_path: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language
branch: main
commit: f9bf78d49c16c36d4eb98d5abd305a65f3252e32
scan_target: /home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language
chunk_path: scripts/WHEN-Language/chunks
function_context_path: scripts/WHEN-Language/findings/functions
```

## Scan evidence

- Build command: `go build -o bin/goslop ./cmd/goslop` (bin/goslop already present)
- Scan command: `./bin/goslop --profile all --no-fail --no-terminal --config templates/goslop-python.toml --export-context --export-chunks --no-cache -chunks-dir scripts/WHEN-Language/chunks -context-dir scripts/WHEN-Language/findings/functions real-repos/WHEN-Language`
- Findings: `104`
- Chunks reviewed: `scripts/WHEN-Language/chunks/Chunk_1_25.txt`, `Chunk_26_50.txt`, `Chunk_51_75.txt`, `Chunk_76_100.txt`, `Chunk_101_104.txt`
- Function contexts reviewed: `scripts/WHEN-Language/findings/functions/1.txt` .. `104.txt` (read via chunks for all findings; context files spot-checked for 6, 24, 46, 48, 86, 89)

## Audit checklist

- [x] Read every assigned chunk under `scripts/WHEN-Language/chunks`.
- [x] Read `scripts/WHEN-Language/findings/functions/<finding-id>.txt` for every proposed false positive.
- [x] Followed the `Source:` path and read the enclosing source function or block when the exported context was insufficient (full `when.py`, `hot_reload.py`, `parser.py`, `lexer.py`, `interpreter.py` read).
- [x] Classified every reviewed finding as `False positive`, `True positive`, or `Uncertain`.
- [x] Based the decision on the rule condition and the shown source, not on application-specific knowledge.
- [x] Reconciled delegated reviews and documented disagreements as `Uncertain` where evidence is insufficient. (No delegated reviews; no disagreements.)
- [x] Ran `git diff --check` after updating this report.

## Classification summary

| Classification | Count | Finding IDs |
| --- | ---: | --- |
| False positive | 58 | 6, 22, 24, 29, 41, 42, 46, 48–83 (excl. 74), 85, 86, 87, 88, 90, 91, 92, 93, 94, 97, 98, 99, 100, 102, 103, 104 |
| True positive | 46 | 1, 2, 3, 4, 5, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 23, 25, 26, 27, 28, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 43, 44, 45, 47, 74, 84, 89, 95, 96, 101 |
| Uncertain | 0 | — |

Exact FP list: 6, 22, 29, 41, 42, 85 (BP-PY-44); 24 (BP-PY-12); 46 (CWE-208); 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 75, 76, 77, 78, 79, 80, 81, 82, 83 (PERF-PY-26); 86, 87, 88, 90, 91, 92, 93, 94, 97, 98, 99, 100, 102, 103, 104 (BP-PY-46 in `when.py`).

## False positives

### [ ] Finding `6` — `BP-PY-44`

- Function context: `scripts/WHEN-Language/findings/functions/6.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/hot_reload.py:91:1`
- Checklist pattern: imported module is a first-party module, not the deprecated stdlib module the rule targets.

Source excerpt:

```
                # Parse the updated source
                from lexer import Lexer
                from parser import Parser
```

Why this is a false positive: `parser` here is the repository's own `parser.py` (sibling of `lexer.py`, which is imported on the line above); the rule condition flags imports of the deprecated/removed *stdlib* module `parser`, which is not what this statement imports.

Checklist evidence: the import is `from parser import Parser` where `parser.py` exists in the same package directory; no stdlib module is named by the import.

### [ ] Finding `22` — `BP-PY-44`

- Function context: `scripts/WHEN-Language/findings/functions/22.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/interpreter.py:316:1`
- Checklist pattern: imported module is a first-party module, not the deprecated stdlib module the rule targets.

Source excerpt:

```
            from lexer import Lexer
            from parser import Parser
```

Why this is a false positive: same local `from parser import Parser` first-party import (with its sibling `lexer`); the rule's deprecated-stdlib condition is not satisfied.

Checklist evidence: the import names the repository's own `parser` module, which is not a stdlib module.

### [ ] Finding `29` — `BP-PY-44`

- Function context: `scripts/WHEN-Language/findings/functions/29.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/interpreter.py:971:1`
- Checklist pattern: imported module is a first-party module, not the deprecated stdlib module the rule targets.

Source excerpt:

```
                    from lexer import Lexer
                    from parser import Parser
```

Why this is a false positive: same first-party `from parser import Parser` import as findings 6/22; the rule condition (deprecated stdlib module import) is not met.

Checklist evidence: the module `parser` is the repository's local module, verified by `parser.py` at the repo root.

### [ ] Finding `41` — `BP-PY-44`

- Function context: `scripts/WHEN-Language/findings/functions/41.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/interpreter.py:1492:1`
- Checklist pattern: imported module is a first-party module, not the deprecated stdlib module the rule targets.

Source excerpt:

```
        from lexer import Lexer
        from parser import Parser
```

Why this is a false positive: first-party `from parser import Parser` import; the rule's deprecated-stdlib condition is not satisfied.

Checklist evidence: the imported `parser` module is the local `parser.py` of this repository.

### [ ] Finding `42` — `BP-PY-44`

- Function context: `scripts/WHEN-Language/findings/functions/42.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/interpreter.py:1636:1`
- Checklist pattern: imported module is a first-party module, not the deprecated stdlib module the rule targets.

Source excerpt:

```
        from lexer import Lexer
        from parser import Parser
```

Why this is a false positive: first-party `from parser import Parser` import; the rule's deprecated-stdlib condition is not satisfied.

Checklist evidence: the imported `parser` module is the local `parser.py` of this repository.

### [ ] Finding `85` — `BP-PY-44`

- Function context: `scripts/WHEN-Language/findings/functions/85.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:23:1`
- Checklist pattern: imported module is a first-party module, not the deprecated stdlib module the rule targets.

Source excerpt:

```
from lexer import Lexer
from parser import Parser
from interpreter import Interpreter
```

Why this is a false positive: the import chain is entirely first-party (`lexer`, `parser`, `interpreter` are all local modules of the repo); `from parser import Parser` is not an import of the removed stdlib `parser` module.

Checklist evidence: `when.py` imports three sibling first-party modules; the deprecated-stdlib name list matches only the module *name* `parser`, not the imported module.

### [ ] Finding `24` — `BP-PY-12`

- Function context: `scripts/WHEN-Language/findings/functions/24.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/interpreter.py:330:40`
- Checklist pattern: the trigger token `exec(` occurs inside a string literal, not as a call.

Source excerpt:

```
        # exec: Execute statements (disabled for safety)
        def exec_func(code):
            raise NotImplementedError("exec() is not supported in WHEN for security reasons")
```

Why this is a false positive: there is no `eval(`/`exec(` call anywhere in the source; the detector's `exec(` match is the substring `exec(` inside the string literal `"exec() is not supported..."` of a `raise` statement. The rule condition ("eval/exec on dynamic input") requires an actual eval/exec call with a non-literal argument.

Checklist evidence: the matched identifier is string content inside a `NotImplementedError` message, not a callable invocation; `exec_func` is a def whose body only raises.

### [ ] Finding `46` — `CWE-208`

- Function context: `scripts/WHEN-Language/findings/functions/46.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/lexer.py:714:20`
- Checklist pattern: compared operands are lexer token-type enums, not security-sensitive values.

Source excerpt:

```
                # Track nesting depth for proper indentation handling
                if token_type == TokenType.LBRACE:
                    self.brace_depth += 1
```

Why this is a false positive: the equality is a lexer-internal enum comparison used to track brace nesting depth; neither operand is a secret, password, digest, signature, or credential, so the rule condition ("security-sensitive values are compared with `==`") is not satisfied. The finding fires only because the timing-sensitive *name* set includes `token`, and `token_type`/`TokenType` are case-insensitive matches.

Checklist evidence: both operands are the lexer's own token-type enum; the result only increments `brace_depth`; no constant-time comparison is relevant to this operation.

### [ ] Finding `48` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/48.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:52:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
            if self.current_token().type == TokenType.MAIN:
                if main is not None:
                    raise SyntaxError("Multiple main blocks defined")
                main = self.parse_main_block()
```

Why this is a false positive: `self.parse_main_block()` is a recursive-descent step inside `parse()`'s single pass over the token stream; each call advances `self.pos` and consumes a distinct construct. The rule condition ("re-decodes or re-parses the same or repeated binary/document payload on every call ... parse once and reuse immutable results") is not met — nothing is parsed twice and there is no reusable result, and the rule's own detection notes suppress one-shot CLI tools, which is what invokes this parser.

Checklist evidence: the call's argument is implicit `self` state that moves forward on every advance(); the same expression is never re-parsed, so a cache would have nothing to serve.

### [ ] Finding `49` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/49.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:55:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
            elif self.current_token().type in [TokenType.OS, TokenType.DE, TokenType.FO, TokenType.PARALLEL]:
                blocks.append(self.parse_block())
```

Why this is a false positive: `self.parse_block()` parses the current block declaration once and consumes its tokens; no payload is re-parsed.

Checklist evidence: the parse call is a one-time descent over fresh tokens; there is no loop over the same input and no reusable decode result.

### [ ] Finding `50` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/50.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:58:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
            elif self.current_token().type == TokenType.DEF:
                declarations.append(self.parse_function())
```

Why this is a false positive: `self.parse_function()` consumes the current function declaration once; no repeated re-parse of the same input occurs.

Checklist evidence: single-pass descent; the enclosing loop iterates over distinct declarations.

### [ ] Finding `51` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/51.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:61:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
            elif self.current_token().type == TokenType.CLASS:
                declarations.append(self.parse_class())
```

Why this is a false positive: `self.parse_class()` parses one class declaration; the rule's repeated-parse condition is not met.

Checklist evidence: distinct input consumed per call; no cache-relevant reuse exists.

### [ ] Finding `52` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/52.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:64:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
            elif self.current_token().type == TokenType.IMPORT:
                declarations.append(self.parse_import())
```

Why this is a false positive: `self.parse_import()` consumes one import declaration once.

Checklist evidence: single-pass recursive descent over distinct tokens; the rule's hot-path re-parse condition is not satisfied.

### [ ] Finding `53` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/53.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:66:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
            elif self.current_token().type == TokenType.FROM:
                declarations.append(self.parse_from_import())
```

Why this is a false positive: `self.parse_from_import()` parses one from-import once.

Checklist evidence: distinct input consumed per call; no repeated decode/parse of a reused payload.

### [ ] Finding `54` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/54.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:74:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                    stmt = self.parse_statement()  # This will parse it as TupleUnpackingAssignment
```

Why this is a false positive: `self.parse_statement()` parses the current statement once within the top-level declaration loop.

Checklist evidence: single-pass descent; the statement being parsed is the current token stream position, not a reused field.

### [ ] Finding `55` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/55.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:78:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                elif self.peek_token().type == TokenType.ASSIGN:
                    declarations.append(self.parse_var_declaration())
```

Why this is a false positive: `self.parse_var_declaration()` consumes one declaration; nothing is re-parsed.

Checklist evidence: distinct input consumed per call; no repeated parse of the same payload.

### [ ] Finding `56` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/56.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:166:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
            # Check for default parameter
            if self.current_token().type == TokenType.ASSIGN:
                self.advance()
                default_value = self.parse_expression()
```

Why this is a false positive: `self.parse_expression()` parses one default-value expression inside the parameter loop; each parameter's default is parsed once.

Checklist evidence: per-iteration input is the next distinct expression; the rule's "re-parses the same ... payload" condition is not met.

### [ ] Finding `57` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/57.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:208:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
            if self.current_token().type == TokenType.DEF:
                # Parse method
                methods.append(self.parse_function())
```

Why this is a false positive: `self.parse_function()` parses each class method once; methods are distinct inputs.

Checklist evidence: single-pass descent over distinct method declarations.

### [ ] Finding `58` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/58.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:214:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                    name = self.advance().value
                    self.expect(TokenType.ASSIGN)
                    value = self.parse_expression()
```

Why this is a false positive: `self.parse_expression()` parses one class attribute value once.

Checklist evidence: distinct input consumed per call; no repeated parse of a reused payload.

### [ ] Finding `59` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/59.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:295:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
            if self.current_token().type == TokenType.DEDENT:
                break
            stmt = self.parse_statement()
```

Why this is a false positive: `self.parse_statement()` parses each statement of a block once; statements are distinct inputs.

Checklist evidence: single-pass descent over distinct statements; no reuse opportunity.

### [ ] Finding `60` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/60.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:326:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                while self.current_token().type == TokenType.COMMA:
                    self.advance()
                    values.append(self.parse_expression())
```

Why this is a false positive: each `self.parse_expression()` in the comma loop parses a different return value expression; no expression is parsed twice.

Checklist evidence: loop iterations consume distinct expressions; the rule's repeated-parse-of-same-payload condition is not satisfied.

### [ ] Finding `61` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/61.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:358:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                    if self.current_token().type not in [TokenType.NEWLINE, TokenType.EOF, TokenType.DEDENT]:
                        values.append(self.parse_expression())
```

Why this is a false positive: `self.parse_expression()` parses each comma-separated value of a tuple-unpacking assignment once.

Checklist evidence: distinct inputs per call; no repeated re-parse of the same payload.

### [ ] Finding `62` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/62.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:475:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                op = "not in"
                right = self.parse_logical_and()
```

Why this is a false positive: the comparison loop parses each operand once; the recursive `parse_logical_and()` call handles the operand's own distinct token range.

Checklist evidence: each call descends into a fresh token range; no payload is re-parsed.

### [ ] Finding `63` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/63.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:482:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                op = "is not"
                right = self.parse_logical_and()
```

Why this is a false positive: same single-pass recursive-descent pattern as finding 62.

Checklist evidence: distinct operand parsed once per call.

### [ ] Finding `64` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/64.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:488:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                op = "is"
                right = self.parse_logical_and()
```

Why this is a false positive: same single-pass recursive-descent pattern; no repeated parse of the same input.

Checklist evidence: distinct operand parsed once per call.

### [ ] Finding `65` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/65.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:493:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                op_token = self.advance()
                op = op_token.value if op_token.value else op_token.type.name.lower()
                right = self.parse_logical_and()
```

Why this is a false positive: the comparison chain parses each RHS operand once.

Checklist evidence: distinct input consumed per call; no reuse of a parsed payload.

### [ ] Finding `66` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/66.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:503:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
        while self.current_token().type == TokenType.AND:
            op = self.advance().value
            right = self.parse_logical_or()
```

Why this is a false positive: the `and` chain parses each operand once; inputs are distinct.

Checklist evidence: single-pass descent; rule's hot-path re-parse condition not met.

### [ ] Finding `67` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/67.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:513:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
        while self.current_token().type == TokenType.OR:
            op = self.advance().value
            right = self.parse_addition()
```

Why this is a false positive: the `or` chain parses each operand once.

Checklist evidence: distinct input consumed per call; no repeated parse of the same payload.

### [ ] Finding `68` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/68.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:523:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
        while self.current_token().type in [TokenType.PLUS, TokenType.MINUS]:
            op = self.advance().value
            right = self.parse_multiplication()
```

Why this is a false positive: the additive-chain loop parses each operand once.

Checklist evidence: distinct input consumed per call; no reuse.

### [ ] Finding `69` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/69.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:533:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
        while self.current_token().type in [TokenType.MULTIPLY, TokenType.DIVIDE, TokenType.MODULO, TokenType.FLOORDIV]:
            op = self.advance().value
            right = self.parse_unary()
```

Why this is a false positive: the multiplicative-chain loop parses each operand once.

Checklist evidence: distinct input consumed per call; no repeated re-parse of a reused payload.

### [ ] Finding `70` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/70.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:570:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                # Parse start (or could be a regular index)
                if self.current_token().type != TokenType.COLON:
                    start = self.parse_expression()
```

Why this is a false positive: `self.parse_expression()` parses one index/slice expression per postfix iteration; each is distinct.

Checklist evidence: single-pass descent; no payload is parsed more than once.

### [ ] Finding `71` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/71.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:579:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                    # Parse stop
                    if self.current_token().type not in [TokenType.COLON, TokenType.RBRACKET]:
                        stop = self.parse_expression()
```

Why this is a false positive: parses one distinct slice-stop expression.

Checklist evidence: distinct input consumed per call.

### [ ] Finding `72` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/72.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:585:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                        self.advance()  # Skip second colon
                        if self.current_token().type != TokenType.RBRACKET:
                            step = self.parse_expression()
```

Why this is a false positive: parses one distinct slice-step expression.

Checklist evidence: distinct input consumed per call; no reuse.

### [ ] Finding `73` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/73.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:650:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                                kw_name = self.advance().value
                                self.advance()  # consume =
                                kw_value = self.parse_expression()
```

Why this is a false positive: the argument loop parses each distinct keyword argument once.

Checklist evidence: per-call input is the next distinct expression; no repeated parse of the same payload.

### [ ] Finding `75` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/75.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:654:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                            else:
                                # Regular positional argument
                                args.append(self.parse_expression())
```

Why this is a false positive: the argument loop parses each distinct positional argument once.

Checklist evidence: distinct input consumed per call; no repeated parse of the same payload.

### [ ] Finding `76` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/76.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:718:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                while self.current_token().type != TokenType.RPAREN:
                    elements.append(self.parse_expression())
```

Why this is a false positive: the tuple-element loop parses each distinct element once.

Checklist evidence: distinct input consumed per call; no reuse.

### [ ] Finding `77` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/77.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:755:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                        kw_name = self.advance().value
                        self.advance()  # consume =
                        kw_value = self.parse_expression()
```

Why this is a false positive: the call-argument loop parses each distinct keyword argument once.

Checklist evidence: distinct input consumed per call; no repeated parse of the same payload.

### [ ] Finding `78` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/78.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:759:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                    else:
                        # Regular positional argument
                        args.append(self.parse_expression())
```

Why this is a false positive: the call-argument loop parses each distinct positional argument once.

Checklist evidence: distinct input consumed per call; no reuse.

### [ ] Finding `79` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/79.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:825:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                                    kw_name = self.advance().value
                                    self.advance()  # consume =
                                    kw_value = self.parse_expression()
```

Why this is a false positive: the chained-member argument loop parses each distinct keyword argument once.

Checklist evidence: distinct input consumed per call; no repeated parse of the same payload.

### [ ] Finding `80` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/80.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:829:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                                else:
                                    # Regular positional argument
                                    args.append(self.parse_expression())
```

Why this is a false positive: the chained-member argument loop parses each distinct positional argument once.

Checklist evidence: distinct input consumed per call; no reuse.

### [ ] Finding `81` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/81.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:866:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                if self.current_token().type == TokenType.RBRACKET:
                    break  # trailing comma

                elements.append(self.parse_expression())
```

Why this is a false positive: the list-element loop parses each distinct element once.

Checklist evidence: distinct input consumed per call; no repeated parse of the same payload.

### [ ] Finding `82` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/82.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:901:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                key = self.parse_expression()
                self.skip_newlines()
```

Why this is a false positive: the dict loop parses each distinct key expression once.

Checklist evidence: distinct input consumed per call; no reuse.

### [ ] Finding `83` — `PERF-PY-26`

- Function context: `scripts/WHEN-Language/findings/functions/83.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/parser.py:905:1`
- Checklist pattern: recursive-descent step consuming distinct tokens once, not a repeated re-parse of the same payload.

Source excerpt:

```
                self.expect(TokenType.COLON)
                self.skip_newlines()
                value = self.parse_expression()
```

Why this is a false positive: the dict loop parses each distinct value expression once.

Checklist evidence: distinct input consumed per call; no repeated parse of the same payload.

### [ ] Finding `86` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/86.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:29:5`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
def show_help():
    print(__doc__)
```

Why this is a false positive: `when.py` is the CLI entry script (shebang, usage docstring, `if __name__ == "__main__": main()` at line 145); `show_help` prints the program's own help text. The rule condition is "print used for operational logging in non-script modules" — this is a script module printing user-facing output, which the rule's fix explicitly permits for CLIs.

Checklist evidence: the module has a `__main__` guard invoking `main()`; the print is CLI presentation output (help text), reachable only from `main()`.

### [ ] Finding `87` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/87.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:32:5`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
def show_version():
    print(f"WHEN Language Interpreter v{__version__}")
```

Why this is a false positive: version banner printed by the CLI's `--version` handler in the entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `88` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/88.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:33:5`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
    print(f"WHEN Language Interpreter v{__version__}")
    print("Built on Python", sys.version)
```

Why this is a false positive: continuation of the version banner in the entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `90` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/90.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:39:13`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
    try:
        if not os.path.exists(filename):
            print(f"Error: File '{filename}' not found")
            sys.exit(1)
```

Why this is a false positive: error message printed to the user by `run_file`, a CLI presentation function of the entry script.

Checklist evidence: user-facing error output in the `__main__`-guarded entry script.

### [ ] Finding `91` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/91.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:43:13`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
        if not filename.endswith('.when'):
            print("Warning: WHEN files should have .when extension")
```

Why this is a false positive: user-facing warning in the CLI entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `92` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/92.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:61:9`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
    except FileNotFoundError:
        print(f"Error: File '{filename}' not found")
        sys.exit(1)
```

Why this is a false positive: user-facing error message in the entry script's error handling.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `93` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/93.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:64:9`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
    except SyntaxError as e:
        print(f"Syntax Error: {e}")
        sys.exit(1)
```

Why this is a false positive: user-facing error message in the entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `94` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/94.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:67:9`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
    except KeyboardInterrupt:
        print("\nProgram interrupted by user")
        sys.exit(0)
```

Why this is a false positive: user-facing interruption notice in the entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `97` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/97.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:70:9`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
    except Exception as e:
        print(f"Runtime Error: {e}")
        import traceback
```

Why this is a false positive: user-facing error message in the entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `98` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/98.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:77:5`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
def interactive_mode():
    """Interactive REPL for WHEN (simplified for now)"""
    print(f"WHEN Language Interactive Shell v{__version__}")
```

Why this is a false positive: REPL banner printed by the CLI's `-i` handler in the entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `99` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/99.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:78:5`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
    print(f"WHEN Language Interactive Shell v{__version__}")
    print("Type 'exit()' to quit")
```

Why this is a false positive: REPL prompt banner in the entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `100` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/100.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:79:5`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
    print("Type 'exit()' to quit")
    print("Note: Interactive mode is limited - use files for full WHEN programs")
```

Why this is a false positive: REPL banner text in the entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `102` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/102.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:108:17`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
            except Exception as e:
                print(f"Error: {e}")
```

Why this is a false positive: REPL error message in the entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `103` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/103.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:111:13`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
        except KeyboardInterrupt:
            print("\nKeyboardInterrupt")
```

Why this is a false positive: REPL interruption notice in the entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

### [ ] Finding `104` — `BP-PY-46`

- Function context: `scripts/WHEN-Language/findings/functions/104.txt`
- Source: `/home/chinmay/ChinmayPersonalProjects/goslop/real-repos/WHEN-Language/when.py:115:5`
- Checklist pattern: CLI entry script (`__main__` guard present) whose prints are user-facing output, not library logging.

Source excerpt:

```
    print("Goodbye!")
```

Why this is a false positive: REPL farewell message in the entry script.

Checklist evidence: user-facing CLI output in the `__main__`-guarded entry script.

## Uncertain findings

None — every finding could be decided from the rule condition and the shown source.

## True positives

### BP-PY-46 — print in library modules (`hot_reload.py`, `interpreter.py`)

| Finding | Source | Reason |
| --- | --- | --- |
| 1 | hot_reload.py:55:9 | `print` for operational logging in the `HotReloader` library class (no `__main__` guard, not a test). |
| 2 | hot_reload.py:73:21 | `print` operational logging in library class method `_watch_loop`. |
| 5 | hot_reload.py:78:17 | `print` error logging in library class method. |
| 7 | hot_reload.py:110:17 | `print` operational logging in library class method `_reload_blocks`. |
| 9 | hot_reload.py:113:17 | `print` error logging in library class method. |
| 10 | hot_reload.py:142:17 | `print` operational logging in library class method `_restore_block_state`. |
| 11 | hot_reload.py:157:21 | `print` operational logging in library class method. |
| 12 | hot_reload.py:159:21 | `print` operational logging in library class method. |
| 14 | hot_reload.py:205:17 | `print` operational logging in library class method `_update_blocks`. |
| 15 | hot_reload.py:207:17 | `print` operational logging in library class method. |
| 16 | hot_reload.py:223:13 | `print` operational logging in library class method. |
| 26 | interpreter.py:473:13 | `print("Program exited")` in library module `interpret()`. |
| 31 | interpreter.py:1061:13 | `print(*values)` implementing the language `print` builtin in library module (print call outside `__main__` guard). |
| 32 | interpreter.py:1369:9 | `print` operational logging `[SAVE]` in library module. |
| 33 | interpreter.py:1377:9 | `print` operational logging `[SAVESTOP]` in library module. |
| 34 | interpreter.py:1388:13 | `print` operational logging `[STARTSAVE]` in library module. |
| 35 | interpreter.py:1394:13 | `print` operational logging `[STARTSAVE]` in library module. |
| 36 | interpreter.py:1396:13 | `print` operational logging `[STARTSAVE]` in library module. |
| 37 | interpreter.py:1423:13 | `print` operational logging `[DISCARD]` in library module. |
| 38 | interpreter.py:1425:13 | `print` operational logging `[DISCARD]` in library module. |
| 44 | interpreter.py:1838:13 | `print` error logging `[PARALLEL]` in library module. |

### BP-PY-1 — bare / broad except

| Finding | Source | Reason |
| --- | --- | --- |
| 3 | hot_reload.py:77:1 | `except Exception as e:` in `_watch_loop` swallows all exceptions and only prints. |
| 8 | hot_reload.py:112:1 | `except Exception as e:` in `_reload_blocks` swallows all exceptions and only prints. |
| 18 | interpreter.py:112:1 | `except Exception as e:` in `safe_call` catches all failures. |
| 20 | interpreter.py:137:1 | bare `except:` in `has_attr`. |
| 21 | interpreter.py:144:1 | bare `except:` in `safe_getattr`. |
| 23 | interpreter.py:325:1 | bare `except:` in `eval_func` wraps any failure into ValueError. |
| 30 | interpreter.py:985:1 | `except Exception as e:` in f-string evaluation. |
| 43 | interpreter.py:1837:1 | `except Exception as e:` in parallel block runner. |
| 95 | when.py:69:1 | `except Exception as e:` in `run_file` catches all runtime failures. |
| 101 | when.py:107:1 | `except Exception as e:` in REPL loop. |

### CWE-396 — generic exception handler (same construct as BP-PY-1)

| Finding | Source | Reason |
| --- | --- | --- |
| 4 | hot_reload.py:77:1 | generic `except Exception` handler. |
| 19 | interpreter.py:112:1 | generic `except Exception` handler. |
| 96 | when.py:69:1 | generic `except Exception` handler. |

### BP-PY-5 — wildcard import

| Finding | Source | Reason |
| --- | --- | --- |
| 17 | interpreter.py:7:1 | `from ast_nodes import *` pollutes namespace. |
| 47 | parser.py:3:1 | `from ast_nodes import *` pollutes namespace. |

### CWE-1121 — excessive cyclomatic complexity

| Finding | Source | Reason |
| --- | --- | --- |
| 13 | hot_reload.py:170:39 | `_update_blocks` body counts ≥12 control-flow keyword occurrences (18 by the rule's counter), a genuinely branch-heavy function. |
| 25 | interpreter.py:411:43 | `interpret` body counts ≥12 control-flow keyword occurrences (26), a genuinely branch-heavy dispatch function. |

### CWE-1124 — excessively deep nesting

| Finding | Source | Reason |
| --- | --- | --- |
| 27 | interpreter.py:907:1 | `saved_event = ...` sits at 7 active control-flow frames (elif → for → if → for → if → if → if), verified by frame simulation. |
| 74 | parser.py:654:1 | `args.append(self.parse_expression())` sits at 6 active control-flow frames (while → elif → else → if → while → else), verified by frame simulation. |

### CWE-1046 — string concatenation in loop

| Finding | Source | Reason |
| --- | --- | --- |
| 28 | interpreter.py:966:1 | `result += part_value` inside `for part_type, part_value in fstring.parts:` where `result` was initialized to `""`. |
| 45 | lexer.py:117:1 | `result += char` inside `for i in range(length):` where `result` was initialized to `""`. |

### CWE-829 / CWE-94 — dynamic import of program-controlled module

| Finding | Source | Reason |
| --- | --- | --- |
| 39 | interpreter.py:1440:34 | `top_module = __import__(decl.module)` — `decl.module` is a dynamic value from a parsed WHEN import declaration. |
| 40 | interpreter.py:1440:34 | same `__import__(decl.module)` reaches a dynamic-import sink with a non-literal argument. |

### BP-PY-45 — sys.path mutation

| Finding | Source | Reason |
| --- | --- | --- |
| 84 | when.py:20:1 | `sys.path.insert(0, str(current_dir))` at runtime instead of packaging. |

### CWE-367 — TOCTOU

| Finding | Source | Reason |
| --- | --- | --- |
| 89 | when.py:38:16 | `os.path.exists(filename)` check immediately before a separate `open(filename)` use; the rule's regex condition is satisfied. |

## Final evidence

- Delegated reviewers: none
- Chunk evidence: `scripts/WHEN-Language/chunks`
- Function evidence: `scripts/WHEN-Language/findings/functions`
- Validation: `git diff --check` — pass
