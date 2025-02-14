package opcodes

import "NintenGo/internal/cpu"

func init() {
	cpu.RegisterOpcode(0x2A, ROLa, cpu.AddressingModeAccumulator, 1, 2, "ROL")
	cpu.RegisterOpcode(0x26, ROL, cpu.AddressingModeZeroPage, 2, 5, "ROL")
	cpu.RegisterOpcode(0x36, ROL, cpu.AddressingModeZeroPageX, 2, 6, "ROL")
	cpu.RegisterOpcode(0x2E, ROL, cpu.AddressingModeAbsolute, 3, 6, "ROL")
	cpu.RegisterOpcode(0x3E, ROL, cpu.AddressingModeAbsoluteX, 3, 7, "ROL")
}

// ROLa Rotate Left Accumulator
func ROLa(c *cpu.CPU, _ uint) {
	value := c.RegisterA

	oldCarry := c.Status.Contains(cpu.CarryFlag)
	c.Status.SetBool(cpu.CarryFlag, value>>7 == 1)

	value <<= 1
	if oldCarry {
		value |= 1
	}

	c.RegisterA = value
	c.SetZeroAndNegativeFlags(value)
}

// ROL Rotate Left
func ROL(c *cpu.CPU, addressingMode uint) {
	addr := c.GetOpAddress(addressingMode)
	value := c.ReadMemory(addr)

	oldCarry := c.Status.Contains(cpu.CarryFlag)
	c.Status.SetBool(cpu.CarryFlag, value>>7 == 1)

	value <<= 1
	if oldCarry {
		value |= 1
	}

	c.WriteMemory(addr, value)
	c.SetZeroAndNegativeFlags(value)
}
