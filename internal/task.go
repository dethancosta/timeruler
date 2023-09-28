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
	Description string `json:"description"`
	StartTime time.Time `json:"startTime"`
	EndTime time.Time `json:"endTime"`
	Tag string `json:"tag"`
} // TODO enforce 5min as unit of increment

// NewTask returns a new Task object with the given description,
// start time, and end time. If end <= start, a nil pointer is returned.
func NewTask(desc string, start, end time.Time) *Task {
	if start.Compare(end) >= 0 {
		return nil
	}

	return &Task{
		Description: desc,
		StartTime: start,
		EndTime: end,
	}
}

// WithTag adds the given tag to the receiver pointer
// and returns the result.
func (t *Task) WithTag(tag string) *Task {
	t.Tag = tag
	return t
}

// Break returns a Task object to be used
// as "free time" in a schedule.
func Break(start, end time.Time) *Task {
	return NewTask(
		"Take a break",
		start,
		end,
		).WithTag(BreakTag)
}

// String returns the string representation of a Task.
func (t *Task) String() string {
	s := t.Description

	if len(strings.TrimSpace(t.Tag)) > 0 {
		s += fmt.Sprintf("\t(%s)\t", t.Tag)
	}

	s += t.EndTime.Format(time.TimeOnly)

	return s
}

func (t Task) IsValid() bool {
	return t.EndTime.Sub(t.StartTime).Minutes() >= 5.0
}
