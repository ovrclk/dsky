package dsky

import (
	"reflect"
	"testing"
)

func TestSectionRows(t *testing.T) {
	d := NewSectionData("id").
		Add("a", "a1").
		Add("a", "a2").
		Add("a", "a3").
		Add("b", "b1").
		Add("b", "b2").
		Add("b", "b3")
	got := d.Rows()
	want := [][]interface{}{
		{"a1", "b1"},
		{"a2", "b2"},
		{"a3", "b3"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, actual  %v", want, got)
	}

	d = NewSectionData("id").
		Add("a", "a1").
		Add("a", "a2").
		Add("b", "b1").
		Add("b", "b2").
		Add("b", "b3")
	got = d.Rows()
	want = [][]interface{}{
		{"a1", "b1"},
		{"a2", "b2"},
		{nil, "b3"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, actual  %v", want, got)
	}

}
