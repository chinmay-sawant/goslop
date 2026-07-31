package cwe_test

import "testing"

func TestCWECryptoHitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		id   string
		src  string
		hit  bool
	}{
		{"CWE-295 requests verify false", "CWE-295", "import requests\nrequests.get('https://service.example', verify=False)\n", true},
		{"CWE-295 default TLS verification", "CWE-295", "import requests\nrequests.get('https://service.example')\n", false},
		{"CWE-328 hashlib MD5", "CWE-328", "import hashlib\ndigest = hashlib.md5(data).hexdigest()\n", true},
		{"CWE-328 SHA-256", "CWE-328", "import hashlib\ndigest = hashlib.sha256(data).hexdigest()\n", false},
		{"CWE-335 fixed random seed", "CWE-335", "import random\nrandom.seed(42)\n", true},
		{"CWE-335 runtime random seed", "CWE-335", "import random\nrandom.seed()\n", false},
		{"CWE-338 random token", "CWE-338", "import random\ntoken = random.getrandbits(128)\n", true},
		{"CWE-338 secrets token", "CWE-338", "import secrets\ntoken = secrets.token_urlsafe(32)\n", false},
		{"CWE-347 disabled JWT signature", "CWE-347", "import jwt\nclaims = jwt.decode(token, options={'verify_signature': False})\n", true},
		{"CWE-347 verified JWT", "CWE-347", "import jwt\nclaims = jwt.decode(token, key=key, algorithms=['HS256'])\n", false},
		{"CWE-1204 fixed AES CBC IV", "CWE-1204", "from Crypto.Cipher import AES\ncipher = AES.new(key, AES.MODE_CBC, iv=b'0' * 16)\n", true},
		{"CWE-1204 fresh AES CBC IV", "CWE-1204", "from Crypto.Cipher import AES\nimport secrets\ncipher = AES.new(key, AES.MODE_CBC, iv=secrets.token_bytes(16))\n", false},
		{"CWE-1240 XOR cipher", "CWE-1240", "def xor_encrypt(data, key):\n    return bytes(a ^ b for a, b in zip(data, key))\n", true},
		{"CWE-1240 reviewed library call", "CWE-1240", "from cryptography.fernet import Fernet\ndef encrypt(data):\n    return Fernet(key).encrypt(data)\n", false},
		{"CWE-1241 predictable session", "CWE-1241", "import random\nsession = random.choice(values)\n", true},
		{"CWE-1241 ordinary random sample", "CWE-1241", "import random\nchoice = random.choice(values)\n", false},
		{"CWE-1392 default password", "CWE-1392", "password = 'admin'\n", true},
		{"CWE-1392 non-default environment password", "CWE-1392", "import os\npassword = os.environ['SERVICE_PASSWORD']\n", false},
		{"comments do not create crypto findings", "CWE-328", "# hashlib.md5(data)\n\"\"\"random.seed(0)\"\"\"\n", false},
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
