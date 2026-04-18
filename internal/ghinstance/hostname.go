package ghinstance

import (
	"fmt"
)

const (
	// DefaultHostname is the default GitHub hostname
	DefaultHostname = "github.com"
	// DefaultAlternativeHostnames are alternative GitHub hostnames
	DefaultAlternativeHostnames = "gist.github.com,github.com"
)

// NormalizeHostname normalizes a GitHub hostname
func NormalizeHostname(hostname string) string {
	if hostname == "" {
		return DefaultHostname
	}
	return hostname
}

// IsGitHub returns true if the hostname is a GitHub hostname
func IsGitHub(hostname string) bool {
	hostname = NormalizeHostname(hostname)
	return hostname == DefaultHostname
}

// Hostname returns the normalized hostname
func Hostname(hostname string) string {
	return NormalizeHostname(hostname)
}

// String returns a string representation of the hostname
func String(hostname string) string {
	return NormalizeHostname(hostname)
}

// Format formats a hostname for display
func Format(hostname string) string {
	return NormalizeHostname(hostname)
}

// New returns a new hostname
func New(hostname string) string {
	return NormalizeHostname(hostname)
}

// Validate validates a hostname
func Validate(hostname string) error {
	hostname = NormalizeHostname(hostname)
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}
	return nil
}