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

	cmp := c.RegisterX - value
	c.Status.SetBool(cpu.CarryFlag, c.RegisterX >= value)
	c.SetZeroAndNegativeFlags(cmp)
}
