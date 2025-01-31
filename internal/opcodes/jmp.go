package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x4C, JMP, cpu.AddressingModeAbsolute, 3, "JMP")
	cpu.RegisterOpcode(0x6C, JMP, cpu.AddressingModeIndirect, 3, "JMP")
}

// JMP Jump
func JMP(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	c.ProgramCounter = uint16(value)
}
