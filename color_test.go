package color

import (
	"strings"
	"testing"
)

func TestSprint_PreservesPercentLiteral(t *testing.T) {
	c := New(FgGreen)
	out := c.Sprint("100%")
	if !strings.Contains(out, "100%") {
		t.Fatalf("expected literal percent, got %q", out)
	}
}

func TestSprintf_FormatsVerb(t *testing.T) {
	c := New(FgGreen)
	out := c.Sprintf("%d", 7)
	if !strings.Contains(out, "7") {
		t.Fatalf("expected formatted number, got %q", out)
	}
}
