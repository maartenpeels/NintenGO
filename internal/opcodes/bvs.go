package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x70, BVS, cpu.AddressingModeRelative, 2, "BVS")
}

// BVS Branch if Overflow Set
func BVS(c *cpu.CPU, addressingMode uint) {
	c.BranchIf(c.IsFlagSet(cpu.OverflowFlag))
}
