// GPU.go
package main

import (
	"math/rand"
)

const BUFFER_WIDTH int = 64
const BUFFER_HEIGHT int = 32
const DISPLAY_SCALE int = 15
const PIXEL_COUNT int = BUFFER_HEIGHT * BUFFER_WIDTH

type C8Color struct {
	Color [3]byte
}

func (c *C8Color) Red() (red byte) {
	red = c.Color[0]
	return
}

func (c *C8Color) Green() (green byte) {
	green = c.Color[1]
	return
}

func (c *C8Color) Blue() (blue byte) {
	blue = c.Color[2]
	return
}

type C8FrameBuffer struct {
	Buffer [PIXEL_COUNT]C8Color
}

func (c *C8FrameBuffer) ClearScreen() {
	for i := range c.Buffer {
		c.Buffer[i] = C8Color{[3]byte{0, 0, 0}}
	}
}

func (c *C8FrameBuffer) RandomNoise() {
	for i := range c.Buffer {
		randPix := rand.Intn(2)
		if randPix == 0 {
			c.Buffer[i] = C8Color{[3]byte{255, 255, 255}}
		} else {
			c.Buffer[i] = C8Color{[3]byte{0, 0, 0}}
		}
	}
}

func (c *C8FrameBuffer) GetPixel(x, y int) (color *C8Color) {
	color = &c.Buffer[(y*BUFFER_WIDTH)+x]
	return
}

func (c *C8FrameBuffer) TestPixel(x, y int) (pixelLit bool) {
	pixel := c.GetPixel(x, y)
	if pixel.Color[0] > 0 || pixel.Color[1] > 0 || pixel.Color[2] > 0 {
		pixelLit = true
	} else {
		pixelLit = false
	}
	return
}

func (c *C8FrameBuffer) TurnPixelOff(x, y int) {
	c.SetPixel(x, y, C8Color{[3]byte{0, 0, 0}})
}

// When Chip8 sprites are drawn, trying to turn on a pixel that is already on sets the Vf
// register, which is used to check for collisions in games.  It also turns the pixel off.
// In this way, sprites can be erased by drawing them to the same location twice.
func (c *C8FrameBuffer) TurnPixelOn(x, y int) (collision bool) {
	if x >= BUFFER_WIDTH {
		x -= BUFFER_WIDTH
	}

	if y >= BUFFER_HEIGHT {
		y -= BUFFER_HEIGHT
	}
	collision = c.TestPixel(x, y)
	if collision {
		c.TurnPixelOff(x, y)
	} else {
		c.SetPixel(x, y, C8Color{[3]byte{255, 255, 255}})
	}
	return
}

func (c *C8FrameBuffer) SetPixel(x, y int, color C8Color) {
	c.Buffer[(y*BUFFER_WIDTH)+x] = color
}

func (c *Chip8Memory) DrawSprite(opcode uint16) {
	collision := false

	x := int(c.Registers[((opcode >> 8) & 0x000F)])
	y := int(c.Registers[((opcode >> 4) & 0x000F)])
	size := int((opcode & 0x000F))

	for i := 0; i < size; i++ {
		spriteLine := c.ReadByte(c.Indexer + uint16(i))
		for currentPixel := 0; currentPixel <= 8; currentPixel++ {
			if (spriteLine & (1 << uint(8-currentPixel))) > 0 {
				collision = c.Buffer.TurnPixelOn((x + currentPixel), y+i)
			}
		}
	}

	c.Flag = collision
}
