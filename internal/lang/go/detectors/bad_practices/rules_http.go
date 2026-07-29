package badpractices

import (
	"strings"

	"github.com/chinmay/codehound/internal/core"
	"github.com/chinmay/codehound/internal/rules"
)

func init() {
	RegisterRule("BP-101", detectBP101)
	RegisterRule("BP-102", detectBP102)
	RegisterRule("BP-104", detectBP104)
	RegisterRule("BP-105", detectBP105)
	RegisterRule("BP-107", detectBP107)
	RegisterRule("BP-109", detectBP109)
	RegisterRule("BP-110", detectBP110)
	RegisterRule("BP-111", detectBP111)
	RegisterRule("BP-116", detectBP116)
	RegisterRule("BP-117", detectBP117)
	RegisterRule("BP-119", detectBP119)
	RegisterRule("BP-120", detectBP120)
	RegisterRule("BP-122", detectBP122)
	RegisterRule("BP-146", detectBP146)
	RegisterRule("BP-147", detectBP147)
	RegisterRule("BP-149", detectBP149)
	RegisterRule("BP-151", detectBP151)
	RegisterRule("BP-155", detectBP155)
	RegisterRule("BP-156", detectBP156)
	RegisterRule("BP-158", detectBP158)
	RegisterRule("BP-159", detectBP159)
	RegisterRule("BP-160", detectBP160)
}

func detectBP101(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-101")
	// Write before WriteHeader
	src := unit.Source
	if !strings.Contains(src, "WriteHeader") && !strings.Contains(src, ".Write(") {
		return
	}
	writePos := strings.Index(src, "w.Write(")
	if writePos < 0 {
		writePos = strings.Index(src, "Write([]byte")
	}
	headerPos := strings.Index(src, "WriteHeader(")
	if writePos >= 0 && headerPos >= 0 && writePos < headerPos {
		pushAt(unit, meta, headerPos, "HTTP header/status written after the response body has started", out)
	}
}

func detectBP102(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-102")
	if !strings.Contains(unit.Source, "http.ResponseWriter") {
		return
	}
	// if err != nil { return } without WriteHeader/http.Error
	if strings.Contains(unit.Source, "err != nil") && strings.Contains(unit.Source, "return") {
		if !strings.Contains(unit.Source, "WriteHeader") && !strings.Contains(unit.Source, "http.Error") &&
			!strings.Contains(unit.Source, "w.Write") {
			// weak: only if handler-shaped
			if pos := strings.Index(unit.Source, "err != nil"); pos >= 0 {
				pushAt(unit, meta, pos, "HTTP error path returns without writing a status response", out)
			}
		}
	}
}

func detectBP104(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-104")
	// duplicate Handle/HandleFunc patterns
	patterns := map[string]int{}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		for _, meth := range []string{"HandleFunc(", "Handle(", ".GET(", ".POST("} {
			if i := strings.Index(t, meth); i >= 0 {
				// extract first string literal arg
				rest := t[i+len(meth):]
				if strings.HasPrefix(rest, `"`) {
					end := strings.Index(rest[1:], `"`)
					if end >= 0 {
						pat := rest[:end+2]
						if prev, ok := patterns[pat]; ok {
							pushAt(unit, meta, line.byte, "duplicate HTTP mux pattern registration", out)
							_ = prev
							return
						}
						patterns[pat] = line.byte
					}
				}
			}
		}
	}
}

func detectBP105(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-105")
	if !strings.Contains(unit.Source, "http.Cookie") && !strings.Contains(unit.Source, "&http.Cookie{") {
		return
	}
	src := unit.Source
	if strings.Contains(src, "http.Cookie{") || strings.Contains(src, "&http.Cookie{") {
		// extract cookie literal roughly
		if !strings.Contains(src, "HttpOnly") || !strings.Contains(src, "Secure") {
			if pos := strings.Index(src, "Cookie{"); pos >= 0 {
				pushAt(unit, meta, pos, "sensitive cookie missing Secure and/or HttpOnly flags", out)
			}
		}
	}
}

