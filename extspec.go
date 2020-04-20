package agenda

import (
	"sort"
	"strings"
	"time"
)

type ScheduleType uint8

// Types of schedules
const (
	SchConstantDelay ScheduleType = iota
	SchDefaultCron
	SchDoY
	SchAt
)

// SpecAgenda specifies a duty cycle (to the second granularity).
// It is computed for every day and re-uses the default cron parser.
type SpecAgenda struct {
	SpecSchedule
	ConstantDelaySchedule
	Type    ScheduleType
	DoY, At string
	Wait    time.Duration
	FirstExec	bool
}

// Next returns the next time this schedule is activated, greater than the given
// time.  If no time can be found to satisfy the schedule, return the zero time.
func (s *SpecAgenda) Next(t time.Time) (next time.Time) {
	//Consider delay variable
	// Add delay one time only
	if s.FirstExec {
		s.FirstExec = false
	} else {
		t = t.Add(-s.Wait)
	}

	switch s.Type {
	case SchDoY:
		next = s.doyNext(t)
	case SchAt:
		date, err := time.Parse(time.RFC3339, s.At)
		if err != nil {
			panic("could not convert command @at, err: " + err.Error())
		}
		if date.Before(time.Now()) {
			return t.AddDate(142, 0, 0) //I hope this schedule gets canceled until that...
		}
		next = date
	case SchConstantDelay:
		//next = s.ConstantDelaySchedule.Next(t.Add(s.Wait))
		next = s.ConstantDelaySchedule.Next(t)
	case SchDefaultCron:
		next = s.SpecSchedule.Next(t)
	}

	next = next.Add(s.Wait)
	return
}

func (s *SpecAgenda) doyNext(t time.Time) time.Time {
	//DoY is a list of strings in the format MM/DD
	dayPieces := strings.Split(s.DoY, " ")
	sort.Strings(dayPieces)

	thisYear, thisMonth, thisDay := t.Date()

	//Find the next day
	for _, date := range dayPieces {
		dateParts := strings.Split(date, "/")
		month, err := mustParseInt(dateParts[0])
		if err != nil {
			panic(err)
		}
		day, err := mustParseInt(dateParts[1])
		if err != nil {
			panic(err)
		}

		s.Dom = getBits(day, day, 1)
		s.Month = getBits(month, month, 1)
		scheduleTime := s.SpecSchedule.Next(t)

		if month == uint(thisMonth) {
			if day == uint(thisDay) {
				//Check if the time is still doable
				if scheduleTime.Before(t.AddDate(0, 0, 1)) {
					//This hour is still able to be executed...
					// fmt.Println("Next 24h " + scheduleTime.String())
					return scheduleTime
				}
			} else if day > uint(thisDay) {
				if scheduleTime.Before(t.AddDate(1, 0, 0)) {
					//This hour is still able to be executed...
					// fmt.Println("Still this year " + scheduleTime.String())
					return scheduleTime
				}
			}
		} else if month > uint(thisMonth) {
			return scheduleTime
		}
	}

	//No date was found... wrap the year!
	t = time.Date(thisYear+1, time.January, 1, 0, 0, 0, 0, s.Location)

	for _, date := range dayPieces {
		dateParts := strings.Split(date, "/")
		month, err := mustParseInt(dateParts[0])
		if err != nil {
			panic(err)
		}
		day, err := mustParseInt(dateParts[1])
		if err != nil {
			panic(err)
		}

		s.Dom = getBits(day, day, 1)
		s.Month = getBits(month, month, 1)
		scheduleTime := s.SpecSchedule.Next(t)
		next1Year := t.AddDate(1, 0, 0)

		if scheduleTime.Before(next1Year) {
			//This hour is still able to be executed...
			// fmt.Println("Next year " + scheduleTime.String())
			return scheduleTime
		}
	}

	panic("no time found")
}
