package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	// Unofficial RRA opcodes (ROR memory then ADC with accumulator)
	cpu.RegisterOpcode(0x63, RRA, cpu.AddressingModeIndirectX, 2, 8, "*RRA")
	cpu.RegisterOpcode(0x67, RRA, cpu.AddressingModeZeroPage, 2, 5, "*RRA")
	cpu.RegisterOpcode(0x6F, RRA, cpu.AddressingModeAbsolute, 3, 6, "*RRA")
	cpu.RegisterOpcode(0x73, RRA, cpu.AddressingModeIndirectY, 2, 8, "*RRA")
	cpu.RegisterOpcode(0x77, RRA, cpu.AddressingModeZeroPageX, 2, 6, "*RRA")
	cpu.RegisterOpcode(0x7B, RRA, cpu.AddressingModeAbsoluteY, 3, 7, "*RRA")
	cpu.RegisterOpcode(0x7F, RRA, cpu.AddressingModeAbsoluteX, 3, 7, "*RRA")
}

// RRA - Rotate Right then ADC (unofficial)
// This instruction rotates the memory location right by one bit (including carry) and then adds the result to the accumulator (with carry)
func RRA(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// Rotate right through carry
	oldCarry := c.Status.Contains(cpu.CarryFlag)
	c.Status.SetBool(cpu.CarryFlag, (value&0x01) != 0)
	value >>= 1
	if oldCarry {
		value |= 0x80
	}

	// Write the rotated value back to memory
	c.WriteMemory(addr, value)

	// Then perform ADC with the rotated value
	// ADC logic: A = A + M + C
	if c.Status.Contains(cpu.DecimalMode) {
		// Decimal mode
		a := uint16(c.RegisterA)
		m := uint16(value)
		carry := uint16(0)
		if c.Status.Contains(cpu.CarryFlag) {
			carry = 1
		}

		// Add as BCD
		lo := (a & 0x0F) + (m & 0x0F) + carry
		hi := (a & 0xF0) + (m & 0xF0)

		// Handle low digit overflow
		if lo > 0x09 {
			hi += 0x10
			lo += 0x06
		}

		// Handle high digit overflow
		if hi > 0x90 {
			hi += 0x60
		}

		// Set carry flag
		c.Status.SetBool(cpu.CarryFlag, hi > 0xFF)

		// Set overflow flag (signed overflow)
		c.Status.SetBool(cpu.OverflowFlag, ((a^m)&0x80) == 0 && ((a^(hi+lo))&0x80) != 0)

		// Set the result
		c.RegisterA = uint8((hi + (lo & 0x0F)) & 0xFF)
	} else {
		// Binary mode
		a := uint16(c.RegisterA)
		m := uint16(value)
		carry := uint16(0)
		if c.Status.Contains(cpu.CarryFlag) {
			carry = 1
		}

		result := a + m + carry

		// Set carry flag
		c.Status.SetBool(cpu.CarryFlag, result > 0xFF)

		// Set overflow flag (signed overflow)
		c.Status.SetBool(cpu.OverflowFlag, ((a^result)&0x80) != 0 && ((m^result)&0x80) != 0)

		// Set the result
		c.RegisterA = uint8(result & 0xFF)
	}

	// Set zero and negative flags based on result
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
