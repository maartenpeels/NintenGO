package cpu

type OpcodeHandler func(*CPU, uint)

type Opcode struct {
	Name           string
	AddressingMode uint
	Length         uint
	Handler        OpcodeHandler
}

var opcodeMap = map[uint8]Opcode{}

func RegisterOpcode(opcode uint8, handler OpcodeHandler, addressingMode uint, length uint, name string) {
	opcodeMap[opcode] = Opcode{Name: name, AddressingMode: addressingMode, Length: length, Handler: handler}
}

func Dispatch(opcode uint8) (Opcode, bool) {
	_opcode, exists := opcodeMap[opcode]
	return _opcode, exists
}
