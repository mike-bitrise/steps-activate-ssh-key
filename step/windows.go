package step

import (
	"os"
)

// We need the ssh-agent supplied with Git for Windows,
// so we set the PATH variable to include the Git for Windows bin directory.
func init() {
	path := os.Getenv("PATH")
	newPath := `C:\Program Files\Git\usr\bin` + string(os.PathListSeparator) + path
	_ = os.Setenv("PATH", newPath)
}
