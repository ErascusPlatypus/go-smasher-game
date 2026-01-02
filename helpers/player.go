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

func (p *Player) selectSprite() {
	switch {
	case p.isDead:
		p.sprite = p.deathSprites[p.deathPos]
	case p.dashing:

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
