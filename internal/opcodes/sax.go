package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	// Unofficial SAX opcodes (Store A AND X)
	cpu.RegisterOpcode(0x83, SAX, cpu.AddressingModeIndirectX, 2, 6, "*SAX")
	cpu.RegisterOpcode(0x87, SAX, cpu.AddressingModeZeroPage, 2, 3, "*SAX")
	cpu.RegisterOpcode(0x8F, SAX, cpu.AddressingModeAbsolute, 3, 4, "*SAX")
	cpu.RegisterOpcode(0x97, SAX, cpu.AddressingModeZeroPageY, 2, 4, "*SAX")
}

// SAX - Store A AND X (unofficial)
// This instruction ANDs the contents of the A and X registers and stores the result in memory
func SAX(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.RegisterA & c.RegisterX
	c.WriteMemory(addr, value)
}
