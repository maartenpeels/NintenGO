package main

import (
	"NintenGo/internal/cpu"
	_ "NintenGo/internal/opcodes" // Ugly, but cannot call from cpu package
	"fmt"
	"github.com/AllenDang/giu"
	"image"
	"image/color"
	"math/rand/v2"
)

var Cpu *cpu.CPU
var Rand = rand.NewPCG(42, 1024)

func main() {
	// The following program is a simple snake game
	// Memory layout:
	// 0xFE: Random number
	// 0xFF: Last pressed key
	// 0x0200-0x0600: Screen buffer
	// 0x0600-...: Snake code
	program := []byte{
		0x20, 0x06, 0x06, 0x20, 0x38, 0x06, 0x20, 0x0d, 0x06, 0x20, 0x2a, 0x06, 0x60, 0xa9, 0x02,
		0x85, 0x02, 0xa9, 0x04, 0x85, 0x03, 0xa9, 0x11, 0x85, 0x10, 0xa9, 0x10, 0x85, 0x12, 0xa9,
		0x0f, 0x85, 0x14, 0xa9, 0x04, 0x85, 0x11, 0x85, 0x13, 0x85, 0x15, 0x60, 0xa5, 0xfe, 0x85,
		0x00, 0xa5, 0xfe, 0x29, 0x03, 0x18, 0x69, 0x02, 0x85, 0x01, 0x60, 0x20, 0x4d, 0x06, 0x20,
		0x8d, 0x06, 0x20, 0xc3, 0x06, 0x20, 0x19, 0x07, 0x20, 0x20, 0x07, 0x20, 0x2d, 0x07, 0x4c,
		0x38, 0x06, 0xa5, 0xff, 0xc9, 0x77, 0xf0, 0x0d, 0xc9, 0x64, 0xf0, 0x14, 0xc9, 0x73, 0xf0,
		0x1b, 0xc9, 0x61, 0xf0, 0x22, 0x60, 0xa9, 0x04, 0x24, 0x02, 0xd0, 0x26, 0xa9, 0x01, 0x85,
		0x02, 0x60, 0xa9, 0x08, 0x24, 0x02, 0xd0, 0x1b, 0xa9, 0x02, 0x85, 0x02, 0x60, 0xa9, 0x01,
		0x24, 0x02, 0xd0, 0x10, 0xa9, 0x04, 0x85, 0x02, 0x60, 0xa9, 0x02, 0x24, 0x02, 0xd0, 0x05,
		0xa9, 0x08, 0x85, 0x02, 0x60, 0x60, 0x20, 0x94, 0x06, 0x20, 0xa8, 0x06, 0x60, 0xa5, 0x00,
		0xc5, 0x10, 0xd0, 0x0d, 0xa5, 0x01, 0xc5, 0x11, 0xd0, 0x07, 0xe6, 0x03, 0xe6, 0x03, 0x20,
		0x2a, 0x06, 0x60, 0xa2, 0x02, 0xb5, 0x10, 0xc5, 0x10, 0xd0, 0x06, 0xb5, 0x11, 0xc5, 0x11,
		0xf0, 0x09, 0xe8, 0xe8, 0xe4, 0x03, 0xf0, 0x06, 0x4c, 0xaa, 0x06, 0x4c, 0x35, 0x07, 0x60,
		0xa6, 0x03, 0xca, 0x8a, 0xb5, 0x10, 0x95, 0x12, 0xca, 0x10, 0xf9, 0xa5, 0x02, 0x4a, 0xb0,
		0x09, 0x4a, 0xb0, 0x19, 0x4a, 0xb0, 0x1f, 0x4a, 0xb0, 0x2f, 0xa5, 0x10, 0x38, 0xe9, 0x20,
		0x85, 0x10, 0x90, 0x01, 0x60, 0xc6, 0x11, 0xa9, 0x01, 0xc5, 0x11, 0xf0, 0x28, 0x60, 0xe6,
		0x10, 0xa9, 0x1f, 0x24, 0x10, 0xf0, 0x1f, 0x60, 0xa5, 0x10, 0x18, 0x69, 0x20, 0x85, 0x10,
		0xb0, 0x01, 0x60, 0xe6, 0x11, 0xa9, 0x06, 0xc5, 0x11, 0xf0, 0x0c, 0x60, 0xc6, 0x10, 0xa5,
		0x10, 0x29, 0x1f, 0xc9, 0x1f, 0xf0, 0x01, 0x60, 0x4c, 0x35, 0x07, 0xa0, 0x00, 0xa5, 0xfe,
		0x91, 0x00, 0x60, 0xa6, 0x03, 0xa9, 0x00, 0x81, 0x10, 0xa2, 0x00, 0xa9, 0x01, 0x81, 0x10,
		0x60, 0xa6, 0xff, 0xea, 0xea, 0xca, 0xd0, 0xfb, 0x60}

	Cpu = cpu.New(0x0600, program)

	wnd := giu.NewMasterWindow("NES Emulator", 1400, 1000, 0)
	wnd.Run(renderEmulatorWindow)
}

// TODO: This is ugly as hell, refactor.
func handleInput(c *cpu.CPU) {
	if giu.IsKeyDown(giu.KeyW) {
		c.Memory[0xFF] = 0x77
	} else if giu.IsKeyDown(giu.KeyS) {
		c.Memory[0xFF] = 0x73
	} else if giu.IsKeyDown(giu.KeyA) {
		c.Memory[0xFF] = 0x61
	} else if giu.IsKeyDown(giu.KeyD) {
		c.Memory[0xFF] = 0x64
	} else {
		// Do nothing
	}
}

