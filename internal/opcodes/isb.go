package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	// Unofficial ISB opcodes (Increment memory then SBC)
	cpu.RegisterOpcode(0xE3, ISB, cpu.AddressingModeIndirectX, 2, 8, "*ISB")
	cpu.RegisterOpcode(0xE7, ISB, cpu.AddressingModeZeroPage, 2, 5, "*ISB")
	cpu.RegisterOpcode(0xEF, ISB, cpu.AddressingModeAbsolute, 3, 6, "*ISB")
	cpu.RegisterOpcode(0xF3, ISB, cpu.AddressingModeIndirectY, 2, 8, "*ISB")
	cpu.RegisterOpcode(0xF7, ISB, cpu.AddressingModeZeroPageX, 2, 6, "*ISB")
	cpu.RegisterOpcode(0xFB, ISB, cpu.AddressingModeAbsoluteY, 3, 7, "*ISB")
	cpu.RegisterOpcode(0xFF, ISB, cpu.AddressingModeAbsoluteX, 3, 7, "*ISB")
}

// ISB - Increment memory then SBC (unofficial)
// This instruction increments the contents of a memory location and then subtracts
// the result from the accumulator (with borrow)
func ISB(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// Increment the memory value
	value++
	c.WriteMemory(addr, value)

	// Then perform SBC with the incremented value
	// SBC logic: A = A - M - (1 - C)
	if c.Status.Contains(cpu.DecimalMode) {
		// Decimal mode
		a := uint16(c.RegisterA)
		m := uint16(value)
		borrow := uint16(1)
		if c.Status.Contains(cpu.CarryFlag) {
			borrow = 0
		}

		result := a - m - borrow

		// Set carry flag if no borrow required
		c.Status.SetBool(cpu.CarryFlag, result < 0x100)

		// Set overflow flag
		c.Status.SetBool(cpu.OverflowFlag, ((a^m)&0x80 != 0) && ((a^result)&0x80 != 0))

		// Adjust for BCD
		if (a & 0xF) < (m&0xF)+borrow {
			result = (result - 0x6) & 0xFF
		}

		if result > 0x99 {
			result -= 0x60
		}

		c.RegisterA = uint8(result)
	} else {
		// Binary mode
		a := c.RegisterA
		m := value
		borrow := uint8(1)
		if c.Status.Contains(cpu.CarryFlag) {
			borrow = 0
		}

		result := a - m - borrow

		// Set carry flag if no borrow required
		c.Status.SetBool(cpu.CarryFlag, int(a) >= int(m)+int(borrow))

		// Set overflow flag
		c.Status.SetBool(cpu.OverflowFlag, ((a^m)&0x80 != 0) && ((a^result)&0x80 != 0))

		c.RegisterA = result
	}

	// Set zero and negative flags based on result
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
