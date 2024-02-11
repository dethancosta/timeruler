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
	task := NewTask("Task1", time.Now().Add(-10*time.Minute), time.Now().Add(5*time.Minute))
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

	taskBefore := NewTask("Before Break", time.Now().Add(-30*time.Minute), time.Now().Add(-20*time.Minute))
	taskAfter := NewTask("After Break", time.Now().Add(20*time.Minute), time.Now().Add(30*time.Minute))
	tl2, err := NewTaskList(taskBefore, taskAfter)
	if err != nil {
		t.Fatalf(err.Error())
	}
	s2 := NewSchedule(tl2)
	expectedBreak := Break(taskBefore.EndTime, taskAfter.StartTime)
	if s2.GetCurrentTaskStr() != expectedBreak.String() {
		t.Fatalf("Expected: %s, Got: %s", expectedBreak.String(), s2.GetCurrentTaskStr())
	}
	// TODO add more tests
}

func TestBuildFromFile(t *testing.T) {
	mealsWithBreaks, err := BuildFromFile("./test_data/meals_w_breaks.csv")
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := `[09:00:00-09:15:00] Eat Breakfast (food)
	[09:15:00-12:15:00] Break (break)
	[12:15:00-12:45:00] Eat Lunch (food)
	[12:45:00-17:00:00] Break (break)
	[17:00:00-18:00:00] Eat Dinner (food)
	[18:00:00-23:30:00] Break (break)
	[23:30:00-23:45:00] Go To Sleep ()`
	got := strings.Replace(strings.Join(strings.Fields(mealsWithBreaks.String()), ""),
		"->", "", -1)
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
		[09:15:00-12:15:00] Break (break)
		[12:15:00-12:45:00] Eat Lunch (food)
		[12:45:00-17:00:00] Break (break)
		[17:00:00-18:00:00] Eat Dinner (food)
		[18:00:00-23:30:00] Break (break)
		[23:30:00-23:45:00] Go To Sleep ()`,
	), "")
	got = strings.Replace(strings.Join(strings.Fields(quantizedOne.String()), ""),
		"->", "", -1)
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
		[12:45:00-17:00:00] Break (break)
		[17:00:00-18:00:00] Eat Dinner (food)
		[18:00:00-23:30:00] Break (break)
		[23:30:00-23:45:00] Go To Sleep ()`,
	), "")
	got = strings.Replace(strings.Join(strings.Fields(quantizedTwo.String()), ""),
		"->", "", -1)
	if expected != got {
		t.Fatalf("Expected: %s\n Got: %s", expected, quantizedTwo.String())
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

func TestGetTasksWithin(t *testing.T) {
	sched, err := BuildFromFile("./test_data/meals_w_breaks.csv")
	if err != nil {
		t.Fatalf(err.Error())
	}
	now := time.Now()
	for i := range sched.Tasks {
		sched.Tasks[i].StartTime = sched.Tasks[i].StartTime.AddDate(now.Year(), int(now.Month()-1), now.Day()-1)
		sched.Tasks[i].EndTime = sched.Tasks[i].EndTime.AddDate(now.Year(), int(now.Month()-1), now.Day()-1)
	}
	expected := `
	[09:00:00-09:15:00] Eat Breakfast (food)
	[09:15:00-12:15:00] Break (break)
	[12:15:00-12:45:00] Eat Lunch (food)
	[12:45:00-17:00:00] Break (break)
	[17:00:00-18:00:00] Eat Dinner (food)
	[18:00:00-23:30:00] Break (break)
	[23:30:00-23:45:00] Go To Sleep ()`
	got := strings.Replace(strings.Join(strings.Fields(sched.String()), ""),
		"->", "", -1)
	if strings.Join(strings.Fields(expected), "") != got {
		t.Fatalf("Expected: %s\n Got: %s", expected, sched.String())
	}
	set := sched.GetTasksWithin(
		sched.Tasks[0].StartTime.Add(1*time.Minute),
		sched.Tasks[0].EndTime.Add(-1*time.Minute),
	)
	if len(set) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(set))
	}
	set = sched.GetTasksWithin(
		sched.Tasks[1].StartTime.Add(1*time.Hour),
		sched.Tasks[1].EndTime.Add(-1*time.Hour),
	)
	if len(set) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(set))
	}
	set = sched.GetTasksWithin(
		sched.Tasks[2].StartTime.Add(1*time.Minute),
		sched.Tasks[4].StartTime.Add(-1*time.Minute),
	)
	if len(set) != 2 {
		t.Fatalf("Expected 2 tasks, got %d\n%s", len(set), TaskList(set).String())
	}
	set = sched.GetTasksWithin(
		sched.Tasks[6].StartTime,
		sched.Tasks[6].EndTime.Add(10*time.Minute),
	)
	if len(set) != 1 {
		t.Fatalf("Expected 1 task, got %d \n%s", len(set), TaskList(set).String())
	}
	// TODO add more tests?
}

