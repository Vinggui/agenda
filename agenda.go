package agenda

import (
	"errors"
	"strings"
	"time"
)

//Errors
var (
	ErrEmptySpec = errors.New("empty spec argument")
	ErrPastDate  = errors.New("input date is already in the past")
)

func copySchedules(schAgenda *SpecAgenda, schDefault Schedule) {
	defaultSchedule := schDefault.(*SpecSchedule)
	schAgenda.Second = defaultSchedule.Second
	schAgenda.Minute = defaultSchedule.Minute
	schAgenda.Hour = defaultSchedule.Hour
	schAgenda.Dom = defaultSchedule.Dom
	schAgenda.Month = defaultSchedule.Month
	schAgenda.Dow = defaultSchedule.Dow
	schAgenda.Location = defaultSchedule.Location
}

/*
Parse extends the entries provided by cron package
*/
func (c *Cron) Parse(spec string) (sch Schedule, err error) {
	//The +<duration> indicates a delays to start the following specs.
	// <duration> is expressed in time. (readable for time.ParseDuration)

	//@doy accepts values from 1-366 (Feb 29th is accounted)
	//@at is a single date in time using time.String default format
	//    like RFC3339 "2006-01-02T15:04:05Z07:00"

	customSchedule := &SpecAgenda{}
	sch = customSchedule
	if strings.HasPrefix(spec, "+") {
		SpecStartPoint := strings.Index(spec, " ")
		customSchedule.Delay, err = time.ParseDuration(spec[1:SpecStartPoint])
		if err != nil {
			return
		}
		spec = spec[SpecStartPoint+1:]
		// customSchedule.Delay = delay
	}

	var genericSchedule Schedule
	if strings.HasPrefix(spec, "@doy") {
		atStartPoint := strings.Index(spec, "[") + 1
		specStartPoint := strings.Index(spec, "]")

		dates := spec[atStartPoint:specStartPoint]

		spec = spec[specStartPoint+2:]

		// the possible cron formats, including second or not
		genericSchedule, err = c.parser.Parse(spec + " * * *")
		if err != nil {
			return
		}
		copySchedules(customSchedule, genericSchedule)

		customSchedule.DoY = dates

		return customSchedule, nil
	} else if strings.HasPrefix(spec, "@at") {
		firstSpace := strings.Index(spec, " ") + 1
		customSchedule.At = spec[firstSpace:]
		var date time.Time
		date, err = time.Parse(time.RFC3339, customSchedule.At)
		if err != nil {
			return
		}

		if date.Before(time.Now()) {
			err = ErrPastDate
			return
		}

		return

	} else {
		// the possible cron formats, including second or not
		genericSchedule, err = c.parser.Parse(spec)
		if err != nil {
			return
		}
		copySchedules(customSchedule, genericSchedule)
		return
	}
}
