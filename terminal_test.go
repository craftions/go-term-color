package color

import (
	"os"
	"testing"
)

type fakeTerminalDetector struct {
	isTerminal bool
}

func (f fakeTerminalDetector) IsTerminal(fd uintptr) bool {
	return f.isTerminal
}

func TestTerminalDetector(t *testing.T) {
	tests := []struct {
		name        string
		isTerminal  bool
		wantNoColor bool
	}{
		{"Is a Terminal", true, false},
		{"Not a Terminal", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := fakeTerminalDetector{isTerminal: tt.isTerminal}
			NoColor = !detector.IsTerminal(os.Stdout.Fd())
			if NoColor != tt.wantNoColor {
				t.Errorf("Expected NoColor=%v, got %v", tt.wantNoColor, NoColor)
			}
		})
	}
}
