package internal

import (
	"sort"
	"time"
)

type TaskList []*Task // TODO implement queue functionality

func (tl TaskList) GetTaskAtTime(t time.Time) *Task {
	// TODO implement
	// use binary search on time 
	lo := 0
	hi := len(tl) - 1
	mid := (hi - lo) / 2

	for lo < hi {
		if tl[mid].EndTime.After(t) && t.After(tl[mid].StartTime) {
			return tl[mid]
		}

		if t.After(tl[mid].EndTime) {
			lo = mid + 1
			mid = (hi - lo) / 2
			continue
		}

		if tl[mid].StartTime.After(t) {
			hi = mid - 1
			mid = (hi - lo) / 2
			continue
		}
	}

	return nil
}

func (tl TaskList) IsConflict(t Task) bool {
	for _, task := range tl {
		if task.Tag == BreakTag {continue} // No conflicts with breaks
		if t.StartTime.After(task.StartTime) && task.EndTime.After(t.StartTime) {
			return true
		}
		if task.StartTime.After(t.StartTime) && t.EndTime.After(task.StartTime) {
			return true
		}
	}

	return false
}

func (tl TaskList) sort() {
	// TODO ensure this works without pointer receiver (and that logic is correct)
	sort.Slice(tl, func(i, j int) bool {return tl[j].StartTime.After(tl[i].EndTime)})
}
