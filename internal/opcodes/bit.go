package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x24, BIT, cpu.AddressingModeZeroPage, 2, 3, "BIT")
	cpu.RegisterOpcode(0x2C, BIT, cpu.AddressingModeAbsolute, 3, 4, "BIT")
}

// BIT Bit Test
func BIT(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	c.Status.SetBool(cpu.ZeroFlag, c.RegisterA&value == 0)
	c.Status.SetBool(cpu.OverflowFlag, value&0x40 != 0)
	c.Status.SetBool(cpu.NegativeFlag, value&0x80 != 0)
}
