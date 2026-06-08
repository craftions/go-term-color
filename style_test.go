package color

import "testing"

func TestMultiplesAtributos(t *testing.T) {
	originalNoColor := NoColor
	NoColor = false
	defer func() { NoColor = originalNoColor }()

	tests := []struct {
		name     string
		color    *Color
		input    string
		expected string
	}{
		{"Cyan Bold Underline", New(FgCyan, Bold, Underline), "test 1", "\x1b[36;1;4mtest 1\x1b[0m"},
		{"BgRed FgWhite", New(BgRed, FgWhite), "test 2", "\x1b[41;37mtest 2\x1b[0m"},
		{"Blink", New(BlinkSlow), "test 3", "\x1b[5mtest 3\x1b[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.color.Sprintf("%s", tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
