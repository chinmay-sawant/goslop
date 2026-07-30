// Package config loads and merges goslop.toml with CLI options.
//
// Schema mirrors the Rust product subset used by goslop: only/skip, fail_on,
// include/exclude, exclude_tests, baseline, cache, taint, bad_practices, export.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chinmay/goslop/internal/core"
	"github.com/chinmay/goslop/internal/rules"
	toml "github.com/pelletier/go-toml/v2"
)

// FileName is the default project config file name.
const FileName = "goslop.toml"

// Document is the root TOML document.
type Document struct {
	Goslop Section `toml:"goslop"`
}

// Section is the [goslop] table.
type Section struct {
	Languages    []string           `toml:"languages"`
	FailOn       *string            `toml:"fail_on"`
	Skip         []string           `toml:"skip"`
	Only         []string           `toml:"only"`
	Include      []string           `toml:"include"`
	Exclude      []string           `toml:"exclude"`
	ExcludeTests *bool              `toml:"exclude_tests"`
	Baseline     BaselineConfig     `toml:"baseline"`
	Cache        CacheConfig        `toml:"cache"`
	Taint        TaintConfig        `toml:"taint"`
	Typed        TypedConfig        `toml:"typed"`
	BadPractices BadPracticesConfig `toml:"bad_practices"`
	Export       ExportConfig       `toml:"export"`
}

// BaselineConfig is [goslop.baseline].
type BaselineConfig struct {
	Enabled *bool   `toml:"enabled"`
	Path    *string `toml:"path"`
}

// CacheConfig is [goslop.cache].
type CacheConfig struct {
	Enabled          *bool    `toml:"enabled"`
	Path             *string  `toml:"path"`
	MaxSizeMB        *uint64  `toml:"max_size_mb"`
	EvictTargetRatio *float64 `toml:"evict_target_ratio"`
	MaxFileSizeMB    *uint64  `toml:"max_file_size_mb"`
}

// TaintConfig is [goslop.taint].
type TaintConfig struct {
	Enabled   *bool `toml:"enabled"`
	ShowPaths *bool `toml:"show_paths"`
}

// TypedConfig is [goslop.typed].
type TypedConfig struct {
	Enabled *bool `toml:"enabled"`
}

// BadPracticesConfig is [goslop.bad_practices].
type BadPracticesConfig struct {
	Enabled           *bool             `toml:"enabled"`
	Severity          *string           `toml:"severity"`
	SeverityOverrides map[string]string `toml:"severity_overrides"`
}

// ExportConfig is [goslop.export] — on-disk context/chunk export options.
type ExportConfig struct {
	// WholeFunction, when true (default), expands Context in exported
	// findings/chunks to the enclosing FuncDecl/FuncLit. When false, uses a
	// nearby ~4-line window around the hit line.
	WholeFunction *bool `toml:"whole_function"`
}

// Load reads path as TOML. Missing file returns (nil, nil).
func Load(path string) (*Document, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	return Parse(b)
}

// Parse decodes TOML bytes. Unknown fields are rejected.
func Parse(data []byte) (*Document, error) {
	var doc Document
	// Defaults for nested tables when partially specified.
	doc.Goslop.Baseline.Enabled = boolPtr(true)
	doc.Goslop.Cache.Enabled = boolPtr(true)

	dec := toml.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("parse goslop.toml: %w", err)
	}
	if err := validate(&doc.Goslop); err != nil {
		return nil, err
	}
	return &doc, nil
}

func validate(s *Section) error {
	if s.FailOn != nil {
		if _, err := core.ParseFailPolicy(*s.FailOn); err != nil {
			return fmt.Errorf("goslop.fail_on: %w", err)
		}
	}
	if s.BadPractices.Severity != nil {
		if _, err := rules.ParseSeverity(*s.BadPractices.Severity); err != nil {
			return fmt.Errorf("goslop.bad_practices.severity: %w", err)
		}
	}
	for id, sev := range s.BadPractices.SeverityOverrides {
		if _, err := rules.ParseSeverity(sev); err != nil {
			return fmt.Errorf("goslop.bad_practices.severity_overrides[%q]: %w", id, err)
		}
	}
	return nil
}

