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
}
