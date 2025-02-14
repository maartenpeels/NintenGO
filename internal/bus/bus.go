package bus

import (
	"NintenGo/internal/common"
	"NintenGo/internal/ppu"
	"NintenGo/internal/rom"
	"fmt"
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
	cycles  uint64

	prgRom *rom.Rom
	ppu    *ppu.PPU
}

func NewBus(rom *rom.Rom) *BUS {
	return &BUS{
		cpuVram: [VramSize]byte{},
		cycles:  0,

		prgRom: rom,
		ppu:    ppu.NewPPU(rom),
	}
}

func (b *BUS) Tick(addCycles uint) {
	b.cycles += uint64(addCycles)
	b.ppu.Tick(addCycles)
}

func (b *BUS) ReadMemory(address uint16) uint8 {
	switch {
	case address >= RamStart && address <= RamMirrorsEnd:
		mirrorDownAddress := address & 0b00000111_11111111
		return b.cpuVram[mirrorDownAddress]

	case address == 0x2000 || address == 0x2001 || address == 0x2003 ||
		address == 0x2005 || address == 0x2006 || address == 0x4014:
		common.Log.Errorf("Attempt to read from write-only PPU address 0x%04X", address)
		return 0 // TODO: What should we return here? Should we panic?

	case address == 0x2007:
		return b.ppu.ReadData()

	case address >= PpuRegistersStart && address <= PpuRegistersMirrorsEnd:
		mirrorDownAddr := address & 0b00100000_00000111
		return b.ReadMemory(mirrorDownAddr)

	case address >= 0x8000:
		return b.ReadPrgRom(address)

	default:
		common.Log.Debugf("Ignoring mem access at 0x%04X", address)
		return 0
	}
}

func (b *BUS) WriteMemory(address uint16, value uint8) {
	switch {
	case address >= RamStart && address <= RamMirrorsEnd:
		mirrorDownAddress := address & 0b00000111_11111111
		b.cpuVram[mirrorDownAddress] = value

	case address == 0x2000:
		b.ppu.WriteToCTRL(value)

	case address == 0x2006:
		b.ppu.WriteToPPUAddr(value)

	case address == 0x2007:
		b.ppu.WriteToData(value)

	case address >= 0x2008 && address <= PpuRegistersMirrorsEnd:
		mirrorDownAddr := address & 0b00100000_00000111
		b.WriteMemory(mirrorDownAddr, value)

	case address >= 0x8000:
		panic(fmt.Sprintf("Attempt to write to Cartridge ROM space: 0x%04X", address))

	default:
		common.Log.Debugf("Ignoring mem write-access at 0x%04X", address)
	}
}

func (b *BUS) ReadPrgRom(address uint16) uint8 {
	newAddr := address - 0x8000
	if len(b.prgRom.PrgRom) == 0x4000 && newAddr >= 0x4000 {
		// If we have a 16KB PRG ROM, mirror the lower bank
		newAddr %= 0x4000
	}
	return b.prgRom.PrgRom[newAddr]
}
