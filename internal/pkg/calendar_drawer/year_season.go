package calendar_drawer

type yearSeason struct {
	Type   yearSeasonType
	Months [3]*Month
}

func (s *yearSeason) monthsNames() []string {
	var names []string

	for _, m := range s.Months {
		if m != nil {
			names = append(names, m.Name())
		}
	}

	return names
}

type yearSeasonType int

func (t *yearSeasonType) name() string {
	return seasonTypeToName[*t]
}
