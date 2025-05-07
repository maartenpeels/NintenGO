package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xE0, CPX, cpu.AddressingModeImmediate, 2, 2, "CPX")
	cpu.RegisterOpcode(0xE4, CPX, cpu.AddressingModeZeroPage, 2, 3, "CPX")
	cpu.RegisterOpcode(0xEC, CPX, cpu.AddressingModeAbsolute, 3, 4, "CPX")
}

// CPX Compare X Register
func CPX(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// C - Set if X >= M
	c.Status.SetBool(cpu.CarryFlag, c.RegisterX >= value)

	// Z - Set if X = M
	c.Status.SetBool(cpu.ZeroFlag, c.RegisterX == value)

	// N - Set if bit 7 of the result is set
	result := c.RegisterX - value
	c.Status.SetBool(cpu.NegativeFlag, (result&0x80) != 0)
}
