package cwe_test

import "testing"

func TestCWESecretsHitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		id   string
		src  string
		hit  bool
	}{
		{"CWE-798 hard-coded API key", "CWE-798", "API_KEY = 'sk_live_example_1234'\n", true},
		{"CWE-798 environment key is safe", "CWE-798", "import os\nAPI_KEY = os.environ['API_KEY']\n", false},
		{"CWE-256 plaintext password", "CWE-256", "db_password = 'correct-horse-battery'\n", true},
		{"CWE-256 environment password is safe", "CWE-256", "import os\ndb_password = os.environ['DB_PASSWORD']\n", false},
		{"CWE-260 Django database literal", "CWE-260", "DATABASES = {'default': {'PASSWORD': 'database-secret'}}\n", true},
		{"CWE-260 configuration environment lookup is safe", "CWE-260", "DATABASES = {'default': {'PASSWORD': os.environ['DB_PASSWORD']}}\n", false},
		{"CWE-261 Base64 password", "CWE-261", "import base64\ndef encode(password):\n    return base64.b64encode(password.encode())\n", true},
		{"CWE-261 password hash is safe", "CWE-261", "import bcrypt\ndef digest(password):\n    return bcrypt.hashpw(password.encode(), bcrypt.gensalt())\n", false},
		{"CWE-312 source secret", "CWE-312", "SECRET_KEY = 'django-secret-example-1234'\n", true},
		{"CWE-312 environment secret is safe", "CWE-312", "import os\nSECRET_KEY = os.environ['SECRET_KEY']\n", false},
		{"CWE-319 HTTP basic auth", "CWE-319", "import requests\ndef login(password):\n    return requests.get('http://service.example/login', auth=('user', password))\n", true},
		{"CWE-319 HTTPS basic auth is safe", "CWE-319", "import requests\ndef login(password):\n    return requests.get('https://service.example/login', auth=('user', password))\n", false},
		{"CWE-319 HTTP token URL path is safe", "CWE-319", "import requests\nrequests.get('http://service.example/token')\n", false},
		{"CWE-547 insecure Django cookie setting", "CWE-547", "SESSION_COOKIE_SECURE = False\n", true},
		{"CWE-547 secure Django cookie setting is safe", "CWE-547", "SESSION_COOKIE_SECURE = True\n", false},
		{"CWE-523 requests verification disabled", "CWE-523", "import requests\nrequests.get('https://service.example', verify=False)\n", true},
		{"CWE-523 default verification is safe", "CWE-523", "import requests\nrequests.get('https://service.example')\n", false},
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
