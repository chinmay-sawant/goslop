package cwe_test

import "testing"

func TestSSRFRedirectAndChannelHitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		rule string
		src  string
		hit  bool
	}{
		{"request URL reaches requests", "CWE-918", "import requests\ndef fetch(request):\n    return requests.get(request.args['url'])\n", true},
		{"fixed URL is safe", "CWE-918", "import requests\ndef fetch():\n    return requests.get('https://api.example.test/status')\n", false},
		{"request URL reaches redirect", "CWE-601", "from flask import redirect\ndef next_page(request):\n    return redirect(request.args['next'])\n", true},
		{"fixed redirect is safe", "CWE-601", "from flask import redirect\ndef next_page():\n    return redirect('/account')\n", false},
		{"reuse and wildcard bind", "CWE-605", "import socket\ndef serve():\n    sock = socket.socket()\n    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)\n    sock.bind(('0.0.0.0', 8080))\n", true},
		{"reuse without wildcard bind is safe", "CWE-605", "import socket\ndef serve():\n    sock = socket.socket()\n    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)\n    sock.bind(('127.0.0.1', 8080))\n", false},
		{"unsigned webhook body", "CWE-924", "def payment_webhook(request):\n    event = request.get_json()\n    process(event)\n", true},
		{"HMAC verified webhook body is safe", "CWE-924", "import hmac\ndef payment_webhook(request):\n    body = request.get_data()\n    if not hmac.compare_digest(request.headers['X-Signature'], sign(body)):\n        return 'invalid', 400\n    process(body)\n", false},
		{"callback logs in request identity", "CWE-940", "def oauth_callback(request):\n    user = lookup(request.args['user_id'])\n    login_user(user)\n", true},
		{"callback verifies state", "CWE-940", "def oauth_callback(request):\n    verify_state(request.args['state'])\n    user = lookup(request.args['user_id'])\n    login_user(user)\n", false},
		{"request mail recipient", "CWE-941", "from django.core.mail import send_mail\ndef invite(request):\n    return send_mail('Invite', 'secret', 'service@example.test', request.args['email'])\n", true},
		{"fixed mail recipient is safe", "CWE-941", "from django.core.mail import send_mail\ndef invite():\n    return send_mail('Invite', 'secret', 'service@example.test', ['ops@example.test'])\n", false},
		{"docstring examples do not run", "CWE-918", "def docs():\n    \"\"\"requests.get(request.args['url'])\"\"\"\n    return 'ok'\n", false},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := hasRule(runScan(tc.src), tc.rule); got != tc.hit {
				t.Fatalf("%s hit=%v, want %v\\nsrc:\\n%s", tc.rule, got, tc.hit, tc.src)
			}
		})
	}
}
