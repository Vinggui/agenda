Under development! Experimental stage. Agenda package adds a couple of milliseconds to the original cron (20-40ms).

For general cron functions, refer to [![GoDoc](http://godoc.org/github.com/robfig/cron?status.png)](http://godoc.org/github.com/robfig/cron)

# agenda
To download agenda
	go get github.com/Vinggui/agenda

Import it in your program as:

	import "github.com/Vinggui/agenda"

It requires Go 1.11 or later due to usage of Go Modules.

### Background - Cron spec format

There are two cron spec formats in common usage:

- The "standard" cron format, described on [the Cron wikipedia page] and used by
  the cron Linux system utility.

- The cron format used by [the Quartz Scheduler], commonly used for scheduled
  jobs in Java software

[the Cron wikipedia page]: https://en.wikipedia.org/wiki/Cron
[the Quartz Scheduler]: http://www.quartz-scheduler.org/documentation/quartz-2.x/tutorials/crontrigger.html

The original version of this package included an optional "seconds" field, which
made it incompatible with both of these formats. Now, the "standard" format is
the default format accepted, and the Quartz format is opt-in.

### Background - Agenda spec format

In additional, if user selects the agenda parser instead, you can do:

#### Select the days of the year (DoY). A list of MM/DD can be passed inside the []
	@doy [12/04 03/02] * *         <-- Default parser
	@doy [12/04 03/02] * * *       <-- Quartz format which includes the seconds
 Obs: the dates do not need to be sorted.

#### Select the *at* a specific date
	@at 2019-12-02T15:04:05Z
	@at 2020-02-29T15:04:05-03:00
 Obs: Command "at" is now using the RFC3339 format.
      Only one date can be passed

#### Add a delay to the schedule. Delay is give by a "+<time.ParseDuration>"
	+14m21s */10 * * * *                  <-- Default parser
	+14m21s */10 * * * * *                <-- Quartz format which includes the seconds
	+2h @doy [12/04 03/02] * * *
	+12s @at 2020-02-29T15:04:05-03:00
