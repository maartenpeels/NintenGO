package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xC8, INY, cpu.AddressingModeImplicit, 1, 2, "INY")
}

// INY Increment Y Register
func INY(c *cpu.CPU, _ uint) {
	c.RegisterY++
	c.SetZeroAndNegativeFlags(c.RegisterY)
}
