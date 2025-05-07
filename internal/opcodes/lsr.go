package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x4A, LSRa, cpu.AddressingModeAccumulator, 1, 2, "LSR")
	cpu.RegisterOpcode(0x46, LSR, cpu.AddressingModeZeroPage, 2, 5, "LSR")
	cpu.RegisterOpcode(0x56, LSR, cpu.AddressingModeZeroPageX, 2, 6, "LSR")
	cpu.RegisterOpcode(0x4E, LSR, cpu.AddressingModeAbsolute, 3, 6, "LSR")
	cpu.RegisterOpcode(0x5E, LSR, cpu.AddressingModeAbsoluteX, 3, 7, "LSR")
}

// LSRa Logical Shift Right (accumulator)
func LSRa(c *cpu.CPU, addressingMode uint) {
	value := c.RegisterA

	c.Status.SetBool(cpu.CarryFlag, value&0x01 == 1)
	value >>= 1

	c.RegisterA = value
	c.SetZeroAndNegativeFlags(value)
}

// LSR Logical Shift Right
func LSR(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	c.Status.SetBool(cpu.CarryFlag, value&0x01 == 1)
	value >>= 1

	c.WriteMemory(addr, value)
	c.SetZeroAndNegativeFlags(value)
}
