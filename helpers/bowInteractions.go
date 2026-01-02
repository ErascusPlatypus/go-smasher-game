package helpers

func handleBowSpecial(attacker, target *Player) {
    if box, ok := attacker.GetPushbackHitbox(); ok {
        if box.Intersects(target.GetRect()) {
            target.TakeDamage(10)
            pushForce := 150.0
            if !attacker.facingRight {
                pushForce = -150.0
            }
            target.X += pushForce
            target.VelY = -8
            attacker.pushbackHit = true
        }
    }
}

func (g *Game) handleBowSpecialDamage() {
    handleBowSpecial(g.playerOne, g.playerTwo)
    handleBowSpecial(g.playerTwo, g.playerOne)
}

func handleArrows(arrows *[]*Arrow, attacker, target *Player) {
    active := (*arrows)[:0]

    for _, a := range *arrows {
        a.Update()
        if !a.Active {
            continue
        }

        if a.GetRect().Intersects(target.GetRect()) {
            target.TakeDamage(a.GetDamage())
            a.Active = false
            continue
        }

        active = append(active, a)
    }

    *arrows = active
}

func (g *Game) handleBowDamage() {
	handleArrows(&g.ArrowsOne, g.playerOne, g.playerTwo)
	handleArrows(&g.ArrowsTwo, g.playerTwo, g.playerOne)
}