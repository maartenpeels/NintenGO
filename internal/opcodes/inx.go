package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xE8, INX, cpu.AddressingModeImplicit, 1, "INX")
}

// INX Increment X Register
func INX(c *cpu.CPU, _ uint) {
	c.RegisterX++
	c.SetZeroAndNegativeFlags(c.RegisterX)
}
