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

func (s *Schedule) ChangeCurrentTaskUntil(desc, tag string, end time.Time) {
	// TODO implement
}

func (s *Schedule) GetCurrentTaskStr() string {
	return strconv.Itoa(s.CurrentID) + "\t" + s.CurrentTask.String()
}

func (s *Schedule) AddTask(t Task) error {
	// TODO implement
	// Used when no conflicts with current tasks are expected
	// Returns an error in the case of an overlap/conflict
	// or if task could otherwise not be added
	return nil
}

func (s *Schedule) UpdateTimeBlock(tasks ...Task) error {
	// TODO implement
	// Used when conflicts with current tasks may be expected
	// Returns an error if update could not be completed
	return nil
	// NOTE old tasks should be kept in display, and crossed out
	// rather than disappear (by default)
}

func (s *Schedule) UpdateCurrentTask() error {
	// TODO implement
	// For use with timer or change/request from client
	return nil
}

func (s *Schedule) RemoveTask(id int) error {
	//TODO implement
	// id corresponds to index in Tasks queue
	return nil
}

func (s *Schedule) CrossOutTask(id int) error {
	// TODO implement
	// id corresponds to index in Tasks queue
	return nil
}

