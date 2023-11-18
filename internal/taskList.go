package internal

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type TaskList []*Task

// GetTaskAtTime returns the task occupying the
// time block that contains the given time. 
// (nil, -1) is returned if there is no task at time t.
func (tl TaskList) GetTaskAtTime(t time.Time) (*Task, int) {
	for i := range tl {
		if tl[i].StartTime.Compare(t) <= 0 && tl[i].EndTime.Compare(t) >= 0 {
			return tl[i], i
		}
	}

	return nil, -1
}

// IsConflict returns true if there is overlap between
// the given task and any of the tasks currently in
// the TaskList.
func (tl TaskList) IsConflict(t Task) bool {
	for _, task := range tl {
		if task.Tag == BreakTag {
			continue
		} // No conflicts with breaks
		if t.StartTime.Before(task.EndTime) && task.StartTime.Before(t.EndTime) {
			return true
		}
	}

	return false
}

// NewTaskList creates a new TaskList from the given tasks.
// It returns nil and an error if there is a time conflict.
func NewTaskList(tasks ...Task) (TaskList, error) {
	// TODO test
	// add all tasks, sort, then check for conflicts

	for _, t := range tasks {
		if !t.IsValid() {
			return nil, InvalidScheduleError{"Invalid Task was given."}
		}
	}

	tl := TaskList{}
	var taskRef *Task
	for t := 0; t < len(tasks); t++ {
		taskRef = &tasks[t]
		err := taskRef.Quantize()
		if err != nil {
			return nil, fmt.Errorf("Error quantizing task: %v", err)
		}
		tl = append(tl, taskRef)
	}

	tl.sort()
	if !tl.IsConsistent() {
		return nil, InvalidScheduleError{"List of tasks contains a conflict."}
	}

	// Add breaks as needed
	var former time.Time
	var latter time.Time
	for i := 0; i < len(tl)-1; i++ {
		former = tl[i].EndTime
		latter = tl[i+1].StartTime
		if former.Compare(latter) != 0 {
			b := Break(former, latter)
			tl = append(tl[:i+1], append([]*Task{&b}, tl[i+1:]...)...)
		}
	}

	return tl, nil
}

// IsConsistent returns true if the TaskList has no overlapping
// tasks, and false otherwise. It assumes that the TaskList is
// sorted.
func (tl TaskList) IsConsistent() bool {
	for i := 0; i < len(tl)-1; i++ {
		if tl[i].EndTime.After(tl[i+1].StartTime) {
			return false
		}
	}

	return true
}

// ResolveConflicts adjusts the start and end times of the tasks starting at
// the given index to accomodate the new given task. It updates a copy of
// the task list and returns the updated copy.
func (tl TaskList) ResolveConflicts(oldTaskId int, newTask *Task) (TaskList, error) {
	if oldTaskId < 0 || oldTaskId >= len(tl) {
		return nil, IndexOutOfBoundsError{}
	}
	oldTask := tl.get(oldTaskId)

	for oldTaskId < len(tl) && oldTask.Conflicts(*newTask) {

		updated := Resolve(*oldTask, *newTask)
		if oldTaskId < len(tl)-1 {
			post := tl[oldTaskId+1:]
			tl = append(tl[:oldTaskId], updated...)
			tl = append(tl, post...)
		} else {
			tl = append(tl[:oldTaskId], updated...)
		}

		oldTaskId++
		oldTask = tl.get(oldTaskId)
	}

	tl = append(tl, newTask)
	tl.sort()

	return tl, nil
}

func (tl TaskList) sort() {
	sort.Slice(tl, func(i, j int) bool { return tl[i].EndTime.Compare(tl[j].StartTime) <= 0 })
}

func (tl TaskList) get(idx int) *Task {
	return tl[idx]
}

func (tl TaskList) String() string {
	sb := strings.Builder{}
	for t := range tl {
		sb.WriteString(tl[t].String() + "\n")
	}
	return sb.String()
}
