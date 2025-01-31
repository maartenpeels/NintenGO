package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xD8, CLD, cpu.AddressingModeImplicit, 1, "CLD")
}

// CLD Clear Decimal Mode
func CLD(c *cpu.CPU, _ uint) {
	c.SetFlag(cpu.DecimalMode, false)
}
