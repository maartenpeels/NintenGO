package cpu

const (
	CarryFlag        uint8 = 1 << 0 // Bit 0
	ZeroFlag         uint8 = 1 << 1 // Bit 1
	InterruptDisable uint8 = 1 << 2 // Bit 2
	DecimalMode      uint8 = 1 << 3 // Bit 3
	BreakCommand     uint8 = 1 << 4 // Bit 4
	OverflowFlag     uint8 = 1 << 6 // Bit 6
	NegativeFlag     uint8 = 1 << 7 // Bit 7
)

func (c *CPU) SetFlag(flag uint8, set bool) {
	if set {
		c.Status |= flag
	} else {
		c.Status &^= flag
	}
}

func (c *CPU) IsFlagSet(flag uint8) bool {
	return c.Status&flag > 0
}
