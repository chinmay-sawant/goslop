package cwe

import "github.com/chinmay/goslop/internal/rules"

// metaByID maps rule id -> metadata (filled by RegisterRule).
var metaByID map[string]*rules.RuleMetadata