func TestUpdateTimeBlock(t *testing.T) {
	sched, err := BuildFromFile("./test_data/meals_w_breaks.csv")
	if err != nil {
		t.Fatalf(err.Error())
	}
	nTask := NewTask("Nap", sched.Tasks[0].StartTime, sched.Tasks[0].EndTime)
	err = sched.UpdateTimeBlock(nTask)
	if err != nil {
		t.Fatalf(err.Error() + "\n" + nTask.String())
	}
	expected := `
	[09:00:00-09:15:00] Nap ()
	[09:15:00-12:15:00] Break (break)
	[12:15:00-12:45:00] Eat Lunch (food)
	[12:45:00-17:00:00] Break (break)
	[17:00:00-18:00:00] Eat Dinner (food)
	[18:00:00-23:30:00] Break (break)
	[23:30:00-23:45:00] Go To Sleep ()`
	got := strings.Replace(strings.Join(strings.Fields(sched.String()), ""),
		"->", "", -1)
	if strings.Join(strings.Fields(expected), "") != got {
		t.Fatalf("Expected: %s\n Got: %s", expected, sched.String())
	}

	nTask = NewTask("Nap Again", sched.Tasks[1].StartTime.Add(1*time.Hour), sched.Tasks[1].EndTime.Add(-1*time.Hour))
	if nTask.IsEmpty() {
		t.Fatalf("Task should not be empty")
	}

	err = sched.UpdateTimeBlock(nTask)
	if err != nil {
		t.Fatalf(err.Error() + "\n" + nTask.String() + "\n\n" + sched.String())
	}
	expected = `
	[09:00:00-09:15:00] Nap ()
	[09:15:00-10:15:00] Break (break)
	[10:15:00-11:15:00] Nap Again ()
	[11:15:00-12:15:00] Break (break)
	[12:15:00-12:45:00] Eat Lunch (food)
	[12:45:00-17:00:00] Break (break)
	[17:00:00-18:00:00] Eat Dinner (food)
	[18:00:00-23:30:00] Break (break)
	[23:30:00-23:45:00] Go To Sleep ()`
	got = strings.Replace(strings.Join(strings.Fields(sched.String()), ""),
		"->", "", -1)
	if strings.Join(strings.Fields(expected), "") != got {
		t.Fatalf("Expected: %s\n Got: %s", expected, sched.String())
	}
	// TODO add more tests (subset of exisitng time, straddling 2, incorrect, 2 new tasks, etc.)
}

func TestAddTask(t *testing.T) {
	sched, err := BuildFromFile("./test_data/meals_w_breaks.csv")
	if err != nil {
		t.Fatalf(err.Error())
	}
	nTask := NewTask("Nap", sched.Tasks[0].StartTime, sched.Tasks[0].EndTime)
	err = sched.AddTask(nTask)
	if err == nil {
		t.Fatalf("Expected error when adding task that overlaps with existing task")
	}
	nTask = NewTask("Nap", sched.Tasks[1].StartTime, sched.Tasks[1].EndTime)
	err = sched.AddTask(nTask)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected := `
	[09:00:00-09:15:00] Eat Breakfast (food)
	[09:15:00-12:15:00] Nap ()
	[12:15:00-12:45:00] Eat Lunch (food)
	[12:45:00-17:00:00] Break (break)
	[17:00:00-18:00:00] Eat Dinner (food)
	[18:00:00-23:30:00] Break (break)
	[23:30:00-23:45:00] Go To Sleep ()`
	got := strings.Replace(strings.Join(strings.Fields(sched.String()), ""),
		"->", "", -1)
	if strings.Join(strings.Fields(expected), "") != got {
		t.Fatalf("Expected: %s\n Got: \n%s", expected, sched.String())
	}

	nTask = NewTask("Nap Again", sched.Tasks[3].StartTime.Add(30*time.Minute), sched.Tasks[3].EndTime.Add(-20*time.Minute))
	err = sched.AddTask(nTask)
	if err != nil {
		t.Fatalf(err.Error())
	}
	expected = `
		[09:00:00-09:15:00] Eat Breakfast (food)
		[09:15:00-12:15:00] Nap ()
		[12:15:00-12:45:00] Eat Lunch (food)
		[12:45:00-13:15:00] Break (break)
		[13:15:00-16:40:00] Nap Again ()
		[16:40:00-17:00:00] Break (break)
		[17:00:00-18:00:00] Eat Dinner (food)
		[18:00:00-23:30:00] Break (break)
		[23:30:00-23:45:00] Go To Sleep ()`
	got = strings.Replace(strings.Join(strings.Fields(sched.String()), ""),
		"->", "", -1)
	if strings.Join(strings.Fields(expected), "") != got {
		t.Fatalf("Expected: %s\n Got: \n%s", expected, sched.String())
	}
}

func TestChangeCurrentTaskUntil(t *testing.T) {
	sched, err := BuildFromFile("./test_data/meals_w_breaks.csv")
	if err != nil {
		t.Fatalf(err.Error())
	}
	now := time.Now()
	err = sched.ChangeCurrentTaskUntil("Nap", "", now.Add(1*time.Hour))
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = sched.ChangeCurrentTaskUntil("Nap", "", now.Add(-1*time.Hour))
	if err == nil {
		t.Fatalf("Expected error when changing current task to before start time")
	}
	err = sched.ChangeCurrentTaskUntil("Nap", "", now.Add(24*time.Hour))
	if err == nil {
		t.Fatalf("Expected error when setting end time to another day")
	}
}

func TestFixBreaks(t *testing.T) {
	// TODO finish implementing
	tl, err := NewTaskList(NewTask("", time.Now(), time.Now().Add(5 * time.Minute)))
	if err != nil {
		t.Fatalf("NewTaskList should not throw error: %s", err.Error())
	}
	sched := NewSchedule(tl)
	l := len(sched.Tasks)
	sched.FixBreaks()
	if l != len(sched.Tasks) {
		t.Fatalf("FixBreaks should not change length of single-task schedule")
	}
}

func TestNewSchedule(t *testing.T) {
	// TODO implement
}
