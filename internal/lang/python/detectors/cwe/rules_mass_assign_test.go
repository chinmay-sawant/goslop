package cwe_test

import "testing"

func TestCWE915HitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		src  string
		hit  bool
	}{
		{
			name: "Django request data unpacked into model",
			src:  "def create(request):\n    return User.objects.create(**request.data)\n",
			hit:  true,
		},
		{
			name: "request driven setattr loop",
			src:  "def update(obj, request):\n    for key, value in request.data.items():\n        setattr(obj, key, value)\n",
			hit:  true,
		},
		{
			name: "allowlisted fields safe",
			src:  "def update(obj, request):\n    allowed = {'display_name'}\n    for key, value in request.data.items():\n        if key in allowed:\n            setattr(obj, key, value)\n",
			hit:  false,
		},
		{
			name: "internal mapping unpack safe",
			src:  "def create(values):\n    return User(**values)\n",
			hit:  false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := hasRule(runScan(tc.src), "CWE-915"); got != tc.hit {
				t.Fatalf("CWE-915 hit=%v, want %v\nsrc:\n%s", got, tc.hit, tc.src)
			}
		})
	}
}

func TestCWE914HitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		src  string
		hit  bool
	}{
		{
			name: "request selects global variable",
			src:  "def read(request):\n    return globals()[request.args['name']]\n",
			hit:  true,
		},
		{
			name: "request controls setattr name",
			src:  "def update(obj, request, value):\n    setattr(obj, request.form['field'], value)\n",
			hit:  true,
		},
		{
			name: "allowlisted name safe",
			src:  "def update(obj, request, value):\n    field = {'name': 'display_name'}[request.form['field']]\n    setattr(obj, field, value)\n",
			hit:  false,
		},
		{
			name: "static attribute safe",
			src:  "def update(obj, value):\n    setattr(obj, 'display_name', value)\n",
			hit:  false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := hasRule(runScan(tc.src), "CWE-914"); got != tc.hit {
				t.Fatalf("CWE-914 hit=%v, want %v\nsrc:\n%s", got, tc.hit, tc.src)
			}
		})
	}
}

func TestCWE916HitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		src  string
		hit  bool
	}{
		{
			name: "MD5 password hash",
			src:  "import hashlib\n\ndef digest(password):\n    return hashlib.md5(password.encode()).hexdigest()\n",
			hit:  true,
		},
		{
			name: "SHA1 password hash",
			src:  "import hashlib\n\ndef digest(passwd):\n    return hashlib.sha1(passwd.encode()).hexdigest()\n",
			hit:  true,
		},
		{
			name: "MD5 non password safe",
			src:  "import hashlib\n\ndef etag(body):\n    return hashlib.md5(body).hexdigest()\n",
			hit:  false,
		},
		{
			name: "bcrypt password hash safe",
			src:  "import bcrypt\n\ndef digest(password):\n    return bcrypt.hashpw(password.encode(), bcrypt.gensalt())\n",
			hit:  false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := hasRule(runScan(tc.src), "CWE-916"); got != tc.hit {
				t.Fatalf("CWE-916 hit=%v, want %v\nsrc:\n%s", got, tc.hit, tc.src)
			}
		})
	}
}
