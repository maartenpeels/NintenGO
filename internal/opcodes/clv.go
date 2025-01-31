package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xB8, CLV, cpu.AddressingModeImplicit, 1, "CLV")
}

// CLV Clear Overflow Flag
func CLV(c *cpu.CPU, _ uint) {
	c.SetFlag(cpu.OverflowFlag, false)
}
