package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x08, PHP, cpu.AddressingModeImplicit, 1, 3, "PHP")
}

// PHP No Operation
func PHP(c *cpu.CPU, _ uint) {
	c.PushStack(c.Status.Value())
}
