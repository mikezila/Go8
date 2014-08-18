Go8
===

A simple Chip8 emulator in Go

This is a port of my Sharp8 emualtor.  It uses OpenGL for display and GLFW3 for input/windowmaking.  I'm still a total novice at OpenGL, so a large part of the OpenGL init/drawing was lifted wholesale from other projects I found.  It's not really ready yet, but it decodes all opcodes and some are impliemnted.

Building
--------

This is only tested on Mac OS X, so your milage may vary elsewhere.  You'll need the go-gl and glfw3 go packages installed.  In order to get them installed you'll need to install glfw3 and glew manually.  There aren't any pre-built binaries that I could find for OS X, but compiling it isn't too hard.

Once all that is done, just go build and the magic will happen.  Right now the rom that will be run is hard coded as a path in the source, but you should be quick enough to change that if you like, it's in the go8.go file.

License
-------

Fully public domain.  No thanks or mention are required.  Take it, copy it, have fun with it.
