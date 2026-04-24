package color

import (
	"os"
	"testing"
)

func TestHelperStrings(t *testing.T) {
	originalNoColor := NoColor
	NoColor = false
	defer func() { NoColor = originalNoColor }()

	tests := []struct {
		name     string
		fn       func(string, ...interface{}) string
		format   string
		args     []interface{}
		expected string
	}{
		{"RedString", RedString, "red %d", []interface{}{1}, "\x1b[31mred 1\x1b[0m"},
		{"GreenString", GreenString, "green", nil, "\x1b[32mgreen\x1b[0m"},
		{"YellowString", YellowString, "yellow", nil, "\x1b[33myellow\x1b[0m"},
		{"BlueString", BlueString, "blue", nil, "\x1b[34mblue\x1b[0m"},
		{"MagentaString", MagentaString, "magenta", nil, "\x1b[35mmagenta\x1b[0m"},
		{"CyanString", CyanString, "cyan", nil, "\x1b[36mcyan\x1b[0m"},
		{"WhiteString", WhiteString, "white", nil, "\x1b[37mwhite\x1b[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.format, tt.args...)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNoColorActivado(t *testing.T) {
	tests := []struct {
		name    string
		noColor bool
		fn      func(string, ...interface{}) string
		input   string
		want    string
	}{
		{"NoColor True", true, BlueString, "test blue", "test blue"},
		{"NoColor False", false, BlueString, "test blue", "\x1b[34mtest blue\x1b[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalNoColor := NoColor
			NoColor = tt.noColor
			defer func() { NoColor = originalNoColor }()

			result := tt.fn(tt.input)
			if result != tt.want {
				t.Errorf("NoColor=%v: Expected %q, got %q", tt.noColor, tt.want, result)
			}
		})
	}
}

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
			result := tt.color.wrap("%s", tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestEnvNoColor(t *testing.T) {
	originalEnv := os.Getenv("NO_COLOR")
	defer func() {
		if originalEnv == "" {
			os.Unsetenv("NO_COLOR")
		} else {
			os.Setenv("NO_COLOR", originalEnv)
		}
	}()

	tests := []struct {
		name       string
		envValue   string
		shouldBind bool
		want       string
	}{
		{"Env Set", "1", true, "env_test"},
		{"Env Empty", "", false, "\x1b[32menv_test\x1b[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", tt.envValue)
			}

			oldNoColor := NoColor
			NoColor = os.Getenv("NO_COLOR") != ""
			defer func() { NoColor = oldNoColor }()

			s := GreenString("env_test")
			if s != tt.want {
				t.Errorf("NO_COLOR env=%q: Expected %q, got %q", tt.envValue, tt.want, s)
			}
		})
	}
}

func TestMetodosImpresion(t *testing.T) {
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	tests := []struct {
		name    string
		noColor bool
	}{
		{"With Colors", false},
		{"Without Colors", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NoColor = tt.noColor

			// Test printing methods don't panic
			Red("Direct function test (red)")
			Green("Direct function test (green)")
			Yellow("Direct function test (yellow)")
			Blue("Direct function test (blue)")
			Magenta("Direct function test (magenta)")
			Cyan("Direct function test (cyan)")
			White("Direct function test (white)")

			color := New(FgRed, BlinkSlow)
			color.Print("Object print test")
			color.Println("Object println test")
			color.Printf("Object printf test %d\n", 1)
		})
	}
}

func TestCheckIfTerminal(t *testing.T) {
	// Just to test that it does not panic and returns a bool
	result := CheckIfTerminal(os.Stdout.Fd())
	_ = result

	// We test with an invalid FD to force it to return false
	result2 := CheckIfTerminal(^uintptr(0))
	if result2 {
		t.Errorf("expected false for dummy fd")
	}
}

func TestSetupNoColor(t *testing.T) {
	// Save environment variables and state
	oldNoColor := NoColor
	oldEnv := os.Getenv("TERM")
	oldNoColorEnv := os.Getenv("NO_COLOR")
	defer func() {
		NoColor = oldNoColor
		os.Setenv("TERM", oldEnv)
		if oldNoColorEnv == "" {
			os.Unsetenv("NO_COLOR")
		} else {
			os.Setenv("NO_COLOR", oldNoColorEnv)
		}
	}()

	// Test 1: TERM=dumb
	os.Setenv("TERM", "dumb")
	os.Unsetenv("NO_COLOR")
	setupNoColor()
	if !NoColor {
		t.Errorf("Expected NoColor=true when TERM=dumb")
	}

	// Test 2: NO_COLOR=1
	os.Setenv("TERM", "")
	os.Setenv("NO_COLOR", "1")
	setupNoColor()
	if !NoColor {
		t.Errorf("Expected NoColor=true when NO_COLOR=1")
	}

	// Test 3: Default env
	os.Unsetenv("NO_COLOR")
	os.Unsetenv("TERM")
	setupNoColor()
}
