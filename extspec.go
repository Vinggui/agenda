package agenda

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

// SpecAgenda specifies a duty cycle (to the second granularity).
// It is computed for every day and re-uses the default cron parser.
type SpecAgenda struct {
	SpecSchedule
	DoY, At string
	Delay   time.Duration
}

// Next returns the next time this schedule is activated, greater than the given
// time.  If no time can be found to satisfy the schedule, return the zero time.
func (s *SpecAgenda) Next(t time.Time) (next time.Time) {
	//Consider delay variable

	if s.Delay != 0 {
		t = t.Add(-s.Delay)
	}

	if len(s.DoY) > 0 {
		next = s.doyNext(t)
	} else if len(s.At) > 0 {
		date, err := time.Parse(time.RFC3339, s.At)
		if err != nil {
			panic("could not convert command @at, err: " + err.Error())
		}
		if date.Before(time.Now()) {
			return t.AddDate(142, 0, 0) //I hope this schedule gets canceled until that...
		}
		next = date
	} else {
		next = s.SpecSchedule.Next(t)
	}

	if s.Delay != 0 {
		next = next.Add(s.Delay)
	}
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
			log.Fatal(err)
		}
		day, err := mustParseInt(dateParts[1])
		if err != nil {
			log.Fatal(err)
		}

		if month >= uint(thisMonth) {
			s.Dom = getBits(day, day, 1)
			s.Month = getBits(month, month, 1)
			scheduleTime := s.SpecSchedule.Next(t)
			if day == uint(thisDay) {
				//Check if the time is still doable
				if scheduleTime.Before(t.AddDate(0, 0, 1)) {
					//This hour is still able to be executed...
					fmt.Println("Next 24h " + scheduleTime.String())
					return scheduleTime
				}
			} else if day > uint(thisDay) {
				if scheduleTime.Before(t.AddDate(1, 0, 0)) {
					//This hour is still able to be executed...
					fmt.Println("Still this year " + scheduleTime.String())
					return scheduleTime
				}
			}
		}
	}

	//No date was found... wrap the year!
	t = time.Date(thisYear+1, time.January, 1, 0, 0, 0, 0, s.Location)

	for _, date := range dayPieces {
		dateParts := strings.Split(date, "/")
		month, err := mustParseInt(dateParts[0])
		if err != nil {
			log.Fatal(err)
		}
		day, err := mustParseInt(dateParts[1])
		if err != nil {
			log.Fatal(err)
		}

		s.Dom = getBits(day, day, 1)
		s.Month = getBits(month, month, 1)
		scheduleTime := s.SpecSchedule.Next(t)
		next1Year := t.Add(time.Hour * 24 * 365)

		if scheduleTime.Before(next1Year) {
			//This hour is still able to be executed...
			fmt.Println("Next year " + scheduleTime.String())
			return scheduleTime
		}
	}

	panic("no time found")
}
