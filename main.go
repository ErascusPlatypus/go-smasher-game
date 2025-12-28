package main

import (
	"embed"
	"image/color"
	"pro12_fighter/helpers"

	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

)

//go:embed assets/**
var assets embed.FS

type GameState int 
var choices = []string{"Fist", "Sword", "Pistol"}

const (
	StateChoice GameState = iota
	StatePlaying
)

type Game struct {
	player    *helpers.Player
	state GameState
	choice string
	platforms []helpers.Platform

	choiceIndex int
}

func (g *Game) updateChoice() {
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.choiceIndex = (g.choiceIndex + 1) % len(choices)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.choiceIndex--
		if g.choiceIndex < 0 {
			g.choiceIndex = len(choices) - 1
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.choice = choices[g.choiceIndex]
		g.state = StatePlaying
	}
}

func (g *Game) Update() error {
	switch g.state {
	case StateChoice:
		g.updateChoice()
	case StatePlaying:
		g.player.Update(g.platforms)
	}
	return nil
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	if g.state == StateChoice {
		g.loadChoiceScreen(screen)
		return
	}

	for _, p := range g.platforms {
		helpers.DrawPlatform(screen, p)
	}

	drawHitbox(screen, g.player.X, g.player.Y, g.player.Width, g.player.Height)
	g.player.Draw(screen)
}


func (g *Game) loadChoiceScreen(screen *ebiten.Image) {
	startY := 220
	lineHeight := 42

	text.Draw(
		screen,
		"Choose Your Weapon",
		helpers.DefaultFont,
		460,
		160,
		color.Black,
	)

	for i, c := range choices {
		y := startY + i*lineHeight

		if i == g.choiceIndex {
			bg := ebiten.NewImage(260, lineHeight)
			bg.Fill(color.RGBA{30, 30, 30, 220})

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(450, float64(y-28))
			screen.DrawImage(bg, op)
		}

		col := color.Black
		if i == g.choiceIndex {
			col = color.White
		}

		text.Draw(
			screen,
			c,
			helpers.DefaultFont,
			480,
			y,
			col,
		)
	}
}


func drawHitbox(screen *ebiten.Image, x, y, w, h float64) {
	img := ebiten.NewImage(int(w), int(h))
	img.Fill(color.RGBA{255, 0, 0, 100})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(img, op)
}


func main() {
	helpers.Init(assets)

	g := &Game{
		player: helpers.NewPlayer(),
		state: StateChoice,
	}

	g.platforms = []helpers.Platform{
		{X: 0, Y: 700, Width: 1200, Height: 10},

		{X: 150, Y: 550, Width: 140, Height: 10},
		{X: 500, Y: 400, Width: 150, Height: 10},
		{X: 850, Y: 550, Width: 140, Height: 10},
	}

	ebiten.SetWindowSize(1200, 800)

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
