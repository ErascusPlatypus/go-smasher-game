package helpers

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
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

	VelY, VelX    float64
	onGround      bool
	choice        string
	sprite        *ebiten.Image
	idleSprite    *ebiten.Image
	walkSprites   []*ebiten.Image
	attackSprites []*ebiten.Image
	damageSprites []*ebiten.Image
	deathSprites  []*ebiten.Image

	walkPos              int
	walkTimer, walkDelay int
	isWalking            bool
	isDead               bool

	deathPos     int
	facingRight  bool
	attackPos    int
	damagePos    int
	attacking    bool
	Health       int
	hitThisSwing bool
	takingDamage bool
	hasHit       bool
	dashing      bool
	dashHit      bool
	dashVel      float64

	charging bool
	chargeTime int
	pushbackActive bool 
	pushbackHit bool 

	controls Controls

	shootTimer   *Timer
	damageTimer  *Timer
	deathTimer   *Timer
	bombCooldown *Timer
	dashTimer    *Timer
	dashCooldown *Timer
	pushbackTimer *Timer
	pushbackCooldown *Timer
}

var pistolIdlePath = "assets/pistol_idle.png"
var pistolWalkPath = "assets/pistol_run_*.png"
var pistolDamagePath = "assets/pistol_hit_*.png"
var pistolDeathPath = "assets/pistol_death_*.png"

var swordIdlePath = "assets/sword_idle.png"
var swordWalkPath = "assets/sword_run_*.png"
var swordAttackPath = "assets/sword_combo_*.png"
var swordDamagePath = "assets/sword_hit_*.png"
var swordDeathPath = "assets/sword_death_*.png"

var bowIdlePath = "assets/bow_idle_01.png"
var bowWalkPath = "assets/bow_walk_*.png"
var bowAttackPath = "assets/bow_attack_*.png"
var bowDamagePath = "assets/bow_hit_*.png"
var bowDeathPath = "assets/bow_death_*.png"


const PlayerScale = 0.75

func NewPlayer(choice string, c Controls) *Player {
	var idleSprite *ebiten.Image
	var walkSprites []*ebiten.Image
	var attackSprites []*ebiten.Image
	var damageSprites []*ebiten.Image
	var deathSprites []*ebiten.Image

	var health int

	if choice == "Sword" {
		idleSprite = LoadImage(swordIdlePath)
		walkSprites = LoadImages(swordWalkPath)
		attackSprites = LoadImages(swordAttackPath)
		damageSprites = LoadImages(swordDamagePath)
		deathSprites = LoadImages(swordDeathPath)
		health = 150
	} else if choice == "Pistol" {
		idleSprite = LoadImage(pistolIdlePath)
		walkSprites = LoadImages(pistolWalkPath)
		attackSprites = walkSprites
		damageSprites = LoadImages(pistolDamagePath)
		deathSprites = LoadImages(pistolDeathPath)
		health = 100
	} else if choice == "Bow" {
		idleSprite = LoadImage(bowIdlePath)
		walkSprites = LoadImages(bowWalkPath)
		attackSprites = LoadImages(bowAttackPath)
		damageSprites = LoadImages(bowDamagePath)
		deathSprites = LoadImages(bowDeathPath)
		health = 100
	}

	w := float64(idleSprite.Bounds().Dx()) * PlayerScale
	h := float64(idleSprite.Bounds().Dy()) * PlayerScale

	var timer = NewTimer(500 * time.Millisecond)
	if choice == "Sword" {
		timer = NewTimer(250 * time.Millisecond)
	} else if choice == "Bow" {
		timer = NewTimer(1000 * time.Millisecond)
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
		damageSprites: damageSprites,
		deathSprites:  deathSprites,
		Health:        health,
		choice:        choice,
		attackPos:     0,
		damagePos:     0,
		hasHit:        false,
		walkPos:       0,
		walkDelay:     8,
		facingRight:   true,
		shootTimer:    timer,
		damageTimer:   NewTimer(120 * time.Millisecond),
		bombCooldown:  NewTimer(10 * time.Second),
		dashTimer:     NewTimer(180 * time.Millisecond),
		dashCooldown:  NewTimer(6 * time.Second),
		pushbackTimer: NewTimer(150 * time.Millisecond),
		pushbackCooldown: NewTimer(8 * time.Second),
		deathPos:      0,
		isDead:        false,
		deathTimer:    NewTimer(120 * time.Millisecond),
		controls:      c,
		charging: false,
		chargeTime: 0,
	}
}

func (p *Player) updateCooldowns() {
	if p.dashCooldown.IsActive() && p.dashCooldown.IsReady() {
		p.dashCooldown.Stop()
	}
	if p.bombCooldown.IsActive() && p.bombCooldown.IsReady() {
		p.bombCooldown.Stop()
	}

	if p.pushbackCooldown.IsActive() && p.pushbackCooldown.IsReady() {
		p.pushbackCooldown.Stop()
	}
}

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

