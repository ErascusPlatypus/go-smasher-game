package helpers

import (
    "image/color"
    "math"

    "github.com/hajimehoshi/ebiten/v2"
)

type Arrow struct {
    X, Y       float64
    VelX, VelY float64
    sprite     *ebiten.Image
    direction  bool
    Active     bool
    Charged    bool
    glowFrame  int // for pulsing effect
}

const (
    arrowSpeed   = 12.0
    arrowGravity = 0.3
)

var arrowSprite = "assets/arrow.png"

func NewArrow(x, y float64, dir bool, charged bool) *Arrow {
    sprite := LoadImage(arrowSprite)

    velY := -3.0
    if charged {
        velY = -5.0
    }

    return &Arrow{
        X:         x,
        Y:         y,
        VelY:      velY,
        sprite:    sprite,
        Active:    true,
        direction: dir,
        Charged:   charged,
        glowFrame: 0,
    }
}

func (a *Arrow) Update() {
    if a.direction {
        a.X += arrowSpeed
    } else {
        a.X -= arrowSpeed
    }

    a.VelY += arrowGravity
    a.Y += a.VelY

    if a.Charged {
        a.glowFrame++
    }

    if a.Y < -50 || a.Y > 1200 || a.X < -50 || a.X > 1400 {
        a.Active = false
    }
}

func (a *Arrow) Draw(screen *ebiten.Image, playerOne bool) {
    if !a.Active {
        return
    }

    speed := arrowSpeed
    if !a.direction {
        speed = -arrowSpeed
    }
    angle := math.Atan2(a.VelY, speed)

    if a.Charged {
        a.drawTrail(screen, playerOne)
    }

    opts := &ebiten.DrawImageOptions{}
    w, h := a.sprite.Bounds().Dx(), a.sprite.Bounds().Dy()
    opts.GeoM.Translate(-float64(w)/2, -float64(h)/2)
    opts.GeoM.Rotate(angle)
    opts.GeoM.Translate(a.X, a.Y)

    if playerOne {
        opts.ColorScale.Scale(0.4, 0.4, 1.0, 1.0)
    } else {
        opts.ColorScale.Scale(1.0, 0.4, 0.4, 1.0)
    }

    if a.Charged {
        opts.ColorScale.Scale(1.3, 1.3, 1.3, 1.0)
    }

    screen.DrawImage(a.sprite, opts)
}

func (a *Arrow) drawTrail(screen *ebiten.Image, playerOne bool) {
    pulse := 0.8 + 0.2*math.Sin(float64(a.glowFrame)*0.3)

    for i := 1; i <= 8; i++ {
        trailX := a.X
        trailY := a.Y

        if a.direction {
            trailX -= float64(i) * 14
        } else {
            trailX += float64(i) * 14
        }
        trailY -= a.VelY * float64(i) * 0.2

        size := 16 - i*2
        if size < 4 {
            size = 4
        }

        particle := ebiten.NewImage(size, size)

        var particleColor color.RGBA
        if playerOne {
            particleColor = color.RGBA{100, 180, 255, 255} 
        } else {
            particleColor = color.RGBA{255, 120, 80, 255} 
        }
        particle.Fill(particleColor)

        opts := &ebiten.DrawImageOptions{}
        opts.GeoM.Translate(-float64(size)/2, -float64(size)/2)
        opts.GeoM.Translate(trailX, trailY)

        alpha := float32((1.0 - float64(i)*0.1) * pulse)
        opts.ColorScale.Scale(1.0, 1.0, 1.0, alpha)

        screen.DrawImage(particle, opts)
    }
}

func (a *Arrow) GetRect() Rect {
    w, h := a.sprite.Bounds().Dx(), a.sprite.Bounds().Dy()
    return NewRect(a.X, a.Y, float64(w), float64(h))
}

func (a *Arrow) GetDamage() int {
    if a.Charged {
        return 20
    }
    return 12
}