package cwe_test

import "testing"

func TestCWECodeDynamicHitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		id   string
		src  string
		hit  bool
	}{
		{
			name: "CWE-749 route exposes eval", id: "CWE-749", hit: true,
			src: "from flask import Flask\napp = Flask(__name__)\n\n@app.route('/run')\ndef run(expression):\n    return eval(expression)\n",
		},
		{
			name: "CWE-749 internal literal eval is not exposed", id: "CWE-749", hit: false,
			src: "def constant():\n    return eval('1 + 1')\n",
		},
		{
			name: "CWE-829 dynamic import", id: "CWE-829", hit: true,
			src: "import importlib\ndef load(name):\n    return importlib.import_module(name)\n",
		},
		{
			name: "CWE-829 literal import is package controlled", id: "CWE-829", hit: false,
			src: "import importlib\ndef load():\n    return importlib.import_module('json')\n",
		},
		{
			name: "CWE-695 ctypes native library", id: "CWE-695", hit: true,
			src: "import ctypes\nlib = ctypes.CDLL('libc.so.6')\n",
		},
		{
			name: "CWE-695 ctypes import alone is safe", id: "CWE-695", hit: false,
			src: "import ctypes\nvalue = ctypes.c_int(1)\n",
		},
		{
			name: "CWE-214 password in subprocess argv", id: "CWE-214", hit: true,
			src: "import subprocess\ndef upload(password):\n    subprocess.run(['tool', '--password', password])\n",
		},
		{
			name: "CWE-214 password file is safe", id: "CWE-214", hit: false,
			src: "import subprocess\ndef upload():\n    subprocess.run(['tool', '--password-file', '/run/secret'])\n",
		},
		{
			name: "CWE-214 dynamic password file path is safe", id: "CWE-214", hit: false,
			src: "import subprocess\ndef upload(path):\n    subprocess.run(['tool', '--password-file', path])\n",
		},
		{
			name: "CWE-215 debug logs token", id: "CWE-215", hit: true,
			src: "import logging\ndef request(token):\n    logging.debug('token=%s', token)\n",
		},
		{
			name: "CWE-215 redacted debug literal is safe", id: "CWE-215", hit: false,
			src: "import logging\ndef request():\n    logging.debug('token redacted')\n",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := hasRule(runScan(tc.src), tc.id)
			if got != tc.hit {
				t.Fatalf("%s hit=%v, want %v", tc.id, got, tc.hit)
			}
		})
	}
}
