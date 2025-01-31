package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x49, EOR, cpu.AddressingModeImmediate, 2, "EOR")
	cpu.RegisterOpcode(0x45, EOR, cpu.AddressingModeZeroPage, 2, "EOR")
	cpu.RegisterOpcode(0x55, EOR, cpu.AddressingModeZeroPageX, 2, "EOR")
	cpu.RegisterOpcode(0x4D, EOR, cpu.AddressingModeAbsolute, 3, "EOR")
	cpu.RegisterOpcode(0x5D, EOR, cpu.AddressingModeAbsoluteX, 3, "EOR")
	cpu.RegisterOpcode(0x59, EOR, cpu.AddressingModeAbsoluteY, 3, "EOR")
	cpu.RegisterOpcode(0x41, EOR, cpu.AddressingModeIndirectX, 2, "EOR")
	cpu.RegisterOpcode(0x51, EOR, cpu.AddressingModeIndirectY, 2, "EOR")
}

// EOR Exclusive OR
func EOR(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	c.RegisterA ^= value
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
