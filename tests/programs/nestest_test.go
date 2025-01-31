package programs

import (
	"NintenGo/internal/cpu"
	"io"
	"os"
	"testing"
)

func TestNESTest(t *testing.T) {
	// NESTest is a test ROM that runs a series of tests on the CPU
	dir, _ := os.Getwd()
	t.Logf("Current working directory: %s", dir)
	// Load the test ROM
	file, err := os.Open("../../roms/nestest.nes")
	if err != nil || file == nil {
		t.Error("Error opening test ROM, err: ", err)
		return
	}

	program, err := io.ReadAll(file)
	if err != nil {
		t.Error("Error reading test ROM")
		return
	}

	// Run the test ROM
	c := cpu.New(0x8000, program)
	c.Run()
}
