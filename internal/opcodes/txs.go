package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x9A, TXS, cpu.AddressingModeImplicit, 1, 2, "TXS")
}

// TXS Transfer X to Stack Pointer
func TXS(c *cpu.CPU, _ uint) {
	c.StackPointer = c.RegisterX
}
