package opcodes

import (
	"NintenGo/internal/cpu"
)

func init() {
	// RRA (ROR + ADC) - Unofficial opcode
	cpu.RegisterOpcode(0x77, RRA, cpu.AddressingModeZeroPageX, 2, "RRA")
	// Add other addressing modes if needed
}

// RRA Rotate Right then Add with Carry (Unofficial)
func RRA(c *cpu.CPU, addressingMode uint) {
	address := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(address)

	// ROR part
	oldCarry := c.IsFlagSet(cpu.CarryFlag)
	c.SetFlag(cpu.CarryFlag, value&0x01 != 0)
	value = value >> 1
	if oldCarry {
		value |= 0x80
	}
	c.WriteMemory(address, value)

	// ADC part
	c.AddToRegisterA(value)
	c.SetZeroAndNegativeFlags(c.RegisterA)
}
