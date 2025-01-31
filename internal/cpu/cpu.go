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

	Memory [0xFFFF]uint8
}

type Error struct {
	PC      uint16
	OpCode  uint8
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("CPU Error at PC:%04X (OpCode:%02X): %s", e.PC, e.OpCode, e.Message)
}

func New(start uint16, program []byte) *CPU {
	c := &CPU{
		RegisterA: 0,
		RegisterX: 0,
		RegisterY: 0,

		Status: 0,

		ProgramCounter: 0,
		StackPointer:   StackReset,

		Memory: [0xFFFF]uint8{},
	}
	c.loadProgram(start, program)
	c.Reset()
	return c
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

	// Format status flags for better readability
	statusFlags := fmt.Sprintf("%c%c%c%c%c%c%c",
		map[bool]rune{true: 'N', false: '-'}[c.IsFlagSet(NegativeFlag)],
		map[bool]rune{true: 'V', false: '-'}[c.IsFlagSet(OverflowFlag)],
		map[bool]rune{true: 'B', false: '-'}[c.IsFlagSet(BreakCommand)],
		map[bool]rune{true: 'D', false: '-'}[c.IsFlagSet(DecimalMode)],
		map[bool]rune{true: 'I', false: '-'}[c.IsFlagSet(InterruptDisable)],
		map[bool]rune{true: 'Z', false: '-'}[c.IsFlagSet(ZeroFlag)],
		map[bool]rune{true: 'C', false: '-'}[c.IsFlagSet(CarryFlag)])

	// Build the log message based on addressing mode
	switch opcode.AddressingMode {
	case AddressingModeImplicit, AddressingModeAccumulator:
		common.Log.Debugf("PC:%04X | OP:%02X %-3s | A:%02X X:%02X Y:%02X | SP:%02X | %s",
			programCounterStart,
			opCodeByte,
			opcode.Name,
			c.RegisterA,
			c.RegisterX,
			c.RegisterY,
			c.StackPointer,
			statusFlags)
	default:
		addrMode := getAddressingModeString(opcode.AddressingMode, address)
		common.Log.Debugf("PC:%04X | OP:%02X %-3s | %s | VAL:%02X | A:%02X X:%02X Y:%02X | SP:%02X | %s",
			programCounterStart,
			opCodeByte,
			opcode.Name,
			addrMode,
			value,
			c.RegisterA,
			c.RegisterX,
			c.RegisterY,
			c.StackPointer,
			statusFlags)
	}

	// Execute the instruction
	opcode.Handler(c, opcode.AddressingMode)

	// Only increment PC after execution if it wasn't modified by the instruction
	if c.ProgramCounter == programCounterState {
		c.ProgramCounter += uint16(opcode.Length - 1)
		common.Log.Debugf("PC:%04X | Next instruction", c.ProgramCounter)

	} else {
		common.Log.Debugf("PC:%04X | Jumped to %04X", programCounterStart, c.ProgramCounter)
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
	return c.Memory[address]
}

func (c *CPU) WriteMemory(address uint16, value uint8) {
	c.Memory[address] = value
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

func (c *CPU) loadProgram(start uint16, program []byte) {
	common.Log.Debugf("Loading program of %d bytes at address %04X", len(program), start)
	for i, b := range program {
		c.Memory[start+uint16(i)] = b
	}

	// Verify the load
	common.Log.Debug("Program load verification:")
	for i := uint16(0); i < uint16(len(program)); i += 16 {
		var bytes []string
		for j := uint16(0); j < 16 && i+j < uint16(len(program)); j++ {
			bytes = append(bytes, fmt.Sprintf("%02X", c.Memory[start+i+j]))
		}
		common.Log.Debugf("%04X: %s", start+i, strings.Join(bytes, " "))
	}

	c.WriteMemoryU16(0xFFFC, start)
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
