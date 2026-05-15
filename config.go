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

// resolveMode allows tests to determine NoColor properly
func resolveMode(env map[string]string, isTerminal bool, current Mode) Mode {
	if current == ModeNever {
		return ModeNever
	}
	if current == ModeAlways {
		return ModeAlways
	}
	if val, ok := env["NO_COLOR"]; ok && val != "" {
		return ModeNever
	}
	if val, ok := env["TERM"]; ok && val == "dumb" {
		return ModeNever
	}
	if !isTerminal {
		return ModeNever
	}
	return ModeAuto
}

// setupNoColor initializes the NoColor variable based on CurrentMode and environment.
func setupNoColor() {
	env := map[string]string{
		"NO_COLOR": os.Getenv("NO_COLOR"),
		"TERM":     os.Getenv("TERM"),
	}
	isTerminal := globalDetector.IsTerminal(os.Stdout.Fd())
	mode := resolveMode(env, isTerminal, CurrentMode)

	if mode == ModeNever {
		NoColor = true
		if CurrentMode == ModeNever {
			colorDisabledReason = ReasonForced
		} else if env["NO_COLOR"] != "" {
			colorDisabledReason = ReasonNoColorEnv
		} else if env["TERM"] == "dumb" {
			colorDisabledReason = ReasonDumbTerm
		} else if !isTerminal {
			colorDisabledReason = ReasonNotTerminal
		}
		return
	}

	if mode == ModeAlways {
		NoColor = false
		colorDisabledReason = ReasonForced
		return
	}

	NoColor = false
	colorDisabledReason = ReasonAutoTerminal
}
