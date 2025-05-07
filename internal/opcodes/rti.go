package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x40, RTI, cpu.AddressingModeImplicit, 1, 6, "RTI")
}

// RTI Return from Interrupt
func RTI(c *cpu.CPU, _ uint) {
	// When pulling status from stack during RTI:
	// - Break flag (bit 4) should be cleared
	// - Unused flag (bit 5) should be set
	value := c.PopStack()
	value &= ^uint8(cpu.BreakCommand) // Clear Break flag
	value |= cpu.NotUsedFlag          // Set NotUsed flag
	c.Status.Update(value)

	c.ProgramCounter = c.PopStackU16()
}
