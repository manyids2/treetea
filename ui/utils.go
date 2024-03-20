package ui

type (
	errMsg error
)

func max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if b >= a {
		return a
	} else {
		return b
	}
}

// InsertItemIntoSlice Insert an item into slice of items at the given index.
func InsertItemIntoSlice(items []string, index int, item string) []string {
	if items == nil {
		return []string{item}
	}
	if index >= len(items) {
		return append(items, item)
	}

	index = max(0, index)

	items = append(items, "")
	copy(items[index+1:], items[index:])
	items[index] = item
	return items
}

// Remove an item from a slice of items at the given index. This runs in O(n).
func RemoveItemFromSlice(i []string, index int) []string {
	if index >= len(i) {
		return i // noop
	}
	copy(i[index:], i[index+1:])
	i[len(i)-1] = ""
	return i[:len(i)-1]
}
