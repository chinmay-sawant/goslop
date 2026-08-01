package perf

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/core"
	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestPhase2Rules(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		rule string
		src  string
		want bool
	}{
		{"materialized sort", "PERF-PY-1", "async def percentile(session):\n    rows = (await session.execute(stmt)).scalars().all()\n    return sorted(rows)[0]\n", true},
		{"database aggregate", "PERF-PY-1", "async def percentile(session):\n    return await session.scalar(select(func.percentile_cont(0.95)))\n", false},
		{"batch create", "PERF-PY-3", "def import_rows(rows):\n    for row in rows:\n        Event.objects.create(name=row['name'])\n", true},
		{"dependent create", "PERF-PY-3", "def import_rows(rows):\n    for row in rows:\n        event = Event.objects.create(name=row['name'])\n        EventLog.objects.create(event=event)\n", false},
		{"counter save", "PERF-PY-4", "def reserve(stock, amount):\n    stock.reserved_quantity += amount\n    stock.save()\n", true},
		{"atomic counter", "PERF-PY-4", "def reserve(stock, amount):\n    Stock.objects.filter(pk=stock.pk).update(reserved_quantity=models.F('reserved_quantity') + amount)\n", false},
		{"parse then dump", "PERF-PY-9", "def ingest(request, session):\n    payload = request.get_json()\n    session.add(Event(body=json.dumps(payload)))\n", true},
		{"dump for log unrelated create", "PERF-PY-9", "def ingest(request, session):\n    payload = request.get_json()\n    logger.info(json.dumps(payload))\n    session.add(Event(name='x'))\n", false},
		{"raw payload", "PERF-PY-9", "def ingest(request, session):\n    session.add(Event(body=request.data))\n", false},
		{"sleep after work", "PERF-PY-10", "def worker():\n    while True:\n        processed = process_batch()\n        if processed:\n            logger.info('processed')\n        time.sleep(5)\n", true},
		{"continue after work", "PERF-PY-10", "def worker():\n    while True:\n        processed = process_batch()\n        if processed:\n            continue\n        time.sleep(5)\n", false},
		{"per row mutation", "PERF-PY-11", "def redrive():\n    for job in Job.objects.filter(state='dead'):\n        job.state = 'queued'\n        job.save()\n", true},
		{"set update", "PERF-PY-11", "def redrive():\n    Job.objects.filter(state='dead').update(state='queued')\n", false},
		{"sqlalchemy per row mutation", "PERF-PY-11", "def redrive(session):\n    for job in session.query(Job).filter(Job.state == 'dead'):\n        job.state = 'queued'\n    session.commit()\n", true},
		{"non ORM save", "PERF-PY-11", "def persist(records):\n    for record in records:\n        record.state = 'queued'\n        record.save()\n", false},
		{"unbounded json", "PERF-PY-12", "@app.post('/ingest')\ndef ingest(request):\n    payload = request.get_json()\n    return payload\n", true},
		{"limited json", "PERF-PY-12", "@app.post('/ingest')\ndef ingest(request):\n    if request.content_length > MAX_BODY_BYTES:\n        raise ValueError('too large')\n    payload = request.get_json()\n    return payload\n", false},
		{"lookup then create", "PERF-PY-14", "def create(key):\n    existing = Event.objects.filter(idempotency_key=key).first()\n    if existing:\n        return existing\n    return Event.objects.create(idempotency_key=key)\n", true},
		{"helper lookup then insert", "PERF-PY-14", "async def create(request, session):\n    existing = await self._find_batch_by_key(request.idempotency_key)\n    if existing:\n        return existing\n    session.add(Event(idempotency_key=request.idempotency_key))\n", true},
		{"upsert", "PERF-PY-14", "def create(key):\n    return Event.objects.get_or_create(idempotency_key=key)\n", false},
		{"top-level lookup after function", "PERF-PY-14", "def unrelated():\n    return 1\n\nexisting = Event.objects.filter(idempotency_key=key).first()\nEvent.objects.create(idempotency_key=key)\n", false},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := core.DefaultScanContext()
			ctx.Only = []string{tc.rule}
			unit := core.NewParsedUnit(core.LanguagePython, "app.py", tc.src)
			var out []rules.Finding
			NewPythonPerfScan().Run(ctx, unit, &out)
			got := false
			for _, finding := range out {
				if finding.RuleID == tc.rule {
					got = true
				}
			}
			if got != tc.want {
				t.Fatalf("%s got finding=%v want=%v: %#v", tc.rule, got, tc.want, out)
			}
		})
	}
}
