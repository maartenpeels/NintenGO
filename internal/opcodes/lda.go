package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0xA9, LDA, cpu.AddressingModeImmediate, 2, 2, "LDA")
	cpu.RegisterOpcode(0xA5, LDA, cpu.AddressingModeZeroPage, 2, 3, "LDA")
	cpu.RegisterOpcode(0xB5, LDA, cpu.AddressingModeZeroPageX, 2, 4, "LDA")
	cpu.RegisterOpcode(0xAD, LDA, cpu.AddressingModeAbsolute, 3, 4, "LDA")
	cpu.RegisterOpcode(0xBD, LDA, cpu.AddressingModeAbsoluteX, 3, 4, "LDA")
	cpu.RegisterOpcode(0xB9, LDA, cpu.AddressingModeAbsoluteY, 3, 4, "LDA")
	cpu.RegisterOpcode(0xA1, LDA, cpu.AddressingModeIndirectX, 2, 6, "LDA")
	cpu.RegisterOpcode(0xB1, LDA, cpu.AddressingModeIndirectY, 2, 5, "LDA")
}

// LDA loads an immediate value into the accumulator.
func LDA(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	c.RegisterA = value
	c.SetZeroAndNegativeFlags(value)
}