func detectBP107(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-107")
	// middleware func(http.Handler) http.Handler without next.ServeHTTP
	if strings.Contains(unit.Source, "http.Handler") && strings.Contains(unit.Source, "func(") {
		if strings.Contains(unit.Source, "http.HandlerFunc") || strings.Contains(unit.Source, "http.Handler)") {
			if !strings.Contains(unit.Source, "ServeHTTP") && !strings.Contains(unit.Source, "next.ServeHTTP") &&
				!strings.Contains(unit.Source, "next(") {
				// only if looks like middleware
				if strings.Contains(unit.Source, "next http.Handler") || strings.Contains(unit.Source, "http.Handler)") {
					if pos := strings.Index(unit.Source, "http.Handler"); pos >= 0 {
						pushAt(unit, meta, pos, "HTTP middleware never invokes the next handler", out)
					}
				}
			}
		}
	}
}

func detectBP109(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-109")
	if !strings.Contains(unit.Source, "gin.") && !strings.Contains(unit.Source, "*gin.Context") {
		return
	}
	// c.JSON(4xx/5xx) without c.Abort
	if (strings.Contains(unit.Source, "c.JSON(") || strings.Contains(unit.Source, "c.AbortWithStatus")) &&
		(strings.Contains(unit.Source, "http.StatusBadRequest") || strings.Contains(unit.Source, "400") ||
			strings.Contains(unit.Source, "http.StatusInternalServerError") || strings.Contains(unit.Source, "500")) {
		if !strings.Contains(unit.Source, "c.Abort") && !strings.Contains(unit.Source, "AbortWithStatus") {
			if pos := strings.Index(unit.Source, "c.JSON("); pos >= 0 {
				pushAt(unit, meta, pos, "Gin error response without Abort; handler may continue", out)
			}
		}
	}
}

func detectBP110(unit *core.ParsedUnit, facts *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-110")
	if !strings.Contains(unit.Source, "ShouldBind") && !strings.Contains(unit.Source, "BindJSON") && !strings.Contains(unit.Source, "Bind(") {
		return
	}
	msg := "Gin bind error is discarded; check the error and abort the request"
	bindNames := []string{"ShouldBindJSON", "ShouldBind", "BindJSON", "MustBindWith", "ShouldBindQuery"}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		for _, b := range bindNames {
			if strings.Contains(t, b+"(") {
				// bare statement: c.ShouldBindJSON(&x)
				if !strings.Contains(t, ":=") && !strings.Contains(t, "err") && !strings.HasPrefix(t, "if ") {
					pushAt(unit, meta, line.byte, msg, out)
				}
			}
		}
	}
	_ = facts
}

func detectBP111(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-111")
	if !strings.Contains(unit.Source, "go ") || !strings.Contains(unit.Source, "gin.Context") && !strings.Contains(unit.Source, "*gin.Context") {
		return
	}
	// go func using c. after go
	if strings.Contains(unit.Source, "go func") && (strings.Contains(unit.Source, "c.JSON") || strings.Contains(unit.Source, "c.Request") || strings.Contains(unit.Source, "c.")) {
		// only if c is gin context in signature
		if strings.Contains(unit.Source, "*gin.Context") || strings.Contains(unit.Source, "gin.Context") {
			if pos := strings.Index(unit.Source, "go "); pos >= 0 {
				pushAt(unit, meta, pos, "Gin context used from a goroutine; copy values before leaving the handler", out)
			}
		}
	}
}

func detectBP116(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-116")
	if !strings.Contains(unit.Source, "echo.") && !strings.Contains(unit.Source, "echo.Context") {
		return
	}
	// c.JSON then return err
	if strings.Contains(unit.Source, "c.JSON(") && strings.Contains(unit.Source, "return err") {
		if pos := strings.Index(unit.Source, "return err"); pos >= 0 {
			pushAt(unit, meta, pos, "Echo handler writes an error response and also returns the raw error", out)
		}
	}
}

func detectBP117(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-117")
	if !strings.Contains(unit.Source, "echo") {
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, ".Bind(") && !strings.Contains(t, "err") && !strings.Contains(t, ":=") && !strings.HasPrefix(t, "if ") {
			pushAt(unit, meta, line.byte, "Echo Bind error is discarded", out)
		}
	}
}

