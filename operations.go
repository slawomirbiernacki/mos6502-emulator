package main

func (c *Cpu) adc(value byte) {
	//TODO maybe replace masks with shifts?
	sum := uint16(value) + uint16(c.A) + uint16(c.C)
	if (sum & 0b10000000) != uint16(c.A&0b10000000) {
		c.V = 1
	} else {
		c.V = 0
	}
	if sum == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
	c.N = c.A >> 7
	//TODO BCD
	if sum > 255 {
		c.C = 1
	} else {
		c.C = 0
	}
	c.A = byte(sum & 0xFF)
}

func (c *Cpu) sbc(value byte) {
	//TODO maybe replace masks with shifts?
	sub := int16(value) - int16(c.A) - ^int16(c.C) // TODO not sure if whole carry should be negated or just last bit, not sure if it should be int16 or uint16
	if (sub > 127) || (sub < -127) {
		c.V = 1
	} else {
		c.V = 0
	}
	if sub == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
	c.N = byte(sub) >> 7
	//TODO BCD
	if sub >= 0 {
		c.C = 1
	} else {
		c.C = 0
	}
	c.A = byte(sub & 0xFF)
}

func (c *Cpu) and(value byte) {
	c.A = c.A & value
	c.N = c.A >> 7
	if c.A == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) ora(value byte) {
	c.A = c.A | value
	c.N = c.A >> 7
	if c.A == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) eor(value byte) {
	c.A = c.A ^ value
	c.N = c.A >> 7
	if c.A == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) cmp(value byte) {
	compared := c.A - value
	c.N = compared >> 7
	if c.A >= value {
		c.C = 1
	} else {
		c.C = 0
	}
	if compared == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) lda(value byte) {
	c.A = value
	if value == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
	c.N = value >> 7
}

func (c *Cpu) asl(value byte) byte {
	c.C = value >> 7
	shifted := byte((uint16(value) << 1) & 0xFE)
	c.N = shifted >> 7
	if shifted == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
	return shifted
}

func (c *Cpu) rol(value byte) byte {
	carry := value >> 7

	rolled := byte((uint16(value) << 1) & 0xFE)
	rolled = rolled | c.C
	c.C = carry
	c.N = rolled >> 7
	if rolled == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
	return rolled
}

// Logical shift right
func (c *Cpu) lsr(value byte) byte {
	c.N = 0
	c.C = value | 0b00000001
	shifted := byte((uint16(value) >> 1) & 0x7F)
	if shifted == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
	return shifted
}

func (c *Cpu) ror(value byte) byte {
	carry := value | 0b00000001

	rolled := byte((uint16(value) >> 1) & 0x7F)
	if c.C == 1 {
		rolled = rolled | 0x80
	} else {
		rolled = rolled | 0x00
	}

	c.C = carry
	c.N = rolled >> 7
	if rolled == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
	return rolled
}

func (c *Cpu) ldx(value byte) {
	c.X = value
	c.N = value >> 7
	if value == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) dec(value byte) byte {
	value = (value - 1) & 0xFF
	c.N = value >> 7
	if value == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
	return value
}

func (c *Cpu) inc(value byte) byte {
	value = (value + 1) & 0xFF
	c.N = value >> 7
	if value == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
	return value
}

func (c *Cpu) bit(value byte) {
	test := c.A & value
	c.N = test >> 7
	c.V = (test & 0b01000000) >> 6
	if value == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) ldy(value byte) {
	c.Y = value
	c.N = value >> 7
	if value == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) cpy(value byte) {
	compared := c.Y - value
	c.N = compared >> 7
	if c.Y >= value {
		c.C = 1
	} else {
		c.C = 0
	}
	if compared == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) cpx(value byte) {
	compared := c.X - value
	c.N = compared >> 7
	if c.X >= value {
		c.C = 1
	} else {
		c.C = 0
	}
	if compared == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) pla(value byte) {
	c.N = value >> 7
	if value == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) dey() {
	c.Y = c.Y - 1
	c.N = c.Y >> 7
	if c.Y == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) dex() {
	c.X = c.X - 1
	c.N = c.X >> 7
	if c.X == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) iny() {
	c.Y = c.Y + 1
	c.N = c.Y >> 7
	if c.Y == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) tay() {
	c.Y = c.A
	c.N = c.Y >> 7
	if c.Y == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) tya() {
	c.A = c.Y
	c.N = c.A >> 7
	if c.A == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) inx() {
	c.X = c.X + 1
	c.N = c.X >> 7
	if c.X == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) txa() {
	c.A = c.X
	c.N = c.A >> 7
	if c.A == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) tax() {
	c.X = c.A
	c.N = c.X >> 7
	if c.X == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *Cpu) tsx() {
	c.X = c.S
	c.N = c.X >> 7
	if c.X == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}
