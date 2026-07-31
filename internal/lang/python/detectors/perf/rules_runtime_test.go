package perf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/fixture"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestRuntimeRulesHitAndMiss(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		rule string
		path string
		hit  string
		miss string
	}{
		{
			rule: "PERF-PY-5", path: "workers/delivery.py",
			hit:  "func", // replaced below to keep each Python example readable
			miss: "func",
		},
		{
			rule: "PERF-PY-7", path: "app/middleware.py",
			hit: `from starlette.middleware.base import BaseHTTPMiddleware
class AuditMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request, call_next):
        return await call_next(request)
app.add_middleware(AuditMiddleware)
`,
			miss: `from starlette.types import ASGIApp
class AuditMiddleware:
    async def __call__(self, scope, receive, send):
        await self.app(scope, receive, send)
app.add_middleware(AuditMiddleware)
`,
		},
		{
			rule: "PERF-PY-17", path: "config/database.py",
			hit: `import os
ENVIRONMENT = "production"
DATABASE_URL = "postgresql://db/service"
engine = create_engine(DATABASE_URL)
`,
			miss: `import os
ENVIRONMENT = "production"
DATABASE_URL = "postgresql://db/service"
engine = create_engine(DATABASE_URL, pool_pre_ping=True, pool_recycle=1800, pool_timeout=10)
`,
		},
		{
			rule: "PERF-PY-18", path: "app/normalization.py",
			hit: `import re
def normalize_payload(payload):
    payload = re.sub(r"\\s+", " ", payload)
    payload = re.sub(r"[^a-z ]", "", payload)
    return payload
`,
			miss: `import re
def normalize_payload(payload):
    compact = re.sub(r"\\s+", " ", payload)
    return re.sub(r"[^a-z ]", "", compact)
`,
		},
		{
			rule: "PERF-PY-22", path: "config/database.py",
			hit: `ENVIRONMENT = "production"
DATABASES = {"default": {"ENGINE": "django.db.backends.sqlite3", "NAME": "service.db"}}
CELERY_WORKER_CONCURRENCY = 4
`,
			miss: `ENVIRONMENT = "production"
DATABASES = {"default": {"ENGINE": "django.db.backends.postgresql", "NAME": "service"}}
CELERY_WORKER_CONCURRENCY = 4
`,
		},
	} {
		tc := tc
		t.Run(tc.rule, func(t *testing.T) {
			t.Parallel()
			if tc.rule == "PERF-PY-5" {
				tc.hit = `def worker():
    batch = claim_pending_batch()
    for job in batch:
        deliver_webhook(job)
`
				tc.miss = `import asyncio
async def worker():
    batch = claim_pending_batch()
    await asyncio.gather(*(deliver_webhook(job) for job in batch))
`
			}
			assertPerfRule(t, tc.rule, tc.path, tc.hit, true)
			assertPerfRule(t, tc.rule, tc.path, tc.miss, false)
		})
	}
}

func TestRuntimeConfigRulesSuppressTestsAndLocalConfigurations(t *testing.T) {
	t.Parallel()
	assertPerfRule(t, "PERF-PY-17", "tests/test_database.py", `ENVIRONMENT = "production"
engine = create_engine("postgresql://db/service")
`, false)
	assertPerfRule(t, "PERF-PY-17", "config/local_database.py", `ENVIRONMENT = "production"
engine = create_engine("sqlite:///local.db")
`, false)
	assertPerfRule(t, "PERF-PY-22", "tests/test_settings.py", `ENVIRONMENT = "production"
DATABASES = {"default": {"ENGINE": "django.db.backends.sqlite3"}}
CELERY_WORKER_CONCURRENCY = 4
`, false)
	assertPerfRule(t, "PERF-PY-22", "config/local_settings.py", `ENVIRONMENT = "production"
DATABASES = {"default": {"ENGINE": "django.db.backends.sqlite3"}}
CELERY_WORKER_CONCURRENCY = 4
`, false)
}

func TestRuntimeRuleFixturePairs(t *testing.T) {
	t.Parallel()
	for _, rule := range []string{"PERF-PY-5", "PERF-PY-7", "PERF-PY-17", "PERF-PY-18", "PERF-PY-22"} {
		rule := rule
		t.Run(rule, func(t *testing.T) {
			t.Parallel()
			for _, vulnerable := range []bool{true, false} {
				suffix := "safe"
				if vulnerable {
					suffix = "vulnerable"
				}
				fixturePath := filepath.Join("..", "..", "..", "..", "..", "tests", "fixtures", "python", "perf", rule+"-"+suffix+".txt")
				contents, err := os.ReadFile(fixturePath)
				if err != nil {
					t.Fatalf("read %s: %v", fixturePath, err)
				}
				fx, err := fixture.ParseFixture(string(contents), filepath.Base(fixturePath))
				if err != nil {
					t.Fatalf("parse %s: %v", fixturePath, err)
				}
				assertPerfRule(t, rule, fx.Filename, fx.Source, vulnerable)
			}
		})
	}
}

func assertPerfRule(t *testing.T, rule, path, source string, want bool) {
	t.Helper()
	ctx := core.DefaultScanContext()
	ctx.Only = []string{rule}
	unit := core.NewParsedUnit(core.LanguagePython, path, source)
	var findings []rules.Finding
	NewPythonPerfScan().Run(ctx, unit, &findings)
	got := false
	for _, finding := range findings {
		if finding.RuleID == rule {
			got = true
			break
		}
	}
	if got != want {
		t.Fatalf("%s on %s: got fire=%v, want %v; findings=%#v\nsource:\n%s", rule, path, got, want, findings, source)
	}
}
