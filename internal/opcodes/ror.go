package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x6A, RORa, cpu.AddressingModeAccumulator, 1, "ROR")
	cpu.RegisterOpcode(0x66, ROR, cpu.AddressingModeZeroPage, 2, "ROR")
	cpu.RegisterOpcode(0x76, ROR, cpu.AddressingModeZeroPageX, 2, "ROR")
	cpu.RegisterOpcode(0x6E, ROR, cpu.AddressingModeAbsolute, 3, "ROR")
	cpu.RegisterOpcode(0x7E, ROR, cpu.AddressingModeAbsoluteX, 3, "ROR")
}

// RORa Rotate Right Accumulator
func RORa(c *cpu.CPU, _ uint) {
	value := c.RegisterA

	oldCarry := c.IsFlagSet(cpu.CarryFlag)
	c.SetFlag(cpu.CarryFlag, value&1 == 1)

	value >>= 1
	if oldCarry {
		value |= 1
	}

	c.RegisterA = value
	c.SetZeroAndNegativeFlags(value)
}

// ROR Rotate Right
func ROR(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	oldCarry := c.IsFlagSet(cpu.CarryFlag)
	c.SetFlag(cpu.CarryFlag, value&1 == 1)

	value >>= 1
	if oldCarry {
		value |= 1
	}

	c.WriteMemory(addr, value)
	c.SetZeroAndNegativeFlags(value)
}
