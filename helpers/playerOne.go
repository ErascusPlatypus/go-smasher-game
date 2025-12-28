package helpers

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	JumpVelocity = -12.0
	Gravity      = 0.4
	GroundY      = 560
	VelX         = 5.0
)

type Player struct {
	X, Y float64

	Width, Height    float64
	OffsetX, OffsetY float64

	VelY, VelX           float64
	onGround             bool
	choice               string
	sprite               *ebiten.Image
	idleSprite           *ebiten.Image
	walkSprites          []*ebiten.Image
	attackSprites        []*ebiten.Image
	walkPos              int
	walkTimer, walkDelay int
	facingRight          bool
	attackPos            int
	attacking            bool

	shootTimer *Timer
}

var fistIdlePath = "assets/idle_fig.png"
var swordIdlePath = "assets/sword_idle.png"
var pistolIdlePath = "assets/pistol_idle.png"
var pistolWalkPath = "assets/pistol_run_*.png"
var swordWalkPath = "assets/sword_run_*.png"
var fistWalkPath = "assets/walk_*.png"
var swordAttackPath = "assets/sword_combo_*.png"

const PlayerScale = 0.75

func NewPlayer(choice string) *Player {
	var idleSprite *ebiten.Image
	var walkSprites []*ebiten.Image
	var attackSprites []*ebiten.Image

	if choice == "Sword" {
		idleSprite = LoadImage(swordIdlePath)
		walkSprites = LoadImages(swordWalkPath)
		attackSprites = LoadImages(swordAttackPath)
	} else if choice == "Pistol" {
		idleSprite = LoadImage(pistolIdlePath)
		walkSprites = LoadImages(pistolWalkPath)
		attackSprites = walkSprites
	} else {
		idleSprite = LoadImage(fistIdlePath)
		walkSprites = LoadImages(fistWalkPath)
		attackSprites = walkSprites
	}

	w := float64(idleSprite.Bounds().Dx()) * PlayerScale
	h := float64(idleSprite.Bounds().Dy()) * PlayerScale

	var timer = NewTimer(500 * time.Millisecond)
	if choice == "Sword" {
		timer = NewTimer(250 * time.Millisecond)
	}

	return &Player{
		X:             100,
		Y:             0,
		Width:         w * 0.75,
		Height:        h,
		sprite:        idleSprite,
		idleSprite:    idleSprite,
		walkSprites:   walkSprites,
		attackSprites: attackSprites,
		choice:        choice,
		attackPos:     0,
		walkPos:       0,
		walkDelay:     8,
		facingRight:   true,
		shootTimer:    timer,
	}
}

func (p *Player) Update(platforms []Platform, choice string, bulletList *[]*Bullet) {
	if p.attacking {
		p.sprite = p.attackSprites[p.attackPos]
	} else if p.VelX != 0 {
		p.sprite = p.walkSprites[p.walkPos]
	} else {
		p.sprite = p.idleSprite
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if choice == "Pistol" {
			if !p.shootTimer.IsActive() {
				p.shootTimer.Start()
			}

			if p.shootTimer.IsReady() {
				bx := p.X
				if p.facingRight {
					bx += p.Width
				}

				by := p.Y + p.Height*0.4
				*bulletList = append(*bulletList, NewBullet(bx, by, p.facingRight))
				p.shootTimer.Reset()
			}
		}

		if choice == "Sword" && inpututil.IsKeyJustPressed(ebiten.KeySpace) && !p.attacking {
			p.VelY = JumpVelocity + 4
			p.onGround = false
			p.attacking = true
			p.attackPos = 0
			p.shootTimer.Start()
		}
	}

	if p.attacking && p.shootTimer.IsReady() {
		p.attackPos++

		if p.attackPos >= len(p.attackSprites) {
			p.attacking = false
			p.attackPos = 0
			p.shootTimer.Stop()
		} else {
			p.shootTimer.Reset()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) && p.onGround {
		p.VelY = JumpVelocity
		p.onGround = false
	}

	if !p.attacking && ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.VelX = VelX
		p.walkTimer++

		p.facingRight = true

		if p.walkTimer >= p.walkDelay {
			p.walkTimer = 0
			p.walkPos = (p.walkPos + 1) % len(p.walkSprites)
		}

		p.sprite = p.walkSprites[p.walkPos]
	}

	if !p.attacking && ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.VelX = -VelX
		p.walkTimer++

		p.facingRight = false

		if p.walkTimer >= p.walkDelay {
			p.walkTimer = 0
			p.walkPos--
			if p.walkPos < 0 {
				p.walkPos = len(p.walkSprites) - 1
			}
		}

		p.sprite = p.walkSprites[p.walkPos]
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.VelY += Gravity
	}

	p.VelY += Gravity
	nextY := p.Y + p.VelY

	p.onGround = false

	for _, plat := range platforms {
		if p.X+p.Width <= plat.X || p.X >= plat.X+plat.Width {
			continue
		}

		if p.VelY > 0 {
			if p.Y+p.Height <= plat.Y && nextY+p.Height >= plat.Y {
				nextY = plat.Y - p.Height
				p.VelY = 0
				p.onGround = true
			}
		}

		if p.VelY < 0 {
			platBottom := plat.Y + plat.Height

			if p.Y >= platBottom && nextY <= platBottom {
				nextY = platBottom
				p.VelY = 0
			}
		}
	}

	p.Y = nextY
	p.X += p.VelX
	p.VelX = 0
}

func (p *Player) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	w := float64(p.sprite.Bounds().Dx())
	opts.GeoM.Translate(-w/2, 0)

	if !p.facingRight {
		opts.GeoM.Scale(-PlayerScale, PlayerScale)
	} else {
		opts.GeoM.Scale(PlayerScale, PlayerScale)
	}

	opts.GeoM.Translate(p.X, p.Y)

	screen.DrawImage(p.sprite, opts)
}
