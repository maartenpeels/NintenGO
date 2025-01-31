package opcodes

import (
	"NintenGo/internal/common"
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x60, RTS, cpu.AddressingModeImplicit, 1, "RTS")
}

// RTS Return from Subroutine
func RTS(c *cpu.CPU, _ uint) {
	returnAddr := c.PopStackU16() + 1
	common.Log.Debugf("RTS: Returning to %04X", returnAddr)
	c.ProgramCounter = returnAddr
}