func detectBP119(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-119")
	if !strings.Contains(unit.Source, "fiber") && !strings.Contains(unit.Source, "*fiber.Ctx") {
		return
	}
	if strings.Contains(unit.Source, "go ") && strings.Contains(unit.Source, "*fiber.Ctx") {
		if pos := strings.Index(unit.Source, "go "); pos >= 0 {
			pushAt(unit, meta, pos, "Fiber context used from a goroutine after the handler may return", out)
		}
	}
}

func detectBP120(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-120")
	if !strings.Contains(unit.Source, "BodyParser") {
		return
	}
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, "BodyParser(") && !strings.Contains(t, "err") && !strings.Contains(t, ":=") && !strings.HasPrefix(t, "if ") {
			pushAt(unit, meta, line.byte, "Fiber BodyParser error is discarded", out)
		}
	}
}

func detectBP122(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-122")
	if !strings.Contains(unit.Source, "chi.") && !strings.Contains(unit.Source, "go-chi") {
		return
	}
	if strings.Contains(unit.Source, "func(") && strings.Contains(unit.Source, "http.Handler") {
		if !strings.Contains(unit.Source, "next.ServeHTTP") && !strings.Contains(unit.Source, "next(") {
			if pos := strings.Index(unit.Source, "http.Handler"); pos >= 0 {
				pushAt(unit, meta, pos, "Chi middleware never invokes the next handler", out)
			}
		}
	}
}

func detectBP146(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-146")
	sensitive := []string{"password", "secret", "token", "api_key", "apikey", "ssn", "credit_card"}
	srcLower := strings.ToLower(unit.Source)
	if !(strings.Contains(unit.Source, "log.") || strings.Contains(unit.Source, "slog.") || strings.Contains(unit.Source, "logger.")) {
		return
	}
	for _, s := range sensitive {
		if strings.Contains(srcLower, s) && (strings.Contains(unit.Source, "log.") || strings.Contains(unit.Source, "slog.")) {
			// only if log line nearby
			for _, line := range codeLines(unit.Source) {
				lt := strings.ToLower(line.text)
				if (strings.Contains(lt, "log.") || strings.Contains(lt, "slog.")) && strings.Contains(lt, s) {
					pushAt(unit, meta, line.byte, "sensitive field may be logged; redact secrets before logging", out)
					return
				}
			}
		}
	}
}

func detectBP147(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-147")
	if packageName(unit.Source) == "main" {
		return
	}
	if strings.Contains(unit.Source, "log.") && (strings.Contains(unit.Source, "slog.") || strings.Contains(unit.Source, "zap.") || strings.Contains(unit.Source, "logrus")) {
		if pos := strings.Index(unit.Source, "log."); pos >= 0 {
			pushAt(unit, meta, pos, "service package mixes standard log with structured logging", out)
		}
	}
}

func detectBP149(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-149")
	// err != nil { slog.Error("msg") } without err attribute
	lines := codeLines(unit.Source)
	for i, line := range lines {
		t := strings.TrimSpace(line.text)
		if strings.Contains(t, "err != nil") {
			for j := i + 1; j < len(lines) && j < i+5; j++ {
				n := strings.TrimSpace(lines[j].text)
				if n == "}" {
					break
				}
				if (strings.Contains(n, "slog.Error") || strings.Contains(n, "log.Error") || strings.Contains(n, ".Error(")) &&
					!strings.Contains(n, "err") {
					pushAt(unit, meta, lines[j].byte, "error log omits the error value as an attribute", out)
					return
				}
			}
		}
	}
}

func detectBP151(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-151")
	if !strings.Contains(unit.Source, "os.Getenv") && !strings.Contains(unit.Source, "os.LookupEnv") {
		return
	}
	// getenv of secret-like then log
	for _, line := range codeLines(unit.Source) {
		t := strings.TrimSpace(line.text)
		upper := strings.ToUpper(t)
		if (strings.Contains(t, "os.Getenv") || strings.Contains(t, "os.LookupEnv")) &&
			(strings.Contains(upper, "SECRET") || strings.Contains(upper, "PASSWORD") || strings.Contains(upper, "TOKEN") || strings.Contains(upper, "API_KEY")) {
			if strings.Contains(unit.Source, "log.") || strings.Contains(unit.Source, "slog.") {
				// if same ident logged
				pushAt(unit, meta, line.byte, "secret environment value may be logged", out)
				return
			}
		}
	}
}

