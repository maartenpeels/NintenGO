package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x00, BRK, cpu.AddressingModeImplicit, 1, 7, "BRK")
}

// BRK Force Interrupt
func BRK(c *cpu.CPU, _ uint) {
	c.Status.Set(cpu.BreakCommand)
}
