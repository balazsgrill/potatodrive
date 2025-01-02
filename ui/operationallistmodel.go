package ui

import "github.com/lxn/walk"

type OperationalListModelBase struct {
	walk.ListModelBase
	items []interface{}
}

var _ walk.ListModel = (*OperationalListModelBase)(nil)

func (t *OperationalListModelBase) ItemCount() int {
	return len(t.items)
}

func (t *OperationalListModelBase) Value(index int) interface{} {
	return t.items[index]
}

func (t *OperationalListModelBase) InsertItem(index int, item interface{}) {
	t.items = append(t.items[:index], append([]interface{}{item}, t.items[index:]...)...)
	t.PublishItemsInserted(index, index)
}

func (t *OperationalListModelBase) RemoveItem(index int) {
	t.items = append(t.items[:index], t.items[index+1:]...)
	t.PublishItemsRemoved(index, index)
}

func (t *OperationalListModelBase) ChangeItem(index int, item interface{}) {
	t.items[index] = item
	t.PublishItemChanged(index)
}
