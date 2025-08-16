package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  int = 800
	windowHeight int = 600
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Snooke", nil, nil)
	if err != nil {
		log.Fatalln("create window:", err)
	}
	defer window.Destroy()

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	x := (mode.Width - windowWidth) / 2
	y := (mode.Height - windowHeight) / 2

	window.SetPos(x, y)

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize glow:", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Match viewport to framebuffer size now and when resized
	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	// Enable vsync
	glfw.SwapInterval(1)

	// Background color
	gl.ClearColor(0.10, 0.10, 0.12, 1.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var vertexShader = `
#version 330 core

layout(location = 0) in vec2 aPos;
uniform vec2 uViewport;

void main() {
    vec2 ndc = (aPos / uViewport) * 2.0 - 1.0;
    gl_Position = vec4(ndc * vec2(1, -1), 0.0, 1.0);
}
` + "\x00"

var fragmentShader = `
#version 330 core

uniform vec3 uColor;
out vec4 FragColor;

void main() {
    FragColor = vec4(uColor, 1.0);
}
` + "\x00"
