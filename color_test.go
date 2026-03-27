package color

import (
	"fmt"
	"os"
	"testing"
)

func TestRedString(t *testing.T) {
	originalNoColor := NoColor
	NoColor = false // Forzamos a aplicar colores para test
	defer func() { NoColor = originalNoColor }()

	s := RedString("test rojo")
	expected := "\x1b[31mtest rojo\x1b[0m"
	if s != expected {
		t.Errorf("Esperaba %q, obtuve %q", expected, s)
	}

	mg := MagentaString("magenta")
	expectedMg := "\x1b[35mmagenta\x1b[0m"
	if mg != expectedMg {
		t.Errorf("Esperaba %q para magenta, obtuve %q", expectedMg, mg)
	}
}

func TestNoColorActivado(t *testing.T) {
	originalNoColor := NoColor
	NoColor = true // Simulamos que estamos en un archivo / no-tty
	defer func() { NoColor = originalNoColor }()

	s := BlueString("test azul")
	expected := "test azul"
	if s != expected {
		t.Errorf("Se esperaba que no tuviera atributos ANSI, obtuve %q", s)
	}
}

func TestMultiplesAtributos(t *testing.T) {
	originalNoColor := NoColor
	NoColor = false
	defer func() { NoColor = originalNoColor }()

	c := New(FgCyan, Bold, Underline)
	s := c.wrap("cyan bold underline")
	// 36 es cyan, 1 es bold, 4 es underline
	// el orden de .format() es igual al de añadido
	expected := "\x1b[36;1;4mcyan bold underline\x1b[0m"
	if s != expected {
		t.Errorf("Esperaba múltiples atributos %q, obtuve %q", expected, s)
	}
}

func TestEnvNoColor(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	oldNoColor := NoColor
	NoColor = os.Getenv("NO_COLOR") != ""
	defer func() { NoColor = oldNoColor }()

	s := GreenString("env_test")
	if s != "env_test" {
		t.Errorf("NO_COLOR no está siendo respetado, obtuve %q", s)
	}
}

func TestMetodosImpresion(t *testing.T) {
	// Testeamos que no haya panic
	originalNoColor := NoColor
	NoColor = false
	defer func() { NoColor = originalNoColor }()

	fmt.Print("\n--- Salida de pruebas locales ---\n")
	Yellow("Direct function test (debe ser amarillo)")
	New(FgRed, BlinkSlow).Println("Objeto color test (rojo parpadeante lento)")
	fmt.Print("--- Fin pruebas locales ---\n")
}
