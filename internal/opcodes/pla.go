package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x68, PLA, cpu.AddressingModeImplicit, 1, 4, "PLA")
}

// PLA Pull Accumulator
func PLA(c *cpu.CPU, _ uint) {
	c.RegisterA = c.PopStack()
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
