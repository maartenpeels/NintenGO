package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xC0, CPY, cpu.AddressingModeImmediate, 2, "CPY")
	cpu.RegisterOpcode(0xC4, CPY, cpu.AddressingModeZeroPage, 2, "CPY")
	cpu.RegisterOpcode(0xCC, CPY, cpu.AddressingModeAbsolute, 3, "CPY")
}

// CPY Compare Y Register
func CPY(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	cmp := c.RegisterY - value
	c.SetFlag(cpu.CarryFlag, c.RegisterY >= value)
	c.SetZeroAndNegativeFlags(cmp)
}
