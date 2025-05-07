package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x6A, RORa, cpu.AddressingModeAccumulator, 1, 2, "ROR")
	cpu.RegisterOpcode(0x66, ROR, cpu.AddressingModeZeroPage, 2, 5, "ROR")
	cpu.RegisterOpcode(0x76, ROR, cpu.AddressingModeZeroPageX, 2, 6, "ROR")
	cpu.RegisterOpcode(0x6E, ROR, cpu.AddressingModeAbsolute, 3, 6, "ROR")
	cpu.RegisterOpcode(0x7E, ROR, cpu.AddressingModeAbsoluteX, 3, 7, "ROR")
}

// RORa Rotate Right Accumulator
func RORa(c *cpu.CPU, _ uint) {
	value := c.RegisterA

	// Save the current carry flag
	oldCarry := c.Status.Contains(cpu.CarryFlag)

	// Set carry flag to contents of old bit 0
	c.Status.SetBool(cpu.CarryFlag, value&1 != 0)

	// Shift right
	value >>= 1

	// Put old carry into bit 7
	if oldCarry {
		value |= 0x80
	}

	c.RegisterA = value
	c.SetZeroAndNegativeFlags(value)
}

// ROR Rotate Right
func ROR(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// Save the current carry flag
	oldCarry := c.Status.Contains(cpu.CarryFlag)

	// Set carry flag to contents of old bit 0
	c.Status.SetBool(cpu.CarryFlag, value&1 != 0)

	// Shift right
	value >>= 1

	// Put old carry into bit 7
	if oldCarry {
		value |= 0x80
	}

	c.WriteMemory(addr, value)
	c.SetZeroAndNegativeFlags(value)
}
