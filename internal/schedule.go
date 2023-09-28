package internal

import (
	"strconv"
	"time"
)


type Schedule struct {
	Tasks TaskList
	CurrentTask *Task
	CurrentID int // ID of the current task
	OldTasks TaskList // Tasks that have been crossed out or completed
}

// ChangeCurrentTaskUntil updates the schedule's current task 
// to have the given description, tag, and end time. It does not
// assume that there will not be a conflict.
func (s *Schedule) ChangeCurrentTaskUntil(desc, tag string, end time.Time) error {
	// TODO implement
	if end.Compare(time.Now()) <= 1 {
		return InvalidTimeError{"Task ends before the current time."}
	}

	newCurrent := NewTask(desc, time.Now(), end).WithTag(tag)
	if !s.Tasks.IsConflict(*newCurrent) {
		s.Tasks = append(s.Tasks, newCurrent)
		s.Tasks.sort()
	} else {
		// TODO implement update of affected tasks
	}

	return nil
}

// GetCurrentTaskStr returns the schedule's current task 
// as a formatted string.
func (s *Schedule) GetCurrentTaskStr() string {
	return strconv.Itoa(s.CurrentID) + "\t" + s.CurrentTask.String()
}

// AddTask is used when no conflicts with the schedule's current 
// tasks are expected. It returns an error in the case of an
// overlap/conflict or if the tasks could otherwise not be added
func (s *Schedule) AddTask(t Task) error {
	// TODO test
	
	if s.Tasks.IsConflict(t) {
		return InvalidScheduleError{"Task conflicts with schedule."}
	}
	s.Tasks = append(s.Tasks, &t)
	s.Tasks.sort()

	return nil
}

// UpdateTimeBlock updates the schedule's task list with
// the given collection of tasks. It is not assumed that
// no conflicts will exist, and will alter the existing
// tasks as needed. It returns an error if the update 
// could not be completed. If hardRemove is false,
// the affected (portions of) tasks will be crossed out
// rather than completely removed.
func (s *Schedule) UpdateTimeBlock(hardRemove bool, tasks ...Task) error {
	// TODO implement
	return nil
}

// UpdateCurrentTask checks the schedule's task list for the 
// task scheduled for the time that the function is called.
// Note that the task will be nil if there is no currently
// scheduled task.
func (s *Schedule) UpdateCurrentTask() error {
	// TODO test
	// For use with timer or change/request from client
	s.CurrentTask = s.Tasks.GetTaskAtTime(time.Now())

	return nil
}

// RemoveTask removes the ith task in the task list.
// It returns an error if the id is invalid or the 
// task could otherwise not be removed.
func (s *Schedule) RemoveTask(id int) error {
	//TODO implement
	// id corresponds to index in Tasks queue
	// return an error if the task is a break?
	return nil
}

// CrossOutTask crosses out the ith task in the task 
// list. It returns an error if the id is invalid or
// the task could otherwise not be crossed off.
func (s *Schedule) CrossOutTask(id int) error {
	// TODO implement
	// id corresponds to index in Tasks queue
	// return an error if the task is a break
	return nil
}

