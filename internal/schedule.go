package internal

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Schedule struct {
	Tasks       TaskList
	CurrentTask *Task
	CurrentID   int // ID of the current task
}

// GertTasksWithin returns all tasks that occur within a given time frame
func (s *Schedule) GetTasksWithin(after time.Time, before time.Time) []*Task {
	_, before_idx := s.Tasks.GetTaskAtTime(after)
	_, after_idx := s.Tasks.GetTaskAtTime(before)
	return s.Tasks[after_idx : before_idx+1]
}

// ChangeCurrentTaskUntil updates the schedule's current task
// to have the given description, tag, and end time. It is not
// assumed that there will not be a conflict.
func (s *Schedule) ChangeCurrentTaskUntil(desc, tag string, end time.Time) error {
	// TODO test
	// TODO append to log
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
		if err != nil {
			return err
		}

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
	// TODO append to log
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
// could not be completed.
func (s *Schedule) UpdateTimeBlock(tasks ...Task) error {
	// TODO test
	// TODO append to log
	for _, t := range tasks {
		if !t.IsValid() {
			return InvalidTimeError{"One or more tasks has an invalid time."}
		}
		if t.StartTime.Before(time.Now()) {
			return InvalidTimeError{"Task cannot start before the current time."}
		}
		todayY, todayM, todayD := time.Now().Date()
		if y, m, d := t.StartTime.Date(); y != todayY || m != todayM || d != todayD {
			return InvalidTimeError{"Task cannot start before the current day."}
		}
		if y, m, d := t.EndTime.Date(); y != todayY || m != todayM || d != todayD {
			return InvalidTimeError{"Task cannot end after the current day."}
		}

		conflictBlock := s.GetTasksWithin(t.StartTime, t.EndTime)
		var newBlocks []*Task
		for i := range conflictBlock {
			newBlocks = append(newBlocks, Resolve(conflictBlock[i], &t)...)
		}
		_, n := s.Tasks.GetTaskAtTime(t.StartTime)
		_, m := s.Tasks.GetTaskAtTime(t.EndTime)
		// TODO check if n is 0, and if m is len(...)-1
		s.Tasks = append(s.Tasks[:n-1], append(conflictBlock, s.Tasks[m+1:]...)...)
		s.Tasks = append(s.Tasks, &t)
		s.Tasks.sort()
	}

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
	//TODO append to log
	// id corresponds to index in Tasks queue
	// return an error if the task is a break?
	if id < 0 || id >= len(s.Tasks) {
		return IndexOutOfBoundsError{}
	}

	if s.Tasks[id].IsBreak() {
		return InvalidScheduleError{"Can't remove a break (considered empty)."}
	}

	if id == len(s.Tasks)-1 {
		s.Tasks = s.Tasks[:id]
		return nil
	}

	s.Tasks = append(s.Tasks[:id], s.Tasks[id+1:]...)
	return nil
}

// Print prints the schedule in a barebones format.
// Intended for debugging.
func (s Schedule) Print() {
	sb := strings.Builder{}
	for i, t := range s.Tasks {
		if i == s.CurrentID {
			sb.WriteString("->")
		} else {
			sb.WriteString("  ")
		}

		sb.WriteString("[" + t.StartTime.Format(time.TimeOnly))
		sb.WriteString("-" + t.EndTime.Format(time.TimeOnly) + "] ")
		sb.WriteString(t.Description + " (" + t.Tag + ")\n")
	}
}

// BuildFromFile creates a schedule from a csv file with filename schedule.csv
func BuildFromFile() (*Schedule, error) {
	// TODO test
	// TODO log?
	f, err := os.Open("schedule.csv")
	if err != nil {
		return nil, fmt.Errorf("BuildFromFile: %w", err)
	}
	r := csv.NewReader(f)
	taskList := TaskList{}
	lc := 0
	var desc string
	var start time.Time
	var end time.Time
	var tag string

	for line, err := r.Read(); err == nil; {
		lc++
		if len(line) < 3 {
			return nil,
				errors.New("BuildFromFile: Field missing from line " + strconv.Itoa(lc))
		}
		desc = line[0]
		start, err = time.Parse(time.TimeOnly, line[1])
		if err != nil {
			return nil, errors.New("BuildFromFile: time value improperly formatted on line " + strconv.Itoa(lc))
		}
		end, err = time.Parse(time.TimeOnly, line[2])
		if err != nil {
			return nil, errors.New("BuildFromFile: time value improperly formatted on line " + strconv.Itoa(lc))
		}
		if len(line) == 4 {
			tag = line[3]
		}

		task := NewTask(desc, start, end).WithTag(tag)
		taskList = append(taskList, task)
	}
	if err != io.EOF {
		return nil, fmt.Errorf("BuildFromFile: %w", err)
	}

	taskList.sort()
	current, index := taskList.GetTaskAtTime(time.Now())
	return &Schedule{
		Tasks:       taskList,
		CurrentTask: current,
		CurrentID:   index,
	}, nil
}
