package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0xA2, LDX, cpu.AddressingModeImmediate, 2, "LDX")
	cpu.RegisterOpcode(0xA6, LDX, cpu.AddressingModeZeroPage, 2, "LDX")
	cpu.RegisterOpcode(0xB6, LDX, cpu.AddressingModeZeroPageY, 2, "LDX")
	cpu.RegisterOpcode(0xAE, LDX, cpu.AddressingModeAbsolute, 3, "LDX")
	cpu.RegisterOpcode(0xBE, LDX, cpu.AddressingModeAbsoluteY, 3, "LDX")
}

// LDX Load X Register
func LDX(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	c.RegisterX = value
	c.SetZeroAndNegativeFlags(value)
}
