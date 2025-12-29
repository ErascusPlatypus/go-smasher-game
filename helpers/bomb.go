package helpers

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BombGravity = 0.4
	BombRadius  = 80.0

	BombMaxTravel     = 600.0
	BombBounceLossY   = 0.6
	BombBounceLossX   = 0.75
	BombMinBounceVelY = 2.0
)

var bombSprite = "assets/bomb_anim.png"

type Bomb struct {
	X, Y         float64
	travelled    float64
	VelX, VelY   float64
	Active       bool
	Exploded     bool
	HasDamaged   bool
	ExplosionTTL int

	bounces int

	timer  *Timer
	sprite *ebiten.Image
}

func NewBomb(x, y float64, facingRight bool) *Bomb {
	vx := 6.0
	if !facingRight {
		vx = -vx
	}

	sprite := LoadImage(bombSprite)

	return &Bomb{
		X:            x,
		Y:            y,
		VelX:         vx,
		VelY:         -8,
		Active:       true,
		Exploded:     false,
		ExplosionTTL: 0,
		timer:        NewTimer(1200 * time.Millisecond),
		sprite:       sprite,
	}

}

func (b *Bomb) Update() {
	if !b.Active {
		return
	}

	dx := b.VelX
	b.X += dx
	b.travelled += abs(dx)

	b.VelY += BombGravity
	b.Y += b.VelY

	if b.travelled >= BombMaxTravel {
		b.Exploded = true
		b.Active = false
		return
	}

	if b.Y >= 690.0 {
		b.Y = 680.0

		b.VelY = -b.VelY * BombBounceLossY
		b.VelX *= BombBounceLossX
		b.bounces++

		if abs(b.VelY) < BombMinBounceVelY {
			b.Exploded = true
			b.Active = false
			return
		}
	}

	if b.timer.IsReady() {
		b.Exploded = true
		b.Active = false
	}
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func (b *Bomb) Draw(screen *ebiten.Image) {
	if !b.Active {
		return
	}

	op := &ebiten.DrawImageOptions{}

	w := float64(b.sprite.Bounds().Dx())
	h := float64(b.sprite.Bounds().Dy())

	op.GeoM.Translate(-w/2, -h/2)

	op.GeoM.Scale(0.25, 0.25)

	op.GeoM.Translate(
		b.X+w*0.25/2,
		b.Y+h*0.25/2,
	)

	screen.DrawImage(b.sprite, op)
}

func (b *Bomb) HitsPlayer(p *Player) bool {
	dx := (p.X) - b.X
	dy := (p.Y + p.Height/2) - b.Y
	return dx*dx+dy*dy <= BombRadius*BombRadius
}

func (b *Bomb) GetRect() Rect {
	w, h := b.sprite.Bounds().Dx(), b.sprite.Bounds().Dy()
	return NewRect(b.X, b.Y, float64(w), float64(h))
}
