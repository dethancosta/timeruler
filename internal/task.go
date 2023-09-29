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
} 

// NewTask returns a new Task object with the given description,
// start time, and end time. If end <= start, a nil pointer is returned.
func NewTask(desc string, start, end time.Time) *Task {
	t := &Task{
		Description: desc,
		StartTime: start,
		EndTime: end,
	}

	err := t.Quantize()
	if err != nil {return nil}

	return t
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
	b := NewTask(
		"Take a break",
		start,
		end,
		).WithTag(BreakTag)

	err := b.Quantize()
	if err != nil {return nil}

	return b
}

// String returns the string representation of a Task.
func (t Task) String() string {
	s := t.Description

	if len(strings.TrimSpace(t.Tag)) > 0 {
		s += fmt.Sprintf("\t(%s)\t", t.Tag)
	}

	s += t.EndTime.Format(time.TimeOnly)

	return s
}

// IsValid returns true if a Task's start and end times are
// at least five minutes apart, and returns false otherwise.
func (t Task) IsValid() bool {
	return t.EndTime.Sub(t.StartTime).Minutes() >= 5.0
}


// Helper functions

// Quantize rounds a task's start time and end time
// to 5-minute increments.
func (t *Task) Quantize() error {
	if !t.IsValid() {return InvalidTimeError{"Invalid task time."}}

	t.StartTime = t.StartTime.Round(5 * time.Minute)
	t.EndTime = t.EndTime.Round(5 * time.Minute)

	return nil
}
