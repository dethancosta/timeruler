package internal

import (
	"testing"
	"time"
)

var (
	task1 = NewTask("Task 1", time.Now(), time.Now().Add(30 * time.Minute))
	task2 = NewTask("Task 2", time.Now().Add(35 * time.Minute), time.Now().Add(45 * time.Minute))
	task3 = NewTask("Task 3", time.Now().Add(50 * time.Minute), time.Now().Add(55 * time.Minute))
	task4 = NewTask("Task 4", time.Now().Add(100 * time.Minute), time.Now().Add(105 * time.Minute))
	task2_5 = NewTask("Task 2.5", time.Now().Add(40 * time.Minute), time.Now().Add(50 * time.Minute))
	task5 = NewTask("Task 5", time.Now().Add(500 * time.Minute), time.Now().Add(550 * time.Minute))
	task6 = NewTask("Task 6", time.Now().Add(-10 * time.Minute), time.Now().Add(-5 * time.Minute))

	validList = TaskList{
		task1,
		task2,
		task3,
	}
	invalidList = TaskList{
		task1,
		task2,
		task3,
		task2_5,
	}
	emptyList = TaskList{}
)

func TestGetTaskAtTime(t *testing.T) {
	t1 := validList.GetTaskAtTime(time.Now().Add(10 * time.Minute))
	if t1 == nil {
		t.Fatalf("Wanted \"Task 1\", got nil,")
	}
	if t1.Description != "Task 1" {
		t.Fatalf("Wanted \"Task 1\", got: %s\n", t1.Description)
	}

	t2 := validList.GetTaskAtTime(time.Now())
	if t2 == nil {
		t.Fatalf("Wanted \"Task 1\", got nil,")
	}
	if t2.Description != "Task 1" {
		t.Fatalf("Wanted \"Task 1\", got: %s\n", t2.Description)
	}

	t3 := validList.GetTaskAtTime(time.Now().Add(40 * time.Minute))
	if t3 == nil {
		t.Fatalf("Wanted \"Task 2\", got nil,")
	}
	if t3.Description != "Task 2" {
		t.Fatalf("Wanted \"Task 2\", got: %s\n", t3.Description)
	}

	t4 := validList.GetTaskAtTime(time.Now().Add(50 * time.Minute))
	if t4 == nil {
		t.Fatalf("Wanted \"Task 3\", got nil,")
	}
	if t4.Description != "Task 3" {
		t.Fatalf("Wanted \"Task 3\", got: %s\n", t4.Description)
	}

	t5 := validList.GetTaskAtTime(time.Now().Add(100 * time.Minute))
	if t5 != nil {
		t.Fatalf("Wanted nil, got: %s", t5.Description)
	}

	t6 := validList.GetTaskAtTime(time.Now().Add(- 10 * time.Minute))
	if t6 != nil {
		t.Fatalf("Wanted nil, got: %s", t6.Description)
	}
}

func TestIsConflict(t *testing.T) {
	test1 := validList.IsConflict(*task2_5)
	if !test1 {
		t.Fatalf("Wanted true, got: false")
	}

	test2 := validList.IsConflict(*task5)
	if test2 {
		t.Fatalf("Wanted false, got true")
	}

	test3 := validList.IsConflict(*task6)
	if test3 {
		t.Fatalf("Wanted false, got true")
	}
}