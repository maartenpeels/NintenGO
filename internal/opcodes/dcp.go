package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	// Unofficial DCP opcodes (Decrement memory then CMP)
	cpu.RegisterOpcode(0xC3, DCP, cpu.AddressingModeIndirectX, 2, 8, "*DCP")
	cpu.RegisterOpcode(0xC7, DCP, cpu.AddressingModeZeroPage, 2, 5, "*DCP")
	cpu.RegisterOpcode(0xCF, DCP, cpu.AddressingModeAbsolute, 3, 6, "*DCP")
	cpu.RegisterOpcode(0xD3, DCP, cpu.AddressingModeIndirectY, 2, 8, "*DCP")
	cpu.RegisterOpcode(0xD7, DCP, cpu.AddressingModeZeroPageX, 2, 6, "*DCP")
	cpu.RegisterOpcode(0xDB, DCP, cpu.AddressingModeAbsoluteY, 3, 7, "*DCP")
	cpu.RegisterOpcode(0xDF, DCP, cpu.AddressingModeAbsoluteX, 3, 7, "*DCP")
}

// DCP - Decrement memory then Compare (unofficial)
// This instruction decrements the contents of a memory location and then compares
// the result with the accumulator
func DCP(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// Decrement the memory value
	value--
	c.WriteMemory(addr, value)

	// Then perform CMP with the decremented value
	result := c.RegisterA - value

	// Set carry flag if A >= M
	c.Status.SetBool(cpu.CarryFlag, c.RegisterA >= value)

	// Set zero flag if A == M
	c.Status.SetBool(cpu.ZeroFlag, c.RegisterA == value)

	// Set negative flag based on result
	c.Status.SetBool(cpu.NegativeFlag, (result&0x80) != 0)
}
