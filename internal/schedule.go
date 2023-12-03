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
func (s *Schedule) GetTasksWithin(before time.Time, after time.Time) []*Task {
	endTime := s.Tasks[len(s.Tasks)-1].EndTime
	startTime := s.Tasks[0].StartTime
	_, before_idx := s.Tasks.GetTaskAtTime(before.Add(1 * time.Minute))
	_, after_idx := s.Tasks.GetTaskAtTime(after)
	if after.Compare(endTime) >= 0 {
		after_idx = len(s.Tasks) - 1
	}
	if before.Compare(startTime) <= 0 {
		before_idx = 0
	}
	if before_idx == -1 || after_idx == -1 {
		return nil
	}
	return s.Tasks[before_idx : after_idx+1]
}

// ChangeCurrentTaskUntil updates the schedule's current task
// to have the given description, tag, and end time. It is not
// assumed that there will not be a conflict.
func (s *Schedule) ChangeCurrentTaskUntil(desc, tag string, end time.Time) error {
	if end.Compare(time.Now()) <= 0 {
		return InvalidTimeError{"Task ends before the current time."}
	}
	now := time.Now()
	if end.Day() != now.Day() || end.Month() != now.Month() || end.Year() != now.Year() {
		return InvalidTimeError{"Task must end during the current day."}
	}

	newCurrent := NewTask(desc, time.Now(), end).WithTag(tag)
	err := newCurrent.Quantize()
	if err != nil {
		return err
	}

	if len(s.Tasks) == 0 || s.Tasks.get(len(s.Tasks)-1).EndTime.Before(time.Now()) {
		s.Tasks = append(s.Tasks, &newCurrent)
	} else if s.Tasks[0].StartTime.After(time.Now()) {
		if !s.Tasks.IsConflict(newCurrent) {
			s.Tasks = append([]*Task{&newCurrent}, s.Tasks...)
		} else {
			newList, err := s.Tasks.ResolveConflicts(newCurrent)
			if err != nil {
				return err
			}
			s.Tasks = newList
		}

	} else if !s.Tasks.IsConflict(newCurrent) {
		_, idx := s.Tasks.GetTaskAtTime(time.Now())
		s.Tasks[idx].EndTime = time.Now() // Should be the break
		err := s.Tasks[idx].Quantize()
		if err != nil {
			return err
		}

		s.Tasks = append(s.Tasks[:idx+1], append([]*Task{&newCurrent}, s.Tasks[idx+1:]...)...)
	} else {
		newList, err := s.Tasks.ResolveConflicts(newCurrent)
		if err != nil {
			return err
		}
		s.Tasks = newList
	}
	s.FixBreaks()

	return nil
}

// GetCurrentTaskStr returns the schedule's current task
// as a formatted string.
func (s *Schedule) GetCurrentTaskStr() string {
	return s.CurrentTask.String()
}

// AddTask is used when no conflicts with the schedule's current
// tasks are expected. It returns an error in the case of an
// overlap/conflict or if the tasks could otherwise not be added
func (s *Schedule) AddTask(t Task) error {
	if !t.IsValid() {
		return InvalidTimeError{"Task times are invalid."}
	}

	if s.Tasks.IsConflict(t) {
		return InvalidScheduleError{"Task conflicts with schedule."}
	}
	//_, idx := s.Tasks.GetTaskAtTime(t.StartTime)
	//newTasks, err := s.Tasks.ResolveConflicts(idx, t)
	err := s.UpdateTimeBlock(t)
	if err != nil {
		return fmt.Errorf("AddTask: %w", err)
	}
	//s.Tasks = newTasks

	return nil
}

// UpdateTimeBlock updates the schedule's task list with
// the given collection of tasks. It is not assumed that
// no conflicts will exist, and will alter the existing
// tasks as needed. It returns an error if the update
// could not be completed.
func (s *Schedule) UpdateTimeBlock(tasks ...Task) error {
	for _, t := range tasks {
		if !t.IsValid() {
			return InvalidTimeError{"One or more tasks has an invalid time."}
		}
		todayY, todayM, todayD := time.Now().Date()
		if y, m, d := t.StartTime.Date(); y != todayY || m != todayM || d != todayD {
			return InvalidTimeError{"Task must start during the current day."}
		}
		if y, m, d := t.EndTime.Date(); y != todayY || m != todayM || d != todayD {
			return InvalidTimeError{"Task must end during the current day."}
		}

		/*
		newTasks, err := NewTaskList()
		if err != nil {
			return fmt.Errorf("UpdateTimeBlock: %w", err)
		}

		for i := range s.Tasks {
			if s.Tasks[i].EndTime.Before(t.StartTime) {
				newTasks = append(newTasks, s.Tasks[i])
			} else if s.Tasks[i].StartTime.After(t.EndTime) {
				newTasks = append(newTasks, &t)
				newTasks = append(newTasks, s.Tasks[i:]...)
				break
			} else {
				if s.Tasks[i].StartTime.Before(t.StartTime) {
					nt := NewTask(s.Tasks[i].Description, s.Tasks[i].StartTime, t.StartTime).WithTag(s.Tasks[i].Tag)
					newTasks = append(newTasks, &nt)
				}
				if s.Tasks[i].EndTime.After(t.EndTime) {
					nt := NewTask(s.Tasks[i].Description, t.EndTime, s.Tasks[i].EndTime).WithTag(s.Tasks[i].Tag)
					newTasks = append(newTasks, &nt)
				}
			}
		}
		*/

		//newTasks.sort()
		
		newTasks, err := s.Tasks.ResolveConflicts(t)
		if err != nil {
			return fmt.Errorf("UpdateTimeBlock: %w", err)
		}
		s.Tasks = newTasks
		s.FixBreaks()
	}

	return nil
}

