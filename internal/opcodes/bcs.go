package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0xB0, BCS, cpu.AddressingModeRelative, 2, "BCS")
}

// BCS Branch if Carry Set
func BCS(c *cpu.CPU, addressingMode uint) {
	c.BranchIf(c.IsFlagSet(cpu.CarryFlag))
}
