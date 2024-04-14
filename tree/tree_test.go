package tree

import (
	"testing"
)

type node struct {
	key      Key
	children []Key
}

func (n node) Key() Key {
	return n.key
}

func (n node) Children() []Key {
	return n.children
}

func TestRender(t *testing.T) {
	// Immutable items
	items := []Item{
		node{"a", []Key{}},
		node{"b", []Key{}},
		node{"c", []Key{}},
	}
	total := len(items)

	// List
	m := New(&items)
	if m.Len() != total {
		t.Fatalf(`Length messed up: %d, %v`, total, m.Order)
	}

	// Print to check order
	out := Render(&m, &items)
	want := `
0 ==> 0:  a
1 ->  1:  b
2 ->  2:  c
`
	if out != want {
		t.Fatalf("\noutput:\n%s\n---\nwant:\n%s\n", out, want)
	}
}

func TestLevels(t *testing.T) {
	// Immutable items
	items := []Item{
		node{"b", []Key{}},
		node{"a", []Key{"b"}},
		node{"c", []Key{}},
	}
	total := len(items)

	// List
	m := New(&items)
	if m.Len() != total {
		t.Fatalf(`Length messed up: %d, %v`, total, m.Order)
	}

	// Print to check order
	out := Render(&m, &items)
	want := `
0 ==> 1:  a
1 ->  0:    b
2 ->  2:  c
`
	if out != want {
		t.Fatalf("\noutput:\n%s\n---\nwant:\n%s\n", out, want)
	}
}

func TestMissingChildren(t *testing.T) {
	// Immutable items
	items := []Item{
		node{"a", []Key{"b"}},
		node{"c", []Key{"d"}},
	}
	total := len(items)

	// List
	m := New(&items)
	if m.Len() != total {
		t.Fatalf(`Length messed up: %d, %v`, total, m.Order)
	}

	// Print to check order
	out := Render(&m, &items)
	want := `
0 ==> 0:  a
1 ->  1:  c
`
	if out != want {
		t.Fatalf("\noutput:\n%s\n---\nwant:\n%s\n", out, want)
	}
}

func TestCurrent(t *testing.T) {
	// Immutable items
	items := []Item{
		node{"b", []Key{}},
		node{"a", []Key{"b"}},
		node{"c", []Key{}},
	}

	// Default current
	m := New(&items)
	if m.Current() != 0 {
		t.Fatalf(`Current messed up: %d`, m.current)
	}

	// Within bounds
	m.SetCurrent(1)
	if m.Current() != 1 {
		t.Fatalf(`Current messed up: %d`, m.current)
	}

	// Print to check order
	out := Render(&m, &items)
	want := `
0 ->  1:  a
1 ==> 0:    b
2 ->  2:  c
`
	if out != want {
		t.Fatalf("\noutput:\n%s\n---\nwant:\n%s\n", out, want)
	}

	// Out of bounds
	err := m.SetCurrent(5)
	if err == nil {
		t.Fatalf(`Expected error: %s`, err)
	}
}
