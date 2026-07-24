package adapters

import (
	"errors"
	"testing"
)

func TestWrapErr_NilError_ReturnsNil(t *testing.T) {
	got := wrapErr(nil, "some operation", 42)
	if got != nil {
		t.Errorf("wrapErr(nil, ...) = %v, want nil", got)
	}
}

func TestWrapErr_WithEntityID(t *testing.T) {
	original := errors.New("connection refused")
	got := wrapErr(original, "updating product", 7)

	if got == nil {
		t.Fatal("wrapErr returned nil for non-nil error")
	}

	want := "updating product (id=7): connection refused"
	if got.Error() != want {
		t.Errorf("got %q, want %q", got.Error(), want)
	}

	if !errors.Is(got, original) {
		t.Error("wrapped error does not unwrap to original")
	}
}

func TestWrapErr_ZeroEntityID(t *testing.T) {
	original := errors.New("timeout")
	got := wrapErr(original, "getting config", 0)

	if got == nil {
		t.Fatal("wrapErr returned nil for non-nil error")
	}

	want := "getting config: timeout"
	if got.Error() != want {
		t.Errorf("got %q, want %q", got.Error(), want)
	}

	if !errors.Is(got, original) {
		t.Error("wrapped error does not unwrap to original")
	}
}

func TestWrapErr_NegativeEntityID(t *testing.T) {
	original := errors.New("not found")
	got := wrapErr(original, "finding user", -1)

	if got == nil {
		t.Fatal("wrapErr returned nil for non-nil error")
	}

	// Negative IDs should not include id context (only > 0 triggers it)
	want := "finding user: not found"
	if got.Error() != want {
		t.Errorf("got %q, want %q", got.Error(), want)
	}
}
