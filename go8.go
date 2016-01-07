// go8.go
package main

import (
	"fmt"
<<<<<<< HEAD
	"math"
	"runtime"

	gl "github.com/go-gl/gl/all-core/gl"
=======
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
>>>>>>> origin/master
	glfw "github.com/go-gl/glfw/v3.1/glfw"
)

const VERSION string = "v0.6"
const TITLE_STRING string = "Go8 " + VERSION
const DEBUG bool = false
const DEBUG_VERBOSE bool = false

func main() {
	fmt.Println("Go8 ", VERSION)

	chip8 := Chip8Memory{}
	chip8.Reset()
	chip8.PackFonts()
	chip8.LoadRom("roms/maze.c8")

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
	gl.Viewport(0, 0, int32(BUFFER_WIDTH*DISPLAY_SCALE), int32(BUFFER_HEIGHT*DISPLAY_SCALE))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(-0.5, float64(BUFFER_WIDTH)-0.5, float64(BUFFER_HEIGHT)-0.5, -0.5, -1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	for !window.ShouldClose() {
		chip8.execute()

		if chip8.Buffer.Dirty {
			gl.Clear(gl.COLOR_BUFFER_BIT)
			gl.Begin(gl.POINTS)
			for y := 0; y < BUFFER_HEIGHT; y++ {
				for x := 0; x < BUFFER_WIDTH; x++ {
					pixel := chip8.Buffer.GetPixel(x, y)
					if pixel {
						gl.Color4i(math.MaxInt32, math.MaxInt32, math.MaxInt32, math.MaxInt32)
					} else {
						gl.Color4i(0, 0, 0, math.MaxInt32)
					}
					gl.Vertex2i(int32(x), int32(y))
				}
			}
			gl.End()
			window.SwapBuffers()
		}
		chip8.Buffer.Dirty = false

		glfw.PollEvents()

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}