func mapColor(byteValue uint8) color.RGBA {
	switch byteValue {
	case 0:
		return color.RGBA{A: 255} // BLACK
	case 1:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // WHITE
	case 2, 9:
		return color.RGBA{R: 128, G: 128, B: 128, A: 255} // GREY
	case 3, 10:
		return color.RGBA{R: 255, A: 255} // RED
	case 4, 11:
		return color.RGBA{G: 255, A: 255} // GREEN
	case 5, 12:
		return color.RGBA{B: 255, A: 255} // BLUE
	case 6, 13:
		return color.RGBA{R: 255, B: 255, A: 255} // MAGENTA
	case 7, 14:
		return color.RGBA{R: 255, G: 255, A: 255} // YELLOW
	default:
		return color.RGBA{G: 255, B: 255, A: 255} // CYAN
	}
}

func renderEmulatorWindow() {
	type FlagState struct {
		N bool
		V bool
		B bool
		D bool
		I bool
		Z bool
		C bool
	}
	var state = FlagState{
		Cpu.IsFlagSet(cpu.NegativeFlag),
		Cpu.IsFlagSet(cpu.OverflowFlag),
		Cpu.IsFlagSet(cpu.BreakCommand),
		Cpu.IsFlagSet(cpu.DecimalMode),
		Cpu.IsFlagSet(cpu.InterruptDisable),
		Cpu.IsFlagSet(cpu.ZeroFlag),
		Cpu.IsFlagSet(cpu.CarryFlag),
	}

	var splitPos float32 = 528

	giu.SingleWindow().Layout(
		giu.SplitLayout(giu.DirectionVertical, &splitPos,
			// Left column (display)
			giu.Layout{
				giu.Column(
					giu.Layout{
						// Display
						giu.Custom(func() {
							canvas := giu.GetCanvas()
							pos := giu.GetCursorScreenPos()

							// Draw the display content
							for i := 0x0200; i < 0x0600; i++ {
								x := (i - 0x0200) % 32
								y := (i - 0x0200) / 32

								colorIdx := Cpu.ReadMemory(uint16(i))
								rgba := mapColor(colorIdx)

								// Calculate pixel positions with proper scaling
								pMin := image.Pt(x*8, y*8)
								pMax := image.Pt(pMin.X+8, pMin.Y+8)

								// Draw the pixel
								canvas.AddRectFilled(pos.Add(pMin), pos.Add(pMax), rgba, 0, 0)
							}
						}),

						// Add a fixed-size container around the display
						giu.Dummy(256, 256),
					},

					// Controls
					giu.Row(
						giu.Button("Run").OnClick(func() {
							Cpu.RunWithCallback(func(c *cpu.CPU) {
								handleInput(c)

								// Set random number
								c.WriteMemory(0xFE, uint8(Rand.Uint64()))

								giu.Update()
							})
						}),

						giu.Button("Step").OnClick(func() {
							_ = Cpu.Step()
						}),

						giu.Button("Reset").OnClick(func() {
							Cpu.Reset()
						}),
					),
				),
			},

			// Right column (debug)
			giu.Layout{
				// CPU Registers
				giu.Label("Registers"),
				giu.Child().Size(giu.Auto, 100).Layout(
					giu.Row(
						giu.Label(fmt.Sprintf("A: $%02X", Cpu.RegisterA)),
						giu.Label(fmt.Sprintf("X: $%02X", Cpu.RegisterX)),
						giu.Label(fmt.Sprintf("Y: $%02X", Cpu.RegisterY)),
					),
					giu.Row(
						giu.Label(fmt.Sprintf("PC: $%04X", Cpu.ProgramCounter)),
						giu.Label(fmt.Sprintf("SP: $%02X", Cpu.StackPointer)),
					),
				),

				// Status Flags
				giu.Label("Status Flags"),
				giu.Child().Size(giu.Auto, 300).Layout(
					giu.Checkbox("N (negative)", &state.N),
					giu.Checkbox("V (overflow)", &state.V),
					giu.Checkbox("B (break command)", &state.B),
					giu.Checkbox("D (decimal mode)", &state.D),
					giu.Checkbox("I (interrupt disable)", &state.I),
					giu.Checkbox("Z (zero)", &state.Z),
					giu.Checkbox("C (carry)", &state.C),
				),

				// RAM Viewer
				giu.Label("RAM Viewer"),
				giu.Child().Layout(
					giu.Custom(func() {
						for addr := 0; addr < len(Cpu.Memory); addr += 16 {
							// Address
							rowStr := fmt.Sprintf("%04X:", addr)

							// Hex values
							for i := 0; i < 16; i++ {
								if addr+i < len(Cpu.Memory) {
									rowStr += fmt.Sprintf(" %02X", Cpu.Memory[addr+i])
								}
							}

							// ASCII representation
							rowStr += "  "
							for i := 0; i < 16; i++ {
								if addr+i < len(Cpu.Memory) {
									b := Cpu.Memory[addr+i]
									if b >= 32 && b <= 126 {
										rowStr += string(b)
									} else {
										rowStr += "."
									}
								}
							}

							giu.Label(rowStr).Build()
						}
					}),
				),
			},
		),
	)
}
