package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	// Unofficial SRE opcodes (LSR memory then EOR with accumulator)
	cpu.RegisterOpcode(0x43, SRE, cpu.AddressingModeIndirectX, 2, 8, "*SRE")
	cpu.RegisterOpcode(0x47, SRE, cpu.AddressingModeZeroPage, 2, 5, "*SRE")
	cpu.RegisterOpcode(0x4F, SRE, cpu.AddressingModeAbsolute, 3, 6, "*SRE")
	cpu.RegisterOpcode(0x53, SRE, cpu.AddressingModeIndirectY, 2, 8, "*SRE")
	cpu.RegisterOpcode(0x57, SRE, cpu.AddressingModeZeroPageX, 2, 6, "*SRE")
	cpu.RegisterOpcode(0x5B, SRE, cpu.AddressingModeAbsoluteY, 3, 7, "*SRE")
	cpu.RegisterOpcode(0x5F, SRE, cpu.AddressingModeAbsoluteX, 3, 7, "*SRE")
}

// SRE - Shift Right then EOR (unofficial)
// This instruction shifts the memory location right by one bit and then EORs the result with the accumulator
func SRE(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// Shift right
	c.Status.SetBool(cpu.CarryFlag, (value&0x01) != 0)
	value >>= 1

	// Write the shifted value back to memory
	c.WriteMemory(addr, value)

	// EOR with accumulator
	c.RegisterA ^= value

	// Set zero and negative flags based on result
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
