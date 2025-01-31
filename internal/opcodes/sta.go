package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x85, STA, cpu.AddressingModeZeroPage, 2, "STA")
	cpu.RegisterOpcode(0x95, STA, cpu.AddressingModeZeroPageX, 2, "STA")
	cpu.RegisterOpcode(0x8D, STA, cpu.AddressingModeAbsolute, 3, "STA")
	cpu.RegisterOpcode(0x9D, STA, cpu.AddressingModeAbsoluteX, 3, "STA")
	cpu.RegisterOpcode(0x99, STA, cpu.AddressingModeAbsoluteY, 3, "STA")
	cpu.RegisterOpcode(0x81, STA, cpu.AddressingModeIndirectX, 2, "STA")
	cpu.RegisterOpcode(0x91, STA, cpu.AddressingModeIndirectY, 2, "STA")
}

// STA Store Accumulator
func STA(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	c.WriteMemory(addr, c.RegisterA)
}
