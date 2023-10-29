package internal

import (
	"testing"
	"time"
)

var (
	tl1 = TaskList{}
)

func TestGetCurrentTask(t *testing.T) {
	task := NewTask("Task1", time.Now().Add(-10 * time.Minute), time.Now().Add(5 * time.Minute))
	tl, err := NewTaskList(
		task,
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	currentTask, currentIdx := tl.GetTaskAtTime(time.Now())

	s1 := Schedule{
		Tasks: tl,
	}
	err = s1.UpdateCurrentTask()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if s1.GetCurrentTaskStr() != task.String() || currentTask.String() != task.String() {
		t.Fatalf("Expected: %s, Got: %s", task.String(), s1.GetCurrentTaskStr())
	}
	if currentIdx != 0 {
		t.Fatalf("Expected: 0, Got: %d", currentIdx)
	}

	// TODO test that gap is a break
	taskBefore := NewTask("Before Break", time.Now().Add(-30 * time.Minute), time.Now().Add(-20 * time.Minute))
	taskAfter := NewTask("Before Break", time.Now().Add(20 * time.Minute), time.Now().Add(30 * time.Minute))
	tl2, err := NewTaskList(taskBefore, taskAfter)
	if err != nil {
		t.Fatalf(err.Error())
	}
	s2 := NewSchedule(tl2)
	expectedBreak := Break(taskBefore.EndTime, taskAfter.StartTime)
	if s2.GetCurrentTaskStr() != expectedBreak.String() {
		t.Fatalf("Expected: %s, Got: %s", expectedBreak.String(), s2.GetCurrentTaskStr())
	}
}

func TestBuildFromFile(t *testing.T) {
	// TODO implement
}

func TestRemoveTask(t *testing.T) {
	// TODO implement
}

func TestUpdateCurrentTask(t *testing.T) {
	// TODO implement
}

func TestUpdateTimeBlock(t *testing.T) {
	// TODO implement
}

func TestAddTask(t *testing.T) {
	// TODO implement
}

func TestChangeCurrentTaskUntil(t *testing.T) {
	// TODO implement
}

func TestGetTasksWithin(t *testing.T) {
	// TODO implement
}
