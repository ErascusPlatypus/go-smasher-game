package helpers

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