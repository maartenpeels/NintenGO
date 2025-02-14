package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x38, SEC, cpu.AddressingModeImplicit, 1, 2, "SEC")
}

// SEC Set Carry Flag
func SEC(c *cpu.CPU, _ uint) {
	c.Status.Set(cpu.CarryFlag)
}
