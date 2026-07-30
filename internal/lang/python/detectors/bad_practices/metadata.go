package badpractices

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

// Hand-written metadata for implemented BP-PY rules (source of truth for detectors).
// Full catalogue lives at ruleset/python/bad-practices.json; do not use Go metadata_gen.go.
var metaByID = map[string]*rules.RuleMetadata{
	"BP-PY-1": {
		ID: "BP-PY-1", Title: "Bare Except Clause",
		Description: "A bare `except:` or `except Exception` without handling or re-raise swallows failures and hides bugs.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Catch specific exception types and handle or re-raise.",
	},
	"BP-PY-2": {
		ID: "BP-PY-2", Title: "Except Pass",
		Description: "An exception handler body is only `pass`, discarding the failure silently.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Log, re-raise, or handle the exception instead of bare pass.",
	},
	"BP-PY-4": {
		ID: "BP-PY-4", Title: "Mutable Default Argument",
		Description: "A function default uses a mutable value (`[]`, `{}`, `set()`) shared across calls.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Use None as default and assign a new list/dict/set inside the body.",
	},
	"BP-PY-6": {
		ID: "BP-PY-6", Title: "assert Used For Runtime Validation",
		Description: "`assert` is used for input or security checks that disappear with `python -O`.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Prefer if + raise ValueError/HTTPException for runtime validation.",
	},
	"BP-PY-7": {
		ID: "BP-PY-7", Title: "open Without Context Manager",
		Description: "A file is opened without a `with` statement, risking resource leaks.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Use `with open(...) as f:` so the file is closed reliably.",
	},
	"BP-PY-8": {
		ID: "BP-PY-8", Title: "subprocess With shell=True",
		Description: "subprocess is invoked with shell=True, enabling shell injection when args are dynamic.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Pass an argv list and keep shell=False.",
	},
	"BP-PY-9": {
		ID: "BP-PY-9", Title: "os.system Or os.popen",
		Description: "os.system / os.popen run a shell command and are hard to secure.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Prefer subprocess with a list argv and shell=False.",
	},
	"BP-PY-10": {
		ID: "BP-PY-10", Title: "pickle Loads Untrusted Data",
		Description: "pickle.load/loads can execute arbitrary code on malicious payloads.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Prefer json or other safe formats for untrusted data; avoid pickle across trust boundaries.",
	},
	"BP-PY-11": {
		ID: "BP-PY-11", Title: "yaml.load Without SafeLoader",
		Description: "yaml.load without SafeLoader can construct arbitrary Python objects.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Use yaml.safe_load or yaml.load(..., Loader=yaml.SafeLoader).",
	},
	"BP-PY-12": {
		ID: "BP-PY-12", Title: "eval Or exec On Dynamic Input",
		Description: "eval/exec on dynamic input enables arbitrary code execution.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Avoid eval/exec; use ast.literal_eval for literals or structured parsers.",
	},
	"BP-PY-13": {
		ID: "BP-PY-13", Title: "Hardcoded Secret In Source",
		Description: "A secret-like name is assigned a non-empty string literal in source.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Load secrets from environment or a secret manager; never commit real credentials.",
	},
	"BP-PY-16": {
		ID: "BP-PY-16", Title: "Flask DEBUG True In Production Code",
		Description: "Flask DEBUG is enabled in non-test code, exposing an interactive debugger.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Keep debug off in production; gate via environment-specific config.",
	},
	"BP-PY-17": {
		ID: "BP-PY-17", Title: "Flask SECRET_KEY Hardcoded",
		Description: "Flask SECRET_KEY / secret_key is a string literal in source.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Load SECRET_KEY from the environment or a secret manager.",
	},
	"BP-PY-21": {
		ID: "BP-PY-21", Title: "Django DEBUG True In Settings",
		Description: "Django DEBUG = True in settings modules risks information disclosure.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Set DEBUG from environment; keep False in production settings.",
	},
	"BP-PY-29": {
		ID: "BP-PY-29", Title: "FastAPI Depends On Mutable Global",
		Description: "FastAPI/Starlette dependencies or routes mutate module-level global state.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Prefer request-scoped dependencies and a proper store instead of mutating globals.",
	},
	"BP-PY-30": {
		ID: "BP-PY-30", Title: "FastAPI Blocking I/O In Async Route",
		Description: "An `async def` route calls blocking I/O (time.sleep, requests, sync ORM) on the event loop.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Use await asyncio.sleep, httpx.AsyncClient, run_in_executor, or a sync def route.",
	},
	"BP-PY-31": {
		ID: "BP-PY-31", Title: "FastAPI response_model Disabled Unsafely",
		Description: "Response returns ORM/internal objects without `response_model`, leaking fields.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Declare response_model with a Pydantic schema or return an explicit DTO.",
	},
	"BP-PY-32": {
		ID: "BP-PY-32", Title: "Starlette FileResponse User Path",
		Description: "`FileResponse` / static file helpers use a path from user input without confinement.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Resolve the path and enforce a trusted prefix before FileResponse.",
	},
	"BP-PY-33": {
		ID: "BP-PY-33", Title: "Jinja2 autoescape Disabled",
		Description: "Jinja2 Environment is created with `autoescape=False` or missing autoescape for HTML.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Enable autoescape=True or select_autoescape for HTML/XML templates.",
	},
	"BP-PY-34": {
		ID: "BP-PY-34", Title: "Jinja2 Markup Or safe Filter On Variables",
		Description: "Template or Python marks dynamic values as safe HTML.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Only mark trusted, sanitized literals as safe; escape untrusted input.",
	},
	"BP-PY-35": {
		ID: "BP-PY-35", Title: "SQLAlchemy text With F-String",
		Description: "`sqlalchemy.text` / execute builds SQL with f-strings or format instead of bound params.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Use text(\"... WHERE id = :id\") with bindparams instead of f-strings.",
	},
	"BP-PY-36": {
		ID: "BP-PY-36", Title: "SQLAlchemy Session Not Closed",
		Description: "A Session/sessionmaker session is created without context manager or close in finally.",
		Severity:    rules.SeverityMedium, Pack: rules.PackBadPractice,
		Fix: "Use `with Session() as session:` or call session.close() on all exit paths.",
	},
	"BP-PY-37": {
		ID: "BP-PY-37", Title: "DB-API Cursor Execute With Percent Format",
		Description: "DB-API `cursor.execute` uses Python `%` formatting on the SQL string.",
		Severity:    rules.SeverityHigh, Pack: rules.PackBadPractice,
		Fix: "Pass SQL with placeholders and a params sequence: execute(sql, params).",
	},
}

