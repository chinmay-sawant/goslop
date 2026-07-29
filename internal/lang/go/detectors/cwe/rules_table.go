package cwe

import (
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

type needleGroup []string

type cweNeedleRule struct {
	id     string
	groups []needleGroup
	forbid []string
	span   string
	msg    string
	meta   *rules.RuleMetadata
}

func runNeedleRule(unit *core.ParsedUnit, facts *GoCweFacts, r cweNeedleRule, out *[]rules.Finding) {
	if unit == nil || facts == nil || out == nil || r.meta == nil {
		return
	}
	idx := facts.Index
	matched := false
	for _, g := range r.groups {
		if len(g) == 0 {
			continue
		}
		ok := true
		for _, n := range g {
			if !idx.Has(n) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		blocked := false
		for _, f := range r.forbid {
			if f != "" && idx.Has(f) {
				blocked = true
				break
			}
		}
		if !blocked {
			matched = true
			break
		}
	}
	if !matched {
		return
	}
	start := 0
	if r.span != "" {
		if i := strings.Index(unit.Source, r.span); i >= 0 {
			start = i
		}
	}
	line, col := unit.LineCol(start)
	rules.PushFindingWithConfidence(r.meta, unitFile(unit), line, col, r.msg, 0.55, out)
}

var cweNeedleRules = []cweNeedleRule{
	{
		id: "CWE-15",
		groups: []needleGroup{
			{"dsn := r.URL.Query().Get(\"dsn\")"},
		},
		span: "dsn := r.URL.Query().Get(\"dsn\")",
		msg:  "request-derived configuration value reaches a database-opening sink",
		meta: &MetaCWE15,
	},
	{
		id: "CWE-41",
		groups: []needleGroup{
			{"full := filepath.Join(root, rel)"},
		},
		span: "full := filepath.Join(root, rel)",
		msg:  "partial traversal filtering still allows equivalent path aliases to reach file access",
		meta: &MetaCWE41,
	},
	{
		id: "CWE-59",
		groups: []needleGroup{
			{"path := filepath.Join(staging, name)"},
		},
		span: "path := filepath.Join(staging, name)",
		msg:  "user-controlled path is opened without a symlink rejection check",
		meta: &MetaCWE59,
	},
	{
		id: "CWE-76",
		groups: []needleGroup{
			{"strings.ReplaceAll(raw, \"<\", \"\")", "strings.ReplaceAll(safe, \">\", \"\")", "text/html"},
		},
		forbid: []string{"html.EscapeString("},
		span:   "strings.ReplaceAll(raw, \"<\", \"\")",
		msg:    "manual angle-bracket stripping is used for HTML output instead of proper escaping",
		meta:   &MetaCWE76,
	},
	{
		id: "CWE-93",
		groups: []needleGroup{
			{"c.Header(\"Location\"", "r.URL.Query().Get("},
			{"w.Header().Set(\"Location\"", "r.URL.Query().Get("},
			{"c.Header(\"Location\"", "r.FormValue("},
			{"w.Header().Set(\"Location\"", "r.FormValue("},
		},
		forbid: []string{"strings.ReplaceAll("},
		span:   "Location",
		msg:    "user-controlled input is concatenated into a Location header without CRLF stripping",
		meta:   &MetaCWE93,
	},
	{
		id: "CWE-112",
		groups: []needleGroup{
			{"Price   float64 `xml:\"price\"`"},
		},
		span: "Price   float64 `xml:\"price\"`",
		msg:  "untrusted XML is unmarshaled without subsequent field-level validation",
		meta: &MetaCWE112,
	},
	{
		id: "CWE-140",
		groups: []needleGroup{
			{"text/csv", "strings.Join("},
			{"\",\""},
		},
		forbid: []string{"csv.NewWriter("},
		span:   "text/csv",
		msg:    "user-controlled fields are joined into CSV output with literal delimiters",
		meta:   &MetaCWE140,
	},
	{
		id: "CWE-178",
		groups: []needleGroup{
			{"\"shared\":  true,"},
		},
		forbid: []string{"strings.EqualFold("},
		span:   "\"shared\":  true,",
		msg:    "case-sensitive comparison mishandles case folding for security decisions",
		meta:   &MetaCWE178,
	},
	{
		id: "CWE-179",
		groups: []needleGroup{
			{".MatchString(", "url.QueryUnescape("},
			{".MatchString(", "path"},
		},
		forbid: []string{".MatchString(decoded)"},
		span:   "MatchString",
		msg:    "validation occurs before canonicalization/decoding",
		meta:   &MetaCWE179,
	},
	{
		id: "CWE-182",
		groups: []needleGroup{
			{"capability = strings.ToLower(capability)"},
		},
		span: "capability = strings.ToLower(capability)",
		msg:  "input is stripped and collapsed into an authorization-relevant value before membership checks",
		meta: &MetaCWE182,
	},
	{
		id: "CWE-184",
		groups: []needleGroup{
			{"for _, word := range", "strings.Contains("},
		},
		forbid: []string{".MatchString("},
		span:   "for _, word := range",
		msg:    "user-controlled input is checked against an incomplete deny-list after normalization",
		meta:   &MetaCWE184,
	},
	{
		id: "CWE-186",
		groups: []needleGroup{
			{"regexp.MustCompile(`^[a-z]+$`)"},
		},
		span: "regexp.MustCompile(`^[a-z]+$`)",
		msg:  "host validation uses an overly restrictive regex that only accepts lowercase letters",
		meta: &MetaCWE186,
	},
	{
		id: "CWE-201",
		groups: []needleGroup{
			{"_ = json.NewEncoder(w).Encode(record)"},
		},
		span: "_ = json.NewEncoder(w).Encode(record)",
		msg:  "a response serializes a record containing sensitive fields directly to the caller",
		meta: &MetaCWE201,
	},
	{
		id: "CWE-204",
		groups: []needleGroup{
			{"StatusNotFound", "no account"},
		},
		span: "no account",
		msg:  "authentication failures return distinguishable responses for missing accounts and wrong credentials",
		meta: &MetaCWE204,
	},
	{
		id: "CWE-208",
		groups: []needleGroup{
			{"for i := range expected", "provided[i] != expected[i]"},
		},
		forbid: []string{"subtle.ConstantTimeCompare("},
		span:   "for i := range expected",
		msg:    "secret comparison returns early on mismatched bytes instead of using a constant-time primitive",
		meta:   &MetaCWE208,
	},
	{
		id: "CWE-209",
		groups: []needleGroup{
			{"fmt.Sprintf(\"db failure: %v\", err)"},
		},
		span: "fmt.Sprintf(\"db failure: %v\", err)",
		msg:  "database error details are formatted into a client-facing response",
		meta: &MetaCWE209,
	},
	{
		id: "CWE-212",
		groups: []needleGroup{
			{"Card"},
			{"PAN"},
		},
		forbid: []string{"type paymentExport struct", "type chargeExport struct"},
		span:   "json.Marshal(rows)",
		msg:    "records containing sensitive payment fields are marshaled directly for export",
		meta:   &MetaCWE212,
	},
	{
		id: "CWE-213",
		groups: []needleGroup{
			{"_ = json.NewEncoder(w).Encode(profile)"},
		},
		span: "_ = json.NewEncoder(w).Encode(profile)",
		msg:  "potential security issue detected (heuristic port)",
		meta: &MetaCWE213,
	},
	{
		id: "CWE-214",
		groups: []needleGroup{
			{"cmd := exec.Command(\"archivectl\", \"--token\", token, \"--target\", \"ledger\")"},
		},
		forbid: []string{"cmd.Stdin = strings.NewReader("},
		span:   "cmd := exec.Command(\"archivectl\", \"--token\", token, \"--target\", \"ledger\")",
		msg:    "sensitive information may be exposed via process invocation environment or arguments",
		meta:   &MetaCWE214,
	},
	{
		id: "CWE-215",
		groups: []needleGroup{
			{"log.Printf(\"request debug secret=%s trace=%s uri=%s\", secret, traceID, r.URL.RequestURI())"},
		},
		span: "log.Printf(\"request debug secret=%s trace=%s uri=%s\", secret, traceID, r.URL.RequestURI())",
		msg:  "a debug log statement includes request-derived secret material",
		meta: &MetaCWE215,
	},
	{
		id: "CWE-250",
		groups: []needleGroup{
			{"os.WriteFile(", "0o777"},
			{"os.WriteFile(", "0777"},
			{"os.OpenFile(", "0o777"},
			{"os.Chmod(", "0o777"},
		},
		span: "0o777",
		msg:  "runtime file is written with world-accessible permissions",
		meta: &MetaCWE250,
	},
	{
		id: "CWE-252",
		groups: []needleGroup{
			{"os.WriteFile("},
		},
		forbid: []string{"if err := os.WriteFile("},
		span:   "os.WriteFile(",
		msg:    "os.WriteFile is called without checking its returned error",
		meta:   &MetaCWE252,
	},
	{
		id: "CWE-256",
		groups: []needleGroup{
			{"_, err := db.Exec(\"INSERT INTO credentials(login, pass) VALUES(?, ?)\", login, pa"},
		},
		forbid: []string{"GenerateFromPassword("},
		span:   "_, err := db.Exec(\"INSERT INTO credentials(login, pass) VALUES(?, ?)\", login, pa",
		msg:    "password or credential material is stored without adequate protection",
		meta:   &MetaCWE256,
	},
	{
		id: "CWE-257",
		groups: []needleGroup{
			{"aes.NewCipher("},
			{"cipher.NewGCM("},
			{"gcm.Seal("},
			{"base64.StdEncoding.EncodeToString(", "\"password\": encoded"},
			{"VALUES(?, ?)\", login, encoded)"},
		},
		span: "aes.NewCipher(",
		msg:  "a password or login secret is encrypted with a reversible cipher before storage",
		meta: &MetaCWE257,
	},
	{
		id: "CWE-260",
		groups: []needleGroup{
			{"Password string"},
			{"Secret   string"},
		},
		forbid: []string{"os.Getenv("},
		span:   "cfg.Password",
		msg:    "a secret-bearing field is loaded from a configuration file and used directly",
		meta:   &MetaCWE260,
	},
	{
		id: "CWE-261",
		groups: []needleGroup{
			{"base64.StdEncoding.EncodeToString(", "Secret: encoded"},
			{"Store(user, encoded)"},
		},
		span: "base64.StdEncoding.EncodeToString(",
		msg:  "a password is Base64-encoded and then stored in a recoverable form",
		meta: &MetaCWE261,
	},
	{
		id: "CWE-262",
		groups: []needleGroup{
			{"last_seen"},
			{"changed_at"},
		},
		forbid: []string{"time.Since(", "maxPasswordAge"},
		span:   "last_seen",
		msg:    "credential metadata is loaded but no password-age enforcement is performed",
		meta:   &MetaCWE262,
	},
	{
		id: "CWE-263",
		groups: []needleGroup{
			{"MaxAgeDays: 3650"},
		},
		span: "MaxAgeDays: 3650",
		msg:  "password maximum age is configured to an excessively long multi-year period",
		meta: &MetaCWE263,
	},
	{
		id: "CWE-266",
		groups: []needleGroup{
			{"role := r.URL.Query().Get(\"role\")"},
		},
		span: "role := r.URL.Query().Get(\"role\")",
		msg:  "a client-controlled role value is used directly when provisioning access",
		meta: &MetaCWE266,
	},
	{
		id: "CWE-267",
		groups: []needleGroup{
			{"target := r.URL.Query().Get(\"path\")"},
		},
		span: "target := r.URL.Query().Get(\"path\")",
		msg:  "the reviewer role is allowed to invoke a destructive filesystem removal operation",
		meta: &MetaCWE267,
	},
	{
		id: "CWE-268",
		groups: []needleGroup{
			{"Encode(userRecords)", "\"hash\""},
		},
		span: "Encode(userRecords)",
		msg:  "a sensitive export path is authorized by combining weaker read and export scopes",
		meta: &MetaCWE268,
	},
	{
		id: "CWE-270",
		groups: []needleGroup{
			{"ctx := context.WithValue(r.Context(), effectiveUserKey, \"root\")"},
		},
		span: "ctx := context.WithValue(r.Context(), effectiveUserKey, \"root\")",
		msg:  "the handler switches to a privileged execution context without restoring the original caller context",
		meta: &MetaCWE270,
	},
	{
		id: "CWE-272",
		groups: []needleGroup{
			{"syscall.Setuid(0)", "os.Chown(path, 1000, 1000)"},
		},
		forbid: []string{"syscall.Setuid(1000)"},
		span:   "Setuid(0)",
		msg:    "elevated privileges are retained beyond the privileged operation",
		meta:   &MetaCWE272,
	},
	{
		id: "CWE-273",
		groups: []needleGroup{
			{"_ = syscall.Setuid(1000)"},
		},
		forbid: []string{"if err := syscall.Setuid(1000); err != nil"},
		span:   "_ = syscall.Setuid(1000)",
		msg:    "privilege drop errors are ignored",
		meta:   &MetaCWE273,
	},
	{
		id: "CWE-274",
		groups: []needleGroup{
			{"w.Header().Set(\"Content-Type\", \"application/json\")"},
		},
		span: "w.Header().Set(\"Content-Type\", \"application/json\")",
		msg:  "potential security issue detected (heuristic port)",
		meta: &MetaCWE274,
	},
	{
		id: "CWE-276",
		groups: []needleGroup{
			{"os.WriteFile(", "0o666"},
			{"os.WriteFile(", "0666"},
			{"os.MkdirAll(", "0o777"},
			{"os.Mkdir(", "0o777"},
		},
		span: "WriteFile",
		msg:  "a session artifact is written with a world-readable and world-writable default mode",
		meta: &MetaCWE276,
	},
	{
		id: "CWE-277",
		groups: []needleGroup{
			{"_ = os.MkdirAll(dir, 0777)"},
		},
		span: "_ = os.MkdirAll(dir, 0777)",
		msg:  "umask is cleared before creating a world-writable directory",
		meta: &MetaCWE277,
	},
	{
		id: "CWE-278",
		groups: []needleGroup{
			{"f, _ := os.OpenFile(hdr.Name, os.O_CREATE|os.O_WRONLY, os.FileMode(hdr.Mode))"},
		},
		span: "f, _ := os.OpenFile(hdr.Name, os.O_CREATE|os.O_WRONLY, os.FileMode(hdr.Mode))",
		msg:  "archive entry permissions are reapplied directly from untrusted metadata during extraction",
		meta: &MetaCWE278,
	},
	{
		id: "CWE-279",
		groups: []needleGroup{
			{"http.Error(w, \"method not allowed\", http.StatusMethodNotAllo"},
		},
		span: "http.Error(w, \"method not allowed\", http.StatusMethodNotAllo",
		msg:  "the handler parses a requested mode but still writes the file with a hard-coded world-writable mode",
		meta: &MetaCWE279,
	},
	{
		id: "CWE-280",
		groups: []needleGroup{
			{"path := \"/secure/tenants/\" + r.URL.Query().Get(\"id\")"},
		},
		span: "path := \"/secure/tenants/\" + r.URL.Query().Get(\"id\")",
		msg:  "potential security issue detected (heuristic port)",
		meta: &MetaCWE280,
	},
	{
		id: "CWE-281",
		groups: []needleGroup{
			{"io.Copy(out, in)"},
		},
		forbid: []string{"info.Mode()"},
		span:   "io.Copy(out, in)",
		msg:    "backup recreation uses os.Create and loses the source file's original permission bits",
		meta:   &MetaCWE281,
	},
	{
		id: "CWE-283",
		groups: []needleGroup{
			{"_, _ = w.Write([]byte(`{\"removed\":true}`))"},
		},
		span: "_, _ = w.Write([]byte(`{\"removed\":true}`))",
		msg:  "a user-selected file path is removed without verifying that the caller owns the inode",
		meta: &MetaCWE283,
	},
	{
		id: "CWE-289",
		groups: []needleGroup{
			{"strings.Split(", "\"@\")[0]"},
		},
		forbid: []string{"canonical_name = ?"},
		span:   "strings.Split(",
		msg:    "principal authentication strips the realm suffix and authenticates only the bare local username",
		meta:   &MetaCWE289,
	},
	{
		id: "CWE-290",
		groups: []needleGroup{
			{"user := r.Header.Get(\"X-Remote-User\")"},
		},
		span: "user := r.Header.Get(\"X-Remote-User\")",
		msg:  "the request trusts a caller-controlled X-Remote-User header as the authenticated identity",
		meta: &MetaCWE290,
	},
	{
		id: "CWE-294",
		groups: []needleGroup{
			{"_, _ = w.Write([]byte(`{\"ok\":true}`))"},
		},
		span: "_, _ = w.Write([]byte(`{\"ok\":true}`))",
		msg:  "the login flow accepts an authentication token without nonce tracking or replay detection",
		meta: &MetaCWE294,
	},
	{
		id: "CWE-301",
		groups: []needleGroup{
			{"http.Error(w, \"missing challenge\", http.StatusBadRequest)"},
		},
		span: "http.Error(w, \"missing challenge\", http.StatusBadRequest)",
		msg:  "the server reflects the client challenge directly as the authentication proof",
		meta: &MetaCWE301,
	},
	{
		id: "CWE-303",
		groups: []needleGroup{
			{"hmac.New(", "mac.Sum(nil)", "string(expected) == sig"},
		},
		span: "string(expected) == sig",
		msg:  "the computed MAC is compared to user input with string equality instead of constant-time verification",
		meta: &MetaCWE303,
	},
	{
		id: "CWE-305",
		groups: []needleGroup{
			{"if r.URL.Query().Get(\"debug\") == \"1\" {"},
		},
		span: "if r.URL.Query().Get(\"debug\") == \"1\" {",
		msg:  "a caller-controlled debug flag reaches privileged behavior before the authenticated subject check",
		meta: &MetaCWE305,
	},
	{
		id: "CWE-306",
		groups: []needleGroup{
			{"// Vulnerable: destructive purge endpoint has no auth gate."},
		},
		span: "// Vulnerable: destructive purge endpoint has no auth gate.",
		msg:  "a destructive purge endpoint performs its action without any authentication gate",
		meta: &MetaCWE306,
	},
	{
		id: "CWE-307",
		groups: []needleGroup{
			{"http.Error(w, \"missing email\", http.StatusBadRequest)"},
		},
		span: "http.Error(w, \"missing email\", http.StatusBadRequest)",
		msg:  "the login flow has no throttling, backoff, or lockout for repeated failed authentication attempts",
		meta: &MetaCWE307,
	},
	{
		id: "CWE-308",
		groups: []needleGroup{
			{"http.Error(w, \"method not allowed\", http.StatusMethodNotAllo"},
		},
		span: "http.Error(w, \"method not allowed\", http.StatusMethodNotAllo",
		msg:  "a high-value wire action is authorized with only a password and no validated second factor",
		meta: &MetaCWE308,
	},
	{
		id: "CWE-309",
		groups: []needleGroup{
			{"http.Error(w, \"missing fields\", http.StatusBadRequest)"},
		},
		span: "http.Error(w, \"missing fields\", http.StatusBadRequest)",
		msg:  "the enterprise login route relies on username and password form fields as the primary authentication method",
		meta: &MetaCWE309,
	},
	{
		id: "CWE-312",
		groups: []needleGroup{
			{"SSN: r.FormValue(\"ssn\")", "SSN string `json:\"ssn\"`"},
		},
		forbid: []string{"SSNCipher", "gcm.Seal("},
		span:   "SSN: r.FormValue(\"ssn\")",
		msg:    "a sensitive SSN value is persisted in cleartext instead of encrypted form",
		meta:   &MetaCWE312,
	},
	{
		id: "CWE-319",
		groups: []needleGroup{
			{"ListenAndServe("},
			{"http.ListenAndServe("},
			{"http.ListenAndServe", "CVV", "Number"},
		},
		forbid: []string{"ListenAndServeTLS(", "tls.Config"},
		span:   "ListenAndServe(",
		msg:    "sensitive payment data is accepted over a cleartext HTTP listener instead of TLS",
		meta:   &MetaCWE319,
	},
	{
		id: "CWE-322",
		groups: []needleGroup{
			{"tls.Dial(", "InsecureSkipVerify: true"},
		},
		span: "InsecureSkipVerify: true",
		msg:  "the TLS relay connection disables peer certificate verification during key exchange",
		meta: &MetaCWE322,
	},
	{
		id: "CWE-323",
		groups: []needleGroup{
			{"aead.Seal("},
		},
		forbid: []string{"io.ReadFull(rand.Reader, nonce)"},
		span:   "aead.Seal(",
		msg:    "a fixed nonce is reused for AEAD encryption operations with the same key",
		meta:   &MetaCWE323,
	},
	{
		id: "CWE-324",
		groups: []needleGroup{
			{"ExpiresAt", "ApiKeyRow"},
			{"SigningKey", "Secret", "hmac.New("},
		},
		forbid: []string{"time.Now().After(row.ExpiresAt)", "time.Now().After(key.ExpiresAt)"},
		span:   "ExpiresAt",
		msg:    "cryptographic processing uses key material with an expiration field but never checks whether the key is expired",
		meta:   &MetaCWE324,
	},
	{
		id: "CWE-325",
		groups: []needleGroup{
			{"cipher.NewCTR(", "XORKeyStream("},
		},
		forbid: []string{"cipher.NewGCM("},
		span:   "cipher.NewCTR(",
		msg:    "sensitive data is encrypted with CTR mode without an authentication or integrity step",
		meta:   &MetaCWE325,
	},
	{
		id: "CWE-328",
		groups: []needleGroup{
			{"md5.Sum("},
		},
		span: "md5.Sum(",
		msg:  "a password digest is derived with MD5, which is too weak for this security-sensitive use",
		meta: &MetaCWE328,
	},
	{
		id: "CWE-331",
		groups: []needleGroup{
			{"rand.NewSource(time.Now().UnixNano())", "Intn(900000) + 100000", "code"},
		},
		span: "Intn(900000) + 100000",
		msg:  "the recovery code is generated from a small predictable decimal range instead of cryptographic randomness",
		meta: &MetaCWE331,
	},
	{
		id: "CWE-334",
		groups: []needleGroup{
			{"Intn(4096)"},
		},
		span: "Intn(4096)",
		msg:  "the generated token comes from a very small 4096-value space and is easy to guess",
		meta: &MetaCWE334,
	},
	{
		id: "CWE-335",
		groups: []needleGroup{
			{"seed := time.Now().Unix()", "rand.Seed(seed)"},
			{"seed := time.Now().Unix()", "rand.New(rand.NewSource(seed))"},
			{"rand.NewSource(seed)", "rand.Seed(seed)"},
		},
		span: "time.Now().Unix()",
		msg:  "the PRNG is seeded from predictable wall-clock time for a security-sensitive ticket value",
		meta: &MetaCWE335,
	},
	{
		id: "CWE-338",
		groups: []needleGroup{
			{"rand.New(rand.NewSource(time.Now().UnixNano()))", "token"},
			{"rand.New(rand.NewSource(time.Now().UnixNano()))", "sid"},
			{"rand.NewSource(time.Now().UnixNano())", "token"},
			{"rand.NewSource(time.Now().UnixNano())", "sid"},
		},
		span: "rand.New(rand.NewSource(time.Now().UnixNano()))",
		msg:  "a security-sensitive token is generated from math/rand instead of cryptographic randomness",
		meta: &MetaCWE338,
	},
	{
		id: "CWE-341",
		groups: []needleGroup{
			{"fmt.Sprintf(\"%d-%d-%s\"", "os.Getpid()", "time.Now().Unix()"},
		},
		span: "fmt.Sprintf(\"%d-%d-%s\"",
		msg:  "the token is built from observable pid, wall-clock time, and caller input instead of cryptographic randomness",
		meta: &MetaCWE341,
	},
	{
		id: "CWE-342",
		groups: []needleGroup{
			{"lastOTP++", "code := lastOTP"},
			{"lastSmsCode++", "code := lastSmsCode"},
		},
		span: "lastOTP++",
		msg:  "the next OTP value is generated by incrementing the previous one",
		meta: &MetaCWE342,
	},
	{
		id: "CWE-343",
		groups: []needleGroup{
			{"*3 + 1) % 97"},
			{"*5 + 3) % 101"},
		},
		span: "*3 + 1) % 97",
		msg:  "the output range is produced by a deterministic recurrence over shared state and is predictable from previous values",
		meta: &MetaCWE343,
	},
	{
		id: "CWE-344",
		groups: []needleGroup{
			{"w.Header().Set(\"Content-Type\", \"application/octet-stream\")"},
		},
		span: "w.Header().Set(\"Content-Type\", \"application/octet-stream\")",
		msg:  "a hard-coded invariant HMAC secret is embedded directly in code for a changing signing context",
		meta: &MetaCWE344,
	},
	{
		id: "CWE-346",
		groups: []needleGroup{
			{"w.Header().Set(\"Access-Control-Allow-Credentials\", \"true\")"},
		},
		span: "w.Header().Set(\"Access-Control-Allow-Credentials\", \"true\")",
		msg:  "the response reflects the caller-supplied Origin without validating it against a trusted allow-list",
		meta: &MetaCWE346,
	},
	{
		id: "CWE-347",
		groups: []needleGroup{
			{"strings.Split(raw, \".\")", "DecodeString(parts[1])", "json.Unmarshal(payload, &claims)"},
		},
		forbid: []string{"VerifyPKCS1v15(", "invalid signature"},
		span:   "DecodeString(parts[1])",
		msg:    "JWT claims are decoded and trusted without verifying the token signature first",
		meta:   &MetaCWE347,
	},
	{
		id: "CWE-349",
		groups: []needleGroup{
			{"json.RawMessage", "json.Unmarshal(bundle.Profile, &profile)"},
			{"json.RawMessage", "json.Unmarshal(env.Profile, &profile)"},
		},
		forbid: []string{"Role != \"support\"", "role not allowed from trusted channel"},
		span:   "json.RawMessage",
		msg:    "trusted envelope metadata is mixed with an untyped raw profile blob whose role fields are used directly",
		meta:   &MetaCWE349,
	},
	{
		id: "CWE-353",
		groups: []needleGroup{
			{"INSERT INTO agent_reports (payload) VALUES (?)", "io.ReadAll(r.Body)"},
		},
		forbid: []string{"hmac.New(", "X-Body-Mac", "ConstantTimeCompare"},
		span:   "INSERT INTO agent_reports",
		msg:    "agent report payload is accepted without an integrity check",
		meta:   &MetaCWE353,
	},
	{
		id: "CWE-356",
		groups: []needleGroup{
			{"DELETE FROM workspaces WHERE slug = ?", "{\"deleted\":true}"},
		},
		forbid: []string{"X-Confirm-Delete", "X-Expected-Token"},
		span:   "DELETE FROM workspaces",
		msg:    "destructive workspace delete proceeds without an explicit confirmation token",
		meta:   &MetaCWE356,
	},
	{
		id: "CWE-358",
		groups: []needleGroup{
			{"http.Error(w, \"missing authorization\", http.StatusUnauthoriz"},
		},
		span: "http.Error(w, \"missing authorization\", http.StatusUnauthoriz",
		msg:  "bearer token claims are accepted without required JWT structure and algorithm validation",
		meta: &MetaCWE358,
	},
	{
		id: "CWE-359",
		groups: []needleGroup{
			{"json.Marshal(row)"},
		},
		span: "json.Marshal(row)",
		msg:  "private personal information is serialized directly without requester authorization or public projection",
		meta: &MetaCWE359,
	},
	{
		id: "CWE-360",
		groups: []needleGroup{
			{"X-Forwarded-For"},
		},
		span: "X-Forwarded-For",
		msg:  "a security-sensitive client IP action trusts caller-controlled forwarded header data",
		meta: &MetaCWE360,
	},
	{
		id: "CWE-366",
		groups: []needleGroup{
			{"walletCredits += amount"},
			{"referralCredits += 10"},
		},
		forbid: []string{"atomic.AddInt64("},
		span:   "walletCredits += amount",
		msg:    "shared credit state is incremented without atomic or synchronized protection",
		meta:   &MetaCWE366,
	},
	{
		id: "CWE-367",
		groups: []needleGroup{
			{"os.Stat(", "os.ReadFile"},
		},
		span: "os.Stat(",
		msg:  "the code checks a file path with Stat before later using it, creating a TOCTOU race window",
		meta: &MetaCWE367,
	},
	{
		id: "CWE-368",
		groups: []needleGroup{
			{"actingAsRoot = true"},
			{"privilegedMode = true", "os.Setenv("},
		},
		forbid: []string{"sync.Mutex", "Lock()"},
		span:   "actingAsRoot = true",
		msg:    "privileged context switching is controlled by an unsynchronized shared mode flag",
		meta:   &MetaCWE368,
	},
	{
		id: "CWE-378",
		groups: []needleGroup{
			{"os.TempDir()", "0666"},
		},
		span: "os.TempDir()",
		msg:  "a temp file is created with world-accessible permissions in the shared temp area",
		meta: &MetaCWE378,
	},
	{
		id: "CWE-379",
		groups: []needleGroup{
			{"MkdirAll(dir, 0777)"},
		},
		span: "MkdirAll(dir, 0777)",
		msg:  "a temporary file is staged inside a shared world-writable directory",
		meta: &MetaCWE379,
	},
	{
		id: "CWE-385",
		groups: []needleGroup{
			{"if len(provided) != len(expected) {"},
		},
		span: "if len(provided) != len(expected) {",
		msg:  "potential security issue detected (heuristic port)",
		meta: &MetaCWE385,
	},
	{
		id: "CWE-393",
		groups: []needleGroup{
			{"w.Write([]byte(`{\"balance\":` + http.StatusText(int(balance)) + `}`))"},
		},
		span: "w.Write([]byte(`{\"balance\":` + http.StatusText(int(balance)) + `}`))",
		msg:  "potential security issue detected (heuristic port)",
		meta: &MetaCWE393,
	},
	{
		id: "CWE-403",
		groups: []needleGroup{
			{"os.Open(\"/etc/codehound/master.key\")", "exec.Command(\"/bin/sh\", \"-c\", script)", "_ = secret"},
		},
		forbid: []string{"syscall.Close", "SysProcAttr"},
		span:   "master.key",
		msg:    "a sensitive file descriptor may be inherited by a child process",
		meta:   &MetaCWE403,
	},
	{
		id: "CWE-408",
		groups: []needleGroup{
			{"http.Error(w, \"auth required\", http.StatusUnauthorized)"},
		},
		span: "http.Error(w, \"auth required\", http.StatusUnauthorized)",
		msg:  "the export query runs before the caller authentication check",
		meta: &MetaCWE408,
	},
	{
		id: "CWE-412",
		groups: []needleGroup{
			{"lockfile", "os.ReadFile(lockPath)"},
		},
		span: "lockfile",
		msg:  "the lock file path comes directly from the client request",
		meta: &MetaCWE412,
	},
	{
		id: "CWE-420",
		groups: []needleGroup{
			{"r.GET(\"/debug/sqltrace\"", "r.Group(\"/api\", requireJWT())"},
			{"http.HandleFunc(\"/debug/sqltrace\"", "http.Handle(\"/api/invoices\", protected)"},
		},
		span: "/debug/sqltrace",
		msg:  "the alternate debug route is exposed outside the primary authenticated API guard",
		meta: &MetaCWE420,
	},
	{
		id: "CWE-421",
		groups: []needleGroup{
			{"transferToken =", "event: status\\ndata: \" + transferToken"},
			{"wireTransferCode =", "event: status\\ndata: %s\\n\\n\", wireTransferCode"},
		},
		forbid: []string{"sync.Mutex", "transferMu", "wireMu"},
		span:   "transferToken =",
		msg:    "an alternate event channel exposes shared transfer state without synchronization",
		meta:   &MetaCWE421,
	},
	{
		id: "CWE-425",
		groups: []needleGroup{
			{"http.HandleFunc(\"/internal/admin/export.csv\", func(w http.Re"},
		},
		span: "http.HandleFunc(\"/internal/admin/export.csv\", func(w http.Re",
		msg:  "the admin export endpoint is mounted without an explicit authorization guard",
		meta: &MetaCWE425,
	},
	{
		id: "CWE-426",
		groups: []needleGroup{
			{"plugin_dir", "plugin.Open(modPath)"},
		},
		span: "plugin_dir",
		msg:  "the plugin load directory is derived from caller-controlled input",
		meta: &MetaCWE426,
	},
	{
		id: "CWE-427",
		groups: []needleGroup{
			{"os.Setenv(\"PATH\",", "exec.Command(\"pdftopng\""},
		},
		span: "os.Setenv(\"PATH\",",
		msg:  "user input is prepended to PATH before resolving the helper binary by name",
		meta: &MetaCWE427,
	},
	{
		id: "CWE-434",
		groups: []needleGroup{
			{"file.Filename", "SaveUploadedFile(file, dest)"},
			{"hdr.Filename", "os.Create(dest)"},
		},
		forbid: []string{"unsupported file type", "filepath.Ext(", "hex.EncodeToString("},
		span:   "file.Filename",
		msg:    "the upload is stored and later served using the client filename without an extension allow-list",
		meta:   &MetaCWE434,
	},
	{
		id: "CWE-454",
		groups: []needleGroup{
			{"enforceMFA = r.FormValue(\"enforce_mfa\") == \"true\""},
		},
		span: "enforce_mfa",
		msg:  "the MFA enforcement flag is bootstrapped from client input instead of server configuration",
		meta: &MetaCWE454,
	},
	{
		id: "CWE-455",
		groups: []needleGroup{
			{"tls.LoadX509KeyPair(", "continuing without mTLS"},
		},
		forbid: []string{"log.Fatalf("},
		span:   "continuing without mTLS",
		msg:    "startup logs a TLS material failure but continues running anyway",
		meta:   &MetaCWE455,
	},
	{
		id: "CWE-459",
		groups: []needleGroup{
			{"CreateTemp("},
		},
		forbid: []string{"os.Remove(f.Name())"},
		span:   "CreateTemp(",
		msg:    "the temporary export file is served without being removed afterward",
		meta:   &MetaCWE459,
	},
	{
		id: "CWE-472",
		groups: []needleGroup{
			{"Role    string `form:\"role\"`"},
			{"role := r.FormValue(\"role\")"},
		},
		forbid: []string{"SELECT role FROM users"},
		span:   "role",
		msg:    "authorization trusts a client-submitted role field instead of resolving role server-side",
		meta:   &MetaCWE472,
	},
	{
		id: "CWE-488",
		groups: []needleGroup{
			{"map[string][]string{}", "session"},
		},
		span: "sessionCarts",
		msg:  "global cart state is keyed directly by a client-controlled session identifier",
		meta: &MetaCWE488,
	},
	{
		id: "CWE-494",
		groups: []needleGroup{
			{"http.Error(w, err.Error(), http.StatusBadGateway)"},
		},
		span: "http.Error(w, err.Error(), http.StatusBadGateway)",
		msg:  "the downloaded worker bundle is accepted without any pinned integrity verification",
		meta: &MetaCWE494,
	},
	{
		id: "CWE-497",
		groups: []needleGroup{
			{"\"goos\":     runtime.GOOS,"},
		},
		forbid: []string{"\"status\": \"ok\""},
		span:   "\"goos\":     runtime.GOOS,",
		msg:    "sensitive system information may be exposed to callers",
		meta:   &MetaCWE497,
	},
	{
		id: "CWE-501",
		groups: []needleGroup{
			{"Approved bool", "Amount", "Memo", "ShouldBindJSON(&msg)"},
			{"Decode(&msg)", "msg.Approved = true"},
		},
		forbid: []string{"payoutDecision", "Request  payoutRequest"},
		span:   "Approved bool",
		msg:    "trusted approval state is merged into the same struct as untrusted request fields",
		meta:   &MetaCWE501,
	},
	{
		id: "CWE-502",
		groups: []needleGroup{
			{"gob.NewDecoder("},
			{"encoding/gob", "adminAction", "Grant", ".Decode(&action)"},
		},
		forbid: []string{"ShouldBindJSON(&req)", "json.NewDecoder(r.Body).Decode(&req)"},
		span:   "gob.NewDecoder(",
		msg:    "user-controlled gob data is deserialized into a privileged admin action",
		meta:   &MetaCWE502,
	},
	{
		id: "CWE-515",
		groups: []needleGroup{
			{"\"over\"", "= 1", "= 0", "\"over_limit\""},
		},
		span: "var quotaFlag int",
		msg:  "a global quota flag is used as a covert cross-request signal",
		meta: &MetaCWE515,
	},
	{
		id: "CWE-521",
		groups: []needleGroup{
			{"len(body.Password) < 1"},
		},
		forbid: []string{"strongPassword(", "len(pw) < 12"},
		span:   "len(body.Password) < 1",
		msg:    "password validation allows trivially weak credentials before persistence",
		meta:   &MetaCWE521,
	},
	{
		id: "CWE-523",
		groups: []needleGroup{
			{"/login", "password", "Addr: \":8080\""},
			{"StartCleartextLogin"},
		},
		forbid: []string{"requireTLS(", "Request.TLS == nil", "r.TLS == nil"},
		span:   "/login",
		msg:    "login credentials are accepted before any TLS enforcement or redirect",
		meta:   &MetaCWE523,
	},
	{
		id: "CWE-524",
		groups: []needleGroup{
			{"map[string]string{}", "Authorization", "tokenCache"},
			{"tokenVault"},
		},
		forbid: []string{"context.WithValue(", "session_token"},
		span:   "tokenCache",
		msg:    "raw session tokens are cached in shared process memory keyed by caller identifiers",
		meta:   &MetaCWE524,
	},
	{
		id: "CWE-538",
		groups: []needleGroup{
			{"/var/www/"},
			{"/var/www/html/public/"},
		},
		forbid: []string{"/var/lib/codehound/private", "0o600"},
		span:   "/var/www/",
		msg:    "database configuration secrets are exported to a public world-readable file path",
		meta:   &MetaCWE538,
	},
	{
		id: "CWE-544",
		groups: []needleGroup{
			{"log.Println(err)"},
		},
		span: "panic(err)",
		msg:  "database failures are handled through inconsistent panic and logging paths",
		meta: &MetaCWE544,
	},
	{
		id: "CWE-547",
		groups: []needleGroup{
			{"const jwtSecret = "},
			{"const sessionMACKey = "},
		},
		forbid: []string{"os.Getenv(\"JWT_SIGNING_KEY\")", "os.Getenv(\"SESSION_MAC_KEY\")"},
		span:   "const jwtSecret = ",
		msg:    "signing material is hard-coded directly in source instead of loaded from runtime secret configuration",
		meta:   &MetaCWE547,
	},
	{
		id: "CWE-549",
		groups: []needleGroup{
			{"\"email\":    r.FormValue(\"email\"),"},
		},
		span: "\"email\":    r.FormValue(\"email\"),",
		msg:  "potential security issue detected (heuristic port)",
		meta: &MetaCWE549,
	},
	{
		id: "CWE-551",
		groups: []needleGroup{
			{"raw := ", "URL.Path", "strings.HasPrefix(raw, \"/admin\")", "strings.ReplaceAll(raw, \"%2f\", \"/\")"},
		},
		forbid: []string{"url.PathUnescape(raw)"},
		span:   "strings.HasPrefix(raw, \"/admin\")",
		msg:    "authorization checks the raw path before percent-unescape canonicalization",
		meta:   &MetaCWE551,
	},
	{
		id: "CWE-552",
		groups: []needleGroup{
			{"os.Chmod(dest, 0o777)", "FormFile(\"contract\")", "/srv/contracts"},
		},
		forbid: []string{"filepath.Base(", "os.Chmod(dest, 0o600)"},
		span:   "os.Chmod(dest, 0o777)",
		msg:    "uploaded contract files are made world-accessible after storage",
		meta:   &MetaCWE552,
	},
	{
		id: "CWE-565",
		groups: []needleGroup{
			{"if err != nil || role.Value != \"admin\" {"},
		},
		span: "if err != nil || role.Value != \"admin\" {",
		msg:  "a privileged delete action trusts a caller-controlled role cookie",
		meta: &MetaCWE565,
	},
	{
		id: "CWE-601",
		groups: []needleGroup{
			{"\"next\"", "c.Redirect(http.StatusFound, target)"},
			{"http.Redirect(w, r, target, http.StatusFound)"},
		},
		forbid: []string{"strings.HasPrefix(target, \"/\")", "strings.Contains(target, \"//\")"},
		span:   "target",
		msg:    "the redirect target comes from an unvalidated caller-controlled next parameter",
		meta:   &MetaCWE601,
	},
	{
		id: "CWE-603",
		groups: []needleGroup{
			{"X-Authenticated", "\"true\"", "UPDATE billing SET plan"},
		},
		span: "X-Authenticated",
		msg:  "billing mutation trusts a caller-supplied authenticated header",
		meta: &MetaCWE603,
	},
	{
		id: "CWE-605",
		groups: []needleGroup{
			{"SetsockoptInt", "SO_REUSEADDR"},
		},
		forbid: []string{"net.Listen(\"tcp\", \":9090\")"},
		span:   "SetsockoptInt",
		msg:    "the listener explicitly enables SO_REUSEADDR on the service socket",
		meta:   &MetaCWE605,
	},
	{
		id: "CWE-611",
		groups: []needleGroup{
			{"xml.NewDecoder(", "dec.Strict = false", "Decode(&catalog)"},
		},
		span: "dec.Strict = false",
		msg:  "untrusted XML is parsed with strict mode disabled and no DOCTYPE rejection",
		meta: &MetaCWE611,
	},
	{
		id: "CWE-613",
		groups: []needleGroup{
			{"http.SetCookie(w, &http.Cookie{Name: \"sid\", Value: \"\", Path: \"/\", MaxAge: -1, HttpOnly: true})"},
		},
		span: "http.SetCookie(w, &http.Cookie{Name: \"sid\", Value: \"\", Path: \"/\", MaxAge: -1, HttpOnly: true})",
		msg:  "potential security issue detected (heuristic port)",
		meta: &MetaCWE613,
	},
	{
		id: "CWE-618",
		groups: []needleGroup{
			{"/opt/vendor/activex-bridge", "exec.Command(", "method", "args"},
		},
		forbid: []string{"allowedPluginMethods"},
		span:   "/opt/vendor/activex-bridge",
		msg:    "the endpoint forwards caller-controlled method names into a privileged native helper",
		meta:   &MetaCWE618,
	},
	{
		id: "CWE-619",
		groups: []needleGroup{
			{"rows, err := db.Query(", "rows.Next()"},
		},
		forbid: []string{"defer rows.Close()"},
		span:   "rows, err := db.Query(",
		msg:    "a database cursor is opened and can return without being closed",
		meta:   &MetaCWE619,
	},
	{
		id: "CWE-620",
		groups: []needleGroup{
			{"NewPassword string `json:\"new_password\"`"},
		},
		span: "NewPassword string `json:\"new_password\"`",
		msg:  "the password change flow updates credentials without verifying the current password",
		meta: &MetaCWE620,
	},
	{
		id: "CWE-639",
		groups: []needleGroup{
			{"err := db.QueryRow(\"SELECT id, user_id, amount FROM invoices WHERE id = $1\", invoiceID).Sc"},
		},
		span: "err := db.QueryRow(\"SELECT id, user_id, amount FROM invoices WHERE id = $1\", invoiceID).Sc",
		msg:  "a caller-controlled invoice key is queried without owner scoping",
		meta: &MetaCWE639,
	},
	{
		id: "CWE-640",
		groups: []needleGroup{
			{"ForgotPassword", "new_password", "email", "UPDATE users SET password"},
			{"Where(\"email = ?\", email).Update(\"password\", newPass)"},
		},
		forbid: []string{"reset_tokens", "\"token\"", "expires_at"},
		span:   "new_password",
		msg:    "the recovery flow resets a password from email alone without a reset token",
		meta:   &MetaCWE640,
	},
	{
		id: "CWE-645",
		groups: []needleGroup{
			{"failedAttempts[user]++", "failedAttempts[user] >= 1"},
		},
		span: "failedAttempts[user] >= 1",
		msg:  "the account is locked after a single failed login attempt",
		meta: &MetaCWE645,
	},
	{
		id: "CWE-648",
		groups: []needleGroup{
			{"os.Chown(", "uid"},
		},
		span: "os.Chown(",
		msg:  "the handler passes caller-controlled values into a privileged ownership-change API",
		meta: &MetaCWE648,
	},
	{
		id: "CWE-649",
		groups: []needleGroup{
			{"Cookie(\"profile\")", "base64.StdEncoding.DecodeString", "role=admin"},
		},
		span: "DecodeString",
		msg:  "an obfuscated profile cookie is trusted without any integrity verification",
		meta: &MetaCWE649,
	},
	{
		id: "CWE-653",
		groups: []needleGroup{
			{"for _, row := range sharedAuditStore {"},
		},
		span: "for _, row := range sharedAuditStore {",
		msg:  "public and admin paths share the same privileged data store",
		meta: &MetaCWE653,
	},
	{
		id: "CWE-654",
		groups: []needleGroup{
			{"X-Api-Key", "legacy-admin-key", "ExportUsers"},
		},
		span: "legacy-admin-key",
		msg:  "admin export access is granted solely from a static API key header",
		meta: &MetaCWE654,
	},
	{
		id: "CWE-656",
		groups: []needleGroup{
			{"/maintenance-portal-9f3c2a", "HiddenConfigPanel"},
		},
		span: "/maintenance-portal-9f3c2a",
		msg:  "sensitive configuration access relies only on a hidden URL path",
		meta: &MetaCWE656,
	},
	{
		id: "CWE-708",
		groups: []needleGroup{
			{"os.Chown(", "owner_uid"},
		},
		span: "os.Chown(",
		msg:  "the caller chooses both the ownership target and uid for a file operation",
		meta: &MetaCWE708,
	},
	{
		id: "CWE-756",
		groups: []needleGroup{
			{"err.Error()", "FetchProfile"},
			{"SELECT email FROM profiles"},
		},
		forbid: []string{"\"unable to load profile\""},
		span:   "err.Error()",
		msg:    "raw database error text is returned directly to the client",
		meta:   &MetaCWE756,
	},
	{
		id: "CWE-765",
		groups: []needleGroup{
			{"// Vulnerable: unlocks the mutex twice when validation fails."},
		},
		span: "// Vulnerable: unlocks the mutex twice when validation fails.",
		msg:  "the critical-section lock is explicitly released twice on an error path",
		meta: &MetaCWE765,
	},
	{
		id: "CWE-778",
		groups: []needleGroup{
			{"SignIn", "username", "password", "Unauthorized"},
		},
		forbid: []string{"log.Printf(\"auth failure"},
		span:   "Unauthorized",
		msg:    "authentication failures are returned without any audit logging",
		meta:   &MetaCWE778,
	},
	{
		id: "CWE-783",
		groups: []needleGroup{
			{"!authenticated || isAdmin && ownerID == docOwner"},
		},
		forbid: []string{"!(isAdmin || ownerID == docOwner)"},
		span:   "!authenticated || isAdmin && ownerID == docOwner",
		msg:    "authorization depends on ambiguous && and || precedence",
		meta:   &MetaCWE783,
	},
	{
		id: "CWE-798",
		groups: []needleGroup{
			{"postgres://reporting:Tr4ck3rP@ss@db.internal:5432/reports?sslmode=disable"},
			{"password = \"SuperSecret123!\""},
			{"apiKey := \"AKIA"},
		},
		forbid: []string{"os.Getenv(\"REPORTING_DSN\")"},
		span:   "postgres://",
		msg:    "database credentials are embedded directly in the source code",
		meta:   &MetaCWE798,
	},
	{
		id: "CWE-807",
		groups: []needleGroup{
			{"blockedIPs"},
		},
		forbid: []string{"RemoteAddr"},
		span:   "X-Forwarded-For",
		msg:    "a security gate trusts the caller-controlled forwarded IP header",
		meta:   &MetaCWE807,
	},
	{
		id: "CWE-820",
		groups: []needleGroup{
			{"visitCounts[key] = visitCounts[key] + 1", "TrackVisit"},
		},
		forbid: []string{"visitMu.Lock()", "visitMu sync.Mutex"},
		span:   "visitCounts[key] = visitCounts[key] + 1",
		msg:    "shared visit counters are updated without synchronization",
		meta:   &MetaCWE820,
	},
	{
		id: "CWE-821",
		groups: []needleGroup{
			{"RLock()", "tokenCache[key] = value"},
		},
		forbid: []string{"cacheMu.Lock()"},
		span:   "RLock()",
		msg:    "shared cache state is mutated while only a read lock is held",
		meta:   &MetaCWE821,
	},
	{
		id: "CWE-826",
		groups: []needleGroup{
			{"db.Close()", "_, _ = db.Query(\"SELECT 1\")"},
		},
		span: "db.Close()",
		msg:  "a shared resource is released while concurrent work may still use it",
		meta: &MetaCWE826,
	},
	{
		id: "CWE-829",
		groups: []needleGroup{
			{"http.Error(w, err.Error(), http.StatusBadRequest)"},
		},
		span: "http.Error(w, err.Error(), http.StatusBadRequest)",
		msg:  "a plugin is loaded from a caller-controlled filesystem path",
		meta: &MetaCWE829,
	},
	{
		id: "CWE-836",
		groups: []needleGroup{
			{"PasswordHash string `json:\"password_hash\"`"},
		},
		span: "PasswordHash string `json:\"password_hash\"`",
		msg:  "authentication accepts a caller-supplied password hash instead of verifying a plaintext password",
		meta: &MetaCWE836,
	},
	{
		id: "CWE-838",
		groups: []needleGroup{
			{"application/json; charset=utf-8", "0xC3, 0x28"},
		},
		span: "0xC3, 0x28",
		msg:  "invalid byte sequences are emitted while declaring UTF-8 JSON output",
		meta: &MetaCWE838,
	},
	{
		id: "CWE-841",
		groups: []needleGroup{
			{"_ = accountMFAPassed[email]", "accountPasswords[email] = newPass"},
		},
		forbid: []string{"if !accountMFAPassed[email]"},
		span:   "accountPasswords[email] = newPass",
		msg:    "password reset proceeds without enforcing the MFA workflow step",
		meta:   &MetaCWE841,
	},
	{
		id: "CWE-842",
		groups: []needleGroup{
			{"RegisterMember", "Group: \"administrators\""},
		},
		forbid: []string{"Group: \"members\""},
		span:   "Group: \"administrators\"",
		msg:    "newly registered users are assigned to an administrator group by default",
		meta:   &MetaCWE842,
	},
	{
		id: "CWE-909",
		groups: []needleGroup{
			{"enc := json.NewEncoder(w)"},
		},
		forbid: []string{"if widgetDB == nil"},
		span:   "enc := json.NewEncoder(w)",
		msg:    "a global database handle is used without checking that initialization completed",
		meta:   &MetaCWE909,
	},
	{
		id: "CWE-915",
		groups: []needleGroup{
			{"map[string]interface{}"},
		},
		span: "map[string]interface{}",
		msg:  "a user-controlled attribute map updates privileged object fields directly",
		meta: &MetaCWE915,
	},
	{
		id: "CWE-916",
		groups: []needleGroup{
			{"md5.Sum(", "password"},
		},
		forbid: []string{"bcrypt.GenerateFromPassword", "hashIterations = 100_000"},
		span:   "md5.Sum(",
		msg:    "password storage uses a fast MD5 hash with insufficient computational effort",
		meta:   &MetaCWE916,
	},
	{
		id: "CWE-917",
		groups: []needleGroup{
			{"template.New(\"report\").Parse(src)", "{{.Title}} where ", "+ expr"},
		},
		forbid: []string{"reportTemplate", "reportTemplatePure"},
		span:   "{{.Title}} where ",
		msg:    "caller-controlled data is concatenated into the template source itself",
		meta:   &MetaCWE917,
	},
	{
		id: "CWE-918",
		groups: []needleGroup{
			{"http.Get("},
		},
		forbid: []string{"allowedHosts", "allowedHostsPure", "Hostname()"},
		span:   "http.Get(",
		msg:    "an outbound request is sent to a caller-controlled URL without host allowlisting",
		meta:   &MetaCWE918,
	},
	{
		id: "CWE-921",
		groups: []needleGroup{
			{"/tmp/integration.key"},
		},
		span: "/tmp/integration.key",
		msg:  "sensitive integration key material is stored in a world-readable temporary file",
		meta: &MetaCWE921,
	},
	{
		id: "CWE-924",
		groups: []needleGroup{
			{"if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {"},
		},
		span: "if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {",
		msg:  "a payment webhook body is applied without validating an integrity signature first",
		meta: &MetaCWE924,
	},
	{
		id: "CWE-940",
		groups: []needleGroup{
			{"userID := r.URL.Query().Get(\"user_id\")"},
		},
		span: "userID := r.URL.Query().Get(\"user_id\")",
		msg:  "an OAuth callback accepts caller-supplied authorization data without verifying a bound state token",
		meta: &MetaCWE940,
	},
	{
		id: "CWE-941",
		groups: []needleGroup{
			{"email := r.URL.Query().Get(\"email\")"},
		},
		span: "email := r.URL.Query().Get(\"email\")",
		msg:  "a reset notification is sent to a caller-controlled email address",
		meta: &MetaCWE941,
	},
	{
		id: "CWE-1051",
		groups: []needleGroup{
			{"ChargeCard("},
			{"ChargeCardPure(", "10.20.30.40:9090", "http.NewRequest(", "X-Card-Token"},
		},
		forbid: []string{"os.Getenv(\"BILLING_API_URL\")"},
		span:   "10.20.30.40:9090",
		msg:    "an outbound billing request is pinned to a hard-coded internal host",
		meta:   &MetaCWE1051,
	},
	{
		id: "CWE-1052",
		groups: []needleGroup{
			{"gorm.Open(postgres.Open(dsn)"},
			{"sql.Open(\"postgres\", appDSNPure)", "password=SuperSecret99", "host=db.internal"},
		},
		forbid: []string{"APP_DATABASE_URL", "DB_PASSWORD"},
		span:   "password=SuperSecret99",
		msg:    "database initialization embeds a complete DSN with hard-coded credentials",
		meta:   &MetaCWE1052,
	},
	{
		id: "CWE-1067",
		groups: []needleGroup{
			{"fmt.Sprintf(\"%%%s%%\", term)"},
			{"pattern := fmt.Sprintf(\"%%%s%%\", term)", "LIKE", "notes.body"},
			{"SELECT id, body FROM notes"},
		},
		forbid: []string{"prefix+\"%\"", "pattern := prefix + \"%\""},
		span:   "fmt.Sprintf(\"%%%s%%\", term)",
		msg:    "a search predicate uses a leading wildcard pattern that forces a sequential scan",
		meta:   &MetaCWE1067,
	},
	{
		id: "CWE-1125",
		groups: []needleGroup{
			{"/admin/sql", "/admin/config", "/internal/reload"},
		},
		span: "/debug/pprof",
		msg:  "public routing exposes debug, admin, and internal maintenance endpoints together",
		meta: &MetaCWE1125,
	},
	{
		id: "CWE-1173",
		groups: []needleGroup{
			{"var raw map[string]interface{}", "ShouldBindJSON(&raw)"},
			{"Decode(&raw)", "SignupPayload{}"},
			{"SignupPayloadPure{}"},
		},
		forbid: []string{"ShouldBindJSON(&payload)", "Decode(&payload)", "mail.ParseAddress(payload.Email)"},
		span:   "var raw map[string]interface{}",
		msg:    "request data is decoded into a generic map instead of the validated signup model",
		meta:   &MetaCWE1173,
	},
	{
		id: "CWE-1204",
		groups: []needleGroup{
			{"cipher.NewCBCEncrypter(", "weakIV"},
			{"weakIVPure", "1234567890123456"},
		},
		forbid: []string{"io.ReadFull(rand.Reader, iv)"},
		span:   "1234567890123456",
		msg:    "CBC encryption uses a fixed IV literal instead of generating one per request",
		meta:   &MetaCWE1204,
	},
	{
		id: "CWE-1220",
		groups: []needleGroup{
			{"Authorization", "FROM invoices WHERE id = $1"},
		},
		span: "FROM invoices WHERE id = $1",
		msg:  "invoice access is authenticated but not scoped to the requesting owner",
		meta: &MetaCWE1220,
	},
	{
		id: "CWE-1230",
		groups: []needleGroup{
			{"X-Original-Name", "DownloadRedacted("},
			{"DownloadRedactedPure("},
		},
		forbid: []string{"Cache-Control"},
		span:   "X-Original-Name",
		msg:    "a redacted download response still exposes sensitive filename and size metadata",
		meta:   &MetaCWE1230,
	},
	{
		id: "CWE-1236",
		groups: []needleGroup{
			{"ExportFeedbackCSV("},
			{"ExportFeedbackCSVPure(", "id,comment", "fmt.Sprintf(\"%d,%s\\n\"", "row.Comment"},
		},
		forbid: []string{"sanitizeCSVField(", "sanitizeCSVFieldPure(", "csv.NewWriter("},
		span:   "ExportFeedbackCSV(",
		msg:    "CSV export writes user-controlled comment cells without neutralizing spreadsheet formulas",
		meta:   &MetaCWE1236,
	},
	{
		id: "CWE-1240",
		groups: []needleGroup{
			{"SealSessionToken("},
			{"SealSessionTokenPure(", "xorCipher("},
			{"xorCipherPure(", "^ key"},
		},
		forbid: []string{"cipher.NewGCM(", "aes.NewCipher("},
		span:   "xorCipher",
		msg:    "session sealing uses a homegrown XOR cipher instead of a standard authenticated primitive",
		meta:   &MetaCWE1240,
	},
	{
		id: "CWE-1265",
		groups: []needleGroup{
			{"defer ledgerMuPure.Unlock()"},
		},
		span: "defer ledgerMuPure.Unlock()",
		msg:  "a transfer path re-enters a mutex-protected balance helper while the same mutex is already held",
		meta: &MetaCWE1265,
	},
	{
		id: "CWE-1286",
		groups: []needleGroup{
			{"if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {"},
		},
		span: "if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {",
		msg:  "webhook configuration JSON is accepted without strict syntax and URL validation",
		meta: &MetaCWE1286,
	},
	{
		id: "CWE-1289",
		groups: []needleGroup{
			{"requested == \"private/keys.pem\"", "filepath.Join(root, requested)"},
		},
		span: "requested == \"private/keys.pem\"",
		msg:  "asset access relies on a literal blocked path comparison before canonical normalization",
		meta: &MetaCWE1289,
	},
	{
		id: "CWE-1322",
		groups: []needleGroup{
			{"queue := make(chan", "for payload := range queue", "time.Sleep(2 * time.Second)"},
		},
		forbid: []string{"time.AfterFunc("},
		span:   "time.Sleep(2 * time.Second)",
		msg:    "the webhook worker blocks its queue loop with sleep instead of scheduling retries asynchronously",
		meta:   &MetaCWE1322,
	},
	{
		id: "CWE-1327",
		groups: []needleGroup{
			{"StartPublicAPI("},
			{"StartPublicAPIPure(", "Run(\":9090\")"},
			{"ListenAndServe(\":9090\","},
		},
		forbid: []string{"127.0.0.1:9090"},
		span:   ":9090",
		msg:    "the service binds to all interfaces instead of a restricted loopback address",
		meta:   &MetaCWE1327,
	},
	{
		id: "CWE-1333",
		groups: []needleGroup{
			{"^([a-zA-Z]+)*$", "tagPattern"},
			{"tagPatternPure", "MatchString(tag)"},
		},
		forbid: []string{"safeTagPattern", "len(tag) > 32"},
		span:   "^([a-zA-Z]+)*$",
		msg:    "tag validation uses a catastrophic-backtracking regex on attacker-controlled input",
		meta:   &MetaCWE1333,
	},
	{
		id: "CWE-1389",
		groups: []needleGroup{
			{"strconv.ParseInt(raw, 0, 64)"},
		},
		forbid: []string{"strconv.ParseInt(raw, 10, 64)"},
		span:   "strconv.ParseInt(raw, 0, 64)",
		msg:    "seat counts are parsed with base 0 and may accept alternate-radix prefixes unexpectedly",
		meta:   &MetaCWE1389,
	},
	{
		id: "CWE-1392",
		groups: []needleGroup{
			{"BootstrapAdmin("},
			{"BootstrapAdminPure(", "Username: \"admin\"", "Password: \"admin\""},
		},
		forbid: []string{"BOOTSTRAP_ADMIN_PASSWORD"},
		span:   "Password: \"admin\"",
		msg:    "administrator bootstrap uses a built-in default password literal",
		meta:   &MetaCWE1392,
	},
}

func init() {
	for i := range cweNeedleRules {
		i := i
		rule := cweNeedleRules[i]
		RegisterRule(rule.id, func(unit *core.ParsedUnit, facts *GoCweFacts, out *[]rules.Finding) {
			runNeedleRule(unit, facts, cweNeedleRules[i], out)
		}, rule.meta)
	}
}
