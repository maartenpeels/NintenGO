package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x50, BVC, cpu.AddressingModeRelative, 2, 2, "BVC")
}

// BVC Branch if Overflow Clear
func BVC(c *cpu.CPU, addressingMode uint) {
	c.BranchIf(!c.Status.Contains(cpu.OverflowFlag))
}
