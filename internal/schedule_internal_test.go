package internal

import (
	"strings"
	"testing"
	"time"
)

var (
	svalidList, _ = NewTaskList(
		NewTask("Task 0", time.Now().Add(-45*time.Minute), time.Now().Add(-20*time.Minute)),
		NewTask("Task 1", time.Now().Add(-10*time.Minute), time.Now().Add(30*time.Minute)),
		NewTask("Task 2", time.Now().Add(35*time.Minute), time.Now().Add(45*time.Minute)),
		NewTask("Task 3", time.Now().Add(45*time.Minute), time.Now().Add(55*time.Minute)),
	)
	semptyList = TaskList{}
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

	_, err = BuildFromFile("./test_data/overlap1.csv")
	if err == nil {
		t.Fatalf("File with overlapping tasks should not compile")
	}

	quantizedOne, err := BuildFromFile("./test_data/quantize_test1.csv")
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected = strings.Join(strings.Fields(
		`[09:00:00-09:15:00] Eat Breakfast (food)
		[09:15:00-12:15:00] Take a break (break)
		[12:15:00-12:45:00] Eat Lunch (food)
		[12:45:00-17:00:00] Take a break (break)
		[17:00:00-18:05:00] Eat Dinner (food)
		[18:05:00-23:30:00] Take a break (break)
		[23:30:00-23:45:00] Go To Sleep ()`,
	), "")
	got = strings.Join(strings.Fields(quantizedOne.String()), "")
	if expected != got {
		t.Fatalf("Expected: %s\n Got: %s", expected, quantizedOne.String())
	}

	quantizedTwo, err := BuildFromFile("./test_data/quantize_test2.csv")
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected = strings.Join(strings.Fields(
		`[09:00:00-09:15:00] Eat Breakfast (food)
		[09:15:00-12:45:00] Eat Lunch (food)
		[12:45:00-17:00:00] Take a break (break)
		[17:00:00-18:05:00] Eat Dinner (food)
		[18:05:00-23:30:00] Take a break (break)
		[23:30:00-23:45:00] Go To Sleep ()`,
	), "")
	got = strings.Join(strings.Fields(quantizedTwo.String()), "")
	if expected != got {
		t.Fatalf("Expected: %s\n Got: %s", expected, quantizedTwo.String())
	}
}

func TestRemoveTask(t *testing.T) {
	quantizedOne, err := BuildFromFile("./test_data/quantize_test1.csv")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = quantizedOne.RemoveTask(-1)
	if err == nil {
		t.Fatalf("Expected error when removing task at negative index")
	}
	err = quantizedOne.RemoveTask(len(quantizedOne.Tasks))
	if err == nil {
		t.Fatalf("Expected error when removing task at index equal to length")
	}
	err = quantizedOne.RemoveTask(0)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := strings.Join(strings.Fields(
		`[12:15:00-12:45:00] Eat Lunch (food)
		[12:45:00-17:00:00] Take a break (break)
		[17:00:00-18:05:00] Eat Dinner (food)
		[18:05:00-23:30:00] Take a break (break)
		[23:30:00-23:45:00] Go To Sleep ()`,
	), "")
	got := strings.Join(strings.Fields(quantizedOne.String()), "")
	if expected != got {
		t.Fatalf("Expected: %s\n Got: %s", expected, quantizedOne.String())
	}

	err = quantizedOne.RemoveTask(2)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected = strings.Join(strings.Fields(
		`[12:15:00-12:45:00] Eat Lunch (food)
		[12:45:00-23:30:00] Take a break (break)
		[23:30:00-23:45:00] Go To Sleep ()`,
	), "")
	got = strings.Join(strings.Fields(quantizedOne.String()), "")
	if expected != got {
		t.Fatalf("Expected: %s\n Got: %s", expected, quantizedOne.String())
	}

	err = quantizedOne.RemoveTask(1)
	if err == nil {
		t.Fatalf("Expected error when removing break")
	}
}

func TestUpdateCurrentTask(t *testing.T) {
	// TODO Add more tests?
	validSchedule := NewSchedule(svalidList)
	validSchedule.CurrentID = 0
	validSchedule.CurrentTask = validSchedule.Tasks[0]
	currentTask := *validSchedule.CurrentTask
	currentId := validSchedule.CurrentID
	if !strings.HasPrefix(currentTask.String(), "Task 0") {
		t.Fatalf("Expected: Task 0, Got: %s", currentTask.String())
	}
	if currentId != 0 {
		t.Fatalf("Expected: 0, Got: %d", currentId)
	}

	validSchedule.UpdateCurrentTask()
	currentTask = *validSchedule.CurrentTask
	currentId = validSchedule.CurrentID
	if !strings.HasPrefix(currentTask.String(), "Task 1") {
		t.Fatalf("Expected: Task 1, Got: %s", currentTask.String())
	}
	if currentId != 2 {
		t.Fatalf("Expected: 1, Got: %d", currentId)
	}
}

func TestUpdateTimeBlock(t *testing.T) {
	/*
	sched, err := BuildFromFile("./test_data/meals_w_breaks.csv")
	if err != nil {
		t.Fatalf(err.Error())
	}
	now := time.Now()
	for i := range sched.Tasks {
		sched.Tasks[i].StartTime = sched.Tasks[i].StartTime.AddDate(now.Year(), int(now.Month()-1), now.Day()-1)
		sched.Tasks[i].EndTime = sched.Tasks[i].EndTime.AddDate(now.Year(), int(now.Month()-1), now.Day()-1)
	}
	nTask := NewTask("Nap", sched.Tasks[0].StartTime, sched.Tasks[0].EndTime)
	err = sched.UpdateTimeBlock(
		nTask,
	) 
	if err != nil {
		t.Fatalf(err.Error() + "\n" + nTask.String())
	}
	expected := `
	[09:00:00-09:15:00] Nap ()
	[09:15:00-12:15:00] Take a break (break)
	[12:15:00-12:45:00] Eat Lunch (food)
	[12:45:00-17:00:00] Take a break (break)
	[17:00:00-18:00:00] Eat Dinner (food)
	[18:00:00-23:30:00] Take a break (break)
	[23:30:00-23:45:00] Go To Sleep ()`
	got := strings.Join(strings.Fields(sched.String()), "")
	if strings.Join(strings.Fields(expected), "") != got {
		t.Fatalf("Expected: %s\n Got: %s", expected, sched.String())
	}
	// TODO add more tests (subset of exisitng time, straddling 2, incorrect, etc.)
	*/
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
