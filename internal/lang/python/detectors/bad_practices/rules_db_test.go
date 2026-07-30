package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY35SQLAlchemyTextFString(t *testing.T) {
	t.Parallel()
	vulnF := `from sqlalchemy import text
def q(uid):
    return text(f"SELECT * FROM users WHERE id = {uid}")
`
	vulnFmt := `from sqlalchemy import text
def q(uid):
    return text("SELECT * FROM users WHERE id = {}".format(uid))
`
	safe := `from sqlalchemy import text
def q(uid):
    return text("SELECT * FROM users WHERE id = :id")
`
	assertRule(t, "BP-PY-35", "db.py", vulnF, true)
	assertRule(t, "BP-PY-35", "db.py", vulnFmt, true)
	assertRule(t, "BP-PY-35", "db.py", safe, false)
	findings := runBP(t, nil, vulnF, "db.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-35" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-35 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY36SessionNotClosed(t *testing.T) {
	t.Parallel()
	vuln := `from sqlalchemy.orm import Session
SessionLocal = sessionmaker()
def work():
    session = SessionLocal()
    session.query(User).all()
    return 1
`
	// Need sessionmaker name present — adjust
	vuln2 := `def work():
    session = SessionLocal()
    rows = session.query(User).all()
    return rows
`
	safeWith := `def work():
    with SessionLocal() as session:
        return session.query(User).all()
`
	safeClose := `def work():
    session = SessionLocal()
    try:
        return session.query(User).all()
    finally:
        session.close()
`
	assertRule(t, "BP-PY-36", "db.py", vuln2, true)
	assertRule(t, "BP-PY-36", "db.py", safeWith, false)
	assertRule(t, "BP-PY-36", "db.py", safeClose, false)
	_ = vuln
}

func TestBPPY37ExecutePercentFormat(t *testing.T) {
	t.Parallel()
	vulnF := `def q(cursor, uid):
    cursor.execute(f"SELECT * FROM t WHERE id = {uid}")
`
	vulnPct := `def q(cursor, uid):
    cursor.execute("SELECT * FROM t WHERE id = %s" % (uid,))
`
	safe := `def q(cursor, uid):
    cursor.execute("SELECT * FROM t WHERE id = %s", (uid,))
`
	assertRule(t, "BP-PY-37", "db.py", vulnF, true)
	assertRule(t, "BP-PY-37", "db.py", vulnPct, true)
	assertRule(t, "BP-PY-37", "db.py", safe, false)
	findings := runBP(t, nil, vulnF, "db.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-37" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-37 severity = %v, want high", f.Severity)
		}
	}
}

func TestDBAndTemplateRulesRegistered(t *testing.T) {
	t.Parallel()
	// Registration smoke: empty-ish sources should not panic / false-positive heavily.
	assertRule(t, "BP-PY-33", "x.py", "x = 1\n", false)
	assertRule(t, "BP-PY-35", "x.py", "x = 1\n", false)
	assertRule(t, "BP-PY-37", "x.py", "x = 1\n", false)
}
