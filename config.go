package color

import (
	"os"
)

// Mode represents the color rendering mode
type Mode string

const (
	ModeAuto   Mode = "auto"
	ModeAlways Mode = "always"
	ModeNever  Mode = "never"
)

// Reason represents the internal decision reason
type Reason string

const (
	ReasonAutoTerminal Reason = "auto_terminal"
	ReasonNoColorEnv   Reason = "no_color_env"
	ReasonDumbTerm     Reason = "term_dumb"
	ReasonForced       Reason = "forced"
	ReasonNotTerminal  Reason = "not_terminal"
)

// NoColor determines if the use of colors should be omitted.
var NoColor bool

// CurrentMode allows querying or overriding the operational mode.
// By default, it operates in ModeAuto.
var CurrentMode Mode = ModeAuto

// colorDisabledReason tracks the internal reasoning for the diagnostic tool.
var colorDisabledReason Reason

func init() {
	setupNoColor()
}

// globalDetector is used for testing purposes
var globalDetector TerminalDetector = defaultDetector{}

// setupNoColor initializes the NoColor variable based on CurrentMode and environment.
func setupNoColor() {
	if CurrentMode == ModeNever {
		NoColor = true
		colorDisabledReason = ReasonForced
		return
	}
	if CurrentMode == ModeAlways {
		NoColor = false
		colorDisabledReason = ReasonForced
		return
	}

	if os.Getenv("NO_COLOR") != "" {
		NoColor = true
		colorDisabledReason = ReasonNoColorEnv
		return
	}

	if os.Getenv("TERM") == "dumb" {
		NoColor = true
		colorDisabledReason = ReasonDumbTerm
		return
	}

	if !globalDetector.IsTerminal(os.Stdout.Fd()) {
		NoColor = true
		colorDisabledReason = ReasonNotTerminal
		return
	}

	NoColor = false
	colorDisabledReason = ReasonAutoTerminal
}
