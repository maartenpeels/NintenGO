package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xC0, CPY, cpu.AddressingModeImmediate, 2, 2, "CPY")
	cpu.RegisterOpcode(0xC4, CPY, cpu.AddressingModeZeroPage, 2, 3, "CPY")
	cpu.RegisterOpcode(0xCC, CPY, cpu.AddressingModeAbsolute, 3, 4, "CPY")
}

// CPY Compare Y Register
func CPY(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// C - Set if Y >= M
	c.Status.SetBool(cpu.CarryFlag, c.RegisterY >= value)

	// Z - Set if Y = M
	c.Status.SetBool(cpu.ZeroFlag, c.RegisterY == value)

	// N - Set if bit 7 of the result is set
	result := c.RegisterY - value
	c.Status.SetBool(cpu.NegativeFlag, (result&0x80) != 0)
}
