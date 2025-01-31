package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x8A, TXA, cpu.AddressingModeImplicit, 1, "TXA")
}

// TXA Transfer X to Accumulator
func TXA(c *cpu.CPU, _ uint) {
	c.RegisterA = c.RegisterX
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
