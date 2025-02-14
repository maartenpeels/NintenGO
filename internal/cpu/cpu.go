package cpu

import (
	"NintenGo/internal/common"
	"fmt"
	"strconv"
	"strings"
)

const (
	AddressingModeImplicit = iota
	AddressingModeAccumulator
	AddressingModeImmediate
	AddressingModeZeroPage
	AddressingModeZeroPageX
	AddressingModeZeroPageY
	AddressingModeRelative
	AddressingModeAbsolute
	AddressingModeAbsoluteX
	AddressingModeAbsoluteY
	AddressingModeIndirect
	AddressingModeIndirectX
	AddressingModeIndirectY
)

const (
	StackStart = 0x100
	StackReset = 0xfd
)

type CPU struct {
	RegisterA uint8
	RegisterX uint8
	RegisterY uint8

	Status uint8

	ProgramCounter uint16
	StackPointer   uint8

	Bus *BUS
}

type Error struct {
	PC      uint16
	OpCode  uint8
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("CPU Error at PC:%04X (OpCode:%02X): %s", e.PC, e.OpCode, e.Message)
}

func New(bus *BUS) *CPU {
	return &CPU{
		RegisterA: 0,
		RegisterX: 0,
		RegisterY: 0,

		Status: 0,

		ProgramCounter: 0,
		StackPointer:   StackReset,

		Bus: bus,
	}
}

func (c *CPU) Reset() {
	c.RegisterA = 0
	c.RegisterX = 0
	c.RegisterY = 0

	c.Status = 0

	c.StackPointer = StackReset

	// Load the reset vector from 0xFFFC
	c.ProgramCounter = c.ReadMemoryU16(0xFFFC)
}

func (c *CPU) Step() bool {
	if c.IsFlagSet(BreakCommand) {
		common.Log.Debug("Break command flag set, halting execution")
		return false
	}

	opCodeByte := c.ReadMemory(c.ProgramCounter)
	programCounterStart := c.ProgramCounter
	c.ProgramCounter += 1
	programCounterState := c.ProgramCounter

	opcode := Dispatch(opCodeByte)

	// Build instruction bytes string
	var instructionBytes []string
	instructionBytes = append(instructionBytes, fmt.Sprintf("%02X", opCodeByte))
	for i := uint16(1); i < uint16(opcode.Length); i++ {
		instructionBytes = append(instructionBytes, fmt.Sprintf("%02X", c.ReadMemory(programCounterStart+i)))
	}
	// Pad with spaces to always have 3 bytes worth of space
	for len(instructionBytes) < 3 {
		instructionBytes = append(instructionBytes, "  ")
	}

	// Store operand address before executing
	var address uint16
	var value uint8

	if opcode.AddressingMode != AddressingModeImplicit &&
		opcode.AddressingMode != AddressingModeAccumulator &&
		opcode.AddressingMode != AddressingModeRelative {
		address = c.GetOpAddress(opcode.AddressingMode)
		value = c.ReadMemory(address)
	} else if opcode.AddressingMode == AddressingModeRelative {
		// Branch instructions
		jump := int8(c.ReadMemory(c.ProgramCounter))
		address = c.ProgramCounter + 1 + uint16(jump)
	}

	// Build the instruction string with proper formatting
	var instruction string
	switch opcode.AddressingMode {
	case AddressingModeImplicit:
		instruction = opcode.Name
	case AddressingModeAccumulator:
		instruction = opcode.Name + " A"
	case AddressingModeImmediate:
		instruction = fmt.Sprintf("%s #$%02X", opcode.Name, value)
	case AddressingModeZeroPage:
		instruction = fmt.Sprintf("%s $%02X = %02X", opcode.Name, address, value)
	case AddressingModeZeroPageX:
		instruction = fmt.Sprintf("%s $%02X,X", opcode.Name, address)
	case AddressingModeZeroPageY:
		instruction = fmt.Sprintf("%s $%02X,Y", opcode.Name, address)
	case AddressingModeAbsolute:
		instruction = fmt.Sprintf("%s $%04X", opcode.Name, address)
	case AddressingModeAbsoluteX:
		instruction = fmt.Sprintf("%s $%04X,X", opcode.Name, address)
	case AddressingModeAbsoluteY:
		instruction = fmt.Sprintf("%s $%04X,Y", opcode.Name, address)
	case AddressingModeIndirect:
		instruction = fmt.Sprintf("%s ($%04X)", opcode.Name, address)
	case AddressingModeIndirectX:
		instruction = fmt.Sprintf("%s ($%02X,X)", opcode.Name, address)
	case AddressingModeIndirectY:
		instruction = fmt.Sprintf("%s ($%02X),Y", opcode.Name, address)
	case AddressingModeRelative:
		instruction = fmt.Sprintf("%s $%04X", opcode.Name, address)
	}

	// Pad instruction to 32 characters
	for len(instruction) < 32 {
		instruction += " "
	}

	// Format the log message
	common.Log.Debugf("%04X  %s  %s A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:  0,  0 CYC:0",
		programCounterStart,
		strings.Join(instructionBytes, " "),
		instruction,
		c.RegisterA,
		c.RegisterX,
		c.RegisterY,
		c.Status,
		c.StackPointer)

	// Execute the instruction
	opcode.Handler(c, opcode.AddressingMode)

	// Only increment PC after execution if it wasn't modified by the instruction
	if c.ProgramCounter == programCounterState {
		c.ProgramCounter += uint16(opcode.Length - 1)
	}

	return true
}

func (c *CPU) Run() {
	for {
		if !c.Step() {
			break
		}
	}
}

