package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xC6, DEC, cpu.AddressingModeZeroPage, 2, 5, "DEC")
	cpu.RegisterOpcode(0xD6, DEC, cpu.AddressingModeZeroPageX, 2, 6, "DEC")
	cpu.RegisterOpcode(0xCE, DEC, cpu.AddressingModeAbsolute, 3, 6, "DEC")
	cpu.RegisterOpcode(0xDE, DEC, cpu.AddressingModeAbsoluteX, 3, 7, "DEC")
}

// DEC Decrement Memory
func DEC(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	newValue := value - 1
	c.WriteMemory(addr, newValue)
	c.SetZeroAndNegativeFlags(newValue)
}
