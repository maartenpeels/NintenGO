package opcodes

import (
	"NintenGo/internal/cpu"
	_ "NintenGo/internal/opcodes" // Ugly, but cannot call from cpu package
	"testing"
)

func TestAND(t *testing.T) {
	program := []byte{
		0x29, 0x0f, // AND #$0f
	}

	c := cpu.New(program)
	c.RegisterA = 0xFF

	c.Step()

	if c.RegisterA != 0x0F {
		t.Error("Register A should be 0x0F, got ", c.RegisterA)
	}
}

func TestANDIndirectX(t *testing.T) {
	program := []byte{
		0x21, 0x10, // AND #$10
		0x00, 0x20, // Low and high bytes of address
		0x0F, // Value at address
	}

	c := cpu.New(program)
	c.RegisterA = 0xFF
	c.RegisterX = 0x01

	c.Step()

	if c.RegisterA != 0x0F {
		t.Error("Register A should be 0x0F, got ", c.RegisterA)
	}
}
