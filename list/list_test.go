package list

import (
	"testing"
)

func TestRender(t *testing.T) {
	// Immutable items
	items := []string{"a", "b", "c"}
	total := len(items)

	// List
	m := New(total)
	if m.Len() != total {
		t.Fatalf(`Length messed up: %d, %v`, total, m.Order)
	}

	// Print to check order
	out := Render(&m, &items)
	want := `
0 ==> 0: a
1 ->  1: b
2 ->  2: c
`
	if out != want {
		t.Fatalf("\noutput:\n%s\n---\nwant:\n%s\n", out, want)
	}
}

func TestCurrent(t *testing.T) {
	// Immutable items
	items := []string{"a", "b", "c"}
	total := len(items)

	// Default current
	m := New(total)
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
0 ->  0: a
1 ==> 1: b
2 ->  2: c
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

func TestInsert(t *testing.T) {
	// Immutable items
	items := []string{"a", "b", "c"}
	total := len(items)

	// List
	m := New(total)

	// Modify items and list separately
	items = append(items, "d")
	m.Insert(3, 1)

	// Check
	out := Render(&m, &items)
	want := `
0 ==> 0: a
1 ->  3: d
2 ->  1: b
3 ->  2: c
`
	if out != want {
		t.Fatalf("\noutput:\n%s\n---\nwant:\n%s\n", out, want)
	}
}

func TestRemove(t *testing.T) {
	// Immutable items
	items := []string{"a", "b", "c"}
	total := len(items)

	// List
	m := New(total)

	// Remove from view, not from items
	m.Remove(1)

	// Check
	out := Render(&m, &items)
	want := `
0 ==> 0: a
1 ->  2: c
`
	if out != want {
		t.Fatalf("\noutput:\n%s\n---\nwant:\n%s\n", out, want)
	}
}

func TestMove(t *testing.T) {
	// Immutable items
	items := []string{"a", "b", "c"}
	total := len(items)

	// List
	m := New(total)

	// Move, shortcut to remove and insert
	m.Move(1, 2)

	// Check
	out := Render(&m, &items)
	want := `
0 ==> 0: a
1 ->  2: c
2 ->  1: b
`
	if out != want {
		t.Fatalf("\noutput:\n%s\n---\nwant:\n%s\n", out, want)
	}
}
