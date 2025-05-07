package opcodes

import "NintenGo/internal/cpu"

func init() {
	// Official NOP
	cpu.RegisterOpcode(0xEA, NOP, cpu.AddressingModeImplicit, 1, 2, "NOP")

	// Unofficial NOPs with zero page addressing
	cpu.RegisterOpcode(0x04, NOP, cpu.AddressingModeZeroPage, 2, 3, "*NOP")
	cpu.RegisterOpcode(0x44, NOP, cpu.AddressingModeZeroPage, 2, 3, "*NOP")
	cpu.RegisterOpcode(0x64, NOP, cpu.AddressingModeZeroPage, 2, 3, "*NOP")

	// Unofficial NOPs with absolute addressing
	cpu.RegisterOpcode(0x0C, NOP, cpu.AddressingModeAbsolute, 3, 4, "*NOP")

	// Unofficial NOPs with zero page X addressing
	cpu.RegisterOpcode(0x14, NOP, cpu.AddressingModeZeroPageX, 2, 4, "*NOP")
	cpu.RegisterOpcode(0x34, NOP, cpu.AddressingModeZeroPageX, 2, 4, "*NOP")
	cpu.RegisterOpcode(0x54, NOP, cpu.AddressingModeZeroPageX, 2, 4, "*NOP")
	cpu.RegisterOpcode(0x74, NOP, cpu.AddressingModeZeroPageX, 2, 4, "*NOP")
	cpu.RegisterOpcode(0xD4, NOP, cpu.AddressingModeZeroPageX, 2, 4, "*NOP")
	cpu.RegisterOpcode(0xF4, NOP, cpu.AddressingModeZeroPageX, 2, 4, "*NOP")

	// Unofficial NOPs with absolute X addressing
	cpu.RegisterOpcode(0x1C, NOP, cpu.AddressingModeAbsoluteX, 3, 4, "*NOP")
	cpu.RegisterOpcode(0x3C, NOP, cpu.AddressingModeAbsoluteX, 3, 4, "*NOP")
	cpu.RegisterOpcode(0x5C, NOP, cpu.AddressingModeAbsoluteX, 3, 4, "*NOP")
	cpu.RegisterOpcode(0x7C, NOP, cpu.AddressingModeAbsoluteX, 3, 4, "*NOP")
	cpu.RegisterOpcode(0xDC, NOP, cpu.AddressingModeAbsoluteX, 3, 4, "*NOP")
	cpu.RegisterOpcode(0xFC, NOP, cpu.AddressingModeAbsoluteX, 3, 4, "*NOP")

	// Unofficial NOPs with implied addressing (1-byte NOPs)
	cpu.RegisterOpcode(0x1A, NOP, cpu.AddressingModeImplicit, 1, 2, "*NOP")
	cpu.RegisterOpcode(0x3A, NOP, cpu.AddressingModeImplicit, 1, 2, "*NOP")
	cpu.RegisterOpcode(0x5A, NOP, cpu.AddressingModeImplicit, 1, 2, "*NOP")
	cpu.RegisterOpcode(0x7A, NOP, cpu.AddressingModeImplicit, 1, 2, "*NOP")
	cpu.RegisterOpcode(0xDA, NOP, cpu.AddressingModeImplicit, 1, 2, "*NOP")
	cpu.RegisterOpcode(0xFA, NOP, cpu.AddressingModeImplicit, 1, 2, "*NOP")

	// Unofficial 2-byte NOPs with immediate addressing
	cpu.RegisterOpcode(0x80, NOP, cpu.AddressingModeImmediate, 2, 2, "*NOP")
	cpu.RegisterOpcode(0x82, NOP, cpu.AddressingModeImmediate, 2, 2, "*NOP")
	cpu.RegisterOpcode(0x89, NOP, cpu.AddressingModeImmediate, 2, 2, "*NOP")
	cpu.RegisterOpcode(0xC2, NOP, cpu.AddressingModeImmediate, 2, 2, "*NOP")
	cpu.RegisterOpcode(0xE2, NOP, cpu.AddressingModeImmediate, 2, 2, "*NOP")
}

// NOP No Operation
func NOP(c *cpu.CPU, _ uint) {
	// Do nothing
}
