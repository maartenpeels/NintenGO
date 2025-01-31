package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xD0, BNE, cpu.AddressingModeRelative, 2, "BNE")
}

// BNE Branch if Not Equal
func BNE(c *cpu.CPU, addressingMode uint) {
	c.BranchIf(!c.IsFlagSet(cpu.ZeroFlag))
}