func (c *CPU) RunWithCallback(callback func(*CPU)) {
	for {
		if !c.Step() {
			break
		}

		callback(c)
	}
}

func (c *CPU) GetOpAddress(addressingMode uint) uint16 {
	switch addressingMode {
	case AddressingModeImmediate:
		return c.ProgramCounter
	case AddressingModeZeroPage:
		return uint16(c.ReadMemory(c.ProgramCounter))
	case AddressingModeZeroPageX:
		pos := c.ReadMemory(c.ProgramCounter)
		return uint16(pos + c.RegisterX)
	case AddressingModeZeroPageY:
		pos := c.ReadMemory(c.ProgramCounter)
		return uint16(pos + c.RegisterY)
	case AddressingModeAbsolute:
		return c.ReadMemoryU16(c.ProgramCounter)
	case AddressingModeAbsoluteX:
		pos := c.ReadMemoryU16(c.ProgramCounter)
		return pos + uint16(c.RegisterX)
	case AddressingModeAbsoluteY:
		pos := c.ReadMemoryU16(c.ProgramCounter)
		return pos + uint16(c.RegisterY)
	case AddressingModeIndirectX:
		base := c.ReadMemory(c.ProgramCounter)
		pos := uint16(base + c.RegisterX)
		lo := c.ReadMemory(pos)
		hi := c.ReadMemory(pos + 1)
		return uint16(hi)<<8 | uint16(lo)
	case AddressingModeIndirectY:
		base := c.ReadMemory(c.ProgramCounter)
		lo := c.ReadMemory(uint16(base))
		hi := c.ReadMemory(uint16(base + 1))
		addr := uint16(hi)<<8 | uint16(lo)
		return addr + uint16(c.RegisterY)
	case AddressingModeIndirect:
		base := c.ReadMemory(c.ProgramCounter)
		lo := c.ReadMemory(uint16(base))
		hi := c.ReadMemory(uint16(base + 1))
		return uint16(hi)<<8 | uint16(lo)
	default:
		panic("Unknown addressing mode: " + strconv.Itoa(int(addressingMode)))
	}
}

func (c *CPU) ReadMemory(address uint16) uint8 {
	return c.Bus.ReadMemory(address)
}

func (c *CPU) WriteMemory(address uint16, value uint8) {
	c.Bus.WriteMemory(address, value)
}

func (c *CPU) ReadMemoryU16(address uint16) uint16 {
	lo := uint16(c.ReadMemory(address))
	hi := uint16(c.ReadMemory(address + 1))
	return (hi << 8) | lo
}

func (c *CPU) WriteMemoryU16(address uint16, value uint16) {
	lo := uint8(value & 0xFF)
	high := uint8(value >> 8)
	c.WriteMemory(address, lo)
	c.WriteMemory(address+1, high)
}

func (c *CPU) PushStack(value uint8) {
	c.WriteMemory(StackStart+uint16(c.StackPointer), value)
	c.StackPointer--
}

func (c *CPU) PopStack() uint8 {
	c.StackPointer++
	return c.ReadMemory(StackStart + uint16(c.StackPointer))
}

func (c *CPU) PushStackU16(value uint16) {
	c.PushStack(uint8(value >> 8))   // Push high byte
	c.PushStack(uint8(value & 0xFF)) // Push low byte
}

func (c *CPU) PopStackU16() uint16 {
	lo := uint16(c.PopStack())
	hi := uint16(c.PopStack())
	return (hi << 8) | lo
}

func (c *CPU) SetZeroAndNegativeFlags(value uint8) {
	c.SetFlag(ZeroFlag, value == 0)
	c.SetFlag(NegativeFlag, value&0x80 != 0)
}

// AddToRegisterA adds a value to the accumulator and updates the flags accordingly.
func (c *CPU) AddToRegisterA(value uint8) {
	sum := uint16(c.RegisterA) + uint16(value)
	if c.IsFlagSet(CarryFlag) {
		sum++
	}

	c.SetFlag(CarryFlag, sum > 0xFF)

	var result = uint8(sum)
	if (value^result)&(result^c.RegisterA)&0x80 != 0 {
		c.SetFlag(OverflowFlag, true)
	} else {
		c.SetFlag(OverflowFlag, false)
	}

	c.RegisterA = uint8(sum)
}

func (c *CPU) BranchIf(condition bool) {
	if condition {
		jump := int8(c.ReadMemory(c.ProgramCounter))
		jumpAddr := c.ProgramCounter + 1 + uint16(jump)
		c.ProgramCounter = jumpAddr
	}
}

// Helper function to get human-readable addressing mode string
func getAddressingModeString(mode uint, address uint16) string {
	switch mode {
	case AddressingModeImmediate:
		return fmt.Sprintf("#$%02X", address)
	case AddressingModeZeroPage:
		return fmt.Sprintf("$%02X", address)
	case AddressingModeZeroPageX:
		return fmt.Sprintf("$%02X,X", address)
	case AddressingModeZeroPageY:
		return fmt.Sprintf("$%02X,Y", address)
	case AddressingModeAbsolute:
		return fmt.Sprintf("$%04X", address)
	case AddressingModeAbsoluteX:
		return fmt.Sprintf("$%04X,X", address)
	case AddressingModeAbsoluteY:
		return fmt.Sprintf("$%04X,Y", address)
	case AddressingModeIndirect:
		return fmt.Sprintf("($%04X)", address)
	case AddressingModeIndirectX:
		return fmt.Sprintf("($%02X,X)", address)
	case AddressingModeIndirectY:
		return fmt.Sprintf("($%02X),Y", address)
	default:
		return fmt.Sprintf("$%04X", address)
	}
}
