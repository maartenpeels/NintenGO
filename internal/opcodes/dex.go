package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xCA, DEX, cpu.AddressingModeImplicit, 1, "DEX")
}

// DEX Decrement X Register
func DEX(c *cpu.CPU, _ uint) {
	c.RegisterX--
	c.SetZeroAndNegativeFlags(c.RegisterX)
}
