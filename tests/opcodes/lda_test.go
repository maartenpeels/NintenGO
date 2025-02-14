package opcodes

import (
	"NintenGo/internal/cpu"
	_ "NintenGo/internal/opcodes" // Ugly, but cannot call from cpu package
	"testing"
)

func TestLDA(t *testing.T) {
	program := []byte{0xa9, 0x05, 0x00}
	c := cpu.New(0x800, program)
	c.Run()

	if c.RegisterA != 0x05 {
		t.Error("register A != 0x05")
	}

	if c.Status.Contains(cpu.ZeroFlag) != false {
		t.Error("Zero flag should be false")
	}

	if c.Status.Contains(cpu.NegativeFlag) != false {
		t.Error("Negative flag should be false")
	}
}

func TestLDAZeroPage(t *testing.T) {
	program := []byte{0xa5, 0x10, 0x00}
	c := cpu.New(0x800, program)

	// Prepare memory
	c.WriteMemory(0x10, 0x55)

	c.Run()

	if c.RegisterA != 0x55 {
		t.Error("register A != 0x55")
	}
}

func TestLDAZero(t *testing.T) {
	program := []byte{0xa9, 0x00, 0x00}
	c := cpu.New(0x800, program)
	c.Run()

	if c.RegisterA != 0x00 {
		t.Error("register A != 0x00")
	}

	if c.Status.Contains(cpu.ZeroFlag) != true {
		t.Error("Zero flag should be true")
	}

	if c.Status.Contains(cpu.NegativeFlag) != false {
		t.Error("Negative flag should be false")
	}
}

func TestLDAZeroNegative(t *testing.T) {
	program := []byte{0xa9, 0x80, 0x00}
	c := cpu.New(0x800, program)
	c.Run()

	if c.RegisterA != 0x80 {
		t.Error("register A != 0x80")
	}

	if c.Status.Contains(cpu.ZeroFlag) != false {
		t.Error("Zero flag should be false")
	}

	if c.Status.Contains(cpu.NegativeFlag) != true {
		t.Error("Negative flag should be true")
	}
}
