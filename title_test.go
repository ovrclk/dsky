package dsky

import (
	"testing"
)

func TestTitle_String(t *testing.T) {
	got := NewTitle("foo").H1().String()
	expect := "foo\n===\n"
	if got != expect {
		t.Fatal("==> expected:\n", expect, "==> got\n", got)
	}
}
