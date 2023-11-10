package internal

import (
	"strings"
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
	mealsWithBreaks, err := BuildFromFile("./test_data/meals_w_breaks.csv")
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := `[09:00:00-09:15:00] Eat Breakfast (food)
	[09:15:00-12:15:00] Take a break (break)
	[12:15:00-12:45:00] Eat Lunch (food)
	[12:45:00-17:00:00] Take a break (break)
	[17:00:00-18:00:00] Eat Dinner (food)
	[18:00:00-23:30:00] Take a break (break)
	[23:30:00-23:45:00] Go To Sleep ()`
	got := strings.Join(strings.Fields(mealsWithBreaks.String()), "")
	expected = strings.Join(strings.Fields(expected), "")
	if expected != got {
		t.Fatalf("Expected: %s\n Got: %s", expected, mealsWithBreaks.String())
	}

	// TODO add more test cases
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
