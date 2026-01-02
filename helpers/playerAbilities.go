package helpers

import (
	"github.com/hajimehoshi/ebiten/v2/inpututil"

)

func (p *Player) handleDash() bool {
	if !p.dashing {
		return false
	}

	p.VelY = 0
	p.X += p.dashVel

	if p.dashTimer.IsReady() {
		p.dashing = false
		p.dashTimer.Stop()
	}

	return true
}

func (p *Player) startPushback() {
	p.pushbackActive = true 
	p.pushbackHit = false 
	p.attacking = true 
	p.attackPos = 0 
	p.pushbackTimer.Start()
	p.pushbackCooldown.Start()
}

func (p *Player) handlePushback() bool {
	if !p.pushbackActive {
		return false 
	}

	if p.pushbackTimer.IsReady() {
		p.pushbackActive = false 
		p.pushbackTimer.Stop()
		p.attacking = false 
	}

	return true 
}

func (p *Player) GetPushbackHitbox() (Rect, bool) {
    if !p.pushbackActive || p.choice != "Bow" || p.pushbackHit {
        return Rect{}, false
    }

    w := 80.0
    h := 60.0

    x := p.X
    if p.facingRight {
        x += p.Width / 2
    } else {
        x -= p.Width/2 + w
    }

    y := p.Y + p.Height*0.2

    return NewRect(x, y, w, h), true
}

func (p *Player) handleSpecial(bombs *[]*Bomb) {
	if inpututil.IsKeyJustPressed(p.controls.SpecialOne) {

		if p.choice == "Sword" &&
			!p.dashing &&
			!p.dashCooldown.IsActive() &&
			!p.takingDamage {

			p.startDash()
		}

		if p.choice == "Pistol" &&
			!p.bombCooldown.IsActive() {

			p.throwBomb(bombs)
		}

		if p.choice == "Bow" && 
			!p.pushbackActive && 
			!p.pushbackCooldown.IsActive() {
				p.startPushback()
			}
	}
}

func (p *Player) startDash() {
	p.dashing = true
	p.dashHit = false
	p.attacking = true
	p.attackPos = 0

	p.VelY = JumpVelocity + 4
	p.onGround = false

	p.dashVel = 18
	if !p.facingRight {
		p.dashVel = -18
	}

	p.shootTimer.Start()
	p.dashTimer.Start()
	p.dashCooldown.Start()
}

func (p *Player) throwBomb(bombs *[]*Bomb) {
	bx := p.X
	if p.facingRight {
		bx += p.Width / 2
	} else {
		bx -= p.Width / 2
	}

	by := p.Y + p.Height*0.4
	*bombs = append(*bombs, NewBomb(bx, by, p.facingRight))
	p.bombCooldown.Start()
}