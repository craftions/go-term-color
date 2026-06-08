package color

import (
	"bytes"
	"strings"
	"testing"
)

func TestDiagnostic(t *testing.T) {
	oldNoColor := NoColor
	oldReason := colorDisabledReason
	oldMode := CurrentMode
	defer func() {
		NoColor = oldNoColor
		colorDisabledReason = oldReason
		CurrentMode = oldMode
	}()

	NoColor = true
	CurrentMode = ModeAuto
	colorDisabledReason = ReasonNotTerminal

	diag := Diagnose()
	if diag.ColorEnabled != false {
		t.Errorf("Expected ColorEnabled false, got %v", diag.ColorEnabled)
	}
	if diag.Reason != ReasonNotTerminal {
		t.Errorf("Expected Reason %v, got %v", ReasonNotTerminal, diag.Reason)
	}

	var buf bytes.Buffer
	err := PrintDiagnostic(&buf)
	if err != nil {
		t.Errorf("PrintDiagnostic returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Color Enabled: false") {
		t.Errorf("Output does not contain Color Enabled: false. Output: %s", output)
	}
	if !strings.Contains(output, "Reason: not_terminal") {
		t.Errorf("Output does not contain Reason: not_terminal. Output: %s", output)
	}
}
