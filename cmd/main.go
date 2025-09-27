package main

import (
	_ "embed"
	"log"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed consolas.ttf
var consolasTTF []byte

func main() {
	// Создаем контекст изображения
	const width, height = 480, 360
	dc := gg.NewContext(width, height)

	// Устанавливаем белый фон
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Загружаем шрифт
	face := consolasFace(40)
	dc.SetFontFace(face)

	// Устанавливаем цвет текста
	dc.SetRGB(0, 0, 0)

	// Добавляем текст
	dc.DrawString("Ноябрь", 20, 50)
	dc.DrawString("Пн Вт Ср Чт Пт Сб Вс", 20, 100)
	dc.DrawString("                1  2", 20, 150)
	dc.DrawString(" 3  4  5  6  7  8  9", 20, 200)
	dc.DrawString("10 11 12 13 14 15 16", 20, 250)
	dc.DrawString("17 18 19 20 21 22 23", 20, 300)
	dc.DrawString("24 25 26 27 28 29 30", 20, 350)

	// Сохраняем изображение
	dc.SavePNG("output_gg.png")
}

func consolasFace(size float64) font.Face {
	parsedFont, err := opentype.Parse(consolasTTF)
	if err != nil {
		log.Fatal(err)
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return face
}
