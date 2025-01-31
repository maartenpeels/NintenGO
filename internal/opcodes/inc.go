package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xE6, INC, cpu.AddressingModeZeroPage, 2, "INC")
	cpu.RegisterOpcode(0xF6, INC, cpu.AddressingModeZeroPageX, 2, "INC")
	cpu.RegisterOpcode(0xEE, INC, cpu.AddressingModeAbsolute, 3, "INC")
	cpu.RegisterOpcode(0xFE, INC, cpu.AddressingModeAbsoluteX, 3, "INC")
}

// INC Decrement Memory
func INC(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	newValue := value + 1
	c.WriteMemory(addr, newValue)
	c.SetZeroAndNegativeFlags(newValue)
}
