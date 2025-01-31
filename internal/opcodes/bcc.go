package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x90, BCC, cpu.AddressingModeRelative, 2, "BCC")
}

// BCC Branch if Carry Clear
func BCC(c *cpu.CPU, addressingMode uint) {
	c.BranchIf(!c.IsFlagSet(cpu.CarryFlag))
}
