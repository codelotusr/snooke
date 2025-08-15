package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(1280, 720, "Snooke", nil, nil)
	if err != nil {
		log.Fatalln("create window:", err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("gl init", err)
	}

	// Match viewport to framebuffer size now and when resized
	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	// Enable vsync
	glfw.SwapInterval(1)

	gl.ClearColor(0.10, 0.10, 0.12, 1.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
