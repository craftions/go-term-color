package color

import (
	"os"
	"testing"
)

func TestNoColorActivado(t *testing.T) {
	tests := []struct {
		name    string
		noColor bool
		fn      func(string, ...any) string
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



func TestSetupNoColorModes(t *testing.T) {
	oldNoColor := NoColor
	oldMode := CurrentMode
	oldEnv := os.Getenv("TERM")
	oldNoColorEnv := os.Getenv("NO_COLOR")
	oldDetector := globalDetector
	defer func() {
		NoColor = oldNoColor
		CurrentMode = oldMode
		globalDetector = oldDetector
		os.Setenv("TERM", oldEnv)
		if oldNoColorEnv == "" {
			os.Unsetenv("NO_COLOR")
		} else {
			os.Setenv("NO_COLOR", oldNoColorEnv)
		}
	}()

	tests := []struct {
		name        string
		mode        Mode
		envTerm     string
		envNoColor  string
		isTerminal  bool
		wantNoColor bool
		wantReason  Reason
	}{
		{"Mode Never forces NoColor=true", ModeNever, "", "", true, true, ReasonForced},
		{"Mode Always forces NoColor=false", ModeAlways, "dumb", "1", true, false, ReasonForced},
		{"Mode Auto, TERM=dumb", ModeAuto, "dumb", "", true, true, ReasonDumbTerm},
		{"Mode Auto, NO_COLOR=1", ModeAuto, "", "1", true, true, ReasonNoColorEnv},
		{"Mode Auto, Not Terminal", ModeAuto, "", "", false, true, ReasonNotTerminal},
		{"Mode Auto, Is Terminal", ModeAuto, "", "", true, false, ReasonAutoTerminal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CurrentMode = tt.mode
			globalDetector = fakeTerminalDetector{isTerminal: tt.isTerminal}
			os.Setenv("TERM", tt.envTerm)
			if tt.envNoColor == "" {
				os.Unsetenv("NO_COLOR")
			} else {
				os.Setenv("NO_COLOR", tt.envNoColor)
			}
			setupNoColor()
			if NoColor != tt.wantNoColor {
				t.Errorf("Expected NoColor=%v, got %v", tt.wantNoColor, NoColor)
			}
			if colorDisabledReason != tt.wantReason {
				t.Errorf("Expected Reason=%v, got %v", tt.wantReason, colorDisabledReason)
			}
		})
	}
}
func TestResolveMode_NoColorEnvWins(t *testing.T) {
	env := map[string]string{"NO_COLOR": "1", "TERM": "xterm"}
	got := resolveMode(env, true, ModeAuto)
	if got != ModeNever {
		t.Fatalf("expected ModeNever, got %v", got)
	}
}

func TestResolveMode_TermDumbDisablesColor(t *testing.T) {
	env := map[string]string{"TERM": "dumb"}
	got := resolveMode(env, true, ModeAuto)
	if got != ModeNever {
		t.Fatalf("expected ModeNever, got %v", got)
	}
}
