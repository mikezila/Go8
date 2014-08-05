Go8
===

A simple Chip8 emulaator in Go

This is a port of my Sharp8 emualtor.  It uses OpenGL for display and GLFW for input/windowmaking.  I'm still a total novice at OpenGL, so a large part of the OpenGL init/drawing was lifted wholesale from other projects I found.  It's not really ready yet, but it decodes all opcodes and some are impliemnted.

If you want to build this you'll need go-gl, glfw3, and glew installed.  Once I have it to a stable point I'll make running it less of a pain I hope.  So far it's much, much faster than the C# version, we'll see how it stacks up once all of the opcodes are in.

As with all my projects, feel free to fork/copy/steal, whatever.  This project is in the public domain.
