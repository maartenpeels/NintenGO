package cpu

type StatusRegister struct {
	value uint8
}

const (
	CarryFlag        uint8 = 1 << 0 // Bit 0
	ZeroFlag         uint8 = 1 << 1 // Bit 1
	InterruptDisable uint8 = 1 << 2 // Bit 2
	DecimalMode      uint8 = 1 << 3 // Bit 3
	BreakCommand     uint8 = 1 << 4 // Bit 4
	NotUsedFlag      uint8 = 1 << 5 // Bit 5
	OverflowFlag     uint8 = 1 << 6 // Bit 6
	NegativeFlag     uint8 = 1 << 7 // Bit 7
)

func NewStatusRegister() *StatusRegister {
	return &StatusRegister{
		value: NotUsedFlag | InterruptDisable,
	}
}

func (s *StatusRegister) Contains(flag uint8) bool {
	return s.value&flag != 0
}

func (s *StatusRegister) Set(flag uint8) {
	s.value |= flag
}

func (s *StatusRegister) SetBool(flag uint8, value bool) {
	if value {
		s.Set(flag)
	} else {
		s.Clear(flag)
	}
}

func (s *StatusRegister) Clear(flag uint8) {
	s.value &^= flag
}

func (s *StatusRegister) Update(value uint8) {
	s.value = value
}

func (s *StatusRegister) Value() uint8 {
	return s.value
}
