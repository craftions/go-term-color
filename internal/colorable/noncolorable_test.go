package colorable

import (
	"bytes"
	"errors"
	"testing"
)

func TestNonColorable_PassthroughWrite(t *testing.T) {
	var dst bytes.Buffer
	w := NewNonColorable(&dst)
	_, err := w.Write([]byte("abc"))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dst.String() != "abc" {
		t.Fatalf("expected passthrough write, got %q", dst.String())
	}
}

func TestNonColorable_StripANSI(t *testing.T) {
	var dst bytes.Buffer
	w := NewNonColorable(&dst)
	_, err := w.Write([]byte("\x1b[31mred\x1b[0m"))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dst.String() != "red" {
		t.Fatalf("expected stripped string, got %q", dst.String())
	}
}

func TestNonColorable_IncompleteESC(t *testing.T) {
	var dst bytes.Buffer
	w := NewNonColorable(&dst)
	_, _ = w.Write([]byte("foo\x1b"))
	if dst.String() != "foo" {
		t.Fatalf("expected foo, got %q", dst.String())
	}
}

func TestNonColorable_NonBracketESC(t *testing.T) {
	var dst bytes.Buffer
	w := NewNonColorable(&dst)
	_, _ = w.Write([]byte("foo\x1bXbar"))
	if dst.String() != "foobar" {
		t.Fatalf("expected foobar, got %q", dst.String())
	}
}

func TestNonColorable_IncompleteANSI(t *testing.T) {
	var dst bytes.Buffer
	w := NewNonColorable(&dst)
	_, _ = w.Write([]byte("foo\x1b["))
	if dst.String() != "foo" {
		t.Fatalf("expected foo, got %q", dst.String())
	}
}

type errWriter struct{}

func (errWriter) Write(_ []byte) (int, error) { return 0, errors.New("write failed") }

func TestNonColorable_WriterError(t *testing.T) {
	w := NewNonColorable(errWriter{})
	// This will trigger plaintext.WriteTo(w.out) error
	_, _ = w.Write([]byte("foo\x1b[m"))
}
