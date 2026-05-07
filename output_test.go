package color

import (
	"bytes"
	"os"
	"testing"
)

func TestHelperStrings(t *testing.T) {
	originalNoColor := NoColor
	NoColor = false
	defer func() { NoColor = originalNoColor }()

	tests := []struct {
		name     string
		fn       func(string, ...any) string
		format   string
		args     []any
		expected string
	}{
		{"RedString", RedString, "red %d", []any{1}, "\x1b[31mred 1\x1b[0m"},
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

func TestFprintMethods(t *testing.T) {
	originalNoColor := NoColor
	NoColor = false
	defer func() { NoColor = originalNoColor }()

	color := New(FgRed)

	tests := []struct {
		name     string
		fn       func(w *bytes.Buffer)
		expected string
	}{
		{
			"Fprint",
			func(w *bytes.Buffer) { color.Fprint(w, "test") },
			"\x1b[31mtest\x1b[0m",
		},
		{
			"Fprintf",
			func(w *bytes.Buffer) { color.Fprintf(w, "test %d", 1) },
			"\x1b[31mtest 1\x1b[0m",
		},
		{
			"Fprintln",
			func(w *bytes.Buffer) { color.Fprintln(w, "test") },
			"\x1b[31mtest\x1b[0m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.fn(&buf)
			if got := buf.String(); got != tt.expected {
				t.Errorf("%s: Expected %q, got %q", tt.name, tt.expected, got)
			}
		})
	}
}

func TestSprintVsSprintf(t *testing.T) {
	originalNoColor := NoColor
	NoColor = false
	defer func() { NoColor = originalNoColor }()

	color := New(FgBlue)

	gotSprintf := color.Sprintf("%d%%", 100)
	wantSprintf := "\x1b[34m100%\x1b[0m"
	if gotSprintf != wantSprintf {
		t.Errorf("Sprintf falló: Expected %q, got %q", wantSprintf, gotSprintf)
	}

	gotSprint := color.Sprint("100%")
	wantSprint := "\x1b[34m100%\x1b[0m"
	if gotSprint != wantSprint {
		t.Errorf("Sprint falló interpretando literalmente: Expected %q, got %q", wantSprint, gotSprint)
	}

	literalStr := "Progreso: %d"
	gotSprintLiteral := color.Sprint(literalStr)
	wantSprintLiteral := "\x1b[34mProgreso: %d\x1b[0m"
	if gotSprintLiteral != wantSprintLiteral {
		t.Errorf("Sprint interpretó el comodín equivocadamente: Expected %q, got %q", wantSprintLiteral, gotSprintLiteral)
	}
}

func TestMetodosImpresionGlobales(t *testing.T) {
	originalStdout := os.Stdout
	originalNoColor := NoColor
	defer func() {
		os.Stdout = originalStdout
		NoColor = originalNoColor
	}()

	NoColor = false

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	Red("red line")
	Green("green line")
	Yellow("yellow line")
	Blue("blue line")
	Magenta("magenta line")
	Cyan("cyan line")
	White("white line")

	color := New(FgRed)
	color.Print("print test")
	color.Println("println test")
	color.Printf("printf %d", 1)

	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !bytes.Contains(buf.Bytes(), []byte("\x1b[31mred line\n\x1b[0m")) {
		t.Errorf("Falta 'red line' en el output: %q", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("\x1b[32mgreen line\n\x1b[0m")) {
		t.Errorf("Falta 'green line' en el output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("\x1b[31mprint test\x1b[0m")) {
		t.Errorf("Falta 'print test' en el output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("\x1b[31mprintln test\x1b[0m\n")) {
		t.Errorf("Falta 'println test' en el output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("\x1b[31mprintf 1\x1b[0m")) {
		t.Errorf("Falta 'printf 1' en el output")
	}
}
