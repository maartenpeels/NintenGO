package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xC9, CMP, cpu.AddressingModeImmediate, 2, 2, "CMP")
	cpu.RegisterOpcode(0xC5, CMP, cpu.AddressingModeZeroPage, 2, 3, "CMP")
	cpu.RegisterOpcode(0xD5, CMP, cpu.AddressingModeZeroPageX, 2, 4, "CMP")
	cpu.RegisterOpcode(0xCD, CMP, cpu.AddressingModeAbsolute, 3, 4, "CMP")
	cpu.RegisterOpcode(0xDD, CMP, cpu.AddressingModeAbsoluteX, 3, 4, "CMP")
	cpu.RegisterOpcode(0xD9, CMP, cpu.AddressingModeAbsoluteY, 3, 4, "CMP")
	cpu.RegisterOpcode(0xC1, CMP, cpu.AddressingModeIndirectX, 2, 6, "CMP")
	cpu.RegisterOpcode(0xD1, CMP, cpu.AddressingModeIndirectY, 2, 5, "CMP")
}

// CMP Compare
func CMP(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// C - Set if A >= M
	c.Status.SetBool(cpu.CarryFlag, c.RegisterA >= value)

	// Z - Set if A = M
	c.Status.SetBool(cpu.ZeroFlag, c.RegisterA == value)

	// N - Set if bit 7 of the result is set
	result := c.RegisterA - value
	c.Status.SetBool(cpu.NegativeFlag, (result&0x80) != 0)
}
