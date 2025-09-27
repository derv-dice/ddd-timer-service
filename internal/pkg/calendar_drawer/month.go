package calendar_drawer

import (
	"bytes"
	"ddd-timer-service/internal/pkg/cells_drawer"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

type Month struct {
	Year  int
	Num   time.Month
	Weeks []*Week
}

func NewMonth(yearNum, monthNum int) *Month {
	month := time.Month(monthNum)

	out := Month{
		Year: yearNum,
		Num:  month,
	}

	var weeks []*Week
	weeks = append(weeks, NewWeek())
	weekNum := 0

	// Получаем последний день месяца
	lastDay := time.Date(yearNum, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	for day := 1; day <= lastDay; day++ {
		date := time.Date(yearNum, month, day, 0, 0, 0, 0, time.UTC)
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

			// Добавляем новую неделю только если это не последний день месяца
			if day < lastDay {
				weeks = append(weeks, NewWeek())
				weekNum += 1
			}
		}
	}

	out.Weeks = weeks

	return &out
}

func (m *Month) Name() string {
	return monthToName[m.Num]
}

func (m *Month) String() string {
	br := strings.Builder{}

	br.WriteString(monthToName[m.Num])
	br.WriteString("\n")

	for _, name := range weekNames {
		br.WriteString(name)
		br.WriteString(" ")
	}
	br.WriteString("\n")

	for _, w := range m.Weeks {
		for _, d := range w {
			if d == nil {
				br.WriteString("   ")
			} else {
				br.WriteString(fmt.Sprintf("%2d ", d.Num))
			}
		}
		br.WriteString("\n")
	}

	return br.String()
}

func (m *Month) PNG() ([]byte, image.Image, error) {
	width := monthPNGWidth
	height := monthPNGHeight

	dc := gg.NewContext(width, height)

	face := consolasFace(40)
	dc.SetFontFace(face)

	dc.SetRGBA255(
		int(cells_drawer.GreyColor.R),
		int(cells_drawer.GreyColor.G),
		int(cells_drawer.GreyColor.B),
		int(cells_drawer.GreyColor.A))

	arr := m.strings()

	offset := 50
	for _, str := range arr {
		dc.DrawString(str, 20, float64(offset))
		offset += 50
	}

	dc.SetLineWidth(3)

	dc.SetRGBA255(
		int(cells_drawer.GreyColor.R),
		int(cells_drawer.GreyColor.G),
		int(cells_drawer.GreyColor.B),
		int(cells_drawer.GreyColor.A))

	dc.DrawRectangle(10, 10, float64(width-20), float64(height-20))
	dc.Stroke()

	buff := new(bytes.Buffer)

	if err := dc.EncodePNG(buff); err != nil {
		return nil, nil, err
	}

	return buff.Bytes(), dc.Image(), nil
}

func (m *Month) progressMask(from, to, now time.Time) ([]byte, image.Image, error) {
	width := monthPNGWidth
	height := monthPNGHeight

	dc := gg.NewContext(width, height)

	face := consolasFace(40)
	dc.SetFontFace(face)

	dc.SetRGBA255(
		int(cells_drawer.GreyColor.R),
		int(cells_drawer.GreyColor.G),
		int(cells_drawer.GreyColor.B),
		int(cells_drawer.GreyColor.A))

	arr := m.stringsOpacityMask(from, to, now)

	offset := 50
	for _, str := range arr {
		dc.DrawString(str, 20, float64(offset))
		offset += 50
	}

	buff := new(bytes.Buffer)

	if err := dc.EncodePNG(buff); err != nil {
		return nil, nil, err
	}

	return buff.Bytes(), dc.Image(), nil
}

func (m *Month) strings() []string {
	var res []string

	res = append(res, monthToName[m.Num])

	var wNames string
	for _, name := range weekNames {
		wNames += name + " "
	}
	res = append(res, wNames)

	for _, w := range m.Weeks {
		var str string
		for _, d := range w {
			if d == nil {
				str += "   "
			} else {
				str += fmt.Sprintf("%2d ", d.Num)
			}
		}

		res = append(res, str)
	}

	return res
}

func (m *Month) stringsOpacityMask(from, to, now time.Time) []string {
	var res []string

	res = append(res, " ", " ") // имя месяца и названия недель

	for _, w := range m.Weeks {
		var str string
		for _, d := range w {
			if d == nil {
				str += "   "
			} else {
				date := time.Date(m.Year, m.Num, d.Num, 0, 0, 0, 0, time.UTC)

				// Сегодняшний день
				if m.Year == now.Year() && m.Num == now.Month() && d.Num == now.Day() {

					fmt.Println(date.String())
					str += fmt.Sprintf("%s ", fontMaskHighlight)
					continue
				}

				// Время с момента начала службы
				if (date.After(from) && date.Before(to)) && date.Before(time.Now().Round(time.Hour*24)) {
					if d.Num < 10 {
						str += fmt.Sprintf("%s ", fontMaskShort)
					} else {
						str += fmt.Sprintf("%s ", fontMaskLong)
					}

					continue
				}

				// Если маска не нужна, оставляем пустое место
				str += fmt.Sprintf("%s ", fontMaskNone)
			}
		}

		res = append(res, str)
	}

	return res
}

func (m *Month) yearSeasonIndex() int {
	switch m.Num {
	case 12, 3, 6, 9:
		return 0
	case 1, 4, 7, 10:
		return 1
	case 2, 5, 8, 11:
		return 2
	}

	return -1
}

func (m *Month) seasonType() yearSeasonType {
	switch m.Num {
	case 12, 1, 2:
		return Winter
	case 3, 4, 5:
		return Spring
	case 6, 7, 8:
		return Summer
	case 9, 10, 11:
		return Autumn
	}

	return -1
}
