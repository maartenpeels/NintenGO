package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xF0, BEQ, cpu.AddressingModeRelative, 2, "BEQ")
}

// BEQ Branch if Equal
func BEQ(c *cpu.CPU, addressingMode uint) {
	c.BranchIf(c.IsFlagSet(cpu.ZeroFlag))
}
