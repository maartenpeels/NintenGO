package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xAA, TAX, cpu.AddressingModeImplicit, 1, "TAX")
}

// TAX Transfer Accumulator to X
func TAX(c *cpu.CPU, _ uint) {
	c.RegisterX = c.RegisterA
	c.SetZeroAndNegativeFlags(c.RegisterX)
}
