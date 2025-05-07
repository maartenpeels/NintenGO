package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	// Unofficial RLA opcodes (ROL memory then AND with accumulator)
	cpu.RegisterOpcode(0x23, RLA, cpu.AddressingModeIndirectX, 2, 8, "*RLA")
	cpu.RegisterOpcode(0x27, RLA, cpu.AddressingModeZeroPage, 2, 5, "*RLA")
	cpu.RegisterOpcode(0x2F, RLA, cpu.AddressingModeAbsolute, 3, 6, "*RLA")
	cpu.RegisterOpcode(0x33, RLA, cpu.AddressingModeIndirectY, 2, 8, "*RLA")
	cpu.RegisterOpcode(0x37, RLA, cpu.AddressingModeZeroPageX, 2, 6, "*RLA")
	cpu.RegisterOpcode(0x3B, RLA, cpu.AddressingModeAbsoluteY, 3, 7, "*RLA")
	cpu.RegisterOpcode(0x3F, RLA, cpu.AddressingModeAbsoluteX, 3, 7, "*RLA")
}

// RLA - Rotate Left then AND (unofficial)
// This instruction rotates the memory location left by one bit (including carry) and then ANDs the result with the accumulator
func RLA(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// Rotate left through carry
	oldCarry := c.Status.Contains(cpu.CarryFlag)
	c.Status.SetBool(cpu.CarryFlag, (value&0x80) != 0)
	value <<= 1
	if oldCarry {
		value |= 0x01
	}

	// Write the rotated value back to memory
	c.WriteMemory(addr, value)

	// AND with accumulator
	c.RegisterA &= value

	// Set zero and negative flags based on result
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
