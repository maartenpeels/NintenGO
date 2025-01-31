package programs

import (
	"NintenGo/internal/cpu"
	"testing"
)

func TestBug(t *testing.T) {
	program := []byte{0x8d, 0x06, 0x20, 0xc3}

	c := cpu.New(0x8000, program)
	c.Run()
}
