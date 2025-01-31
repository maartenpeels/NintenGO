package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x10, BPL, cpu.AddressingModeRelative, 2, "BPL")
}

// BPL Branch if Positive
func BPL(c *cpu.CPU, addressingMode uint) {
	c.BranchIf(!c.IsFlagSet(cpu.NegativeFlag))
}
