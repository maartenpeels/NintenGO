package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x08, PHP, cpu.AddressingModeImplicit, 1, 3, "PHP")
}

// PHP Push Processor Status
func PHP(c *cpu.CPU, _ uint) {
	// When pushing the status register to the stack:
	// - bit 4 (Break flag) is always set to 1
	// - bit 5 (unused) is always set to 1
	statusValue := c.Status.Value()
	statusValue |= cpu.BreakCommand | cpu.NotUsedFlag
	c.PushStack(statusValue)
}
