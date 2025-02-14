package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x84, STY, cpu.AddressingModeZeroPage, 2, 3, "STY")
	cpu.RegisterOpcode(0x94, STY, cpu.AddressingModeZeroPageY, 2, 4, "STY")
	cpu.RegisterOpcode(0x8C, STY, cpu.AddressingModeAbsolute, 3, 4, "STY")
}

// STY Store Y Register
func STY(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)

	c.WriteMemory(addr, c.RegisterY)
}
