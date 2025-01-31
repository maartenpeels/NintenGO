package opcodes

import (
	"NintenGo/internal/common"
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x20, JSR, cpu.AddressingModeAbsolute, 3, "JSR")
}

// JSR Jump to Subroutine
func JSR(c *cpu.CPU, addressingMode uint) {
	returnAddr := c.ProgramCounter - 1
	targetAddr := c.GetOpAddress(addressingMode)
	common.Log.Debugf("JSR: Pushing return address %04X, jumping to %04X", returnAddr, targetAddr)
	c.PushStackU16(returnAddr)
	c.ProgramCounter = targetAddr
}
