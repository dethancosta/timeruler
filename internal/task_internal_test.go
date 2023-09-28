package internal

import (
	"time"
	"testing"
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
