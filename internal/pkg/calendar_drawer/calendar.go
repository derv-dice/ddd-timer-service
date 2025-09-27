package calendar_drawer

import (
	"fmt"
	"strings"
	"time"
)

type Calendar struct {
	Years []*Year
}

func NewCalendar(dateFrom, dateTo time.Time) *Calendar {
	if dateFrom.After(dateTo) {
		dateFrom, dateTo = dateTo, dateFrom
	}

	currentYear, currentMonth := dateFrom.Year(), dateFrom.Month()
	endYear, endMonth := dateTo.Year(), dateTo.Month()

	var years []*Year

	yearsCounter := 0
	years = append(years, &Year{
		Num:    currentYear,
		Months: []*Month{},
	})

	for {
		years[yearsCounter].Months = append(years[yearsCounter].Months, NewMonth(currentYear, int(currentMonth)))

		if currentYear == endYear && currentMonth == endMonth {
			break
		}

		currentMonth++
		if currentMonth > 12 {
			currentMonth = 1
			currentYear++

			yearsCounter++
			years = append(years, &Year{
				Num:    currentYear,
				Months: []*Month{},
			})
		}
	}

	return &Calendar{
		Years: years,
	}
}

func (c *Calendar) String() string {
	br := strings.Builder{}

	for _, y := range c.Years {
		br.WriteString(fmt.Sprintf("        %d        \n", y.Num))
		for _, m := range y.Months {
			br.WriteString(m.String())
			br.WriteString("\n")
		}
	}

	return br.String()
}

func (c *Calendar) Seasons() YearSeasonsList {
	var seasons []*yearSeason

	prevMonthSeasonType := c.Years[0].Months[0].seasonType()
	seasons = append(seasons, &yearSeason{
		Type:   prevMonthSeasonType,
		Months: [3]*Month{},
	})

	for _, y := range c.Years {
		for _, month := range y.Months {
			if month.seasonType() != prevMonthSeasonType {
				prevMonthSeasonType = month.seasonType()
				seasons = append(seasons, &yearSeason{
					Type:   prevMonthSeasonType,
					Months: [3]*Month{},
				})
			}

			seasons[len(seasons)-1].Months[month.yearSeasonIndex()] = month
		}
	}

	return seasons
}
