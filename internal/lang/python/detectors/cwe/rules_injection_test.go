package cwe_test

import "testing"

func TestInjectionExpansionHitMiss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		rule string
		src  string
		hit  bool
	}{
		{"ldap f-string filter", "CWE-90", "import ldap3\ndef find(conn, user):\n    return conn.search(f'(uid={user})')\n", true},
		{"ldap escaped filter", "CWE-90", "from ldap3.utils.conv import escape_filter_chars\ndef find(conn, user):\n    return conn.search(f'(uid={escape_filter_chars(user)})')\n", false},
		{"xpath f-string", "CWE-91", "from lxml import etree\ndef find(root, user):\n    return root.xpath(f\"//user[@name='{user}']\")\n", true},
		{"xpath bound variable", "CWE-91", "from lxml import etree\ndef find(root, user):\n    return root.xpath('//user[@name=$name]', name=user)\n", false},
		{"header f-string", "CWE-93", "def redirect(response, next_url):\n    response.headers['Location'] = f'/next?to={next_url}'\n", true},
		{"header literal", "CWE-93", "def redirect(response):\n    response.headers['Location'] = '/home'\n", false},
		{"header CRLF removal", "CWE-93", "def redirect(response, next_url):\n    response.headers['Location'] = next_url.replace('\\r', '').replace('\\n', '')\n", false},
		{"exec dynamic", "CWE-94", "def run(code):\n    exec(code)\n", true},
		{"eval literal", "CWE-94", "def run():\n    return eval('1 + 1')\n", false},
		{"subprocess dynamic argv", "CWE-88", "import subprocess\ndef run(user):\n    return subprocess.run(['tool', '--file', user], shell=False)\n", true},
		{"subprocess static argv", "CWE-88", "import subprocess\ndef run():\n    return subprocess.run(['tool', '--version'], shell=False)\n", false},
		{"logger f-string", "CWE-117", "import logging\ndef log_user(user):\n    logging.info(f'login: {user}')\n", true},
		{"structured logger call", "CWE-117", "import logging\ndef log_user(user):\n    logging.info('login: %s', user)\n", false},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := hasRule(runScan(tc.src), tc.rule)
			if got != tc.hit {
				t.Fatalf("%s hit=%v want=%v\\nsrc:\\n%s", tc.rule, got, tc.hit, tc.src)
			}
		})
	}
}
