package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x18, CLC, cpu.AddressingModeImplicit, 1, 2, "CLC")
}

// CLC Clear Carry Flag
func CLC(c *cpu.CPU, _ uint) {
	c.Status.Clear(cpu.CarryFlag)
}