// MetadataForID returns catalogue metadata for a BP-PY rule id.
func MetadataForID(ruleID string) *rules.RuleMetadata {
	return metaByID[ruleID]
}

// CatalogueSize returns the number of implemented BP-PY metadata entries.
func CatalogueSize() int { return len(metaByID) }

// --- Optional full-catalogue parse (tests / validation) ---

type catalogueEntry struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Category    string `json:"category"`
}

var (
	fullCatalogueOnce sync.Once
	fullCatalogue     map[string]catalogueEntry
	fullCatalogueErr  error
)

// loadFullCatalogue parses ruleset/python/bad-practices.json when reachable from CWD.
func loadFullCatalogue() (map[string]catalogueEntry, error) {
	fullCatalogueOnce.Do(func() {
		path, err := findPythonBPCatalogue()
		if err != nil {
			fullCatalogueErr = err
			return
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			fullCatalogueErr = err
			return
		}
		var m map[string]catalogueEntry
		if err := json.Unmarshal(raw, &m); err != nil {
			fullCatalogueErr = err
			return
		}
		fullCatalogue = m
	})
	return fullCatalogue, fullCatalogueErr
}

func findPythonBPCatalogue() (string, error) {
	candidates := []string{
		"ruleset/python/bad-practices.json",
		filepath.Join("..", "..", "..", "..", "..", "ruleset", "python", "bad-practices.json"),
		filepath.Join("..", "..", "..", "..", "..", "..", "ruleset", "python", "bad-practices.json"),
	}
	// Also walk up from cwd.
	cwd, _ := os.Getwd()
	for dir := cwd; dir != "" && dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		candidates = append(candidates, filepath.Join(dir, "ruleset", "python", "bad-practices.json"))
		if parent := filepath.Dir(dir); parent == dir {
			break
		}
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c, nil
		}
	}
	return "", os.ErrNotExist
}

func severityFromCatalogue(s string) rules.Severity {
	sev, err := rules.ParseSeverity(s)
	if err != nil {
		return rules.SeverityLow
	}
	return sev
}

// ValidateImplementedMetadata checks implemented meta against ruleset/python/bad-practices.json.
func ValidateImplementedMetadata() error {
	full, err := loadFullCatalogue()
	if err != nil {
		return err
	}
	for id, meta := range metaByID {
		entry, ok := full[id]
		if !ok {
			return errMissingCatalogue(id)
		}
		wantSev := severityFromCatalogue(entry.Severity)
		if meta.Severity != wantSev {
			return errSeverityMismatch(id, entry.Severity, meta.Severity.String())
		}
		if !strings.HasPrefix(id, "BP-PY-") {
			return errBadID(id)
		}
		if entry.Name != "" && meta.Title != entry.Name {
			return simpleError("title mismatch for " + id + ": catalogue=" + entry.Name + " meta=" + meta.Title)
		}
	}
	return nil
}

// FullCatalogueSize returns the number of keys in the on-disk Python BP catalogue.
func FullCatalogueSize() (int, error) {
	full, err := loadFullCatalogue()
	if err != nil {
		return 0, err
	}
	return len(full), nil
}

type simpleError string

func (e simpleError) Error() string { return string(e) }

func errMissingCatalogue(id string) error {
	return simpleError("catalogue missing " + id)
}
func errSeverityMismatch(id, want, got string) error {
	return simpleError("severity mismatch for " + id + ": catalogue=" + want + " meta=" + got)
}
func errBadID(id string) error {
	return simpleError("registered id must be BP-PY-*: " + id)
}
