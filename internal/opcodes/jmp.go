package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x4C, JMP, cpu.AddressingModeAbsolute, 3, 3, "JMP")
	cpu.RegisterOpcode(0x6C, JMP, cpu.AddressingModeIndirect, 3, 5, "JMP")
}

// JMP Jump
func JMP(c *cpu.CPU, addressingMode uint) {
	c.ProgramCounter = c.GetOpAddress(addressingMode)
}
