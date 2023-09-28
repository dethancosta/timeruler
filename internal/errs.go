package internal


type InvalidTimeError struct {
	msg string
}
type TimeConflictError struct {
	msg string
}
type InvalidScheduleError struct {
	msg string
}

func (e InvalidTimeError) Error() string {
	return e.msg
}

func (e TimeConflictError) Error() string {
	return e.msg
}

func (e InvalidScheduleError) Error() string {
	return e.msg
}
