package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x0a, ASLA, cpu.AddressingModeAccumulator, 1, "ASL")
	cpu.RegisterOpcode(0x06, ASL, cpu.AddressingModeZeroPage, 2, "ASL")
	cpu.RegisterOpcode(0x16, ASL, cpu.AddressingModeZeroPageX, 2, "ASL")
	cpu.RegisterOpcode(0x0e, ASL, cpu.AddressingModeAbsolute, 3, "ASL")
	cpu.RegisterOpcode(0x1e, ASL, cpu.AddressingModeAbsoluteX, 3, "ASL")
}

// ASLA Arithmetic Shift Left (accumulator)
func ASLA(c *cpu.CPU, addressingMode uint) {
	value := c.RegisterA

	c.SetFlag(cpu.CarryFlag, value&0x80 == 1)
	value <<= 1

	c.RegisterA = value
}

// ASL Arithmetic Shift Left
func ASL(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	c.SetFlag(cpu.CarryFlag, value&0x80 != 0)
	value <<= 1

	c.WriteMemory(addr, value)
	c.SetZeroAndNegativeFlags(value)
}
