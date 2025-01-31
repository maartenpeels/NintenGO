package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xC9, CMP, cpu.AddressingModeImmediate, 2, "CMP")
	cpu.RegisterOpcode(0xC5, CMP, cpu.AddressingModeZeroPage, 2, "CMP")
	cpu.RegisterOpcode(0xD5, CMP, cpu.AddressingModeZeroPageX, 2, "CMP")
	cpu.RegisterOpcode(0xCD, CMP, cpu.AddressingModeAbsolute, 3, "CMP")
	cpu.RegisterOpcode(0xDD, CMP, cpu.AddressingModeAbsoluteX, 3, "CMP")
	cpu.RegisterOpcode(0xD9, CMP, cpu.AddressingModeAbsoluteY, 3, "CMP")
	cpu.RegisterOpcode(0xC1, CMP, cpu.AddressingModeIndirectX, 2, "CMP")
	cpu.RegisterOpcode(0xD1, CMP, cpu.AddressingModeIndirectY, 2, "CMP")
}

// CMP Compare
func CMP(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	cmp := c.RegisterA - value
	c.SetFlag(cpu.CarryFlag, c.RegisterA >= value)
	c.SetZeroAndNegativeFlags(cmp)
}
