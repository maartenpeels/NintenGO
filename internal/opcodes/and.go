package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x29, AND, cpu.AddressingModeImmediate, 2, 2, "AND")
	cpu.RegisterOpcode(0x25, AND, cpu.AddressingModeZeroPage, 2, 3, "AND")
	cpu.RegisterOpcode(0x35, AND, cpu.AddressingModeZeroPageX, 2, 4, "AND")
	cpu.RegisterOpcode(0x2d, AND, cpu.AddressingModeAbsolute, 3, 4, "AND")
	cpu.RegisterOpcode(0x3d, AND, cpu.AddressingModeAbsoluteX, 3, 4, "AND")
	cpu.RegisterOpcode(0x39, AND, cpu.AddressingModeAbsoluteY, 3, 4, "AND")
	cpu.RegisterOpcode(0x21, AND, cpu.AddressingModeIndirectX, 2, 6, "AND")
	cpu.RegisterOpcode(0x31, AND, cpu.AddressingModeIndirectY, 2, 5, "AND")
}

// AND Logical AND
func AND(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// Perform bitwise AND between accumulator and memory value
	c.RegisterA &= value

	// Update processor flags
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
