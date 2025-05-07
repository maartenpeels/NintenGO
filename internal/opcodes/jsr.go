package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	cpu.RegisterOpcode(0x20, JSR, cpu.AddressingModeAbsolute, 3, 6, "JSR")
}

// JSR Jump to Subroutine
func JSR(c *cpu.CPU, addressingMode uint) {
	// PC currently points to the first byte of the JSR operand (e.g., low byte of absolute address).
	// The JSR instruction is 3 bytes long.
	// The address of the JSR opcode itself is c.ProgramCounter - 1.
	// The address of the last byte of the JSR instruction (the high byte of the target address)
	// is (c.ProgramCounter - 1) + 2 = c.ProgramCounter + 1.
	// This is the value that should be pushed onto the stack.
	returnAddr := c.ProgramCounter + 1
	targetAddr := c.GetOpAddress(addressingMode)
	c.PushStackU16(returnAddr)
	c.ProgramCounter = targetAddr
}
