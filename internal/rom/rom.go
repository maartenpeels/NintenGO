package rom

import "fmt"

type Mirroring int

const (
	NesTag         = "NES\x1a"
	PrgRomPageSize = 16384
	ChrRomPageSize = 8192
)

const (
	Vertical Mirroring = iota
	Horizontal
	FourScreen
)

type Rom struct {
	PrgRom          []byte
	ChrRom          []byte
	Mapper          uint8
	ScreenMirroring Mirroring
}

func NewRom(raw []byte) (*Rom, error) {
	if string(raw[:4]) != NesTag {
		return nil, fmt.Errorf("invalid NES file")
	}

	mapper := (raw[7] & 0b1111_0000) | (raw[6] >> 4)
	inesVersion := (raw[7] >> 2) & 0b11
	if inesVersion != 0 {
		return nil, fmt.Errorf("unsupported iNES version: %d", inesVersion)
	}

	fourScreen := (raw[6] & 0b1000) != 0
	verticalMirroring := (raw[6] & 0b1) != 0

	var mirroring Mirroring
	if fourScreen {
		mirroring = FourScreen
	} else if verticalMirroring {
		mirroring = Vertical
	} else {
		mirroring = Horizontal
	}

	prgRomSize := int(raw[4]) * PrgRomPageSize
	chrRomSize := int(raw[5]) * ChrRomPageSize

	skipTrainer := (raw[6] & 0b100) != 0
	prgRomStart := 16
	if skipTrainer {
		prgRomStart += 512
	}
	chrRomStart := prgRomStart + prgRomSize

	if len(raw) < chrRomStart+chrRomSize {
		return nil, fmt.Errorf("invalid ROM size")
	}

	return &Rom{
		PrgRom:          raw[prgRomStart : prgRomStart+prgRomSize],
		ChrRom:          raw[chrRomStart : chrRomStart+chrRomSize],
		Mapper:          mapper,
		ScreenMirroring: mirroring,
	}, nil
}
