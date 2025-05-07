package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	// Unofficial LAX opcodes (Load Accumulator and X register)
	cpu.RegisterOpcode(0xA3, LAX, cpu.AddressingModeIndirectX, 2, 6, "*LAX")
	cpu.RegisterOpcode(0xA7, LAX, cpu.AddressingModeZeroPage, 2, 3, "*LAX")
	cpu.RegisterOpcode(0xAF, LAX, cpu.AddressingModeAbsolute, 3, 4, "*LAX")
	cpu.RegisterOpcode(0xB3, LAX, cpu.AddressingModeIndirectY, 2, 5, "*LAX")
	cpu.RegisterOpcode(0xB7, LAX, cpu.AddressingModeZeroPageY, 2, 4, "*LAX")
	cpu.RegisterOpcode(0xBF, LAX, cpu.AddressingModeAbsoluteY, 3, 4, "*LAX")
}

// LAX - Load Accumulator and X Register (unofficial)
// This instruction loads both the accumulator and the X register with the same value
func LAX(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	// Set both A and X to the same value
	c.RegisterA = value
	c.RegisterX = value

	// Set the zero and negative flags based on the loaded value
	c.SetZeroAndNegativeFlags(value)
}
