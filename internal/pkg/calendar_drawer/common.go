package calendar_drawer

import (
	_ "embed"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed src/consolas.ttf
var consolasTTF []byte

func consolasFace(size float64) font.Face {
	parsedFont, err := opentype.Parse(consolasTTF)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	return face
}

var weekNames = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}

var weekDayToNameShort = map[time.Weekday]string{
	time.Monday:    weekNames[0],
	time.Tuesday:   weekNames[1],
	time.Wednesday: weekNames[2],
	time.Thursday:  weekNames[3],
	time.Friday:    weekNames[4],
	time.Saturday:  weekNames[5],
	time.Sunday:    weekNames[6],
}

var monthToName = map[time.Month]string{
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

const (
	Spring yearSeasonType = iota + 10
	Summer
	Autumn
	Winter
)

var seasonTypeToName = map[yearSeasonType]string{
	Spring: "Весна",
	Summer: "Лето",
	Autumn: "Осень",
	Winter: "Зима",
}

const (
	monthPNGHeight     = 420
	monthPNGWidth      = 500
	seasonsTitleHeight = 70
)

const (
	fontMaskShort     = " N"
	fontMaskLong      = "NN"
	fontMaskHighlight = "__"
	fontMaskNone      = "  "
)
