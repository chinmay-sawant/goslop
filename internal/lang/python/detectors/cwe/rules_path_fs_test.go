package cwe_test

import "testing"

func TestCWEPathFSHitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		id   string
		src  string
		hit  bool
	}{
		{"CWE-73 request path to open", "CWE-73", "from flask import request\nopen(request.args.get('file'))\n", true},
		{"CWE-73 secure filename", "CWE-73", "from werkzeug.utils import secure_filename\nopen(secure_filename(request.args.get('file')))\n", false},
		{"CWE-59 symlink check then open", "CWE-59", "import os\nif not os.path.islink(path):\n    open(path)\n", true},
		{"CWE-59 descriptor no-follow", "CWE-59", "fd = os.open(path, os.O_RDONLY | os.O_NOFOLLOW)\n", false},
		{"CWE-41 normpath request file", "CWE-41", "from flask import request\nopen(os.path.normpath(request.args.get('file')))\n", true},
		{"CWE-41 canonical path", "CWE-41", "from flask import request\nopen(os.path.realpath(os.path.normpath(request.args.get('file'))))\n", false},
		{"CWE-276 chmod world writable", "CWE-276", "import os\nos.chmod(path, 0o777)\n", true},
		{"CWE-276 private permissions", "CWE-276", "import os\nos.chmod(path, 0o600)\n", false},
		{"CWE-378 mktemp", "CWE-378", "import tempfile\nname = tempfile.mktemp()\n", true},
		{"CWE-378 named temp file", "CWE-378", "import tempfile\nwith tempfile.NamedTemporaryFile() as f:\n    f.write(b'ok')\n", false},
		{"CWE-426 prepend cwd", "CWE-426", "import os, sys\nsys.path.insert(0, os.getcwd())\n", true},
		{"CWE-426 fixed library root", "CWE-426", "import sys\nsys.path.append('/opt/trusted/lib')\n", false},
		{"CWE-250 switches to root", "CWE-250", "import os\nos.seteuid(0)\n", true},
		{"CWE-250 drops to app account", "CWE-250", "import os\nos.seteuid(1001)\n", false},
		{"CWE-494 exec remote response", "CWE-494", "from urllib.request import urlopen\nexec(urlopen(url).read())\n", true},
		{"CWE-494 local code execution", "CWE-494", "exec(compile(local_source, 'plugin.py', 'exec'))\n", false},
		{"comments and docstrings are ignored", "CWE-378", "# tempfile.mktemp()\n\"\"\"os.seteuid(0)\"\"\"\n", false},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := hasRule(runScan(tc.src), tc.id); got != tc.hit {
				t.Fatalf("%s hit=%v, want %v\\nsrc:\\n%s", tc.id, got, tc.hit, tc.src)
			}
		})
	}
}
