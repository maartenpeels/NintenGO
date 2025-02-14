package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xBA, TSX, cpu.AddressingModeImplicit, 1, 2, "TSX")
}

// TSX Transfer Stack Pointer to X
func TSX(c *cpu.CPU, _ uint) {
	c.RegisterX = c.StackPointer
	c.SetZeroAndNegativeFlags(c.RegisterX)
}
