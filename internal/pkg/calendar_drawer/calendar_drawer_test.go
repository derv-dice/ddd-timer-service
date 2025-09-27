package calendar_drawer

import (
	"fmt"
	"testing"
)

func Test_generateMonth(t *testing.T) {

	m := NewMonth(2025, 9)
	fmt.Println(m)
}

func Test_generateCalendar(t *testing.T) {

	c := NewCalendar(testDate1, testDate2)
	fmt.Println(c)
}

func Test_MonthStrings(t *testing.T) {
	m := NewMonth(2025, 9)

	ss := m.strings()

	for i, s := range ss {
		fmt.Println(i, ":", s)
	}
}
