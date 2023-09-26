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

func NewTask(desc string, start, end time.Time) *Task {
	return &Task{
		Description: desc,
		StartTime: start,
		EndTime: end,
	}
}

func (t *Task) WithTag(tag string) *Task {
	t.Tag = tag
	return t
}

func Break(start, end time.Time) *Task {
	return NewTask(
		"Take a break",
		start,
		end,
		).WithTag(BreakTag)
}

func (t *Task) String() string {
	s := t.Description

	if len(strings.TrimSpace(t.Tag)) > 0 {
		s += fmt.Sprintf("\t(%s)\t", t.Tag)
	}

	s += t.EndTime.Format(time.TimeOnly)

	return s
}
