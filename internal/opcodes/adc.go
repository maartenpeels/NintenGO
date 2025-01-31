package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x69, ADC, cpu.AddressingModeImmediate, 2, "ADC")
	cpu.RegisterOpcode(0x65, ADC, cpu.AddressingModeZeroPage, 2, "ADC")
	cpu.RegisterOpcode(0x75, ADC, cpu.AddressingModeZeroPageX, 2, "ADC")
	cpu.RegisterOpcode(0x6d, ADC, cpu.AddressingModeAbsolute, 3, "ADC")
	cpu.RegisterOpcode(0x7d, ADC, cpu.AddressingModeAbsoluteX, 3, "ADC")
	cpu.RegisterOpcode(0x79, ADC, cpu.AddressingModeAbsoluteY, 3, "ADC")
	cpu.RegisterOpcode(0x61, ADC, cpu.AddressingModeIndirectX, 2, "ADC")
	cpu.RegisterOpcode(0x71, ADC, cpu.AddressingModeIndirectY, 2, "ADC")
}

// ADC Add with Carry
func ADC(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	c.AddToRegisterA(value)

	// Update processor flags
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