// Discover walks upward from start looking for goslop.toml.
// Returns "" when none found.
func Discover(start string) string {
	if start == "" {
		start = "."
	}
	abs, err := filepath.Abs(start)
	if err != nil {
		abs = start
	}
	dir := abs
	if fi, err := os.Stat(dir); err == nil && !fi.IsDir() {
		dir = filepath.Dir(dir)
	}
	for {
		cand := filepath.Join(dir, FileName)
		if fi, err := os.Stat(cand); err == nil && !fi.IsDir() {
			return cand
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// MergeInput is CLI-visible state used when applying config.
// Boolean "set" flags distinguish explicit CLI overrides from defaults.
type MergeInput struct {
	Only           []string
	Skip           []string
	IncludeTests   bool // CLI --include-tests
	NoCache        bool // CLI --no-cache
	CacheDir       string
	NoBaseline     bool
	BaselineFile   string
	Taint          bool
	NoTaint        bool
	TaintShowPaths bool
	// FailOnCLI is set when the user passed an explicit fail policy (optional).
	// Currently CLI uses --no-fail only; leave nil for profile defaults.
	NoFail bool
	// ConfigPath when non-empty loads that file instead of discovery.
	ConfigPath string
	// Paths are scan roots (used for discovery when ConfigPath empty).
	Paths []string
}

// Merged is config applied onto CLI defaults for the app layer.
type Merged struct {
	Doc            *Document // nil when no config file
	ConfigPath     string
	Only           []string
	Skip           []string
	Include        []string
	Exclude        []string
	IncludeTests   bool
	NoCache        bool
	CacheDir       string
	NoBaseline     bool
	BaselineFile   string
	Taint          bool
	NoTaint        bool
	TaintShowPaths bool
	FailPolicy     *core.FailPolicy // set when config fail_on present and not overridden by --no-fail
	// Bad practices
	BadPracticesEnabled *bool
	BPSeverity          *rules.Severity
	SeverityOverrides   map[string]rules.Severity
	// Cache sizing (optional)
	CacheMaxSizeMB     *uint64
	CacheMaxFileSizeMB *uint64
	// ExportWholeFunction is nil when unset (export defaults to true).
	ExportWholeFunction *bool
}

// LoadAndMerge discovers/loads goslop.toml and merges with CLI input.
func LoadAndMerge(in MergeInput) (*Merged, error) {
	path := in.ConfigPath
	if path == "" {
		start := "."
		if len(in.Paths) > 0 {
			start = in.Paths[0]
		}
		path = Discover(start)
	}
	out := &Merged{
		Only:           append([]string(nil), in.Only...),
		Skip:           append([]string(nil), in.Skip...),
		IncludeTests:   in.IncludeTests,
		NoCache:        in.NoCache,
		CacheDir:       in.CacheDir,
		NoBaseline:     in.NoBaseline,
		BaselineFile:   in.BaselineFile,
		Taint:          in.Taint,
		NoTaint:        in.NoTaint,
		TaintShowPaths: in.TaintShowPaths,
	}
	if path == "" {
		return out, nil
	}
	doc, err := Load(path)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return out, nil
	}
	out.Doc = doc
	out.ConfigPath = path
	s := doc.Goslop

	// only / skip are additive (union of config + CLI).
	out.Only = unionStable(s.Only, in.Only)
	out.Skip = unionStable(s.Skip, in.Skip)
	out.Include = append([]string(nil), s.Include...)
	out.Exclude = append([]string(nil), s.Exclude...)

	// exclude_tests: config false → include tests; CLI --include-tests wins.
	if !in.IncludeTests && s.ExcludeTests != nil && !*s.ExcludeTests {
		out.IncludeTests = true
	}

	// fail_on from config unless --no-fail.
	if !in.NoFail && s.FailOn != nil {
		p, err := core.ParseFailPolicy(*s.FailOn)
		if err != nil {
			return nil, fmt.Errorf("goslop.fail_on: %w", err)
		}
		out.FailPolicy = &p
	}

	// cache: CLI --no-cache wins; else config enabled/path.
	if !in.NoCache && s.Cache.Enabled != nil && !*s.Cache.Enabled {
		out.NoCache = true
	}
	if in.CacheDir == "" && s.Cache.Path != nil && *s.Cache.Path != "" {
		out.CacheDir = *s.Cache.Path
	}
	out.CacheMaxSizeMB = s.Cache.MaxSizeMB
	out.CacheMaxFileSizeMB = s.Cache.MaxFileSizeMB

	// baseline
	if !in.NoBaseline && s.Baseline.Enabled != nil && !*s.Baseline.Enabled {
		out.NoBaseline = true
	}
	if in.BaselineFile == "" && s.Baseline.Path != nil && *s.Baseline.Path != "" {
		out.BaselineFile = *s.Baseline.Path
	}

	// taint: CLI --taint / --no-taint win over config.
	if !in.Taint && !in.NoTaint && s.Taint.Enabled != nil {
		if *s.Taint.Enabled {
			out.Taint = true
		} else {
			out.NoTaint = true
		}
	}
	if !in.TaintShowPaths && s.Taint.ShowPaths != nil && *s.Taint.ShowPaths {
		out.TaintShowPaths = true
	}

	// bad practices
	if s.BadPractices.Enabled != nil {
		out.BadPracticesEnabled = s.BadPractices.Enabled
	}
	if s.BadPractices.Severity != nil {
		if sev, err := rules.ParseSeverity(*s.BadPractices.Severity); err == nil {
			out.BPSeverity = &sev
		}
	}
	if len(s.BadPractices.SeverityOverrides) > 0 {
		out.SeverityOverrides = make(map[string]rules.Severity, len(s.BadPractices.SeverityOverrides))
		for id, name := range s.BadPractices.SeverityOverrides {
			if sev, err := rules.ParseSeverity(name); err == nil {
				out.SeverityOverrides[id] = sev
			}
		}
	}

	// export context mode
	if s.Export.WholeFunction != nil {
		out.ExportWholeFunction = s.Export.WholeFunction
	}

	return out, nil
}

func unionStable(a, b []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(a)+len(b))
	for _, xs := range [][]string{a, b} {
		for _, x := range xs {
			x = strings.TrimSpace(x)
			if x == "" {
				continue
			}
			if _, ok := seen[x]; ok {
				continue
			}
			seen[x] = struct{}{}
			out = append(out, x)
		}
	}
	return out
}

func boolPtr(v bool) *bool { return &v }
