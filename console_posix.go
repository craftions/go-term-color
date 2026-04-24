//go:build !windows
// +build !windows

package color

func init() {
	// POSIX systems (Linux, macOS, Solaris) natively operate with ANSI escape sequences
	// therefore, no additional OS kernel/console configuration is required here.
}
