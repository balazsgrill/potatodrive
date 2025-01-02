package ui

import (
	"github.com/balazsgrill/potatodrive/core/tasks"
	"github.com/lxn/walk"
)

type TaskListModel struct {
	OperationalListModelBase
	sectionsizes []int
	idtoindex    map[int64]int
}

func NewTaskListModel() *TaskListModel {
	return &TaskListModel{
		sectionsizes: []int{0, 0, 0}, // Done, IP, Pending
		idtoindex:    make(map[int64]int),
	}
}

func (t *TaskListModel) targetSection(state tasks.TaskState) int {
	if state.Progress == 100 {
		return 0 // done
	} else if state.Progress > 0 {
		return 1 // in progress
	} else {
		return 2 // pending
	}
}

func (t *TaskListModel) endOfSection(section int) int {
	sectionend := 0
	for i, sectionsize := range t.sectionsizes {
		sectionend += sectionsize
		if i == section {
			return sectionend
		}
	}
	return -1
}

func (t *TaskListModel) currentSection(index int) int {
	sectionend := 0
	for i, sectionsize := range t.sectionsizes {
		sectionend += sectionsize
		if index < sectionend {
			return i
		}
	}
	return -1
}

func (t *TaskListModel) TaskStateListener(state tasks.TaskState) {
	currentindex, exists := t.idtoindex[state.ID]
	if exists {
		currentsecion := t.currentSection(currentindex)
		targetsection := t.targetSection(state)
		if currentsecion != targetsection {
			t.RemoveItem(currentindex)
			t.InsertItemToSection(targetsection, state)
		} else {
			t.ChangeItem(currentindex, state)
		}
	} else {
		t.InsertItemToSection(t.targetSection(state), state)
	}
}

func (t *TaskListModel) RemoveItem(index int) {
	section := t.currentSection(index)
	if section >= 0 {
		t.OperationalListModelBase.RemoveItem(index)
		t.sectionsizes[section]--
	}
}

func (t *TaskListModel) InsertItemToSection(section int, item interface{}) {
	index := t.endOfSection(section)
	t.sectionsizes[section]++
	t.InsertItem(index, item)
}

var _ walk.ListModel = (*TaskListModel)(nil)
