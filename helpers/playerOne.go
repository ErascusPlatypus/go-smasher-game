package helpers

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	JumpVelocity = -12.0
	Gravity = 0.4
	GroundY = 560
	VelX = 5.0
)

type Player struct {
	X, Y float64

	Width, Height float64   
	OffsetX, OffsetY float64 

	VelY, VelX float64
	onGround bool
	sprite *ebiten.Image
	walkSprites [] *ebiten.Image
	walkPos int
	walkTimer, walkDelay int
	facingRight bool
}


var spritePath = "assets/idle_fig.png"
var walkPath = "assets/walk_*.png"

const PlayerScale = 0.75

func NewPlayer() *Player {
	sprite := LoadImage(spritePath)
	walkSprites := LoadImages(walkPath)

	w := float64(sprite.Bounds().Dx()) * PlayerScale
	h := float64(sprite.Bounds().Dy()) * PlayerScale
	
	return &Player{
		X: 100,
		Y: 0,
		Width:  w*0.85,
		Height: h,
		sprite: sprite,
		walkSprites: walkSprites,
		walkPos: 0,
		walkDelay: 8,
		facingRight: true,
	}
}

func (p *Player) Update(platforms [] Platform) {
	if ebiten.IsKeyPressed(ebiten.KeySpace) && p.onGround{
		p.VelY = JumpVelocity
		p.onGround = false 
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.VelX = VelX
		p.walkTimer++ 

		p.facingRight = true

		if p.walkTimer >= p.walkDelay {
			p.walkTimer = 0 
			p.walkPos = (p.walkPos + 1) % len(p.walkSprites) 
		}

		p.sprite = p.walkSprites[p.walkPos]
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
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
		if p.X + p.Width <= plat.X || p.X >= plat.X + plat.Width {
			continue 
		}
		
		if p.VelY > 0 {
			if p.Y + p.Height <= plat.Y && nextY + p.Height >= plat.Y {
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
	if p.facingRight == true {
		opts.GeoM.Scale(PlayerScale, PlayerScale)
	} else {
		opts.GeoM.Scale(-PlayerScale, PlayerScale)
	}

	opts.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.sprite, opts)
}
