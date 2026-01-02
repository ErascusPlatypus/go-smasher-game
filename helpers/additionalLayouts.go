package helpers

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

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

func (g *Game) loadEndScreenOverlay(screen *ebiten.Image, winnerOne bool) {
    overlay := ebiten.NewImage(1200, 800)
    overlay.Fill(color.RGBA{0, 0, 0, 180})
    screen.DrawImage(overlay, nil)

    winText := "Player 2 is the Winner!!"
    col := color.RGBA{100, 100, 255, 255}

    if winnerOne {
        winText = "Player 1 is the Winner!!!"
        col = color.RGBA{255, 100, 100, 255}
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
        color.White,
    )
}

func drawHitbox(screen *ebiten.Image, p *Player) {
	width := p.Width
	height := p.Height

	if p.choice == "Bow" {
		width = p.Width * 1.3
		height = p.Height * 1.2
	}

	img := ebiten.NewImage(int(width), int(height))
	img.Fill(color.RGBA{255, 0, 0, 120})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		p.X-width/2,
		p.Y,
	)

	screen.DrawImage(img, op)
}
