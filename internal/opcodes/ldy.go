package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0xA0, LDY, cpu.AddressingModeImmediate, 2, "LDY")
	cpu.RegisterOpcode(0xA4, LDY, cpu.AddressingModeZeroPage, 2, "LDY")
	cpu.RegisterOpcode(0xB4, LDY, cpu.AddressingModeZeroPageX, 2, "LDY")
	cpu.RegisterOpcode(0xAC, LDY, cpu.AddressingModeAbsolute, 3, "LDY")
	cpu.RegisterOpcode(0xBC, LDY, cpu.AddressingModeAbsoluteX, 3, "LDY")
}

// LDY Load Y Register
func LDY(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	c.RegisterY = value
	c.SetZeroAndNegativeFlags(value)
}
