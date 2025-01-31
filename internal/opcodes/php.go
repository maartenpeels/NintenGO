package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x08, PHP, cpu.AddressingModeImplicit, 1, "PHP")
}

// PHP No Operation
func PHP(c *cpu.CPU, _ uint) {
	c.PushStack(c.Status)
}
