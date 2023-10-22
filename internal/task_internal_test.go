package internal

import (
	"testing"
	"time"
)

func TestInvalidTimes(t *testing.T) {
	tm1 := time.Now()
	tm2 := time.Now().Add(5 * time.Minute)
	tm3 := time.Now().Add(-5 * time.Minute)

	task1 := NewTask("Task 1", tm1, tm2)
	if task1 == nil {
		t.Fatalf("Wanted valid task, got nil.")
	}

	task2 := NewTask("Task 2", tm1, tm1)
	if task2 != nil {
		t.Fatalf("Equal times not nil, wanted nil.")
	}

	task3 := NewTask("Task3", tm1, tm3)
	if task3 != nil {
		t.Fatalf("EndTime < StartTime not nil, wanted nil.")
	}
}

func TestNewTask(t *testing.T) {
	tm1 := time.Now()
	tm2 := time.Now().Add(5 * time.Minute)
	task1 := NewTask("Task 1", tm1, tm2)

	if task1.Description != "Task 1" {
		t.Fatalf("task1 description incorrect.\nWanted \"Task 1\", got: %s", task1.Description)
	}
}

func TestWithTag(t *testing.T) {
	tm1 := time.Now()
	tm2 := time.Now().Add(5 * time.Minute)
	task1 := NewTask("Task 1", tm1, tm2).WithTag("Tag1")

	if task1.Tag != "Tag1" {
		t.Fatalf("task1 tag incorrect.\nWanted \"Tag1\", got: %s", task1.Tag)
	}
}

func TestBreak(t *testing.T) {
	bTask := Break(time.Now(), time.Now().Add(5*time.Minute))
	if bTask.Description != "Take a break" {
		t.Fatalf("break incorrect description. Wanted \"Take a break\", got: %s", bTask.Description)
	}
	if bTask.Tag != "break" {
		t.Fatalf("break incorrect tag. Wanted\"break\", got: %s", bTask.Tag)
	}
}
