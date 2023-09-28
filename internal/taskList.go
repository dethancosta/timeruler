package internal

import (
	"sort"
	"time"
)

type TaskList []*Task // TODO implement queue functionality

func (tl TaskList) GetTaskAtTime(t time.Time) *Task {
	// TODO test
	// use binary search on time 
	lo := 0
	hi := len(tl) - 1
	mid := (hi - lo) / 2

	for lo <= hi {

		if t.After(tl[mid].EndTime) {
			lo = mid + 1
			mid = (lo + hi) / 2
			continue
		}

		if t.Before(tl[mid].StartTime) {
			hi = mid - 1
			mid = (lo + hi) / 2
			continue
		}

		return tl[mid]
	}

	return nil
}

func (tl TaskList) IsConflict(t Task) bool {
	for _, task := range tl {
		if task.Tag == BreakTag {continue} // No conflicts with breaks
		if t.StartTime.Before(task.EndTime) && task.StartTime.Before(t.EndTime) {
			return true
		}
	}

	return false
}

func CreateList(tasks ...Task) (TaskList, error) {
	// TODO implement
	// add all tasks, sort, then check for conflicts
	return nil, nil
}

func (tl TaskList) sort() {
	// TODO ensure this works without pointer receiver (and that logic is correct)
	sort.Slice(tl, func(i, j int) bool {return tl[j].StartTime.After(tl[i].EndTime)})
}
