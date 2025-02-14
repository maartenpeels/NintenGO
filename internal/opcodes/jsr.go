package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x20, JSR, cpu.AddressingModeAbsolute, 3, 6, "JSR")
}

// JSR Jump to Subroutine
func JSR(c *cpu.CPU, addressingMode uint) {
	returnAddr := c.ProgramCounter - 1
	targetAddr := c.GetOpAddress(addressingMode)
	c.PushStackU16(returnAddr)
	c.ProgramCounter = targetAddr
}
