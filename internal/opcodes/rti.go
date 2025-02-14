package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x40, RTI, cpu.AddressingModeImplicit, 1, 6, "RTI")
}

// RTI Return from Interrupt
func RTI(c *cpu.CPU, _ uint) {
	c.Status.Update(c.PopStack())
	c.ProgramCounter = c.PopStackU16()
}
