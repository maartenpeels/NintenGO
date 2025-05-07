package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x0a, ASLA, cpu.AddressingModeAccumulator, 1, 2, "ASL")
	cpu.RegisterOpcode(0x06, ASL, cpu.AddressingModeZeroPage, 2, 5, "ASL")
	cpu.RegisterOpcode(0x16, ASL, cpu.AddressingModeZeroPageX, 2, 6, "ASL")
	cpu.RegisterOpcode(0x0e, ASL, cpu.AddressingModeAbsolute, 3, 6, "ASL")
	cpu.RegisterOpcode(0x1e, ASL, cpu.AddressingModeAbsoluteX, 3, 7, "ASL")
}

// ASLA Arithmetic Shift Left (accumulator)
func ASLA(c *cpu.CPU, addressingMode uint) {
	value := c.RegisterA

	c.Status.SetBool(cpu.CarryFlag, value&0x80 != 0)
	value <<= 1

	c.RegisterA = value
	c.SetZeroAndNegativeFlags(value)
}

// ASL Arithmetic Shift Left
func ASL(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	c.Status.SetBool(cpu.CarryFlag, value&0x80 != 0)
	value <<= 1

	c.WriteMemory(addr, value)
	c.SetZeroAndNegativeFlags(value)
}
