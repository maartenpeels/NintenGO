package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x78, SEI, cpu.AddressingModeImplicit, 1, 2, "SEI")
}

// SEI Set Interrupt Disable
func SEI(c *cpu.CPU, _ uint) {
	c.Status.Set(cpu.InterruptDisable)
}
