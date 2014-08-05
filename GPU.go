// GPU.go
package main

import (
	"math/rand"
	"time"
)

const BUFFER_WIDTH int = 64
const BUFFER_HEIGHT int = 32
const DISPLAY_SCALE int = 10
const PIXEL_COUNT int = BUFFER_HEIGHT * BUFFER_WIDTH

type C8Color struct {
	//rgb bytes
	Color [3]byte
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
	rand.Seed(int64(time.Now().Nanosecond()))
	for i := range c.Buffer {
		randPix := rand.Intn(2)
		if randPix == 0 {
			c.Buffer[i] = C8Color{[3]byte{255, 255, 255}}
		} else {
			c.Buffer[i] = C8Color{[3]byte{0, 0, 0}}
		}
	}
}

func (c *C8FrameBuffer) GetPixel(x, y int) (color C8Color) {
	color = c.Buffer[(y*BUFFER_WIDTH)+x]
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

func (c *C8FrameBuffer) SetPixel(x, y int, color C8Color) {
	c.Buffer[(y*BUFFER_WIDTH)+x] = color
}
