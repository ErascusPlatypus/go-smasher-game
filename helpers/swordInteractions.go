package helpers

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