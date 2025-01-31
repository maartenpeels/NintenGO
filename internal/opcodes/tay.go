package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xA8, TAY, cpu.AddressingModeImplicit, 1, "TAY")
}

// TAY Transfer Accumulator to Y
func TAY(c *cpu.CPU, _ uint) {
	c.RegisterY = c.RegisterA
	c.SetZeroAndNegativeFlags(c.RegisterY)
}
