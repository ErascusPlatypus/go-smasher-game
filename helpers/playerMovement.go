package helpers

import (
	"github.com/hajimehoshi/ebiten/v2"

)

func (p *Player) handleMovementInput() {
	p.isWalking = false

	if p.takingDamage {
		return
	}

	if ebiten.IsKeyPressed(p.controls.Jump) && p.onGround {
		p.VelY = JumpVelocity
		p.onGround = false
	}

	if !p.attacking && ebiten.IsKeyPressed(p.controls.Right) {
		p.moveHorizontal(true)
		p.isWalking = true
	}

	if !p.attacking && ebiten.IsKeyPressed(p.controls.Left) {
		p.moveHorizontal(false)
		p.isWalking = true
	}
}

func (p *Player) moveHorizontal(right bool) {
	if right {
		p.VelX = VelX
		p.facingRight = true
	} else {
		p.VelX = -VelX
		p.facingRight = false
	}

	p.walkTimer++
	if p.walkTimer >= p.walkDelay {
		p.walkTimer = 0
		if right {
			p.walkPos = (p.walkPos + 1) % len(p.walkSprites)
		} else {
			p.walkPos--
			if p.walkPos < 0 {
				p.walkPos = len(p.walkSprites) - 1
			}
		}
	}
}

func (p *Player) applyPhysics(platforms []Platform) {
	p.VelY += Gravity
	nextY := p.Y + p.VelY
	p.onGround = false

	for _, plat := range platforms {
		if p.X+p.Width <= plat.X || p.X >= plat.X+plat.Width {
			continue
		}

		if p.VelY > 0 && p.Y+p.Height <= plat.Y && nextY+p.Height >= plat.Y {
			nextY = plat.Y - p.Height
			p.VelY = 0
			p.onGround = true
		}

		if p.VelY < 0 {
			bottom := plat.Y + plat.Height
			if p.Y >= bottom && nextY <= bottom {
				nextY = bottom
				p.VelY = 0
			}
		}
	}

	p.Y = nextY
	p.X += p.VelX
	p.VelX = 0
}

func (p *Player) clampWorld() {
	half := p.Width / 2
	if p.X-half < 0 {
		p.X = half
	}
	if p.X+half > 1200 {
		p.X = 1200 - half
	}
}