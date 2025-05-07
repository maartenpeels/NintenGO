package cpu

import (
	"NintenGo/internal/bus"
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

	Status *StatusRegister

	ProgramCounter uint16
	StackPointer   uint8

	Bus *bus.BUS

	CurrentIndex uint
	ExpectedLogs []string
}

type Error struct {
	PC      uint16
	OpCode  uint8
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("CPU Error at PC:%04X (OpCode:%02X): %s", e.PC, e.OpCode, e.Message)
}

func New(bus *bus.BUS, expectedLogs []string) *CPU {
	return &CPU{
		RegisterA: 0,
		RegisterX: 0,
		RegisterY: 0,

		Status: NewStatusRegister(),

		ProgramCounter: 0,
		StackPointer:   StackReset,

		Bus: bus,

		CurrentIndex: 0,
		ExpectedLogs: expectedLogs,
	}
}

func (c *CPU) Reset() {
	c.RegisterA = 0
	c.RegisterX = 0
	c.RegisterY = 0

	c.Status = NewStatusRegister()
	c.StackPointer = StackReset

	// Load the reset vector from 0xFFFC
	c.ProgramCounter = c.ReadMemoryU16(0xFFFC)
}

func (c *CPU) Step() bool {
	if c.Status.Contains(BreakCommand) {
		common.Log.Debug("Break command flag set, halting execution")
		return false
	}

	// Capture PC for trace, before it's incremented past the current opcode
	pcForTrace := c.ProgramCounter

	opCodeByte := c.ReadMemory(c.ProgramCounter)
	// programCounterStart := c.ProgramCounter // pcForTrace serves this purpose now
	c.ProgramCounter += 1
	programCounterState := c.ProgramCounter

	opcode := Dispatch(opCodeByte)
	c.Bus.Tick(opcode.Cycles)

	// Generate the full trace log line using the new Trace function
	// We pass a temporary CPU state with the PC set to the beginning of the instruction
	// because Trace itself will manipulate and then restore cpu.ProgramCounter for GetOpAddress.
	tempCPUStateForTrace := *c
	tempCPUStateForTrace.ProgramCounter = pcForTrace // Ensure Trace sees the PC at the start of the opcode
	traceLog := Trace(&tempCPUStateForTrace)

	// Debug log with the new trace format
	common.Log.Debug(traceLog)

	// Check against expected logs if they exist
	if c.ExpectedLogs != nil && len(c.ExpectedLogs) > 0 && c.CurrentIndex < uint(len(c.ExpectedLogs)) {
		expected := c.ExpectedLogs[c.CurrentIndex]
		// Compare traceLog with expected. Need to ensure formats match or use a flexible comparison.
		// For now, let's assume exact match for simplicity, though this might need refinement
		// if there are subtle differences (e.g. extra spaces, case sensitivity) introduced by the new Trace func.

		// For now we ignore the PPU part at the end
		traceLog = strings.SplitN(traceLog, "PPU", 2)[0]
		expected = strings.SplitN(expected, "PPU", 2)[0]

		if traceLog != expected {
			errorMessage := fmt.Sprintf("\nTrace mismatch at PC %04X (Op: %02X, Index: %d):\nExpected: %s\nGot:      %s",
				pcForTrace, opCodeByte, c.CurrentIndex, expected, traceLog)
			common.Log.Error(errorMessage)
			c.Status.Set(BreakCommand)
			return false
		}
	}
	if c.ExpectedLogs != nil {
		c.CurrentIndex++
	}

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
		// X register wraps in zero page
		pos := uint8(base + c.RegisterX)
		lo := c.ReadMemory(uint16(pos))
		hi := c.ReadMemory(uint16(uint8(pos + 1)))
		return uint16(hi)<<8 | uint16(lo)
	case AddressingModeIndirectY:
		base := c.ReadMemory(c.ProgramCounter)
		// Zero page wrapping for the pointer
		lo := c.ReadMemory(uint16(base))
		hi := c.ReadMemory(uint16(uint8(base + 1)))
		addr := uint16(hi)<<8 | uint16(lo)
		return addr + uint16(c.RegisterY)
	case AddressingModeIndirect:
		base := c.ReadMemoryU16(c.ProgramCounter)
		// The 6502 has a hardware bug where it doesn't correctly fetch the high byte if the low byte is at the end of a page
		lo := c.ReadMemory(base)
		hi := c.ReadMemory((base & 0xFF00) | ((base + 1) & 0xFF)) // Page boundary wrap
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

	c.Status.SetBool(ZeroFlag, value == 0)
	c.Status.SetBool(NegativeFlag, value&0x80 != 0)
}

// AddToRegisterA adds a value to the accumulator and updates the flags accordingly.
func (c *CPU) AddToRegisterA(value uint8) {
	sum := uint16(c.RegisterA) + uint16(value)
	if c.Status.Contains(CarryFlag) {
		sum++
	}

	c.Status.SetBool(CarryFlag, sum > 0xFF)

	var result = uint8(sum)
	if (value^result)&(result^c.RegisterA)&0x80 != 0 {
		c.Status.Set(OverflowFlag)
	} else {
		c.Status.Clear(OverflowFlag)
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
