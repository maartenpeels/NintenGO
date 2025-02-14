package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x86, STX, cpu.AddressingModeZeroPage, 2, 3, "STX")
	cpu.RegisterOpcode(0x96, STX, cpu.AddressingModeZeroPageY, 2, 4, "STX")
	cpu.RegisterOpcode(0x8E, STX, cpu.AddressingModeAbsolute, 3, 4, "STX")
}

// STX Store X Register
func STX(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)

	c.WriteMemory(addr, c.RegisterX)
}
