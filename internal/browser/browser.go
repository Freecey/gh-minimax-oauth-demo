package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Browser represents a web browser interface
type Browser struct{}

// Open opens a URL in the default browser
func (b *Browser) Open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // linux
		cmd = "xdg-open"
	}
	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}

// NewBrowser creates a new browser instance
func NewBrowser() *Browser {
	return &Browser{}
}

// BrowseURL opens a URL in the default browser
func BrowseURL(url string) error {
	b := NewBrowser()
	return b.Open(url)
}

// OpenURL opens a URL (alias for BrowseURL)
func OpenURL(url string) error {
	return BrowseURL(url)
}