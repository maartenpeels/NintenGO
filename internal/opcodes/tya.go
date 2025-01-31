package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x98, TYA, cpu.AddressingModeImplicit, 1, "TYA")
}

// TYA Transfer Y to Accumulator
func TYA(c *cpu.CPU, _ uint) {
	c.RegisterA = c.RegisterY
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
