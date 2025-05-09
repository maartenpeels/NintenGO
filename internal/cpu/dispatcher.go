package cpu

type OpcodeHandler func(*CPU, uint)

type Opcode struct {
	Name           string
	AddressingMode uint
	Length         uint
	Cycles         uint
	Handler        OpcodeHandler
}

var opcodeMap = map[uint8]Opcode{}

func RegisterOpcode(opcode uint8, handler OpcodeHandler, addressingMode uint, length uint, cycles uint, name string) {
	opcodeMap[opcode] = Opcode{Name: name, AddressingMode: addressingMode, Length: length, Cycles: cycles, Handler: handler}
}

func Dispatch(opcode uint8) Opcode {
	_opcode, exists := opcodeMap[opcode]

	if !exists {
		return opcodeMap[0xEA] // NOP
	}

	return _opcode
}
