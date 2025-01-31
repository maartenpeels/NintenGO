package programs

import (
	"NintenGo/internal/cpu"
	_ "NintenGo/internal/opcodes" // Ugly, but cannot call from cpu package
	"testing"
)

func TestIncrementXAfterTransfer(t *testing.T) {
	// LDA: Load 0xc0 into A
	// TAX: Transfer A to X
	// INX: Increment X
	// BRK: Break
	program := []byte{0xa9, 0xc0, 0xaa, 0xe8, 0x00}
	c := cpu.New(0x8000, program)
	c.Run()

	if c.RegisterX != 0xc1 {
		t.Error("Register X should be 0xc1")
	}
}
