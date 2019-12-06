package agenda

import (
	"testing"
	"time"
)

func TestParseAgenda(t *testing.T) {
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	timeNow := time.Now()
	standardClient := New()
	secondClient := New(WithParser(secondParser))

	entries := []struct {
		client   *Cron
		expr     string
		expected Schedule
	}{
		{secondClient, "+1m 0 5 * * * *", every5minWaitMin(time.Local, 1)},
		{secondClient, "@doy [06/19] * * *", nextDate(time.Local, 6, 19)},

		{secondClient, "0 5 * * * *", every5min(time.Local)},
		{standardClient, "5 * * * *", every5min(time.Local)},
		{secondClient, "CRON_TZ=UTC  0 5 * * * *", every5min(time.UTC)},
		{standardClient, "CRON_TZ=UTC  5 * * * *", every5min(time.UTC)},
		{secondClient, "CRON_TZ=Asia/Tokyo 0 5 * * * *", every5min(tokyo)},
		{secondClient, "@every 5m", ConstantDelaySchedule{5 * time.Minute}},
		{secondClient, "@midnight", midnight(time.Local)},
		{secondClient, "TZ=UTC  @midnight", midnight(time.UTC)},
		{secondClient, "TZ=Asia/Tokyo @midnight", midnight(tokyo)},
		{secondClient, "@yearly", annual(time.Local)},
		{secondClient, "@annually", annual(time.Local)},
		{
			client: secondClient,
			expr:   "* 5 * * * *",
			expected: &SpecSchedule{
				Second:   all(seconds),
				Minute:   1 << 5,
				Hour:     all(hours),
				Dom:      all(dom),
				Month:    all(months),
				Dow:      all(dow),
				Location: time.Local,
			},
		},
	}

	for _, c := range entries {
		actualParse, err := c.client.ParseAgenda(c.expr)
		if err != nil {
			t.Errorf("%s => unexpected error %v", c.expr, err)
		}

		actualTime := actualParse.Next(timeNow)
		targetTime := c.expected.Next(timeNow)
		if actualTime != targetTime {
			t.Errorf("\t%s \"%s\" => \nexpected\t%s \ngot\t\t%s", timeNow.String(), c.expr, targetTime.String(), actualTime.String())
		}

		// if !reflect.DeepEqual(actual, c.expected) {
		// 	t.Errorf("%s => expected %b, got %b", c.expr, c.expected, val)
		// }
	}
}

func BenchmarkParseAgenda(b *testing.B) {
	timeNow := time.Now()
	client := New(WithParser(secondParser))

	for n := 0; n < b.N; n++ {
		sch, _ := client.ParseAgenda("CRON_TZ=UTC  0 5 * * * *")
		sch.Next(timeNow)
		sch, _ = client.ParseAgenda("CRON_TZ=Asia/Tokyo 0 5 * * * *")
		sch.Next(timeNow)
		sch, _ = client.ParseAgenda("@annually")
		sch.Next(timeNow)
		sch, _ = client.ParseAgenda("TZ=UTC  @midnight")
		sch.Next(timeNow)
		sch, _ = client.ParseAgenda("@every 5m")
		sch.Next(timeNow)
	}
}

func BenchmarkParseDefault(b *testing.B) {
	timeNow := time.Now()
	client := New(WithParser(secondParser))

	for n := 0; n < b.N; n++ {
		sch, _ := client.Parse("CRON_TZ=UTC  0 5 * * * *")
		sch.Next(timeNow)
		sch, _ = client.Parse("CRON_TZ=Asia/Tokyo 0 5 * * * *")
		sch.Next(timeNow)
		sch, _ = client.Parse("@annually")
		sch.Next(timeNow)
		sch, _ = client.Parse("TZ=UTC  @midnight")
		sch.Next(timeNow)
		sch, _ = client.Parse("@every 5m")
		sch.Next(timeNow)
	}
}

func every5minWaitMin(loc *time.Location, min int) *SpecSchedule {
	return &SpecSchedule{1 << 0, 1 << (5 + min), all(hours), all(dom), all(months), all(dow), loc}
}

func nextDate(loc *time.Location, month int, day int) *SpecSchedule {
	return &SpecSchedule{1 << 0, 1 << 0, all(hours), 1 << day, 1 << month, all(dow), loc}
}
