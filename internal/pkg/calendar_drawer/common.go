package calendar_drawer

import "time"

var WeekNames = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}

var WeekDayToNameShort = map[time.Weekday]string{
	time.Monday:    WeekNames[0],
	time.Tuesday:   WeekNames[1],
	time.Wednesday: WeekNames[2],
	time.Thursday:  WeekNames[3],
	time.Friday:    WeekNames[4],
	time.Saturday:  WeekNames[5],
	time.Sunday:    WeekNames[6],
}

var MonthToName = map[time.Month]string{
	time.January:   "Январь",
	time.February:  "Февраль",
	time.March:     "Март",
	time.April:     "Апрель",
	time.May:       "Май",
	time.June:      "Июнь",
	time.July:      "Июль",
	time.August:    "Август",
	time.September: "Сентябрь",
	time.October:   "Октябрь",
	time.November:  "Ноябрь",
	time.December:  "Декабрь",
}

type Calendar struct {
	Years []*Year
}

type Year struct {
	Num    int
	Months [12]*Month
}

type Month struct {
	Year  int
	Num   time.Month
	Weeks []*Week
}

type Week [7]*Day // starts from monday

type Day struct {
	Num        int
	WeekdayNum int
	IsDayOff   bool
}

func emptyWeek() *Week {
	w := Week([7]*Day{nil, nil, nil, nil, nil, nil, nil})
	return &w
}
