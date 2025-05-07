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

	opCodeByte := c.ReadMemory(c.ProgramCounter)
	programCounterStart := c.ProgramCounter
	c.ProgramCounter += 1
	programCounterState := c.ProgramCounter

	opcode := Dispatch(opCodeByte)
	c.Bus.Tick(opcode.Cycles)

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
		// Get the base address before X is added
		baseAddr := c.ReadMemory(c.ProgramCounter)
		instruction = fmt.Sprintf("%s $%02X,X @ %02X = %02X", opcode.Name, baseAddr, address, value)
	case AddressingModeZeroPageY:
		// Get the base address before Y is added
		baseAddr := c.ReadMemory(c.ProgramCounter)
		instruction = fmt.Sprintf("%s $%02X,Y @ %02X = %02X", opcode.Name, baseAddr, address, value)
	case AddressingModeAbsolute:
		instruction = fmt.Sprintf("%s $%04X = %02X", opcode.Name, address, value)
	case AddressingModeAbsoluteX:
		// Get the base address before X is added
		baseAddr := c.ReadMemoryU16(c.ProgramCounter)
		instruction = fmt.Sprintf("%s $%04X,X @ %04X = %02X", opcode.Name, baseAddr, address, value)
	case AddressingModeAbsoluteY:
		// Get the base address before Y is added
		baseAddr := c.ReadMemoryU16(c.ProgramCounter)
		instruction = fmt.Sprintf("%s $%04X,Y @ %04X = %02X", opcode.Name, baseAddr, address, value)
	case AddressingModeIndirect:
		// Special case for JMP indirect
		if opcode.Name == "JMP" {
			// For JMP indirect, we need to show the target address
			base := c.ReadMemoryU16(c.ProgramCounter)
			// The 6502 has a bug where it doesn't correctly fetch the high byte if the low byte is at the end of a page
			// We need to emulate this behavior
			lo := c.ReadMemory(base)
			hi := c.ReadMemory((base & 0xFF00) | ((base + 1) & 0xFF)) // Page boundary wrap
			targetAddr := uint16(hi)<<8 | uint16(lo)

			instruction = fmt.Sprintf("%s ($%04X) = %04X", opcode.Name, base, targetAddr)
		} else {
			instruction = fmt.Sprintf("%s ($%04X)", opcode.Name, address)
		}
	case AddressingModeIndirectX:
		// For IndirectX, we need to show the zero page address and the final address
		zeroPageAddr := c.ReadMemory(c.ProgramCounter)
		effectiveAddr := uint8(zeroPageAddr + c.RegisterX)
		lo := c.ReadMemory(uint16(effectiveAddr))
		hi := c.ReadMemory(uint16(uint8(effectiveAddr + 1)))
		finalAddr := uint16(hi)<<8 | uint16(lo)

		// Format: LDA ($80,X) @ 85 = 0200 = 5A
		instruction = fmt.Sprintf("%s ($%02X,X) @ %02X = %04X = %02X",
			opcode.Name, zeroPageAddr, effectiveAddr, finalAddr, value)
	case AddressingModeIndirectY:
		// For IndirectY, we need to show the zero page address, the address it points to, and the final value
		zeroPageAddr := c.ReadMemory(c.ProgramCounter)
		lo := c.ReadMemory(uint16(zeroPageAddr))
		hi := c.ReadMemory(uint16(uint8(zeroPageAddr + 1)))
		baseAddr := uint16(hi)<<8 | uint16(lo)
		finalAddr := baseAddr + uint16(c.RegisterY)

		// Format: LDA ($89),Y = 0300 @ 0300 = 89
		instruction = fmt.Sprintf("%s ($%02X),Y = %04X @ %04X = %02X",
			opcode.Name, zeroPageAddr, baseAddr, finalAddr, value)
	case AddressingModeRelative:
		instruction = fmt.Sprintf("%s $%04X", opcode.Name, address)
	}

	// Pad instruction to 31 characters
	padLength := 31
	if strings.HasPrefix(instruction, "*") {
		padLength = 32
	}
	for len(instruction) < padLength {
		instruction += " "
	}

	// Format the log message
	// For unofficial opcodes, we need to adjust the spacing between instruction bytes and instruction name
	instructionBytesStr := strings.Join(instructionBytes, " ")

	// Adjust spacing between instruction bytes and instruction name
	// For unofficial opcodes (*NOP), we need 1 space instead of the usual 2
	spacing := "  "
	if strings.HasPrefix(instruction, "*") {
		spacing = " "
	}

	log := fmt.Sprintf("%04X  %s%s%s A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:  0,  0 CYC:0",
		programCounterStart,
		instructionBytesStr,
		spacing,
		instruction,
		c.RegisterA,
		c.RegisterX,
		c.RegisterY,
		c.Status.value,
		c.StackPointer)
	common.Log.Info(log)

	// // Check if line matches the expected log
	// expectedLogLine := c.ExpectedLogs[c.CurrentIndex]
	// expectedLogLineStart := strings.SplitN(expectedLogLine, "PPU:", 2)[0]
	// currentLogStart := strings.SplitN(log, "PPU:", 2)[0]
	// if expectedLogLineStart != currentLogStart {
	// 	common.Log.Errorf("Log mismatch at index %d:\nexpected:\n\t%s\ngot:\n\t%s", c.CurrentIndex, expectedLogLineStart, currentLogStart)
	// 	c.Status.Set(BreakCommand)
	// 	return false
	// }
	// c.CurrentIndex++

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
