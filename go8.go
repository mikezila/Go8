// go8.go
package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glu"
	"runtime"
)

const VERSION string = "v0.0"
const TITLE_STRING string = "Go8 " + VERSION
const DEBUG bool = false
const DEBUG_VERBOSE bool = false

func main() {
	fmt.Println("Go8 ", VERSION)

	chip8 := Chip8Memory{}
	chip8.Reset()
	chip8.PackFonts()
	chip8.LoadRom("roms/IBM.c8")

	// lock glfw/gl calls to a single thread
	runtime.LockOSThread()

	glfw.Init()
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, 0)

	window, err := glfw.CreateWindow(BUFFER_WIDTH*DISPLAY_SCALE, BUFFER_HEIGHT*DISPLAY_SCALE, TITLE_STRING, nil, nil)
	if err != nil {
		panic(err)
	}

	defer window.Destroy()

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Init()
	gl.Disable(gl.DEPTH_TEST)
	gl.PointSize(float32(DISPLAY_SCALE))
	gl.ClearColor(255, 255, 0, 0)
	gl.Viewport(0, 0, BUFFER_WIDTH*DISPLAY_SCALE, BUFFER_HEIGHT*DISPLAY_SCALE)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(-0.5, float64(BUFFER_WIDTH)-0.5, float64(BUFFER_HEIGHT)-0.5, -0.5, -1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	screen := make(chan *C8FrameBuffer)
	go chip8.execute(screen)

	for !window.ShouldClose() {

		buffer := <-screen

		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Begin(gl.POINTS)
		for y := 0; y < BUFFER_HEIGHT; y++ {
			for x := 0; x < BUFFER_WIDTH; x++ {
				pixel := buffer.GetPixel(x, y)
				gl.Color3ubv(&pixel.Color)
				gl.Vertex2i(x, y)
			}
		}
		gl.End()

		window.SwapBuffers()
		glfw.PollEvents()

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}

func checkGLerror() {
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		string, _ := glu.ErrorString(glerr)
		panic(string)
	}
}
