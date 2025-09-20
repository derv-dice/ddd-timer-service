package cells_drawer

import (
	"bytes"
	"ddd-timer-service/internal/pkg/stats_counter"
	"fmt"
	"image"
	"image/color"
	"image/png"
)

const (
	cellsCountX = 26

	cellSize = 22
)

var (
	whiteColor = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	greyColor  = color.RGBA{R: 0x29, G: 0x29, B: 0x29, A: 0xff}
	greenColor = color.RGBA{R: 0x95, G: 0xd1, B: 0xa9, A: 0xff}
)

type CellsDrawer struct {
}

func NewCellsDrawer() *CellsDrawer {
	return &CellsDrawer{}
}

func (c *CellsDrawer) NewCellsImagePNG(stats stats_counter.Stats) ([]byte, error) {
	x1 := stats.PassedDays()
	x2 := stats.LeftDays()
	n := int(x1 + x2)

	cellsCountY := (n / cellsCountX) + (n % cellsCountX)

	imgSizeX := cellSize * cellsCountX
	imgSizeY := cellSize * cellsCountY

	img := image.NewRGBA(image.Rect(0, 0, imgSizeX, imgSizeY))

	greenCellsLeft := int(x1)

	for yi := 0; yi < cellsCountY; yi++ {
		fillGreenTo := 0

		if greenCellsLeft >= cellsCountX {
			fillGreenTo = cellsCountX
		}

		if greenCellsLeft < cellsCountX && greenCellsLeft > 0 {
			fillGreenTo = greenCellsLeft
		}

		drawRowOfSquares(img, yi, cellsCountX, fillGreenTo)

		greenCellsLeft -= cellsCountX
	}

	b := bytes.Buffer{}
	err := png.Encode(&b, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode png: %v", err)
	}

	return b.Bytes(), nil
}

func drawRowOfSquares(img *image.RGBA, offsetY, rowLen, fillGreenTo int) {
	y := offsetY * cellSize
	x := 0

	for xi := 0; xi < rowLen; xi++ {
		if x < fillGreenTo*cellSize {
			drawSquare(img, x, y, 22, greenColor, greyColor, 1)
		} else {
			drawSquare(img, x, y, 22, whiteColor, greyColor, 1)
		}

		x += cellSize
	}
}

func drawSquare(img *image.RGBA, x, y, size int, fillColor, borderColor color.Color, borderWidth int) {
	// Заливка
	for i := y; i < y+size; i++ {
		for j := x; j < x+size; j++ {
			img.Set(j, i, fillColor)
		}
	}

	// Рамка
	if borderWidth > 0 {
		for bw := 0; bw < borderWidth; bw++ {
			// Верхняя граница
			for j := x + bw; j < x+size-bw; j++ {
				img.Set(j, y+bw, borderColor)
			}
			// Нижняя граница
			for j := x + bw; j < x+size-bw; j++ {
				img.Set(j, y+size-1-bw, borderColor)
			}
			// Левая граница
			for i := y + bw; i < y+size-bw; i++ {
				img.Set(x+bw, i, borderColor)
			}
			// Правая граница
			for i := y + bw; i < y+size-bw; i++ {
				img.Set(x+size-1-bw, i, borderColor)
			}
		}
	}
}
