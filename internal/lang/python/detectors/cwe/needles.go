package cwe

// pyCweNeedles is the pack-local SourceIndex table for Python CWE prefilters.
// Keep FN-safe: only tokens that must appear for a priority rule to fire.
var pyCweNeedles = []string{
	// CWE-502 deserialization
	"pickle.loads",
	"pickle.load",
	"pickle.Unpickler",
	"yaml.load",
	"yaml.unsafe_load",
	// CWE-78 command injection
	"os.system",
	"os.popen",
	"subprocess.",
	"shell=True",
	"commands.",
	// CWE-89 SQL injection
	"execute(",
	"executemany(",
	".execute(",
	".executemany(",
	"raw(",
	// CWE-22 path traversal
	"open(",
	"pathlib",
	"os.path.join",
	"os.remove",
	"os.unlink",
	"Path(",
	// CWE-79 XSS
	"mark_safe",
	"Markup(",
	"render_template_string",
	"|safe",
	"HttpResponse(",
	// CWE-749 exposed dangerous route methods
	"eval(",
	"exec(",
	"compile(",
	"__import__(",
	"importlib.import_module",
	// CWE-829 dynamic code inclusion
	"spec_from_file_location",
	"runpy.run_path",
	// CWE-695 low-level functionality
	"ctypes.",
	"cffi.FFI",
	"mmap.mmap",
	// CWE-215 sensitive debugging output
	"print(",
	"logging.debug",
	".debug(",
	// CWE-90 LDAP injection
	"ldap3",
	"ldap.initialize",
	".search(",
	".search_s(",
	// CWE-91 XML / XPath injection
	".xpath(",
	"XPath(",
	".fromstring(",
	// CWE-93 CRLF injection
	"response.headers[",
	".set_header(",
	".add_header(",
	"HttpResponseRedirect(",
	// CWE-117 log injection
	"logging.",
	"logger.",
	"log.",
	// CWE-915 mass assignment
	".objects.create",
	".objects.update",
	".__dict__.update",
	"setattr(",
	"request.data",
	"request.POST",
	// CWE-914 dynamically-identified variables
	"globals()",
	"locals()",
	"vars(",
	// CWE-916 insufficient password-hash computational effort
	"hashlib.md5",
	"hashlib.sha1",
	"crypt.crypt",
	"md5_crypt",
	// CWE-798 / CWE-256 / CWE-260 hard-coded password and credential storage
	"password",
	"PASSWORD",
	"passwd",
	"api_key",
	"API_KEY",
	"SECRET_KEY",
	"aws_secret_access_key",
	"DATABASES",
	"config[",
	// CWE-261 weak password encoding
	"base64.b64encode",
	"binascii.hexlify",
	"codecs.encode",
	// CWE-312 cleartext sensitive storage
	"access_token",
	"auth_token",
	// CWE-319 cleartext transport
	"http://",
	"ftplib.FTP",
	"smtplib.SMTP",
	// CWE-523 disabled credential-transport protections
	"verify=False",
	"CERT_NONE",
	"check_hostname",
	"_create_unverified_context",
	// CWE-547 hard-coded security settings
	"SECURE_SSL_REDIRECT",
	"SESSION_COOKIE_SECURE",
	"CSRF_COOKIE_SECURE",
	"SESSION_COOKIE_HTTPONLY",
	"SECURE_HSTS_SECONDS",
	"ALLOWED_HOSTS",
}
