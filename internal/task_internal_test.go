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
	if task1.IsEmpty() {
		t.Fatalf("Wanted valid task, got nil.")
	}

	task2 := NewTask("Task 2", tm1, tm1)
	if !task2.IsEmpty() {
		t.Fatalf("Equal times not nil, wanted nil.")
	}

	task3 := NewTask("Task3", tm1, tm3)
	if !task3.IsEmpty() {
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
	if bTask.Description != "Break" {
		t.Fatalf("break incorrect description. Wanted \"Take a break\", got: %s", bTask.Description)
	}
	if bTask.Tag != "break" {
		t.Fatalf("break incorrect tag. Wanted\"break\", got: %s", bTask.Tag)
	}
}

func TestResolve(t *testing.T) {
	task1 := NewTask("Task 1", time.Now(), time.Now().Add(5*time.Minute))
	task2 := NewTask("Task 2", time.Now(), time.Now().Add(5*time.Minute))

	tl := Resolve(task1, task2)
	if tl != nil {
		t.Fatalf("Wanted nil, got: %v", tl)
	}

	task3 := NewTask("Task 3", time.Now().Add(15*time.Minute), time.Now().Add(20*time.Minute))

	tl = Resolve(task1, task3)
	if len(tl) != 1 {
		t.Fatalf("Wanted 1 task, got: %v", tl)
	}
	if tl[0].Description != "Task 1" && tl[0].StartTime != task1.StartTime && tl[0].EndTime != task1.EndTime {
		t.Fatalf("Wanted \"Task 1\", got: %s", tl[0].Description)
	}

	task3 = NewTask("Task 3", time.Now().Add(-15*time.Minute), time.Now().Add(15*time.Minute))
	task4 := NewTask("Task 4", time.Now().Add(-5*time.Minute), time.Now().Add(5*time.Minute))
	tl = Resolve(task3, task4)
	if len(tl) != 2 {
		t.Fatalf("Wanted 3 tasks, got: %v", tl)
	}
	if tl[0].EndTime != task4.StartTime {
		t.Fatalf("Wanted %v, got: %v", task4.StartTime, tl[0].EndTime)
	}
	if tl[1].StartTime != task4.EndTime {
		t.Fatalf("Wanted %v, got: %v", task4.EndTime, tl[1].StartTime)
	}

	lastTask := NewTask("Last", task4.EndTime.Add(-5*time.Minute), task4.EndTime.Add(5*time.Minute))

	tl = Resolve(task4, lastTask)
	if len(tl) != 1 {
		t.Fatalf("Wanted 1 task, got: %v", tl)
	}
}

func TestConflicts(t *testing.T) {
	t1 := NewTask("", time.Now(), time.Now().Add(5 * time.Minute))
	t2 := NewTask("", time.Now().Add(5 * time.Minute), time.Now().Add(15 * time.Minute))
	t3 := NewTask("", time.Now().Add(10*time.Minute), time.Now().Add(20*time.Minute))
	if (t1.Conflicts(t2)) {
		t.Fatalf("t2 should not conflict with t1")
	}
	if (t2.Conflicts(t1)) {
		t.Fatalf("t2 Conflicts() should be reflexive")
	}
	if (!t3.Conflicts(t2)) {
		t.Fatalf("t2 should conflict with t3")
	}
	if (!t2.Conflicts(t3)) {
		t.Fatalf("t3 Conflicts() should be reflexive")
	}

}