func (s *Schedule) FixBreaks() {
	// TODO test
	for i := 0; i < len(s.Tasks)-1; i++ {
		if s.Tasks[i].IsBreak() && s.Tasks[i+1].IsBreak() {
			s.Tasks[i+1].StartTime = s.Tasks[i].StartTime
			s.Tasks = append(s.Tasks[:i], s.Tasks[i+1:]...)
		} else if !s.Tasks[i].EndTime.Equal(s.Tasks[i+1].StartTime) {
			b := Break(s.Tasks[i].EndTime, s.Tasks[i+1].StartTime)
			s.Tasks = append(s.Tasks[:i+1], append([]*Task{&b}, s.Tasks[i+1:]...)...)
		}
	}
}

// UpdateCurrentTask checks the schedule's task list for the
// task scheduled for the time that the function is called,
// and updates the schedule's CurrentTask member accordingly.
// Note that the task will be nil if there is no currently
// scheduled task.
func (s *Schedule) UpdateCurrentTask() error {
	// For use with timer or change/request from client
	s.CurrentTask, s.CurrentID = s.Tasks.GetTaskAtTime(time.Now())
	if s.CurrentID == -1 {
		return InvalidScheduleError{}
	}

	return nil
}

// Print prints the schedule in a barebones format.
// Intended for debugging.
func (s Schedule) String() string {
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

	return sb.String()
}

// NewSchedule returns a new Schedule with the given TaskList as a Tasks
// member. It assumes the TaskList is consistent, and returns an empty
// schedule otherwise.
func NewSchedule(taskList TaskList) Schedule {
	if !taskList.IsConsistent() {
		return Schedule{}
	}
	currentTask, currentIdx := taskList.GetTaskAtTime(time.Now())
	return Schedule{
		Tasks:       taskList,
		CurrentTask: currentTask,
		CurrentID:   currentIdx,
	}
}

// BuildFromFile creates a schedule from a csv file with the given name
func BuildFromFile(fileName string) (*Schedule, error) {
	// TODO log?
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("BuildFromFile: %w", err)
	}
	r := csv.NewReader(f)
	taskList := []Task{}
	lc := 0
	var desc string
	var start time.Time
	var end time.Time
	var tag string

	line, err := r.Read()
	for err != io.EOF {
		if err != nil {
			return nil, fmt.Errorf("BuildFromFile: %w", err)
		}
		lc++
		if len(line) < 3 {
			return nil,
				errors.New("BuildFromFile: Field missing from line " + strconv.Itoa(lc))
		}
		desc = line[0]

		// Set task times to the current day (for now)
		now := time.Now()
		start, err = time.Parse(time.TimeOnly, line[1])
		if err != nil {
			return nil, errors.New("BuildFromFile: time value improperly formatted on line " + strconv.Itoa(lc))
		}
		start = time.Date(now.Year(), now.Month(), now.Day(), start.Hour(), start.Minute(), 0, 0, time.Local)
		end, err = time.Parse(time.TimeOnly, line[2])
		if err != nil {
			return nil, errors.New("BuildFromFile: time value improperly formatted on line " + strconv.Itoa(lc))
		}

		end = time.Date(now.Year(), now.Month(), now.Day(), end.Hour(), end.Minute(), 0, 0, time.Local)
		tag = strings.TrimSpace(line[3])
		var task Task
		if len(tag) > 0 {
			task = NewTask(desc, start, end).WithTag(tag)
		} else {
			task = NewTask(desc, start, end)
		}
		if task.IsEmpty() {
			return nil, errors.New("BuildFromFile: Task could not be created on line " + strconv.Itoa(lc))
		}

		taskList = append(taskList, task)
		line, err = r.Read()
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("BuildFromFile: %w", err)
		}
	}
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("BuildFromFile: %w", err)
	}

	tList, err := NewTaskList(taskList...)
	if err != nil {
		return nil, fmt.Errorf("BuildFromFile: %w", err)
	}

	current, index := tList.GetTaskAtTime(time.Now())
	return &Schedule{
		Tasks:       tList,
		CurrentTask: current,
		CurrentID:   index,
	}, nil
}
