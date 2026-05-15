package color

import (
	"github.com/craftions/go-term-check/terminal"
)

// TerminalDetector abstracts the terminal check logic
type TerminalDetector interface {
	IsTerminal(fd uintptr) bool
}

type defaultDetector struct{}

func (d defaultDetector) IsTerminal(fd uintptr) bool {
	return terminal.Check(fd).Terminal
}
