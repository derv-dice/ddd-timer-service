package calendar_drawer

import (
	"bytes"
	"ddd-timer-service/internal/pkg/cells_drawer"
	"image"
	"time"

	"github.com/fogleman/gg"
)

type YearSeasonsList []*yearSeason

func (s *YearSeasonsList) PNG() ([]byte, image.Image, error) {

	width := monthPNGWidth * len(*s)
	height := monthPNGHeight*3 + seasonsTitleHeight

	dc := gg.NewContext(width, height)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	face := consolasFace(40)
	dc.SetFontFace(face)
	dc.SetRGB(0, 0, 0)

	offsetY := 0
	offsetX := 0

	for _, season := range *s {
		offsetY = 0

		dc.DrawString(season.Type.name(), float64(offsetX+20), float64(offsetY+seasonsTitleHeight-20))

		offsetY = seasonsTitleHeight

		for _, month := range season.Months {
			if month == nil {
				offsetY += monthPNGHeight
				continue
			}

			_, img, err := month.PNG()
			if err != nil {
				return nil, nil, err
			}

			dc.DrawImage(img, offsetX, offsetY)

			offsetY += monthPNGHeight
		}

		offsetX += monthPNGWidth
	}

	dc.SetLineWidth(3)
	dc.SetRGB(0, 0, 0)
	dc.DrawRectangle(2, 2, float64(width-4), float64(height-4))
	dc.Stroke()

	buff := new(bytes.Buffer)

	if err := dc.EncodePNG(buff); err != nil {
		return nil, nil, err
	}

	return buff.Bytes(), dc.Image(), nil
}

func (s *YearSeasonsList) PNGWithProgressMask(from, to, now time.Time) ([]byte, image.Image, error) {
	width := monthPNGWidth * len(*s)
	height := monthPNGHeight*3 + seasonsTitleHeight

	dc := gg.NewContext(width, height)

	_, img, err := s.PNG()
	if err != nil {
		return nil, nil, err
	}

	_, imgMask, err := s.progressMask(from, to, now)
	if err != nil {
		return nil, nil, err
	}

	dc.DrawImage(img, 0, 0)
	dc.DrawImage(imgMask, 0, 0)

	buff := new(bytes.Buffer)

	if err = dc.EncodePNG(buff); err != nil {
		return nil, nil, err
	}

	return buff.Bytes(), dc.Image(), nil
}

func (s *YearSeasonsList) progressMask(from, to, now time.Time) ([]byte, image.Image, error) {
	width := monthPNGWidth * len(*s)
	height := monthPNGHeight*3 + seasonsTitleHeight

	dc := gg.NewContext(width, height)

	face := consolasFace(40)
	dc.SetFontFace(face)

	dc.SetRGBA255(
		int(cells_drawer.GreyColor.R),
		int(cells_drawer.GreyColor.G),
		int(cells_drawer.GreyColor.B),
		int(cells_drawer.GreyColor.A))

	offsetY := 0
	offsetX := 0

	for _, season := range *s {
		offsetY = 0

		offsetY = seasonsTitleHeight

		for _, month := range season.Months {
			if month == nil {
				offsetY += monthPNGHeight
				continue
			}

			_, img, err := month.progressMask(from, to, now)
			if err != nil {
				return nil, nil, err
			}

			dc.DrawImage(img, offsetX, offsetY)

			offsetY += monthPNGHeight
		}

		offsetX += monthPNGWidth
	}

	buff := new(bytes.Buffer)

	if err := dc.EncodePNG(buff); err != nil {
		return nil, nil, err
	}

	return buff.Bytes(), dc.Image(), nil
}
