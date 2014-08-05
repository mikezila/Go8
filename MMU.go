//MMU.go - Chip8 memory controller
package main

import (
	"fmt"
	"io/ioutil"
)

const MAX_MEMORY uint16 = 0xFFF
const REGISTER_COUNT int = 16
const PC_START uint16 = 0x200
const FONTS_BEGIN uint16 = 0x000
const STACK_DEPTH int = 16

type Chip8Memory struct {
	Memory     [MAX_MEMORY]byte
	Buffer     C8FrameBuffer
	Registers  [REGISTER_COUNT]byte
	Stack      [STACK_DEPTH]uint16
	SP         byte
	Flag       byte
	PC         uint16
	Indexer    uint16
	SoundTimer byte
	DelayTimer byte
	Crashed    bool
}

func (c *Chip8Memory) Reset() {
	// Zero memory
	for i := 0; i < len(c.Memory); i++ {
		c.Memory[i] = 0
	}

	// Zero registers
	for i := 0; i < len(c.Registers); i++ {
		c.Registers[i] = 0
	}

	// Zero stack
	for i := 0; i < STACK_DEPTH; i++ {
		c.Stack[i] = 0
	}

	// Zero other registers
	c.SP = 0
	c.Flag = 0
	c.Indexer = 0
	c.SoundTimer = 0
	c.DelayTimer = 0
	c.Crashed = false

	// Reset program counter to starting position
	c.PC = PC_START
}

func (c *Chip8Memory) PushStack(address uint16) {
	c.Stack[c.SP] = address
	c.SP++
}

func (c *Chip8Memory) PopStack() (address uint16) {
	c.SP--
	address = c.Stack[c.SP]
	return
}

func (c Chip8Memory) ReadOpCode(index uint16) (opcode uint16) {
	if index%2 != 0 {
		fmt.Println("Warning : Reading opcode at odd address.\n")
	}
	opcode = uint16(c.Memory[index]) << 8
	opcode |= uint16(c.Memory[index+1])
	return
}

func (c Chip8Memory) NextOpCode() (opcode uint16) {
	opcode = c.ReadOpCode(c.PC)
	return
}

func (c *Chip8Memory) LoadRom(romName string) {
	rom, err := ioutil.ReadFile(romName)
	if err != nil {
		panic(err)
	}
	copy(c.Memory[PC_START:], rom[:])
}

func (c *Chip8Memory) WriteByte(data byte, index uint16) {
	c.Memory[index] = data
}

// This packs the 0-F glyphs into memory.
func (c *Chip8Memory) PackFonts() {
	fonts := [80]byte{0xf0, 0x90, 0x90, 0x90, 0xf0, 0x20, 0x60, 0x20, 0x20, 0x70, 0xf0, 0x10, 0xf0, 0x80, 0xf0, 0xf0, 0x10, 0xf0, 0x10, 0xf0, 0x90, 0x90, 0xf0, 0x10, 0x10, 0xf0, 0x80, 0xf0, 0x10, 0xf0, 0xf0, 0x80, 0xf0, 0x90, 0xf0, 0xf0, 0x10, 0x20, 0x40, 0x40, 0xf0, 0x90, 0xf0, 0x90, 0xf0, 0xf0, 0x90, 0xf0, 0x10, 0xf0, 0xf0, 0x90, 0xf0, 0x90, 0x90, 0xe0, 0x90, 0xe0, 0x90, 0xe0, 0xf0, 0x80, 0x80, 0x80, 0xf0, 0xe0, 0x90, 0x90, 0x90, 0xe0, 0xf0, 0x80, 0xf0, 0x80, 0xf0, 0xf0, 0x80, 0xf0, 0x80, 0x80}
	copy(c.Memory[0:], fonts[:])
}

func (c Chip8Memory) DumpMemory() {
	fmt.Println(c.Memory)
}
