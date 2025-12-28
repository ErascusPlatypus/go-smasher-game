package helpers

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	X, Y      float64
	VelX      float64
	sprite    *ebiten.Image
	direction bool
	Active    bool
}

const (
	bulletSpeed = 10.0
)

var bulletSprite = "assets/bullet.png"

func NewBullet(x, y float64, dir bool) *Bullet {
	sprite := LoadImage(bulletSprite)

	return &Bullet{
		X: x, Y: y,
		sprite:    sprite,
		Active:    true,
		direction: dir,
	}
}

func (b *Bullet) Update() {
	if b.direction == true {
		b.X += bulletSpeed
	} else {
		b.X -= bulletSpeed
	}

	if b.Y < -50 || b.Y > 1200 || b.X < -50 || b.X > 1400 {
		b.Active = false
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if !b.Active {
		return
	}

	opts := &ebiten.DrawImageOptions{}
	w, h := b.sprite.Bounds().Dx(), b.sprite.Bounds().Dy()
	opts.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	opts.GeoM.Rotate(math.Pi / 2)

	opts.GeoM.Translate(
		b.X+float64(w)/2,
		b.Y+float64(h)/2,
	)

	screen.DrawImage(b.sprite, opts)
}

