package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x09, ORA, cpu.AddressingModeImmediate, 2, 2, "ORA")
	cpu.RegisterOpcode(0x05, ORA, cpu.AddressingModeZeroPage, 2, 3, "ORA")
	cpu.RegisterOpcode(0x15, ORA, cpu.AddressingModeZeroPageX, 2, 4, "ORA")
	cpu.RegisterOpcode(0x0D, ORA, cpu.AddressingModeAbsolute, 3, 4, "ORA")
	cpu.RegisterOpcode(0x1D, ORA, cpu.AddressingModeAbsoluteX, 3, 4, "ORA")
	cpu.RegisterOpcode(0x19, ORA, cpu.AddressingModeAbsoluteY, 3, 4, "ORA")
	cpu.RegisterOpcode(0x01, ORA, cpu.AddressingModeIndirectX, 2, 6, "ORA")
	cpu.RegisterOpcode(0x11, ORA, cpu.AddressingModeIndirectY, 2, 5, "ORA")
}

// ORA Logical Inclusive OR
func ORA(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	c.RegisterA |= value

	c.SetZeroAndNegativeFlags(c.RegisterA)
}
