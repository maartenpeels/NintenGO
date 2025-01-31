package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x38, SEC, cpu.AddressingModeImplicit, 1, "SEC")
}

// SEC Set Carry Flag
func SEC(c *cpu.CPU, _ uint) {
	c.SetFlag(cpu.CarryFlag, true)
}
