package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x28, PLP, cpu.AddressingModeImplicit, 1, 4, "PLP")
}

// PLP Pull Processor Status
func PLP(c *cpu.CPU, _ uint) {
	// When pulling from stack:
	// - Break flag (bit 4) should be cleared
	// - Unused flag (bit 5) should be set
	value := c.PopStack()
	value &= ^uint8(cpu.BreakCommand) // Clear Break flag
	value |= cpu.NotUsedFlag          // Set NotUsed flag
	c.Status.Update(value)
}
