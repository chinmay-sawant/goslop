package app

import (
	"errors"
	"fmt"
)

// Conventional exit codes (parity with Rust CodeHound):
//
//	0 — clean (no failing findings, no errors)
//	1 — findings that violate FailPolicy
//	2 — usage / configuration error
//	3 — internal / I/O / engine error (reserved)
const (
	ExitClean    = 0
	ExitFailing  = 1
	ExitConfig   = 2
	ExitInternal = 3
)

// ErrFindings signals that findings exceeded the fail policy.
// Main should map ExitCodeError{Code: ExitFailing} to process exit 1.
// Typically Err is nil so nothing is printed to stderr after findings stdout.
var ErrFindings = errors.New("findings exceed fail policy")

// ExitCodeError carries a process exit code alongside an optional cause.
//
// Design:
//   - Run returns nil on success (exit 0).
//   - Failing findings: &ExitCodeError{Code: ExitFailing} (Err usually nil).
//   - Usage/config: &ExitCodeError{Code: ExitConfig, Err: cause} or plain error
//     (main defaults plain errors to exit 2).
//   - Internal failures: &ExitCodeError{Code: ExitInternal, Err: cause}.
type ExitCodeError struct {
	Code int
	Err  error
}

func (e *ExitCodeError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	if e.Code == ExitFailing {
		return ErrFindings.Error()
	}
	return fmt.Sprintf("exit code %d", e.Code)
}

func (e *ExitCodeError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// ExitCode extracts a process exit code from err.
// nil → 0; ExitCodeError → its Code; otherwise → ExitConfig (2).
func ExitCode(err error) int {
	if err == nil {
		return ExitClean
	}
	var ece *ExitCodeError
	if errors.As(err, &ece) {
		return ece.Code
	}
	return ExitConfig
}
