package helpers

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Platform struct {
	X, Y float64
	Width, Height float64
}

func DrawPlatform(screen *ebiten.Image, p Platform) {
	img := ebiten.NewImage(int(p.Width), int(p.Height))
	img.Fill(color.Black)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X, p.Y)

	screen.DrawImage(img, op)
	
}