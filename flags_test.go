package flags

import (
	"reflect"
	"strings"
	"testing"
)

func same(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func equals(t *testing.T, a, b interface{}) {
	t.Helper()
	if !same(a, b) {
		t.Errorf("\nexpected: %v\n  actual: %v\nto be equal", a, b)
	}
}

func differs(t *testing.T, a, b interface{}) {
	t.Helper()
	if same(a, b) {
		t.Errorf("\nexpected: %v\n  actual: %v\nto be different", a, b)
	}
}

func panics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("expected to panic")
		}
	}()
	f()
}

func TestPositional(t *testing.T) {
	pos := newPositional()
	equals(t, pos.Len(), 0)

	parser := NewParser(pos, nil)
	s := "true 42 foo"
	args := strings.Split(s, " ")

	arg0 := pos.Bool("bool", "boolean value")
	equals(t, pos.Len(), 1)

	arg1 := pos.Int("int", "integer value")
	equals(t, pos.Len(), 2)

	arg2 := pos.String("string", "string value")
	equals(t, pos.Len(), 3)

	if err := parser.Parse(args); err != nil {
		t.Errorf("parser.Parse: %v", err)
		return
	}

	equals(t, *arg0, true)
	equals(t, *arg1, 42)
	equals(t, *arg2, "foo")

	panics(t, func() { pos.Bool("bool", "boolean value") })
	panics(t, func() { pos.Int("int", "integer value") })
	panics(t, func() { pos.String("string", "string value") })

	if err := parser.Parse(nil); err == nil {
		t.Error("parser.Parse(nil) = nil, want error")
		return
	}

	if err := parser.Parse([]string{"foo"}); err == nil {
		t.Error("parser.Parse([]string{\"foo\"}}) = nil, want error")
		return
	}

	if err := parser.Parse([]string{"true", "42"}); err == nil {
		t.Error("parser.Parse([]string{\"true\", \"42\"}}) = nil, want error")
		return
	}
}
