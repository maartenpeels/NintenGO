package opcodes

import (
	"NintenGo/internal/cpu"
	_ "NintenGo/internal/opcodes" // Ugly, but cannot call from cpu package
	"testing"
)

func TestTAX(t *testing.T) {
	// Load 5 into A and call TAX
	program := []byte{0xa9, 0x05, 0xaa}
	c := cpu.New(program)
	c.Run()

	if c.RegisterX != 5 {
		t.Error("Register X should be 5")
	}
}