func detectBP155(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-155")
	if !strings.Contains(unit.Source, "json.NewDecoder") && !strings.Contains(unit.Source, "json.Unmarshal") {
		return
	}
	if strings.Contains(unit.Source, "r.Body") || strings.Contains(unit.Source, "Request.Body") || strings.Contains(unit.Source, "c.Request.Body") {
		if !strings.Contains(unit.Source, "MaxBytesReader") && !strings.Contains(unit.Source, "LimitReader") && !strings.Contains(unit.Source, "http.MaxBytesReader") {
			if pos := strings.Index(unit.Source, "json."); pos >= 0 {
				pushAt(unit, meta, pos, "JSON request body decoded without a size limit", out)
			}
		}
	}
}

func detectBP156(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-156")
	// json omitempty on password/secret fields
	for _, line := range codeLines(unit.Source) {
		t := strings.ToLower(line.text)
		if strings.Contains(t, "json:") && strings.Contains(t, "omitempty") {
			if strings.Contains(t, "password") || strings.Contains(t, "secret") || strings.Contains(t, "token") {
				pushAt(unit, meta, line.byte, "security-sensitive JSON field uses omitempty", out)
			}
		}
	}
}

func detectBP158(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-158")
	if !strings.Contains(unit.Source, "grpc") && !strings.Contains(unit.Source, "status.Error") {
		// grpc service methods often return (resp, error)
		if !strings.Contains(unit.Source, "context.Context") {
			return
		}
	}
	if strings.Contains(unit.Source, "return nil, err") && strings.Contains(unit.Source, "context.Context") {
		if !strings.Contains(unit.Source, "status.Error") && !strings.Contains(unit.Source, "status.Errorf") {
			// only if looks like grpc (Request/Response types)
			if strings.Contains(unit.Source, "Request") && strings.Contains(unit.Source, "Response") {
				if pos := strings.Index(unit.Source, "return nil, err"); pos >= 0 {
					pushAt(unit, meta, pos, "gRPC handler returns a raw error; wrap with status.Error", out)
				}
			}
		}
	}
}

func detectBP159(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-159")
	if !strings.Contains(unit.Source, "flag.") {
		return
	}
	// *flag used before flag.Parse
	parsePos := strings.Index(unit.Source, "flag.Parse(")
	if parsePos < 0 {
		// flag defined but never parsed
		if strings.Contains(unit.Source, "flag.String") || strings.Contains(unit.Source, "flag.Int") {
			// dereference * somewhere
			if strings.Contains(unit.Source, "*") {
				if pos := strings.Index(unit.Source, "flag."); pos >= 0 {
					pushAt(unit, meta, pos, "flag value read without flag.Parse", out)
				}
			}
		}
		return
	}
	// use of *flagVar before parsePos
	for _, line := range codeLines(unit.Source) {
		if line.byte >= parsePos {
			break
		}
		t := strings.TrimSpace(line.text)
		if strings.HasPrefix(t, "flag.") {
			continue
		}
		// crude: *name where name was from flag
		if strings.Contains(t, "*") && !strings.Contains(t, "func") && !strings.Contains(t, "import") {
			// only if flag.String etc. before
			if strings.Contains(unit.Source[:parsePos], "flag.String") || strings.Contains(unit.Source[:parsePos], "flag.Int") {
				// too noisy; require *var form after assignment from flag
			}
		}
	}
}

func detectBP160(unit *core.ParsedUnit, _ *bpFacts, out *[]rules.Finding) {
	meta := MetadataForID("BP-160")
	if !strings.Contains(unit.Source, "cobra.Command") && !strings.Contains(unit.Source, "cobra.") {
		return
	}
	if strings.Contains(unit.Source, "Run:") && !strings.Contains(unit.Source, "RunE:") {
		if pos := strings.Index(unit.Source, "Run:"); pos >= 0 {
			pushAt(unit, meta, pos, "Cobra command uses Run instead of RunE; errors are swallowed", out)
		}
	}
}
