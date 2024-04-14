package list

import (
	"fmt"
)

type List struct {
	Order []int // Order to display

	total   int // Total number of items
	current int // Index of current item
}

type ListFace interface {
	Reset()

	// Properties
	Len() int
	Current() int
	SetCurrent(int) error

	// Basic operations
	Insert(int, int)   // Reference to item index, position to insert at
	Remove(int) int    // Remove item from order, returns item_idx of removed order_idx
	Move(int, int) int // Move item to a position ( remove + insert ), returns item_idx of removed order_idx
}

func New(total int) List {
	m := List{total: total}
	m.Reset()
	return m
}

func (m *List) Reset() {
	m.Order = make([]int, m.total)
	// Initially in order: 0, 1, 2, ... ,total - 1
	for i := 0; i < m.total; i++ {
		m.Order[i] = i
	}
}

func Render[T any](m *List, items *[]T) string {
	out := "\n"
	for i, v := range m.Order {
		ind := "-> "
		if i == m.current {
			ind = "==>"
		}
		out += fmt.Sprintf("%d %s %d: %v\n", i, ind, v, (*items)[v])
	}
	return out
}

func (m List) Len() int {
	return m.total
}

func (m List) Current() int {
	return m.current
}

func (m *List) SetCurrent(i int) error {
	if (i > m.total) || (i < 0) {
		return fmt.Errorf("Out of bounds: i - %d ; len: %d", i, m.total)
	}
	m.current = i
	return nil
}

func (m *List) Insert(item_idx int, order_idx int) {
	// TODO: Check order bounds
	m.Order = append(m.Order[:order_idx], append([]int{item_idx}, m.Order[order_idx:]...)...)
	if m.current > order_idx {
		m.current += 1
	}
}

func (m *List) Remove(order_idx int) (item_idx int) {
	// TODO: Check order bounds
	item_idx = m.Order[order_idx]
	m.Order = append(m.Order[:order_idx], m.Order[order_idx+1:]...)
	if m.current > order_idx {
		m.current -= 1
	}
	return item_idx
}

func (m *List) Move(src_idx int, dst_idx int) (item_idx int) {
	// TODO: Check order bounds
	item_idx = m.Remove(src_idx)
	m.Insert(item_idx, dst_idx)
	return item_idx
}
