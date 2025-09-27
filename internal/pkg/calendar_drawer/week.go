package calendar_drawer

type Week [7]*Day // starts from monday

func NewWeek() *Week {
	w := Week([7]*Day{nil, nil, nil, nil, nil, nil, nil})
	return &w
}
