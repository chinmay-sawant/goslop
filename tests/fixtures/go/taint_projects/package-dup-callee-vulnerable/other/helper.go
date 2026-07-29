package other

// openPath exists only to prove bare-name collision across packages does not
// suppress or steal the same-package sink summary in package app.
func openPath(s string) {
	_ = s
}
