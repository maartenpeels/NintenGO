package opcodes

import (
	"NintenGo/internal/cpu"
	_ "NintenGo/internal/opcodes" // Ugly, but cannot call from cpu package
	"testing"
)

func TestINX(t *testing.T) {
	program := []byte{0xe8, 0xe8}
	c := cpu.New(program)

	c.Run()

	if c.RegisterX != 2 {
		t.Error("Register X should be 2")
	}
}

func TestINXOverflow(t *testing.T) {
	program := []byte{0xe8, 0xe8}
	c := cpu.New(program)

	c.RegisterX = 0xff
	c.Run()

	if c.RegisterX != 1 {
		t.Error("Register X should be 1")
	}
}
