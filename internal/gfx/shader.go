package gfx

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Shader struct {
	handle uint32
}

type Program struct {
	handle  uint32
	shaders []*Shader
}

func (s *Shader) Delete() {
	gl.DeleteShader(s.handle)
}

func (p *Program) Delete() {
	for _, shader := range p.shaders {
		shader.Delete()
	}
	gl.DeleteProgram(p.handle)
}

func (p *Program) Attach(shaders ...*Shader) {
	for _, shader := range shaders {
		gl.AttachShader(p.handle, shader.handle)
		p.shaders = append(p.shaders, shader)
	}
}

func (p *Program) Use() {
	gl.UseProgram(p.handle)
}

func (p *Program) Link() error {
	gl.LinkProgram(p.handle)
	return getGlError(p.handle, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog, "PROGRAM::LINKING_FAILURE")
}

func NewProgram(shaders ...*Shader) (*Program, error) {
	program := &Program{handle: gl.CreateProgram()}
	program.Attach(shaders...)

	if err := program.Link(); err != nil {
		return nil, err
	}

	return program, nil
}

func NewShader(src string, sType uint32) (*Shader, error) {
	return compileShader(src, sType)
}

func NewShaderFromFile(file string, sType uint32) (*Shader, error) {
	src, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return compileShader(string(src), sType)
}

func compileShader(src string, sType uint32) (*Shader, error) {
	handle := gl.CreateShader(sType)
	glSrc, freeFn := gl.Strs(src + "\x00")
	defer freeFn()

	gl.ShaderSource(handle, 1, glSrc, nil)
	gl.CompileShader(handle)

	err := getGlError(handle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "SHADER::COMPILE_FAILURE::")
	if err != nil {
		return nil, err
	}

	return &Shader{handle: handle}, nil
}

type getObjIv func(uint32, uint32, *int32)
type getObjInfoLog func(uint32, int32, *int32, *uint8)

func getGlError(glHandle uint32, checkTrueParam uint32, getObjIvFn getObjIv, getObjInfoLogFn getObjInfoLog, failMsg string) error {
	var success int32
	getObjIvFn(glHandle, checkTrueParam, &success)

	if success == gl.FALSE {
		var logLength int32
		getObjIvFn(glHandle, gl.INFO_LOG_LENGTH, &logLength)

		log := gl.Str(strings.Repeat("\x00", int(logLength)))
		getObjInfoLogFn(glHandle, logLength, nil, log)

		return fmt.Errorf("%s: %s", failMsg, gl.GoStr(log))
	}

	return nil
}
