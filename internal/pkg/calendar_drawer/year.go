package calendar_drawer

type Year struct {
	Num    int
	Months []*Month
}

func (y *Year) Seasons() YearSeasonsList {
	var seasons []*yearSeason

	prevMonthSeasonType := y.Months[0].seasonType()
	seasons = append(seasons, &yearSeason{
		Type:   prevMonthSeasonType,
		Months: [3]*Month{},
	})

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

	return seasons
}
