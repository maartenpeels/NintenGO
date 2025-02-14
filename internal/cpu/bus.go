package cpu

import (
	"NintenGo/internal/common"
	"NintenGo/internal/rom"
)

const (
	VramSize               = 2048
	RamStart               = 0x0000
	RamMirrorsEnd          = 0x1FFF
	PpuRegistersStart      = 0x2000
	PpuRegistersMirrorsEnd = 0x3FFF
)

// BUS represents the system bus
//
// Memory Map:
//
//	_______________ $10000  _______________
//
// | PRG-ROM       |       |               |
// | Upper Bank    |       |               |
// |_ _ _ _ _ _ _ _| $C000 | PRG-ROM       |
// | PRG-ROM       |       |               |
// | Lower Bank    |       |               |
// |_______________| $8000 |_______________|
// | SRAM          |       | SRAM          |
// |_______________| $6000 |_______________|
// | Expansion ROM |       | Expansion ROM |
// |_______________| $4020 |_______________|
// | I/O Registers |       |               |
// |_ _ _ _ _ _ _ _| $4000 |               |
// | Mirrors       |       | I/O Registers |
// | $2000-$2007   |       |               |
// |_ _ _ _ _ _ _ _| $2008 |               |
// | I/O Registers |       |               |
// |_______________| $2000 |_______________|
// | Mirrors       |       |               |
// | $0000-$07FF   |       |               |
// |_ _ _ _ _ _ _ _| $0800 |               |
// | RAM           |       | RAM           |
// |_ _ _ _ _ _ _ _| $0200 |               |
// | Stack         |       |               |
// |_ _ _ _ _ _ _ _| $0100 |               |
// | Zero Page     |       |               |
// |_______________| $0000 |_______________|
type BUS struct {
	cpuVram [VramSize]byte // 2KB of VRAM
	rom     *rom.Rom
}

func NewBus(rom *rom.Rom) *BUS {
	return &BUS{
		cpuVram: [VramSize]byte{},
		rom:     rom,
	}
}

func (b *BUS) ReadMemory(address uint16) uint8 {
	switch {
	case address >= RamStart && address <= RamMirrorsEnd:
		mirrorDownAddress := address & 0b00000111_11111111
		return b.cpuVram[mirrorDownAddress]
	case address >= PpuRegistersStart && address <= PpuRegistersMirrorsEnd:
		mirrorDownAddr := address & 0b00100000_00000111
		_ = mirrorDownAddr // Placeholder for future PPU support
		common.Log.Info("PPU is not supported yet")
		return 0
	case address >= 0x4000 && address <= 0x4017:
		// APU and I/O registers
		common.Log.Debug("APU read not implemented yet")
		return 0
	case address >= 0x4018 && address <= 0x401F:
		// APU and I/O functionality that is normally disabled
		common.Log.Debug("APU read not implemented yet")
		return 0
	case address >= 0x4020 && address < 0x8000:
		// Expansion ROM - return 0 for now
		common.Log.Debug("Expansion ROM read not implemented yet")
		return 0
	case address >= 0x8000:
		return b.ReadPrgRom(address)
	default:
		common.Log.Errorf("Incorrect memory access at %04X", address)
		return 0
	}
}

func (b *BUS) ReadPrgRom(address uint16) uint8 {
	newAddr := address - 0x8000
	if len(b.rom.PrgRom) == 0x4000 && newAddr >= 0x4000 {
		// If we have a 16KB PRG ROM, mirror the lower bank
		newAddr %= 0x4000
	}
	return b.rom.PrgRom[newAddr]
}

func (b *BUS) WriteMemory(address uint16, value uint8) {
	switch {
	case address >= RamStart && address <= RamMirrorsEnd:
		mirrorDownAddress := address & 0b00000111_11111111
		b.cpuVram[mirrorDownAddress] = value
	case address >= PpuRegistersStart && address <= PpuRegistersMirrorsEnd:
		mirrorDownAddr := address & 0b00100000_00000111
		_ = mirrorDownAddr // Placeholder for future PPU support
		common.Log.Info("PPU is not supported yet")
	case address >= 0x4000 && address <= 0x4017:
		// APU and I/O registers
		common.Log.Debug("APU write not implemented yet")
	case address >= 0x4018 && address <= 0x401F:
		// APU and I/O functionality that is normally disabled
		common.Log.Debug("APU write not implemented yet")
	case address >= 0x4020 && address < 0x8000:
		// Expansion ROM
		common.Log.Debug("Expansion ROM write not implemented yet")
	case address >= 0x8000:
		common.Log.Errorf("Cannot write to ROM space at %04X", address)
	default:
		common.Log.Errorf("Incorrect memory access at %04X", address)
	}
}
