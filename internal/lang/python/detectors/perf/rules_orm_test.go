package perf

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

// The detector Run boundary is the public seam for these source-only rules.
func TestORMRules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		rule string
		vuln string
		safe string
	}{
		{"lookup in loop", "PERF-PY-2", `def reserve(items):
    for item in items:
        sku = Sku.objects.get(code=item["code"])
`, `def reserve(items):
    skus = Sku.objects.in_bulk([item["code"] for item in items], field_name="code")
    for item in items:
		 sku = skus[item["code"]]
`},
		{"work claim no lock", "PERF-PY-6", `def claim():
    job = Job.objects.filter(status="pending").first()
    job.status = "running"
    job.save()
`, `def claim():
    job = Job.objects.select_for_update(skip_locked=True).filter(status="pending").first()
    job.status = "running"
    job.save()
`},
		{"lazy sqlalchemy relation", "PERF-PY-8", `def deliver(session):
    items = session.query(DeliveryOutbox).all()
    for item in items:
        send(item.partner_endpoint.url)
`, `def deliver(session):
    items = session.query(DeliveryOutbox).options(joinedload(DeliveryOutbox.partner_endpoint)).all()
    for item in items:
        send(item.partner_endpoint.url)
`},
		{"full hydration projection", "PERF-PY-13", `def names():
    for user in User.objects.filter(active=True):
        emit(user.email)
`, `def names():
    for email in User.objects.filter(active=True).values_list("email", flat=True):
        emit(email)
`},
		{"composite filter visible no index", "PERF-PY-15", `class Metric(models.Model):
    tenant_id = models.IntegerField()
    timestamp = models.DateTimeField()

def read(cutoff):
    return Metric.objects.filter(tenant_id=1, timestamp__gte=cutoff)
`, `class Metric(models.Model):
    tenant_id = models.IntegerField()
    timestamp = models.DateTimeField()
    class Meta:
        indexes = [models.Index(fields=["tenant_id", "timestamp"])]

def read(cutoff):
    return Metric.objects.filter(tenant_id=1, timestamp__gte=cutoff)
`},
		{"retention predicate visible no index", "PERF-PY-16", `class Metric(models.Model):
    created_at = models.DateTimeField()

def purge(cutoff):
    return Metric.objects.filter(created_at__lt=cutoff).delete()
`, `class Metric(models.Model):
    created_at = models.DateTimeField(db_index=True)

def purge(cutoff):
    return Metric.objects.filter(created_at__lt=cutoff).delete()
`},
		{"unbounded lock sweep", "PERF-PY-19", `@transaction.atomic
def release_expired(cutoff):
    for row in Reservation.objects.filter(expires_at__lt=cutoff).select_for_update(skip_locked=True):
        row.release()
`, `@transaction.atomic
def release_expired(cutoff):
    for row in Reservation.objects.filter(expires_at__lt=cutoff).select_for_update(skip_locked=True)[:100]:
        row.release()
`},
		{"sqlalchemy unbounded lock sweep", "PERF-PY-19", `def release(session):
    with session.begin():
        for row in session.query(Reservation).with_for_update():
            row.release()
`, `def release(session):
    with session.begin():
        for row in session.query(Reservation).with_for_update().limit(100):
            row.release()
`},
		{"sort visible no index", "PERF-PY-20", `class Stock(models.Model):
    warehouse_id = models.IntegerField()
    quantity = models.IntegerField()

def available():
    return Stock.objects.filter(warehouse_id=1).order_by("-quantity")
`, `class Stock(models.Model):
    warehouse_id = models.IntegerField()
    quantity = models.IntegerField()
    class Meta:
        indexes = [models.Index(fields=["warehouse_id", "-quantity"])]

def available():
    return Stock.objects.filter(warehouse_id=1).order_by("-quantity")
`},
		{"unbounded maintenance delete", "PERF-PY-21", `def purge(cutoff):
    AuditEvent.objects.filter(created_at__lt=cutoff).delete()
`, `def purge(cutoff):
    while True:
        deleted, _ = AuditEvent.objects.filter(created_at__lt=cutoff)[:500].delete()
        if not deleted:
            break
`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertPERFRule(t, tt.rule, tt.vuln, true)
			assertPERFRule(t, tt.rule, tt.safe, false)
		})
	}
}

func TestIndexRulesSuppressWhenModelDeclarationIsUnavailable(t *testing.T) {
	t.Parallel()
	assertPERFRule(t, "PERF-PY-15", `def read(cutoff):
    return Metric.objects.filter(tenant_id=1, timestamp__gte=cutoff)
`, false)
	assertPERFRule(t, "PERF-PY-16", `def purge(cutoff):
    return Metric.objects.filter(created_at__lt=cutoff).delete()
`, false)
	assertPERFRule(t, "PERF-PY-20", `def available():
    return Stock.objects.filter(warehouse_id=1).order_by("-quantity")
`, false)
	assertPERFRule(t, "PERF-PY-15", `# external migration manages this index
class Metric(models.Model):
    tenant_id = models.IntegerField()
    timestamp = models.DateTimeField()

def read(cutoff):
    return Metric.objects.filter(tenant_id=1, timestamp__gte=cutoff)
`, false)
}

func TestLookupInSmallExplicitLoopIsSuppressed(t *testing.T) {
	t.Parallel()
	assertPERFRule(t, "PERF-PY-2", `def reserve():
    for code in ("A", "B"):
        Sku.objects.get(code=code)
`, false)
}

func assertPERFRule(t *testing.T, id, source string, want bool) {
	t.Helper()
	ctx := core.DefaultScanContext()
	ctx.Only = []string{id}
	unit := core.NewParsedUnit(core.LanguagePython, "service.py", source)
	var findings []rules.Finding
	detector := NewPythonPerfScan()
	detector.Run(ctx, unit, &findings)
	got := false
	for _, finding := range findings {
		if finding.RuleID == id {
			got = true
			break
		}
	}
	if got != want {
		t.Fatalf("%s want finding=%t, got %v\nsource:\n%s", id, want, findings, source)
	}
}
