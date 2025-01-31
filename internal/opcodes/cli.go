package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x58, CLI, cpu.AddressingModeImplicit, 1, "CLI")
}

// CLI Clear Interrupt Disable
func CLI(c *cpu.CPU, _ uint) {
	c.SetFlag(cpu.InterruptDisable, false)
}
