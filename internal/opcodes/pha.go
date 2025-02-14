package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x48, PHA, cpu.AddressingModeImplicit, 1, 3, "PHA")
}

// PHA Push Accumulator
func PHA(c *cpu.CPU, _ uint) {
	c.PushStack(c.RegisterA)
}
