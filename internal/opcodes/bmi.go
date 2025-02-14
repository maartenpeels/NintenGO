package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x30, BMI, cpu.AddressingModeRelative, 2, 2, "BMI")
}

// BMI Branch if Minus
func BMI(c *cpu.CPU, addressingMode uint) {
	c.BranchIf(c.Status.Contains(cpu.NegativeFlag))
}
