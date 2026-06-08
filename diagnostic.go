package color

import (
	"fmt"
	"io"
	"runtime"
)

// Diagnostic provides insights into the operational state of the color engine.
type Diagnostic struct {
	GOOS         string
	GOARCH       string
	ColorEnabled bool
	Mode         Mode
	Reason       Reason
}

// Diagnose evaluates and returns the current operational state.
func Diagnose() Diagnostic {
	return Diagnostic{
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		ColorEnabled: !NoColor,
		Mode:         CurrentMode,
		Reason:       colorDisabledReason,
	}
}

// PrintDiagnostic writes the diagnostic state to the provided writer in plain text.
func PrintDiagnostic(w io.Writer) error {
	d := Diagnose()
	_, err := fmt.Fprintf(w, "--- Go Term Color Diagnostic ---\nOS: %s\nARCH: %s\nColor Enabled: %t\nMode: %s\nReason: %s\n--------------------------------\n",
		d.GOOS, d.GOARCH, d.ColorEnabled, d.Mode, d.Reason)
	return err
}
