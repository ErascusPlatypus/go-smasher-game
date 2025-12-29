package main

import (
	"embed"
	"pro12_fighter/helpers"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/**
var assets embed.FS

const (
	StateChoice helpers.GameState = iota
	StatePlaying
)

func main() {
	helpers.Init(assets)

	g := &helpers.Game{
		State:   StateChoice,
		BulletsOne: []*helpers.Bullet{},
		BulletsTwo: []*helpers.Bullet{},
	}

	g.Platforms = []helpers.Platform{
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
