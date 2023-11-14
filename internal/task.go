package internal

import (
	"fmt"
	"strings"
	"time"
)

const (
	BreakTag = "break"
)

type Task struct {
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Tag         string    `json:"tag"`
}

// NewTask returns a new Task object with the given description,
// start time, and end time. If end <= start, an empty task is returned.
func NewTask(desc string, start, end time.Time) Task {
	t := Task{
		Description: desc,
		StartTime:   start,
		EndTime:     end,
	}

	err := t.Quantize()
	if err != nil {
		return Task{}
	}

	return t
}

// WithTag adds the given tag to the receiver pointer
// and returns the result.
func (t Task) WithTag(tag string) Task {
	t.Tag = tag
	return t
}

// Break returns a Task object to be used
// as "free time" in a schedule.
func Break(start, end time.Time) Task {
	b := NewTask(
		"Take a break",
		start,
		end,
	).WithTag(BreakTag)

	err := b.Quantize()
	if err != nil {
		return Task{}
	}

	return b
}

// IsBreak returns true if the task has a break tag
func (t Task) IsBreak() bool {
	return t.Tag == BreakTag
}

// String returns the string representation of a Task.
func (t Task) String() string {
	s := t.Description + "\t"
	s += t.StartTime.Format(time.DateTime) + "\t"
	s += t.EndTime.Format(time.DateTime) + "\t"

	if len(strings.TrimSpace(t.Tag)) > 0 {
		s += fmt.Sprintf("\t(%s)", t.Tag)
	}

	return s
}

// IsValid returns true if a Task's start and end times are
// at least five minutes apart, and returns false otherwise.
func (t Task) IsValid() bool {
	return t.EndTime.Sub(t.StartTime).Minutes() >= 5.0
}

// Conflicts returns true if the given task's time span overlaps
// with the receiver task's time span.
func (t Task) Conflicts(other Task) bool {
	// TODO test
	return t.StartTime.Before(other.EndTime) && t.StartTime.Before(other.EndTime)
}

// Resolve updates the time of the old task to remove overlap
// between the old task's time span and that of the new task.
// it returns a Task with the updated times of oldTask. If
// newTask's time is a subset of oldTask's, 2 Tasks
// will be returned. It assumes the tasks have a conflict,
// so oldTask may be incorrectly updated if there is none.
func Resolve(oldTask, newTask *Task) []*Task {
	// TODO test
	if oldTask.StartTime.Before(newTask.StartTime) {
		if oldTask.EndTime.Compare(newTask.EndTime) <= 0 {
			oldTask.EndTime = newTask.StartTime
			return []*Task{oldTask}
		}
		postTask := &Task{
			Description: oldTask.Description,
			StartTime:   newTask.EndTime,
			EndTime:     oldTask.EndTime,
			Tag:         oldTask.Tag,
		}
		oldTask.EndTime = newTask.StartTime
		return []*Task{oldTask, postTask}
	} else {
		if oldTask.EndTime.Compare(newTask.EndTime) <= 0 {
			// oldTask will be removed
			return nil
		} else {
			oldTask.StartTime = newTask.EndTime
			return []*Task{oldTask}
		}
	}
}

// Helper functions

// Quantize rounds a task's start time and end time
// to 5-minute increments.
func (t *Task) Quantize() error {
	if !t.IsValid() {
		return InvalidTimeError{"Invalid task time."}
	}

	t.StartTime = t.StartTime.Round(5 * time.Minute)
	t.EndTime = t.EndTime.Round(5 * time.Minute)

	return nil
}

// IsEmpty tests whether t is an empty (default-value) Task
func (t Task) IsEmpty() bool {
	return t.StartTime.IsZero() &&
		t.EndTime.IsZero() &&
		strings.TrimSpace(t.Description) == ""
}
