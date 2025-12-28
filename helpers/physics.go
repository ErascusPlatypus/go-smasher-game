package helpers

func (p *Player) floatUp() {
	speed := 2.0
	
	for i := 2.0 ; i >= 1.0 ; i-=0.2 {
		p.Y -= speed*i 
	}

	for i := 1.0 ; i <= 2.5 ; i+=0.2 {
		p.Y += speed*i
	}
}