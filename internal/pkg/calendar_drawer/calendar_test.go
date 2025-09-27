package calendar_drawer

import (
	"fmt"
	"testing"
)

func TestCalendar_Seasons(t *testing.T) {

	c := NewCalendar(testDate1, testDate2)

	for _, s := range c.Seasons() {
		fmt.Printf("%s %v\n", s.Type.name(), s.monthsNames())
	}
}
