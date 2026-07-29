package rules

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	fingerprintTool    = "codehound"
	fingerprintVersion = 2
)

// FingerprintV2 returns codehound:2:{rule}:{file}:{msg_hash16}.
// The file path is normalized to forward slashes.
func FingerprintV2(ruleID, file, message string) string {
	sum := sha256.Sum256([]byte(message))
	msgHash := hex.EncodeToString(sum[:])
	if len(msgHash) > 16 {
		msgHash = msgHash[:16]
	}
	norm := strings.ReplaceAll(file, "\\", "/")
	return fmt.Sprintf("%s:%d:%s:%s:%s", fingerprintTool, fingerprintVersion, ruleID, norm, msgHash)
}

// Fingerprint is an alias for FingerprintV2.
func Fingerprint(ruleID, file, message string) string {
	return FingerprintV2(ruleID, file, message)
}
