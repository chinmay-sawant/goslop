package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY29MutableGlobal(t *testing.T) {
	t.Parallel()
	vuln := `from fastapi import FastAPI, Depends
app = FastAPI()
STORE = {}

@app.get("/x")
def read_x(k: str):
    global STORE
    STORE[k] = 1
    return STORE
`
	vulnMut := `from fastapi import FastAPI
app = FastAPI()
CACHE = {}

@app.post("/c")
async def put_c(k: str, v: str):
    CACHE[k] = v
    return {"ok": True}
`
	safe := `from fastapi import FastAPI
app = FastAPI()

@app.get("/x")
def read_x():
    local = {}
    local["a"] = 1
    return local
`
	assertRule(t, "BP-PY-29", "app.py", vuln, true)
	assertRule(t, "BP-PY-29", "app.py", vulnMut, true)
	assertRule(t, "BP-PY-29", "app.py", safe, false)

	// Flask @app.route must not arm FastAPI-only rules.
	flaskish := `from flask import Flask
app = Flask(__name__)
STORE = {}

@app.route("/x")
def read_x(k):
    global STORE
    STORE[k] = 1
    return STORE
`
	assertRule(t, "BP-PY-29", "routes.py", flaskish, false)
}

func TestBPPY30BlockingIOAsyncRoute(t *testing.T) {
	t.Parallel()
	vulnSleep := `from fastapi import FastAPI
import time
app = FastAPI()

@app.get("/slow")
async def slow():
    time.sleep(1)
    return {"ok": True}
`
	vulnReq := `from fastapi import FastAPI
import requests
app = FastAPI()

@app.get("/proxy")
async def proxy():
    r = requests.get("https://example.com")
    return {"status": r.status_code}
`
	safeAsyncio := `from fastapi import FastAPI
import asyncio
app = FastAPI()

@app.get("/slow")
async def slow():
    await asyncio.sleep(1)
    return {"ok": True}
`
	safeSync := `from fastapi import FastAPI
import time
app = FastAPI()

@app.get("/slow")
def slow():
    time.sleep(1)
    return {"ok": True}
`
	assertRule(t, "BP-PY-30", "app.py", vulnSleep, true)
	assertRule(t, "BP-PY-30", "app.py", vulnReq, true)
	assertRule(t, "BP-PY-30", "app.py", safeAsyncio, false)
	assertRule(t, "BP-PY-30", "app.py", safeSync, false)
	findings := runBP(t, nil, vulnSleep, "app.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-30" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-30 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY31ResponseModel(t *testing.T) {
	t.Parallel()
	vuln := `from fastapi import FastAPI
from sqlalchemy.orm import Session
app = FastAPI()

@app.get("/u/{id}")
def get_user(id: int, session: Session):
    return session.query(User).get(id)
`
	safe := `from fastapi import FastAPI
from sqlalchemy.orm import Session
app = FastAPI()

@app.get("/u/{id}", response_model=UserOut)
def get_user(id: int, session: Session):
    return session.query(User).get(id)
`
	safeDict := `from fastapi import FastAPI
app = FastAPI()

@app.get("/u/{id}")
def get_user(id: int):
    return {"id": id, "name": "x"}
`
	assertRule(t, "BP-PY-31", "app.py", vuln, true)
	assertRule(t, "BP-PY-31", "app.py", safe, false)
	assertRule(t, "BP-PY-31", "app.py", safeDict, false)
}

func TestBPPY32FileResponseUserPath(t *testing.T) {
	t.Parallel()
	vulnParam := `from fastapi import FastAPI
from starlette.responses import FileResponse
app = FastAPI()

@app.get("/files/{name}")
def get_file(name: str):
    return FileResponse(name)
`
	vulnF := `from fastapi import FastAPI
from starlette.responses import FileResponse
app = FastAPI()

@app.get("/files/{name}")
def get_file(name: str):
    return FileResponse(f"/data/{name}")
`
	safe := `from fastapi import FastAPI
from starlette.responses import FileResponse
app = FastAPI()

@app.get("/report")
def report():
    return FileResponse("/safe/fixed.pdf")
`
	assertRule(t, "BP-PY-32", "app.py", vulnParam, true)
	assertRule(t, "BP-PY-32", "app.py", vulnF, true)
	assertRule(t, "BP-PY-32", "app.py", safe, false)
	findings := runBP(t, nil, vulnParam, "app.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-32" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-32 severity = %v, want high", f.Severity)
		}
	}
}

func TestFastAPIRulesRegistered(t *testing.T) {
	t.Parallel()
	for _, id := range []string{"BP-PY-29", "BP-PY-30", "BP-PY-31", "BP-PY-32"} {
		assertRule(t, id, "empty.py", "# no fastapi\n", false)
	}
}
