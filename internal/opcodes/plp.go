package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x28, PLP, cpu.AddressingModeImplicit, 1, "PLP")
}

// PLP Pull Processor Status
func PLP(c *cpu.CPU, _ uint) {
	c.Status = c.PopStack()
}
