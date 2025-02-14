package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x58, CLI, cpu.AddressingModeImplicit, 1, 2, "CLI")
}

// CLI Clear Interrupt Disable
func CLI(c *cpu.CPU, _ uint) {
	c.Status.Clear(cpu.InterruptDisable)
}
