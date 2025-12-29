package helpers

import (
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	DefaultFont font.Face
	WinnerFont font.Face
)

func InitFonts(fsData []byte) {
	tt, err := opentype.Parse(fsData)
	if err != nil {
		log.Fatal(err)
	}

	DefaultFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	WinnerFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    60,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}
