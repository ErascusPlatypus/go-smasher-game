package helpers 

func (p *Player) handleDamageAnimation() {
	if p.takingDamage && p.damageTimer.IsReady() {
		p.damagePos++
		if p.damagePos >= len(p.damageSprites) {
			p.damagePos = 0
			p.takingDamage = false
			p.damageTimer.Stop()
		} else {
			p.damageTimer.Reset()
		}
	}
}

func (p *Player) handleAttackAnimation() {
	if p.attacking && p.shootTimer.IsReady() {
		p.attackPos++
		if p.attackPos >= len(p.attackSprites) {
			p.attacking = false
			p.hitThisSwing = false
			p.attackPos = 0
			p.shootTimer.Stop()
		} else {
			p.shootTimer.Reset()
		}
	}
}