func (p *Player) handleAttack(bullets *[]*Bullet, arrows *[]*Arrow) {
    if p.takingDamage {
        return
    }

    if p.choice == "Pistol" && ebiten.IsKeyPressed(p.controls.Attack) {
        p.handlePistolAttack(bullets)
    }

    if p.choice == "Sword" && inpututil.IsKeyJustPressed(p.controls.Attack) && !p.attacking {
        p.startSwordAttack()
    }

    if p.choice == "Bow" && (!p.shootTimer.IsActive() || p.shootTimer.IsReady() || p.charging) {
        p.handleBowAttack(arrows)
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

func (p *Player) handleDamageAnimation() {
	if p.takingDamage && p.damageTimer.IsReady() {
		p.damagePos++
		if p.damagePos >= len(p.damageSprites) {
			p.damagePos = 0
			p.takingDamage = false
			p.damageTimer.Stop()
		} else {
			p.damageTimer.Reset()
		}
	}
}

func (p *Player) handleAttackAnimation() {
	if p.attacking && p.shootTimer.IsReady() {
		p.attackPos++
		if p.attackPos >= len(p.attackSprites) {
			p.attacking = false
			p.hitThisSwing = false
			p.attackPos = 0
			p.shootTimer.Stop()
		} else {
			p.shootTimer.Reset()
		}
	}
}

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

func (p *Player) selectSprite() {
	switch {
	case p.isDead:
		p.sprite = p.deathSprites[p.deathPos]
	case p.dashing:
		// dash sprite already set
	case p.takingDamage:
		p.sprite = p.damageSprites[p.damagePos]
	case p.attacking:
		p.sprite = p.attackSprites[p.attackPos]
	case p.isWalking:
		p.sprite = p.walkSprites[p.walkPos]
	default:
		p.sprite = p.idleSprite
	}
}

func (p *Player) Update(platforms []Platform, bullets *[]*Bullet, bombs *[]*Bomb, arrows *[]*Arrow) {
	if p.isDead {
        if !p.deathTimer.IsActive() {
            p.deathTimer.Start()
        }
        if p.deathTimer.IsReady() {
            p.deathPos++
            if p.deathPos >= len(p.deathSprites) {
                p.deathPos = len(p.deathSprites) - 1
            } else {
                p.deathTimer.Reset()
            }
        }
        p.selectSprite()
        return
    }

	p.updateCooldowns()

	if p.handleDash() {
		return
	}

	if p.handlePushback() {
        p.applyPhysics(platforms)
        p.clampWorld()
        p.selectSprite()
        return
    }

	p.handleSpecial(bombs)
	p.handleAttack(bullets, arrows)
	p.handleDamageAnimation()
	p.handleAttackAnimation()
	p.handleMovementInput()

	p.applyPhysics(platforms)
	p.clampWorld()

	p.selectSprite()
}

func (p *Player) Draw(screen *ebiten.Image, playerOne bool) {
	opts := &ebiten.DrawImageOptions{}

	w := float64(p.sprite.Bounds().Dx())
	opts.GeoM.Translate(-w/2, 0)

	if !p.facingRight {
		if p.choice == "Bow" {
			opts.GeoM.Scale(-1, 1)
		} else {
			opts.GeoM.Scale(-PlayerScale, PlayerScale)
		}
	} else {
		if p.choice == "Bow" {
			opts.GeoM.Scale(1, 1)
		} else {
			opts.GeoM.Scale(PlayerScale, PlayerScale)
		}
	}

	yOffset := 0.0
    if p.choice == "Bow" {
        yOffset = -30.0
    }

    opts.GeoM.Translate(p.X, p.Y+yOffset)
	if playerOne {
		opts.ColorScale.Scale(0.4, 0.4, 1.0, 1.0)
	} else {
		opts.ColorScale.Scale(1.0, 0.4, 0.4, 1.0)
	}

	if p.dashing {
		opts.ColorScale.Scale(1.6, 1.4, 0.4, 1.0)
	}

	hp := fmt.Sprintf("HP : %d", p.Health)
	x := 100
	if !playerOne {
		x = 1000
	}

	text.Draw(
		screen,
		hp,
		DefaultFont,
		x, 50,
		color.Black,
	)

	screen.DrawImage(p.sprite, opts)
}

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

func (p *Player) GetSwordHitbox() (Rect, bool) {
	if p.choice != "Sword" || !p.attacking {
		return Rect{}, false
	}

	width := 60.0
	height := 40.0

	x := p.X
	if p.facingRight {
		x += p.Width / 2
	} else {
		x -= p.Width/2 + width
	}

	y := p.Y + p.Height*0.3

	return NewRect(x, y, width, height), true
}

func (p *Player) GetDashHitbox() (Rect, bool) {
	if !p.dashing || p.choice != "Sword" || p.dashHit {
		return Rect{}, false
	}

	w := 90.0
	h := 50.0

	x := p.X
	if p.facingRight {
		x += p.Width / 2
	} else {
		x -= p.Width/2 + w
	}

	y := p.Y + p.Height*0.3

	return NewRect(x, y, w, h), true
}

func (p *Player) GetRect() Rect {
	return NewRect(
		p.X-p.Width/2,
		p.Y,
		p.Width,
		p.Height,
	)
}
