package ppu

import (
	"NintenGo/internal/common"
	"NintenGo/internal/rom"
	"fmt"
)

type PPU struct {
	chrRom       []byte     // CHR ROM
	paletteTable [32]byte   // 32 bytes of palette data
	vRam         [2048]byte // 2KB of VRAM
	oamData      [256]byte  // Sprite data

	mirroring rom.Mirroring   // Mirroring type
	addr      AddrRegister    // Address register
	ctrl      ControlRegister // Control register

	internalDataBuffer uint8

	scanline uint16
	cycles   uint64
}

func NewPPU(rom *rom.Rom) *PPU {
	return &PPU{
		chrRom:    rom.ChrRom,
		mirroring: rom.ScreenMirroring,
		addr:      *NewAddrRegister(),
		ctrl:      *NewControlRegister(),

		vRam:         [2048]byte{},
		oamData:      [256]byte{},
		paletteTable: [32]byte{},
	}
}

func (p *PPU) Tick(addCycles uint) bool {
	p.cycles += uint64(addCycles)
	if p.cycles >= 341 {
		p.cycles = p.cycles - 341
		p.scanline++

		if p.scanline == 241 {
			// Signal VBlank
			// TODO: Stuff and things
		}

		if p.scanline == 262 {
			p.scanline = 0
			// TODO: Reset VBlank
			return true
		}
	}

	return false
}

func (p *PPU) ReadData() uint8 {
	addr := p.addr.Get()
	p.IncrementVRamAddr()

	if addr <= 0x1FFF {
		result := p.internalDataBuffer
		p.internalDataBuffer = p.chrRom[addr]
		return result
	} else if addr <= 0x2FFF {
		result := p.internalDataBuffer
		p.internalDataBuffer = p.vRam[p.MirrorVramAddr(addr)]
		return result
	} else if addr <= 0x3EFF {
		common.Log.Warnf("Addr space 0x3000-0x3EFF is not expected to be used, requested addr: 0x%04X", addr)
	} else if addr <= 0x3FFF {
		return p.paletteTable[addr-0x3F00]
	} else {
		common.Log.Errorf("Invalid PPU read address: 0x%04X", addr)
	}
	return 0
}

func (p *PPU) WriteToData(value uint8) {
	addr := p.addr.Get()

	switch {
	case addr <= 0x1FFF:
		common.Log.Warnf("attempt to write to chr rom space 0x%04X", addr)

	case addr <= 0x2FFF:
		mirroredAddr := p.MirrorVramAddr(addr)
		p.vRam[mirroredAddr] = value

	case addr <= 0x3EFF:
		common.Log.Warnf("addr 0x%04X shouldn't be used in reality", addr)

	// Addresses $3F10/$3F14/$3F18/$3F1C are mirrors of $3F00/$3F04/$3F08/$3F0C
	case addr == 0x3F10 || addr == 0x3F14 || addr == 0x3F18 || addr == 0x3F1C:
		addMirror := addr - 0x10
		p.paletteTable[addMirror-0x3F00] = value

	case addr <= 0x3FFF:
		p.paletteTable[addr-0x3F00] = value

	default:
		panic(fmt.Sprintf("unexpected access to mirrored space 0x%04X", addr))
	}

	p.IncrementVRamAddr()
}

func (p *PPU) WriteToPPUAddr(value uint8) {
	p.addr.Update(value)
}

func (p *PPU) WriteToCTRL(value uint8) {
	p.ctrl.Update(value)
}

func (p *PPU) IncrementVRamAddr() {
	p.addr.Increment(p.ctrl.VramAddrInc())
}

func (p *PPU) MirrorVramAddr(address uint16) uint16 {
	// Mirror down 0x3000-0x3eff to 0x2000-0x2eff
	mirroredVRAM := p.addr.Get() & 0b10111111111111

	vRamIndex := mirroredVRAM - 0x2000
	nameTable := vRamIndex / 0x400

	switch {
	case p.mirroring == rom.Vertical && (nameTable == 2 || nameTable == 3):
		return vRamIndex - 0x800
	case p.mirroring == rom.Horizontal && nameTable == 2:
		return vRamIndex - 0x400
	case p.mirroring == rom.Horizontal && nameTable == 1:
		return vRamIndex - 0x400
	case p.mirroring == rom.Horizontal && nameTable == 3:
		return vRamIndex - 0x400
	default:
		return vRamIndex
	}
}
