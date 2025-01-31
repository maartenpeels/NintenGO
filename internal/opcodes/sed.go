package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xF8, SED, cpu.AddressingModeImplicit, 1, "SED")
}

// SED Set Decimal Flag
func SED(c *cpu.CPU, _ uint) {
	c.SetFlag(cpu.DecimalMode, true)
}
