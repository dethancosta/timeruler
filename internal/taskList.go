package internal

import (
	"sort"
	"time"
)

type TaskList []*Task

// GetTaskAtTime returns the task occupying the
// time block that contains the given time. It 
// returns nil if there is no task at time t.
func (tl TaskList) GetTaskAtTime(t time.Time) (*Task, int) {
	// TODO test
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

		return tl[mid], mid
	}

	return nil, -1
}

// IsConflict returns true if there is overlap between
// the given task and any of the tasks currently in 
// the TaskList.
func (tl TaskList) IsConflict(t Task) bool {
	for _, task := range tl {
		if task.Tag == BreakTag {continue} // No conflicts with breaks
		if t.StartTime.Before(task.EndTime) && task.StartTime.Before(t.EndTime) {
			return true
		}
	}

	return false
}

// CreateList creates a new TaskList from the given tasks.
// It returns nil and an error if there is a time conflict.
func CreateList(tasks ...Task) (TaskList, error) {
	// TODO test
	// add all tasks, sort, then check for conflicts

	for _, t := range tasks {
		if !t.IsValid() {
			return nil, InvalidScheduleError{"Invalid Task was given."}
		}
	}

	tl := TaskList{}
	for _, t := range tasks {
		tl = append(tl, &t)
	}

	tl.sort()
	if !tl.IsConsistent() {
		return nil, InvalidScheduleError{"List of tasks contains a confict."}
	}

	return tl, nil
}

// IsConsistent returns true if the TaskList has no overlapping
// tasks, and false otherwise. It assumes that the TaskList is 
// sorted.
func (tl TaskList) IsConsistent() bool {
	for i := 0; i < len(tl) - 1; i++ {
		if tl[i].EndTime.After(tl[i+1].StartTime) {
			return false
		}
	}

	return true
}

// ResolveConflict adjusts the start and end times of the tasks starting at
// the given index to accomodate the new given task. It updates a copy of 
// the task list and returns the updated copy.
func (tl TaskList) ResolveConflicts(oldTaskId int, newTask *Task) (TaskList, error) {
	// TODO test
	if oldTaskId < 0 || oldTaskId >= len(tl) {
		return nil, IndexOutOfBoundsError{}
	}
	oldTask := tl.get(oldTaskId)

	for oldTaskId < len(tl) && oldTask.Conflicts(*newTask) {
		oldTask = tl.get(oldTaskId)

		updated := Resolve(oldTask, newTask)
		if oldTaskId < len(tl)-1 {
			post := tl[oldTaskId+1:]
			tl = append(tl[:oldTaskId], updated...)
			tl = append(tl, post...)
		} else {
			tl = append(tl[:oldTaskId], updated...)
		}

		oldTaskId++
	}

	tl = append(tl, newTask)
	tl.sort()

	return tl, nil
}

func (tl TaskList) sort() {
	// TODO test
	// ensure this works without pointer receiver (and that logic is correct)
	sort.Slice(tl, func(i, j int) bool {return tl[j].StartTime.After(tl[i].EndTime)})
}

func (tl TaskList) get(idx int) *Task {
	return tl[idx]
}
