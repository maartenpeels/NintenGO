package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xEA, NOP, cpu.AddressingModeImplicit, 1, 2, "NOP")
}

// NOP No Operation
func NOP(c *cpu.CPU, _ uint) {
	// Do nothing
}
