package internal

import (
	"strconv"
	"time"
)


type Schedule struct {
	Tasks TaskList
	CurrentTask *Task
	CurrentID int // ID of the current task
	OldTasks TaskList // Tasks that have been crossed out
}

// ChangeCurrentTaskUntil updates the schedule's current task 
// to have the given description, tag, and end time. It does not
// assume that there will not be a conflict.
func (s *Schedule) ChangeCurrentTaskUntil(desc, tag string, end time.Time) error {
	// TODO test
	if end.Compare(time.Now()) <= 1 {
		return InvalidTimeError{"Task ends before the current time."}
	}

	newCurrent := NewTask(desc, time.Now(), end).WithTag(tag)
	err := newCurrent.Quantize()
	if err != nil {
		return err
	}

	if !s.Tasks.IsConflict(*newCurrent) {
		_, idx := s.Tasks.GetTaskAtTime(time.Now())
		s.Tasks[idx].EndTime = time.Now() // Should be the break
		err := s.Tasks[idx].Quantize()
		if err != nil {return err}

		s.Tasks = append(s.Tasks[:idx+1], append([]*Task{newCurrent}, s.Tasks[idx+1:]...)...)
		
		//s.Tasks = append(s.Tasks, newCurrent)
		//s.Tasks.sort()
	} else {
		_, idx := s.Tasks.GetTaskAtTime(time.Now())
		newList, err := s.Tasks.ResolveConflicts(idx, newCurrent)
		if err != nil {
			return err
		}
		s.Tasks = newList
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
	if t.StartTime.Before(time.Now()) {
		return InvalidTimeError{"Task cannot start in the past."}
	}
	if !t.IsValid() {
		return InvalidTimeError{"Task times are invalid."}
	}
	
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

	for _, t := range tasks {
		if !t.IsValid() {
			return InvalidTimeError{"One or more tasks has an invalid time."}
		}
		if !t.StartTime.Before(time.Now()) {
			return InvalidTimeError{"Task cannot start before the current time."}
		}
		// TODO consider checking if task's endTime is after today.
	}

	// TODO finish implementing

	return nil
}

// UpdateCurrentTask checks the schedule's task list for the 
// task scheduled for the time that the function is called.
// Note that the task will be nil if there is no currently
// scheduled task.
func (s *Schedule) UpdateCurrentTask() error {
	// TODO test
	// For use with timer or change/request from client
	s.CurrentTask, _ = s.Tasks.GetTaskAtTime(time.Now())

	return nil
}

// RemoveTask removes the ith task in the task list.
// It returns an error if the id is invalid or the 
// task could otherwise not be removed.
func (s *Schedule) RemoveTask(id int) error {
	//TODO test
	// id corresponds to index in Tasks queue
	// return an error if the task is a break?
	if id < 0 || id >= len(s.Tasks) {
		return IndexOutOfBoundsError{}
	}

	if s.Tasks[id].IsBreak() {
		return InvalidScheduleError{"Can't remove a break (considered empty)."}
	}
	
	if id == len(s.Tasks) - 1 {
		s.Tasks = s.Tasks[:id]
		return nil
	}
	
	s.Tasks = append(s.Tasks[:id], s.Tasks[id+1:]...)
	return nil
}

// CrossOutTask crosses out the ith task in the task 
// list. It returns an error if the id is invalid or
// the task could otherwise not be crossed off.
func (s *Schedule) CrossOutTask(id int) error {
	// TODO test
	// id corresponds to index in Tasks queue
	// return an error if the task is a break

	if id < 0 || id >= len(s.Tasks) {
		return IndexOutOfBoundsError{}
	}
	if s.Tasks[id].IsBreak() {
		return InvalidScheduleError{"Can't remove a break (considered empty slot)."}
	}

	s.OldTasks = append(s.OldTasks, s.Tasks[id])

	// TODO consider not removing task from s.Tasks, but instead checking if each task 
	// is in OldTasks in other Schedule methods.
	if id == len(s.Tasks) - 1 {
		s.Tasks = s.Tasks[:id]
	} else {
		s.Tasks = append(s.Tasks[:id], s.Tasks[id+1:]...)
	}
	
	return nil
}

