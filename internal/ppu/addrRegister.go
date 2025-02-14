package ppu

// AddrRegister is a 16-bit register used to store addresses in the PPU
type AddrRegister struct {
	hi    uint8
	lo    uint8
	hiPtr bool
}

func NewAddrRegister() *AddrRegister {
	return &AddrRegister{
		hi:    0,
		lo:    0,
		hiPtr: true,
	}
}

func (a *AddrRegister) Get() uint16 {
	return uint16(a.hi)<<8 | uint16(a.lo)
}

func (a *AddrRegister) Set(value uint16) {
	a.hi = uint8(value >> 8)
	a.lo = uint8(value & 0xFF)
}

func (a *AddrRegister) Update(value uint8) {
	if a.hiPtr {
		a.hi = value
	} else {
		a.lo = value
	}

	if a.Get() > 0x3FFF {
		a.Set(a.Get() & 0b11111111111111)
	}
}

func (a *AddrRegister) Increment(inc uint8) {
	lo := a.lo
	a.lo += inc
	if lo > a.lo {
		a.hi++
	}

	if a.Get() > 0x3FFF {
		a.Set(a.Get() & 0b11111111111111)
	}
}

func (a *AddrRegister) ResetLatch() {
	a.hiPtr = true
}
