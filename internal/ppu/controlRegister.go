package ppu

type ControlRegister struct {
	value uint8
}

const (
	NameTable1            uint8 = 1 << 0 // Bit 0
	NameTable2            uint8 = 1 << 1 // Bit 1
	VramAddrInc           uint8 = 1 << 2 // Bit 2
	SpriteTableAddr       uint8 = 1 << 3 // Bit 3
	BackgroundPatternAddr uint8 = 1 << 4 // Bit 4
	SpriteSize            uint8 = 1 << 5 // Bit 5
	MasterSlaveSelect     uint8 = 1 << 6 // Bit 6
	GenerateNMI           uint8 = 1 << 7 // Bit 7
)

func NewControlRegister() *ControlRegister {
	return &ControlRegister{
		value: 0,
	}
}

func (c *ControlRegister) Contains(flag uint8) bool {
	return c.value&flag != 0
}

func (c *ControlRegister) Set(flag uint8) {
	c.value |= flag
}

func (c *ControlRegister) Clear(flag uint8) {
	c.value &^= flag
}

func (c *ControlRegister) VramAddrInc() uint8 {
	if c.Contains(VramAddrInc) {
		return 32
	}
	return 1
}

func (c *ControlRegister) Update(value uint8) {
	c.value = value
}

func (c *ControlRegister) Value() uint8 {
	return c.value
}
