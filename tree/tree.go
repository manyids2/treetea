package tree

import (
	"fmt"
	"strings"
)

type Key any

var zeroKey Key

type Item interface {
	Key() Key
	Children() []Key
}

func GetIndex(key Key, items *[]Item) int {
	for i, item := range *items {
		if item.Key() == key {
			return i
		}
	}
	return -1
}

type Tree struct {
	Order []int // Order to display

	Keys    []Key // Keys to search in
	Levels  []int // Nesting level for indent
	Parents []Key // To walk from leaf to root

	total   int // Total number of items
	current int // Index of current item
}

type TreeFace interface {
	Reset()

	// Properties
	Len() int
	Current() int
	SetCurrent(int) error

	// Basic operations - maps to siblings to keep list behaviour
	Insert(int, int) // Reference to item index, position to insert at as sibling
	Remove(int)      // Remove only requested node
	Move(int, int)   // Move node to position as sibling

	// Extension to tree
	InsertChild(int, int)       // Reference to item index, position to insert at as child
	RemoveTree(int) []int       // Remove node and children, should return list of affected nodes
	MoveToChild(int, int)       // Move to position as child
	MoveTreeToChild(int, int)   // Move entire tree to position as child
	MoveTreeToSibling(int, int) // Move entire tree to position as child
}

func New(items *[]Item) Tree {
	m := Tree{total: len(*items)}
	m.Reset(items)
	return m
}

func (m *Tree) indexLevels(i int, level int, items *[]Item) {
	m.Levels[i] = level
	for _, c := range (*items)[i].Children() {
		if idx := GetIndex(c, items); idx != -1 {
			m.indexLevels(idx, level+1, items)
		}
	}
}

func (m *Tree) indexOrder(i int, items *[]Item) {
	m.Order = append(m.Order, i)
	for _, c := range (*items)[i].Children() {
		if idx := GetIndex(c, items); idx != -1 {
			m.indexOrder(idx, items)
		}
	}
}

func (m *Tree) Reset(items *[]Item) {
	m.total = len(*items)
	m.Keys = make([]Key, m.total)
	m.Levels = make([]int, m.total)
	m.Parents = make([]Key, m.total)
	m.Order = make([]int, 0) // Since we append

	// Index the keys
	for i, item := range *items {
		m.Keys[i] = item.Key()
	}

	// Index the Parents
	for _, item := range *items {
		for _, c := range item.Children() {
			if idx := GetIndex(c, items); idx != -1 {
				m.Parents[idx] = item.Key()
			}
		}
	}

	// Recursively index levels and order
	for i := range *items {
		if m.Parents[i] == zeroKey {
			m.indexLevels(i, 0, items)
			m.indexOrder(i, items)
		}
	}
}

func Render(m *Tree, items *[]Item) string {
	out := "\n"
	for i, v := range m.Order {
		level := m.Levels[v]
		indent := strings.Repeat("  ", level)
		ind := "-> "
		if i == m.current {
			ind = "==>"
		}
		out += fmt.Sprintf("%d %s %d: %s %v\n", i, ind, v, indent, (*items)[v].Key())
	}
	return out
}

func (m Tree) Len() int {
	return m.total
}

func (m Tree) Current() int {
	return m.current
}

func (m *Tree) SetCurrent(i int) error {
	if (i > m.total) || (i < 0) {
		return fmt.Errorf("Out of bounds: i - %d ; len: %d", i, m.total)
	}
	m.current = i
	return nil
}
