package helpers

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type GameState int

var choices = []string{"Fist", "Sword", "Pistol"}

var PlayerOneControls = Controls{
	Left:   ebiten.KeyA,
	Right:  ebiten.KeyD,
	Jump:   ebiten.KeyW,
	Down:   ebiten.KeyS,
	Attack: ebiten.KeySpace,
}

var PlayerTwoControls = Controls{
	Left:   ebiten.KeyLeft,
	Right:  ebiten.KeyRight,
	Jump:   ebiten.KeyUp,
	Down:   ebiten.KeyDown,
	Attack: ebiten.KeyEnter,
}

const (
	StateChoice GameState = iota
	StatePlaying
)

type Game struct {
	playerOne      *Player
	playerTwo      *Player
	State          GameState
	StatePlayerOne GameState
	StatePlayerTwo GameState
	choiceOne      string
	choiceTwo      string
	Platforms      []Platform
	BulletsOne     []*Bullet
	BulletsTwo     []*Bullet

	choiceIndexOne int
	choiceIndexTwo int
}

func (g *Game) updateChoice() {
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.choiceIndexTwo = (g.choiceIndexTwo + 1) % len(choices)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.choiceIndexTwo--
		if g.choiceIndexTwo < 0 {
			g.choiceIndexTwo = len(choices) - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.choiceTwo = choices[g.choiceIndexTwo]
		g.playerTwo = NewPlayer(g.choiceTwo, PlayerTwoControls)
		g.playerTwo.X = 900
		g.StatePlayerTwo = StatePlaying
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.choiceIndexOne--
		if g.choiceIndexOne < 0 {
			g.choiceIndexOne = len(choices) - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.choiceIndexOne = (g.choiceIndexOne + 1) % len(choices)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft) {
		g.choiceOne = choices[g.choiceIndexOne]
		g.playerOne = NewPlayer(g.choiceOne, PlayerOneControls)
		g.StatePlayerOne = StatePlaying
	}

	if g.StatePlayerOne == StatePlaying && g.StatePlayerTwo == StatePlaying {
		g.State = StatePlaying
	}
}

func (g *Game) Update() error {
	switch g.State {
	case StateChoice:
		g.updateChoice()
	case StatePlaying:
		g.playerOne.Update(g.Platforms, &g.BulletsOne)
		g.playerTwo.Update(g.Platforms, &g.BulletsTwo)

		if box, ok := g.playerOne.GetSwordHitbox(); ok {
			if box.Intersects(g.playerTwo.GetRect()) && !g.playerOne.hitThisSwing {
				g.playerTwo.TakeDamage(20)
				g.playerOne.hitThisSwing = true
			}
		}

		if box, ok := g.playerTwo.GetSwordHitbox(); ok {
			if box.Intersects(g.playerOne.GetRect()) && !g.playerTwo.hitThisSwing {
				g.playerOne.TakeDamage(20)
				g.playerTwo.hitThisSwing = true
			}
		}

		activeBullets := []*Bullet{}
		for _, b := range g.BulletsOne {
			b.Update()
			if !b.Active {
				continue
			}

			if b.GetRect().Intersects(g.playerTwo.GetRect()) {
				g.playerTwo.TakeDamage(10)
				b.Active = false
			}

			if b.Active {
				activeBullets = append(activeBullets, b)
			}
		}

		g.BulletsOne = activeBullets

		activeBullets = []*Bullet{}
		for _, b := range g.BulletsTwo {
			b.Update()
			if !b.Active {
				continue
			}

			if b.GetRect().Intersects(g.playerOne.GetRect()) {
				g.playerOne.TakeDamage(10)
				b.Active = false
			}

			if b.Active {
				activeBullets = append(activeBullets, b)
			}
		}

		g.BulletsTwo = activeBullets
	}
	return nil
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	if g.State == StateChoice {
		g.loadChoiceScreen(screen)
		return
	}

	for _, p := range g.Platforms {
		DrawPlatform(screen, p)
	}

	// drawHitbox(screen, g.playerOne)
	// drawHitbox(screen, g.playerTwo)

	g.playerOne.Draw(screen, true)
	g.playerTwo.Draw(screen, false)

	for _, b := range g.BulletsOne {
		b.Draw(screen, true)
	}

	for _, b := range g.BulletsTwo {
		b.Draw(screen, false)
	}

}

func (g *Game) loadChoiceScreen(screen *ebiten.Image) {
	lineHeight := 42

	leftPanel := ebiten.NewImage(500, 420)
	// leftPanel.Fill(color.RGBA{200, 220, 255, 255})

	rightPanel := ebiten.NewImage(500, 420)
	// rightPanel.Fill(color.RGBA{255, 210, 210, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(80, 180)
	screen.DrawImage(leftPanel, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(620, 180)
	screen.DrawImage(rightPanel, op)

	text.Draw(
		screen,
		"PLAYER 1",
		DefaultFont,
		240,
		220,
		color.RGBA{0, 60, 150, 255},
	)

	text.Draw(
		screen,
		"PLAYER 2",
		DefaultFont,
		780,
		220,
		color.RGBA{150, 0, 0, 255},
	)

	text.Draw(
		screen,
		"W / S  Move   SPACE  Select",
		DefaultFont,
		160,
		260,
		color.Black,
	)

	text.Draw(
		screen,
		"↑ / ↓  Move   SHIFT  Select",
		DefaultFont,
		700,
		260,
		color.Black,
	)

	startY := 310
	for i, c := range choices {
		y := startY + i*lineHeight

		if i == g.choiceIndexOne {
			bg := ebiten.NewImage(300, lineHeight)
			bg.Fill(color.RGBA{0, 80, 200, 220})

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(140, float64(y-28))
			screen.DrawImage(bg, op)
		}

		col := color.Black
		if i == g.choiceIndexOne {
			col = color.White
		}

		text.Draw(
			screen,
			c,
			DefaultFont,
			170,
			y,
			col,
		)
	}

	for i, c := range choices {
		y := startY + i*lineHeight

		if i == g.choiceIndexTwo {
			bg := ebiten.NewImage(300, lineHeight)
			bg.Fill(color.RGBA{200, 40, 40, 220})

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(680, float64(y-28))
			screen.DrawImage(bg, op)
		}

		col := color.Black
		if i == g.choiceIndexTwo {
			col = color.White
		}

		text.Draw(
			screen,
			c,
			DefaultFont,
			710,
			y,
			col,
		)
	}
}

func drawHitbox(screen *ebiten.Image, p *Player) {
	img := ebiten.NewImage(int(p.Width), int(p.Height))
	img.Fill(color.RGBA{255, 0, 0, 120})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		p.X-p.Width/2,
		p.Y,
	)

	screen.DrawImage(img, op)
}
