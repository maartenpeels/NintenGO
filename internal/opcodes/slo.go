package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	// Unofficial SLO opcodes (ASL memory then ORA with accumulator)
	cpu.RegisterOpcode(0x03, SLO, cpu.AddressingModeIndirectX, 2, 8, "*SLO")
	cpu.RegisterOpcode(0x07, SLO, cpu.AddressingModeZeroPage, 2, 5, "*SLO")
	cpu.RegisterOpcode(0x0F, SLO, cpu.AddressingModeAbsolute, 3, 6, "*SLO")
	cpu.RegisterOpcode(0x13, SLO, cpu.AddressingModeIndirectY, 2, 8, "*SLO")
	cpu.RegisterOpcode(0x17, SLO, cpu.AddressingModeZeroPageX, 2, 6, "*SLO")
	cpu.RegisterOpcode(0x1B, SLO, cpu.AddressingModeAbsoluteY, 3, 7, "*SLO")
	cpu.RegisterOpcode(0x1F, SLO, cpu.AddressingModeAbsoluteX, 3, 7, "*SLO")
}

// SLO - Shift Left then OR (unofficial)
// This instruction shifts the memory location left by one bit and then ORs the result with the accumulator
func SLO(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// Shift left
	c.Status.SetBool(cpu.CarryFlag, (value&0x80) != 0)
	value <<= 1

	// Write the shifted value back to memory
	c.WriteMemory(addr, value)

	// OR with accumulator
	c.RegisterA |= value

	// Set zero and negative flags based on result
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
