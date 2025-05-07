package cpu

import (
	"fmt"
	"strings"
)

// NonReadableAddr lists memory addresses that might be write-only or have side effects if read during tracing.
// Customize this list based on your emulator's memory map and I/O behavior.
var NonReadableAddr = []uint16{
	0x2001, // PPUMASK
	0x2002, // PPUSTATUS
	0x2003, // OAMADDR
	0x2004, // OAMDATA
	0x2005, // PPUSCROLL
	0x2006, // PPUADDR
	0x2007, // PPUDATA
	0x4014, // OAMDMA
	0x4016, // JOY1
	0x4017, // JOY2
}

func isNonReadable(addr uint16) bool {
	for _, nonReadable := range NonReadableAddr {
		if addr == nonReadable {
			return true
		}
	}
	return false
}

// Trace generates a debug trace string for the current CPU state and instruction.
// It mimics common NES emulator log formats.
func Trace(cpu *CPU) string {
	if opcodeMap == nil {
		return fmt.Sprintf("%04X: Error - OpCodeMap not initialized. CPU: A:%02X X:%02X Y:%02X P:%02X SP:%02X", cpu.ProgramCounter, cpu.RegisterA, cpu.RegisterX, cpu.RegisterY, cpu.Status.value, cpu.StackPointer)
	}

	opCodeByte := cpu.ReadMemory(cpu.ProgramCounter)
	op, ok := opcodeMap[opCodeByte]
	if !ok {
		return fmt.Sprintf("%04X: %02X       Unknown Opcode. CPU: A:%02X X:%02X Y:%02X P:%02X SP:%02X", cpu.ProgramCounter, opCodeByte, cpu.RegisterA, cpu.RegisterX, cpu.RegisterY, cpu.Status.value, cpu.StackPointer)
	}

	pc := cpu.ProgramCounter // Capture PC at the start of the instruction
	var hexDump []string
	hexDump = append(hexDump, fmt.Sprintf("%02X", opCodeByte))

	var operandStr string
	var memAddr uint16   // Effective address for operand
	var storedValue byte // Value at memAddr, if applicable

	// Save and temporarily adjust cpu.ProgramCounter for GetOpAddress, then restore.
	originalCPUProgramCounter := cpu.ProgramCounter
	cpu.ProgramCounter = pc + 1 // Set to operand address for GetOpAddress

	switch op.AddressingMode {
	case AddressingModeImmediate:
		memAddr = 0 // No memory address in the traditional sense
		storedValue = 0
	case AddressingModeImplicit, AddressingModeAccumulator:
		memAddr = 0
		storedValue = 0
	case AddressingModeRelative:
		// Relative mode target address calculation needs the byte at pc+1
		// cpu.ProgramCounter is already pc+1 here, so ReadMemory will fetch the offset.
		offset := int8(cpu.ReadMemory(cpu.ProgramCounter))
		memAddr = (pc + 2) + uint16(offset) // Target address for display
		storedValue = 0                     // No value stored at a branch target in this context
	default:
		// For other modes, use GetOpAddress. cpu.ProgramCounter is already set to pc+1.
		memAddr = cpu.GetOpAddress(uint(op.AddressingMode)) // op.Mode needs to be convertible to uint
		if !isNonReadable(memAddr) {
			storedValue = cpu.ReadMemory(memAddr) // Read from actual bus
		} else {
			storedValue = 0 // Or some indicator for non-readable
		}
	}

	cpu.ProgramCounter = originalCPUProgramCounter // Restore cpu.ProgramCounter to its original value (pc)

	switch op.Length {
	case 1:
		switch opCodeByte { // Some 1-byte opcodes directly reference the accumulator
		case 0x0A, 0x4A, 0x2A, 0x6A: // ASL A, LSR A, ROL A, ROR A
			operandStr = "A"
		default:
			// Implicit or Accumulator, typically no operand string unless it's 'A'
			if op.AddressingMode == AddressingModeAccumulator {
				operandStr = "A"
			} else {
				operandStr = ""
			}
		}
	case 2:
		arg1 := cpu.ReadMemory(pc + 1)
		hexDump = append(hexDump, fmt.Sprintf("%02X", arg1))

		switch op.AddressingMode {
		case AddressingModeImmediate:
			operandStr = fmt.Sprintf("#$%02X", arg1)
		case AddressingModeZeroPage:
			operandStr = fmt.Sprintf("$%02X = %02X", memAddr&0xFF, storedValue)
		case AddressingModeZeroPageX:
			operandStr = fmt.Sprintf("$%02X,X @ %02X = %02X", arg1, memAddr&0xFF, storedValue)
		case AddressingModeZeroPageY:
			operandStr = fmt.Sprintf("$%02X,Y @ %02X = %02X", arg1, memAddr&0xFF, storedValue)
		case AddressingModeIndirectX:
			operandStr = fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X", arg1, (uint16(arg1)+uint16(cpu.RegisterX))&0xFF, memAddr, storedValue)
		case AddressingModeIndirectY:
			lo := cpu.ReadMemory(uint16(arg1))
			hi := cpu.ReadMemory(uint16(uint8(arg1 + 1)))
			intermediatePtr := uint16(hi)<<8 | uint16(lo)
			operandStr = fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X", arg1, intermediatePtr, memAddr, storedValue)
		case AddressingModeRelative:
			operandStr = fmt.Sprintf("$%04X", memAddr)
		default:
			operandStr = fmt.Sprintf("$%02X (mode %d?) A:%02X X:%02X Y:%02X P:%02X SP:%02X CYC:%d", arg1, op.AddressingMode, cpu.RegisterA, cpu.RegisterX, cpu.RegisterY, cpu.Status.Value(), cpu.StackPointer, cpu.Bus.GetCycles())
		}

	case 3:
		arg1 := cpu.ReadMemory(pc + 1) // low byte
		arg2 := cpu.ReadMemory(pc + 2) // high byte
		hexDump = append(hexDump, fmt.Sprintf("%02X", arg1), fmt.Sprintf("%02X", arg2))
		rawAddressArg := uint16(arg2)<<8 | uint16(arg1) // The $HHLL part

		switch op.AddressingMode {
		case AddressingModeAbsolute:
			if op.Name == "JMP" || op.Name == "JSR" { // JMP and JSR don't read a value from the address for this log line
				operandStr = fmt.Sprintf("$%04X", memAddr)
			} else {
				operandStr = fmt.Sprintf("$%04X = %02X", memAddr, storedValue)
			}
		case AddressingModeAbsoluteX:
			operandStr = fmt.Sprintf("$%04X,X @ %04X = %02X", rawAddressArg, memAddr, storedValue)
		case AddressingModeAbsoluteY:
			operandStr = fmt.Sprintf("$%04X,Y @ %04X = %02X", rawAddressArg, memAddr, storedValue)
		case AddressingModeIndirect: // JMP ($HHLL)
			operandStr = fmt.Sprintf("($%04X) = %04X", rawAddressArg, memAddr)
		default:
			operandStr = fmt.Sprintf("$%04X (mode %d?) A:%02X X:%02X Y:%02X P:%02X SP:%02X CYC:%d", rawAddressArg, op.AddressingMode, cpu.RegisterA, cpu.RegisterX, cpu.RegisterY, cpu.Status.Value(), cpu.StackPointer, cpu.Bus.GetCycles())
		}
	default:
		operandStr = "" // Should not happen for valid opcodes
	}

	hexStr := strings.Join(hexDump, " ")
	// Pad hexStr to always be 8 chars for alignment (e.g., "00      " or "00 00   " or "00 00 00")
	fixedHexLen := 8
	if len(hexStr) < fixedHexLen {
		hexStr += strings.Repeat(" ", fixedHexLen-len(hexStr))
	}

	// Combine mnemonic and operand
	mnemonicAndOperandStr := op.Name
	if operandStr != "" {
		// Check for accumulator mode for specific opcodes, where operand is 'A'
		if op.AddressingMode == AddressingModeAccumulator ||
			(op.Length == 1 && (opCodeByte == 0x0A || opCodeByte == 0x4A || opCodeByte == 0x2A || opCodeByte == 0x6A)) {
			mnemonicAndOperandStr += " A"
		} else if !strings.HasPrefix(operandStr, "#S") { // #S is my temp hack for immediate to avoid space, use #$
			mnemonicAndOperandStr += " " + operandStr
		} else {
			mnemonicAndOperandStr += strings.Replace(operandStr, "#S", "#$", 1)
		}
	}

	// Pad mnemonicAndOperandStr to 31 characters
	instructionPartLen := 31
	if strings.HasPrefix(op.Name, "*") {
		instructionPartLen = 32
	}
	if len(mnemonicAndOperandStr) < instructionPartLen {
		mnemonicAndOperandStr += strings.Repeat(" ", instructionPartLen-len(mnemonicAndOperandStr))
	}

	// Format: PC  Hex   Mnemonic+Operands              A:XX X:XX Y:XX P:XX SP:XX PPU:XXX,XXX CYC:CCC
	// Example:C000  4C F5 C5  JMP $C5F5                   A:00 X:00 Y:00 P:24 SP:FD PPU:  0, 21 CYC:7
	padding := "  "
	if strings.HasPrefix(op.Name, "*") {
		padding = " "
	}
	finalLog := fmt.Sprintf("%04X  %s%s%s A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d",
		pc,
		hexStr,
		padding,
		mnemonicAndOperandStr,
		cpu.RegisterA,
		cpu.RegisterX,
		cpu.RegisterY,
		cpu.Status.Value(),
		cpu.StackPointer,
		cpu.Bus.GetPPUScanline(),
		cpu.Bus.GetPPUCycle(),
		cpu.Bus.GetCycles())

	return finalLog
}
