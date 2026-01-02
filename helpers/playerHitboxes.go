package helpers

func (p *Player) GetSwordHitbox() (Rect, bool) {
	if p.choice != "Sword" || !p.attacking {
		return Rect{}, false
	}

	width := 60.0
	height := 40.0

	x := p.X
	if p.facingRight {
		x += p.Width / 2
	} else {
		x -= p.Width/2 + width
	}

	y := p.Y + p.Height*0.3

	return NewRect(x, y, width, height), true
}

func (p *Player) GetDashHitbox() (Rect, bool) {
	if !p.dashing || p.choice != "Sword" || p.dashHit {
		return Rect{}, false
	}

	w := 90.0
	h := 50.0

	x := p.X
	if p.facingRight {
		x += p.Width / 2
	} else {
		x -= p.Width/2 + w
	}

	y := p.Y + p.Height*0.3

	return NewRect(x, y, w, h), true
}

func (p *Player) GetRect() Rect {
    width := p.Width
    height := p.Height

    if p.choice == "Bow" {
        width = p.Width * 1.3
        height = p.Height * 1.2
    }

    return NewRect(
        p.X-width/2,
        p.Y,
        width,
        height,
    )
}