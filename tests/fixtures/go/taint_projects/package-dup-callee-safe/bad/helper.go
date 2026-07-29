package bad

import "os"

// openPath is a sink in this package only.
func openPath(s string) {
	os.Open(s)
}
