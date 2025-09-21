package calendar_drawer

import (
	"fmt"
	"testing"
)

func Test_generateMonth(t *testing.T) {

	m := generateMonth(2025, 11)

	fmt.Printf("%s\n", MonthToName[m.Num])

	for _, name := range WeekNames {
		fmt.Printf("%s ", name)
	}
	fmt.Println()

	for _, w := range m.Weeks {
		for _, d := range w {
			if d == nil {
				fmt.Printf("   ")
			} else {
				fmt.Printf("%2d ", d.Num)
			}
		}
		fmt.Println()
	}
}
