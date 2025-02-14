package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x88, DEY, cpu.AddressingModeImplicit, 1, 2, "DEY")
}

// DEY Decrement Y Register
func DEY(c *cpu.CPU, _ uint) {
	c.RegisterY--
	c.SetZeroAndNegativeFlags(c.RegisterY)
}
