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
		g.choiceIndexOne = (g.choiceIndexOne + 1) % len(choices)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.choiceIndexOne--
		if g.choiceIndexOne < 0 {
			g.choiceIndexOne = len(choices) - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.choiceOne = choices[g.choiceIndexOne]
		g.playerOne = NewPlayer(g.choiceOne, PlayerOneControls)
		g.StatePlayerOne = StatePlaying
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.choiceIndexTwo--
		if g.choiceIndexTwo < 0 {
			g.choiceIndexTwo = len(choices) - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.choiceIndexTwo = (g.choiceIndexTwo + 1) % len(choices)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft) {
		g.choiceTwo = choices[g.choiceIndexTwo]
		g.playerTwo = NewPlayer(g.choiceTwo, PlayerTwoControls)
		g.playerTwo.X = 900
		g.StatePlayerTwo = StatePlaying
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
	startYOne := 220
	lineHeight := 42
	y := 0

	text.Draw(
		screen,
		"Choose Your Weapon - Player 1",
		DefaultFont,
		460,
		160,
		color.Black,
	)

	for i, c := range choices {
		y = startYOne + i*lineHeight

		if i == g.choiceIndexOne {
			bg := ebiten.NewImage(260, lineHeight)
			bg.Fill(color.RGBA{30, 30, 30, 220})

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(450, float64(y-28))
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
			480,
			y,
			col,
		)
	}

	startYTwo := y + lineHeight*2

	text.Draw(
		screen,
		"Choose Your Weapon - Player 2",
		DefaultFont,
		460,
		160,
		color.Black,
	)

	for i, c := range choices {
		y = startYTwo + i*lineHeight

		if i == g.choiceIndexTwo {
			bg := ebiten.NewImage(260, lineHeight)
			bg.Fill(color.RGBA{30, 30, 30, 220})

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(450, float64(y-28))
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
			480,
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
