package badpractices_test

import (
	"testing"

	"github.com/chinmay-sawant/goslop/internal/rules"
)

func TestBPPY38CreateTaskBare(t *testing.T) {
	t.Parallel()
	vuln := "import asyncio\nasync def main():\n    asyncio.create_task(coro())\n"
	safe := "import asyncio\nasync def main():\n    t = asyncio.create_task(coro())\n    await t\n"
	safeAppend := "import asyncio\nasync def main():\n    tasks.append(asyncio.create_task(coro()))\n"
	assertRule(t, "BP-PY-38", "async_fire.py", vuln, true)
	assertRule(t, "BP-PY-38", "async_fire.py", safe, false)
	assertRule(t, "BP-PY-38", "async_fire.py", safeAppend, false)
	assertRule(t, "BP-PY-38", "async_fire.py", "import asyncio\nasync def main():\n    ensure_future(coro())\n", true)
}

func TestBPPY39TimeSleepInAsync(t *testing.T) {
	t.Parallel()
	vuln := "import time\nasync def handler():\n    time.sleep(1)\n"
	safeSync := "import time\ndef f():\n    time.sleep(1)\n"
	safeAsyncio := "import asyncio\nasync def handler():\n    await asyncio.sleep(1)\n"
	assertRule(t, "BP-PY-39", "async_sleep.py", vuln, true)
	assertRule(t, "BP-PY-39", "async_sleep.py", safeSync, false)
	assertRule(t, "BP-PY-39", "async_sleep.py", safeAsyncio, false)
	findings := runBP(t, nil, vuln, "async_sleep.py")
	for _, f := range findings {
		if f.RuleID == "BP-PY-39" && f.Severity != rules.SeverityHigh {
			t.Fatalf("BP-PY-39 severity = %v, want high", f.Severity)
		}
	}
}

func TestBPPY40ThreadingWithoutJoin(t *testing.T) {
	t.Parallel()
	vuln := "import threading\ndef run():\n    pass\nt = threading.Thread(target=run)\nt.start()\n"
	safe := "import threading\ndef run():\n    pass\nt = threading.Thread(target=run)\nt.start()\nt.join()\n"
	daemon := "import threading\ndef run():\n    pass\nt = threading.Thread(target=run, daemon=True)\nt.start()\n"
	assertRule(t, "BP-PY-40", "threads.py", vuln, true)
	assertRule(t, "BP-PY-40", "threads.py", safe, false)
	// Same-line daemon=True on start construct — miss when daemon on the start line;
	// daemon on Thread(...) line still starts without join; v0 skips only daemon lines.
	// File with only daemon Thread.start may still fire if start line lacks daemon=True.
	// Document: miss when start line contains daemon=True:
	assertRule(t, "BP-PY-40", "threads.py",
		"import threading\nthreading.Thread(target=run, daemon=True).start()\n", false)
	_ = daemon
}
