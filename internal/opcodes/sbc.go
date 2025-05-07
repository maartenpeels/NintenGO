package opcodes

import "NintenGo/internal/cpu"

func init() {
	// Official SBC opcodes
	cpu.RegisterOpcode(0xE9, SBC, cpu.AddressingModeImmediate, 2, 2, "SBC")
	cpu.RegisterOpcode(0xE5, SBC, cpu.AddressingModeZeroPage, 2, 3, "SBC")
	cpu.RegisterOpcode(0xF5, SBC, cpu.AddressingModeZeroPageX, 2, 4, "SBC")
	cpu.RegisterOpcode(0xED, SBC, cpu.AddressingModeAbsolute, 3, 4, "SBC")
	cpu.RegisterOpcode(0xFD, SBC, cpu.AddressingModeAbsoluteX, 3, 4, "SBC")
	cpu.RegisterOpcode(0xF9, SBC, cpu.AddressingModeAbsoluteY, 3, 4, "SBC")
	cpu.RegisterOpcode(0xE1, SBC, cpu.AddressingModeIndirectX, 2, 6, "SBC")
	cpu.RegisterOpcode(0xF1, SBC, cpu.AddressingModeIndirectY, 2, 5, "SBC")

	// Unofficial SBC opcode
	cpu.RegisterOpcode(0xEB, SBC, cpu.AddressingModeImmediate, 2, 2, "*SBC")
}

// SBC Subtract with Carry
func SBC(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	carry := uint8(0)
	if c.Status.Contains(cpu.CarryFlag) {
		carry = 1
	}

	result := c.RegisterA - value - (1 - carry)

	c.Status.SetBool(cpu.CarryFlag, int(c.RegisterA)-int(value)-int(1-carry) >= 0)
	c.Status.SetBool(cpu.OverflowFlag, (c.RegisterA^result)&(c.RegisterA^value)&0x80 != 0)
	c.RegisterA = result

	c.SetZeroAndNegativeFlags(c.RegisterA)
}
