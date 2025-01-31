package opcodes

import (
	"NintenGo/internal/cpu"
	_ "NintenGo/internal/opcodes" // Ugly, but cannot call from cpu package
	"testing"
)

func TestBRK(t *testing.T) {
	program := []byte{0x00}
	c := cpu.New(0x800, program)

	c.Run()

	if c.IsFlagSet(cpu.BreakCommand) != true {
		t.Error("Break flag should be true")
	}
}
