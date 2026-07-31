package terminal

import "slices"

type Completer interface {
	Do(r []rune, pos int) ([]*CompletableItem, int)
}

type CompletableItem struct {
	value []rune
	items []*CompletableItem
}

func (c *CompletableItem) String() string {
	return string(c.value)
}

func NewItem(v string, items ...*CompletableItem) *CompletableItem {
	return &CompletableItem{[]rune(v), items}
}

func NewCompleter(items ...*CompletableItem) *CompletableItem {
	return NewItem("", items...)
}

func (c *CompletableItem) Do(r []rune, pos int) ([]*CompletableItem, int) {
	return doInternal(c.items, r, pos)
}

func doInternal(items []*CompletableItem, r []rune, pos int) ([]*CompletableItem, int) {
	candidates := make([]*CompletableItem, 0)
	r = TrimSpaceLeft(r[:pos])
	for _, i := range items {
		if len(i.value) > len(r) {
			if HasPrefix(r, i.value) {
				candidates = append(candidates, i)
			}
		} else {
			if HasPrefix(i.value, r) {
				if len(i.items) > 0 {
					c, _ := doInternal(i.items, r[len(i.value):], pos)
					candidates = append(candidates, c...)
				}
			}
		}
	}
	slices.SortFunc(candidates, func(i, j *CompletableItem) int {
		return slices.Compare(i.value, j.value)
	})
	candidates = slices.CompactFunc(candidates, func(i, j *CompletableItem) bool {
		return slices.Equal(i.value, j.value)
	})
	return candidates, 0
}
