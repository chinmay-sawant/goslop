package badpractices_test

import (
	"testing"
)

func TestBPPY43RequirementsWithoutPins(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-43", "BP-PY-43")
	assertBPFixtureCase(t, "BP-PY-43", "BP-PY-43-tilde-pin")
	assertBPFixtureCase(t, "BP-PY-43", "BP-PY-43-directives")
	assertBPFixtureCase(t, "BP-PY-43", "BP-PY-43-non-requirements")
	assertBPFixtureCase(t, "BP-PY-43", "BP-PY-43-dev")
}

func TestBPPY44DeprecatedStdlibImport(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-44", "BP-PY-44")
	assertBPFixtureCase(t, "BP-PY-44", "BP-PY-44-asyncore")
	assertBPFixtureCase(t, "BP-PY-44", "BP-PY-44-cgi")
}

func TestBPPY45SysPathMutation(t *testing.T) {
	t.Parallel()
	assertBPFixtureCase(t, "BP-PY-45", "BP-PY-45")
	assertBPFixtureCase(t, "BP-PY-45", "BP-PY-45-append")
	assertBPFixtureCase(t, "BP-PY-45", "BP-PY-45-test-path")
	assertBPFixtureCase(t, "BP-PY-45", "BP-PY-45-test-file")
	assertBPFixtureCase(t, "BP-PY-45", "BP-PY-45-readonly")
	assertBPFixtureCase(t, "BP-PY-45", "BP-PY-45-docs-conf")
	assertBPFixtureCase(t, "BP-PY-45", "BP-PY-45-file-bootstrap")
	assertBPFixtureCase(t, "BP-PY-45", "BP-PY-45-guarded-bootstrap")
}
