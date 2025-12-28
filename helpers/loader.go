package helpers

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

var Assets embed.FS

func Init(fs embed.FS) {
	Assets = fs

	fontBytes, err := Assets.ReadFile("assets/Inter_24pt-Bold.ttf")
	if err != nil {
		panic(err)
	}

	InitFonts(fontBytes)
}

func LoadImage(path string) *ebiten.Image {
	data, err := Assets.ReadFile(path)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func LoadImages(path string) [] *ebiten.Image {
	files, err := filepath.Glob("assets/walk_*.png")

	if err != nil {
		log.Fatal(err)
	}

	var res [] *ebiten.Image 

	for _, f := range files {
		res = append(res, LoadImage(f))
	}

	return res 
}