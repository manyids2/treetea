package tree

import (
	"testing"
)

type node struct {
	key      string
	children []string
}

func (n node) string() string {
	return n.key
}

func (n node) Key() string {
	return n.key
}

func (n node) Children() []string {
	return n.children
}

func TestRender(t *testing.T) {
	// Immutable items
	items := []Item{
		node{"a", []string{}},
		node{"b", []string{}},
		node{"c", []string{}},
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
		node{"b", []string{}},
		node{"a", []string{"b"}},
		node{"c", []string{}},
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
		node{"a", []string{"b"}},
		node{"c", []string{"d"}},
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
		node{"b", []string{}},
		node{"a", []string{"b"}},
		node{"c", []string{}},
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
