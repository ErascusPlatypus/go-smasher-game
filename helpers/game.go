package helpers

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

var choices = []string{"Sword", "Pistol", "Bow"}

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
	ArrowsOne	   []*Arrow
	ArrowsTwo      []*Arrow

	choiceIndexOne int
	choiceIndexTwo int

	PlayingDeathAnimation bool
	showEndScreen         bool
	deadPlayer *Player
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

func (g *Game) playDeathAnimation(victim *Player) {
	if victim.deathTimer.IsReady() {
		victim.deathPos++
		if victim.deathPos >= len(victim.deathSprites) {
			victim.deathTimer.Stop()
			g.PlayingDeathAnimation = false
			g.showEndScreen = true
		} else {
			victim.deathTimer.Reset()
		}
	}
}

func (g *Game) Update() error {
	switch g.State {
	case StateChoice:
		g.updateChoice()
	case StatePlaying:
		g.playerOne.Update(g.Platforms, &g.BulletsOne, &g.BombsOne, &g.ArrowsOne)
		g.playerTwo.Update(g.Platforms, &g.BulletsTwo, &g.BombsTwo, &g.ArrowsTwo)

		g.handleSwordDamage()
		g.handleSwordSpecialDamage()

		g.handlePistolDamage()
		g.handlePistolSpecialDamage()

		g.handleBowDamage()
		g.handleBowSpecialDamage()

		if g.playerOne.Health <= 0 || g.playerTwo.Health <= 0 {
			g.State = StateGameOver
		}

	case StateGameOver:
        if !g.PlayingDeathAnimation {
            if g.playerOne.Health <= 0 {
                g.deadPlayer = g.playerOne
            } else {
                g.deadPlayer = g.playerTwo
            }

            g.deadPlayer.isDead = true

            g.deadPlayer.deathPos = 0
            g.deadPlayer.deathTimer.Reset()
            g.deadPlayer.deathTimer.Start()
																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																									
            g.PlayingDeathAnimation = true
        }

        if g.PlayingDeathAnimation && !g.showEndScreen {																																																																																																																																																																																																																																																																																																																												
            g.deadPlayer.Update(g.Platforms, nil, nil, nil)
            
            if g.deadPlayer.deathPos >= len(g.deadPlayer.deathSprites)-1 {
                g.showEndScreen = true
            }
        }

        if g.showEndScreen && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
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

	for _, a := range g.ArrowsOne {
        a.Draw(screen, true)
    }
    for _, a := range g.ArrowsTwo {
        a.Draw(screen, false)
    }

    if g.State == StateGameOver && g.showEndScreen {
        winnerOne := g.deadPlayer == g.playerTwo
        g.loadEndScreenOverlay(screen, winnerOne)
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
	g.ArrowsOne = []*Arrow{}
	g.ArrowsTwo = []*Arrow{}

	g.PlayingDeathAnimation = false 
	g.deadPlayer = &Player{}
	g.showEndScreen = false 

	g.Platforms = []Platform{
		{X: 0, Y: 700, Width: 1200, Height: 10},
		{X: 200, Y: 550, Width: 140, Height: 10},
		{X: 500, Y: 400, Width: 150, Height: 10},
		{X: 850, Y: 550, Width: 140, Height: 10},	
	}
}