#version 330 core

layout(location = 0) in vec2 aPos;
uniform vec2 uViewport;
void main() {
    vec2 ndc = (aPos / uViewport) * 2.0 - 1.0;
    gl_Position = vec4(ndc * vec2(1, -1), 0.0, 1.0)
}
