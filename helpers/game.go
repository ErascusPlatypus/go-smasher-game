package helpers

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type GameState int

var choices = []string{"Sword", "Pistol"}

var PlayerOneControls = Controls{
	Left:       ebiten.KeyA,
	Right:      ebiten.KeyD,
	Jump:       ebiten.KeyW,
	Down:       ebiten.KeyS,
	Attack:     ebiten.KeySpace,
	SpecialOne: ebiten.KeyR,
}

var PlayerTwoControls = Controls{
	Left:       ebiten.KeyLeft,
	Right:      ebiten.KeyRight,
	Jump:       ebiten.KeyUp,
	Down:       ebiten.KeyDown,
	Attack:     ebiten.KeyEnter,
	SpecialOne: ebiten.KeyShiftRight,
}

const (
	StateChoice GameState = iota
	StatePlaying
	StateGameOver
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
	BombsOne       []*Bomb
	BombsTwo       []*Bomb

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

func handleSwords(attacker, target *Player) {
	if box, ok := attacker.GetSwordHitbox(); ok {
		if box.Intersects(target.GetRect()) && !attacker.hitThisSwing {
			target.TakeDamage(15)
			attacker.hitThisSwing = true 
		}
	}
}

func (g *Game) handleSwordDamage() {
	handleSwords(g.playerOne, g.playerTwo)
	handleSwords(g.playerTwo, g.playerOne)
}

func handleSwordSpecial(attacker, target *Player) {
	if box, ok := attacker.GetDashHitbox(); ok {
		if box.Intersects(target.GetRect()) {
			target.TakeDamage(35)
			attacker.dashHit = true 
		}
	}
}

func (g *Game) handleSwordSpecialDamage() {
	handleSwordSpecial(g.playerOne, g.playerTwo)
	handleSwordSpecial(g.playerTwo, g.playerOne)
}

func handlePistolSpecial(bombs *[]*Bomb, target *Player) {
	active := (*bombs)[:0]

	for _, b := range *bombs {

		if b.Active {
			b.Update()

			if b.HitsPlayer(target) {
				b.Active = false
				b.Exploded = true
				b.ExplosionTTL = 20
			}
		}

		if b.Exploded {
			if !b.HasDamaged && b.HitsPlayer(target) {
				target.TakeDamage(35)
				b.HasDamaged = true
			}

			b.ExplosionTTL--
			if b.ExplosionTTL > 0 {
				active = append(active, b)
			}
			continue
		}

		if b.Active {
			active = append(active, b)
		}
	}

	*bombs = active
}

func (g *Game) handlePistolSpecialDamage() {
	handlePistolSpecial(&g.BombsOne, g.playerTwo)
	handlePistolSpecial(&g.BombsTwo, g.playerOne)
}

func handleBullets(bullets *[]*Bullet, attacker, target *Player) {
	active := (*bullets)[:0] 

	for _, b := range *bullets {
		b.Update()
		if !b.Active {
			continue
		}

		if b.GetRect().Intersects(target.GetRect()) {
			target.TakeDamage(10)
			b.Active = false
			continue
		}

		active = append(active, b)
	}

	*bullets = active
}

func (g *Game) handlePistolDamage() {
	handleBullets(&g.BulletsOne, g.playerOne, g.playerTwo)
	handleBullets(&g.BulletsTwo, g.playerTwo, g.playerOne)
}

func (g *Game) Update() error {
	switch g.State {
	case StateChoice:
		g.updateChoice()
	case StatePlaying:
		g.playerOne.Update(g.Platforms, &g.BulletsOne, &g.BombsOne)
		g.playerTwo.Update(g.Platforms, &g.BulletsTwo, &g.BombsTwo)

		g.handleSwordDamage()
		g.handleSwordSpecialDamage()

		g.handlePistolDamage()
		g.handlePistolSpecialDamage()

		if g.playerOne.Health <= 0 || g.playerTwo.Health <= 0 {
			g.State = StateGameOver
		}

	case StateGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.Reset()
		}
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

	if g.State == StateGameOver {
		winnerOne := g.playerTwo.Health <= 0
		g.loadEndScreen(screen, winnerOne)
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

	for _, b := range g.BombsOne {
		b.Draw(screen)
	}
	for _, b := range g.BombsTwo {
		b.Draw(screen)
	}
}

func (g *Game) loadChoiceScreen(screen *ebiten.Image) {
	lineHeight := 42

	leftPanel := ebiten.NewImage(500, 420)

	rightPanel := ebiten.NewImage(500, 420)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(80, 180)
	screen.DrawImage(leftPanel, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(620, 180)
	screen.DrawImage(rightPanel, op)

	r := 0
	b := 40
	if g.choiceOne != "" {
		r = 255
	}

	if g.choiceTwo != "" {
		b = 255
	}

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
		"W / S  Move   SHIFT Select",
		DefaultFont,
		160,
		260,
		color.Black,
	)

	text.Draw(
		screen,
		"↑ / ↓  Move   ENTER Select",
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
			bg.Fill(color.RGBA{uint8(r), 80, 200, 220})

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
			bg.Fill(color.RGBA{200, 40, uint8(b), 220})

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

func (g *Game) Reset() {
	g.State = StateChoice

	g.playerOne = nil
	g.playerTwo = nil

	g.StatePlayerOne = StateChoice
	g.StatePlayerTwo = StateChoice

	g.choiceOne = ""
	g.choiceTwo = ""

	g.choiceIndexOne = 0
	g.choiceIndexTwo = 0

	g.BulletsOne = []*Bullet{}
	g.BulletsTwo = []*Bullet{}
	g.BombsOne = []*Bomb{}
	g.BombsTwo = []*Bomb{}

	g.Platforms = []Platform{
		{X: 0, Y: 700, Width: 1200, Height: 10},
		{X: 200, Y: 550, Width: 140, Height: 10},
		{X: 500, Y: 400, Width: 150, Height: 10},
		{X: 850, Y: 550, Width: 140, Height: 10},
	}
}

func (g *Game) loadEndScreen(screen *ebiten.Image, winnerOne bool) {
	winText := "Player 2 is the Winner!!"
	col := color.RGBA{0, 0, 255, 255}

	if winnerOne {
		winText = "Player 1 is the Winner!!!"
		col = color.RGBA{255, 0, 0, 255}
	}

	x := 600 - len(winText)*6

	text.Draw(
		screen,
		winText,
		WinnerFont,
		x-150,
		300,
		col,
	)

	text.Draw(
		screen,
		"Press Enter to restart the game",
		DefaultFont,
		x,
		500,
		color.Black,
	)
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
