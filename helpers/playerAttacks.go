package helpers 

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (p *Player) TakeDamage(d int) {
	if p.takingDamage {
		return
	}

	p.takingDamage = true
	p.damagePos = 0
	p.damageTimer.Reset()
	p.damageTimer.Start()

	p.Health -= d
	if p.Health < 0 {
		p.Health = 0
	}
}

func (p *Player) handleBowAttack(arrows *[]*Arrow) {
    if ebiten.IsKeyPressed(p.controls.Attack) && 
		!p.attacking && !p.charging {
			p.charging = true 
			p.chargeTime = 0 
			p.attacking = true 
			p.attackPos = 0 
			p.shootTimer.Start()
	}

	if p.charging {
		p.chargeTime++ 

		if p.shootTimer.IsReady() {
			p.attackPos++ 
			if p.attackPos >= len(p.attackSprites)-2 {
				p.attackPos = len(p.attackSprites) - 2
			}
			p.shootTimer.Reset()
		}

		if inpututil.IsKeyJustReleased(p.controls.Attack) {
            bx := p.X
            if p.facingRight {
                bx += p.Width
            } else {
                bx -= p.Width * 0.5
            }

            by := p.Y + p.Height*0.4

            charged := p.chargeTime > 40
            *arrows = append(*arrows, NewArrow(bx, by, p.facingRight, charged))

            p.charging = false
            p.chargeTime = 0
            p.attacking = false
            p.attackPos = 0
            p.shootTimer.Reset()
        }
	}
}

func (p *Player) handlePistolAttack(bullets *[]*Bullet) {
	if !p.shootTimer.IsActive() {
		p.shootTimer.Start()
	}

	if p.shootTimer.IsReady() {
		bx := p.X
		if p.facingRight {
			bx += p.Width
		}
		by := p.Y + p.Height*0.4

		*bullets = append(*bullets, NewBullet(bx, by, p.facingRight))
		p.shootTimer.Reset()
	}
}

func (p *Player) startSwordAttack() {
	p.VelY = JumpVelocity + 4
	p.onGround = false
	p.attacking = true
	p.attackPos = 0
	p.shootTimer.Start()
}