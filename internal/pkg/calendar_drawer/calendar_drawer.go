package calendar_drawer

import (
	"time"
)

type CalendarDrawer struct{}

func NewCalendarDrawer() *CalendarDrawer {
	return &CalendarDrawer{}
}

func (c *CalendarDrawer) GenerateCalendarPNG() {

}

func (c *CalendarDrawer) GenerateCalendar(dateFrom, dateTo time.Time) *Calendar {
	// Округление до дня
	dateFrom = dateFrom.Round(24 * time.Hour)
	dateTo = dateTo.Round(24 * time.Hour)

	return &Calendar{}
}

func generateMonth(yearNum, monthNum int) *Month {
	out := Month{
		Year: yearNum,
		Num:  time.Month(monthNum),
	}

	var weeks []*Week
	weeks = append(weeks, emptyWeek())
	weekNum := 0

	for day := 1; ; day++ {
		date := time.Date(yearNum, time.Month(monthNum), day, 0, 0, 0, 0, time.UTC)
		if date.Month() != time.Month(monthNum) {
			break
		}

		wdNum := date.Weekday()

		wd := &Day{
			Num:        day,
			WeekdayNum: int(wdNum),
			IsDayOff:   wdNum == time.Sunday || wdNum == time.Saturday,
		}

		switch wdNum {
		case time.Monday:
			weeks[weekNum][0] = wd
		case time.Tuesday:
			weeks[weekNum][1] = wd
		case time.Wednesday:
			weeks[weekNum][2] = wd
		case time.Thursday:
			weeks[weekNum][3] = wd
		case time.Friday:
			weeks[weekNum][4] = wd
		case time.Saturday:
			weeks[weekNum][5] = wd
		case time.Sunday:
			weeks[weekNum][6] = wd
			weeks = append(weeks, emptyWeek())
			weekNum += 1
		}
	}

	out.Weeks = weeks

	return &out
}
