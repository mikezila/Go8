// CPU.go - Chip8 CPU emulation functions
package main

import (
	"fmt"
	"math/rand"
)

func (c *Chip8Memory) execute(screen chan *C8FrameBuffer) {
	quit := false

	for !quit {
		updateScreen := true
		opcode := c.NextOpCode()

		var decodedOpcode string

		// Decode opcode
		switch {

		// Clear Screen
		case (opcode & 0x00F0) == 0x00E0:
			decodedOpcode = "Clear Screen"
			updateScreen = true
			c.Buffer.ClearScreen()

		// Return from subroutine
		case (opcode & 0x00FF) == 0x00EE:
			decodedOpcode = "Return from subroutine"
			c.PC = c.PopStack()

		// Direct Jump (no stack push)
		case (opcode & 0xF000) == 0x1000:
			decodedOpcode = "Jump"
			c.PC = opcode & 0x0FFF

		// Call subroutine
		case (opcode & 0xF000) == 0x2000:
			decodedOpcode = "Call subroutine"
			c.PushStack(c.PC)
			c.PC = opcode & 0x0FFF

		// Skip next op if Vx == kk : 3xkk
		case (opcode & 0xF000) == 0x3000:
			decodedOpcode = "Skip next if Vx == kk"
			if c.Registers[(opcode&0x0F00)>>8] == byte(opcode&0x00FF) {
				c.PC += 2
			}

		// Skip next op if Vx != kk : 4xkk
		case (opcode & 0xF000) == 0x4000:
			decodedOpcode = "Skip next if Vx != kk"
			if c.Registers[(opcode&0x0F00)>>8] != byte(opcode&0x00FF) {
				c.PC += 2
			}

		// Skip next op if Vx == Vy : 5xy0
		case (opcode & 0xF000) == 0x5000:
			decodedOpcode = "Skip next if Vx == Vy"
			if c.Registers[(opcode&0x0F00)>>8] == c.Registers[(opcode&0x00F0)>>4] {
				c.PC += 2
			}

		// Load Vx with kk : 6xkk
		case (opcode & 0xF000) == 0x6000:
			decodedOpcode = "Load Vx with kk"
			c.Registers[(opcode&0x0F00)>>8] = byte(opcode & 0x00FF)

		// Add kk to Vx : 7xkk
		case (opcode & 0xF000) == 0x7000:
			decodedOpcode = "Add kk to Vx"
			c.Registers[(opcode&0x0F00)>>8] += byte(opcode & 0x00FF)

		// Load Vx with Vy : 8xy0
		case (opcode & 0xF00F) == 0x8000:
			decodedOpcode = "Load Vx with Vy"
			c.Registers[(opcode&0x0F00)>>8] = c.Registers[(opcode&0x00F0)>>4]

		// Load Vx with Vx OR Vy
		case (opcode & 0xF00F) == 0x8001:
			decodedOpcode = "Load Vx with Vy OR Vx"
			c.Registers[(opcode&0x0F00)>>8] |= c.Registers[(opcode&0x00F0)>>4]

		// Load Vx with Vx AND Vy
		case (opcode & 0xF00F) == 0x8002:
			decodedOpcode = "Load Vx with Vx AND Vy"
			c.Registers[(opcode&0x0F00)>>8] &= c.Registers[(opcode&0x00F0)>>4]

		// Load Vx with Vx XOR Vy
		case (opcode & 0xF00F) == 0x8003:
			decodedOpcode = "Load Vx with Vx XOR Vy"
			c.Registers[(opcode&0x0F00)>>8] ^= c.Registers[(opcode&0x00F0)>>4]

		// Load Vx with Vx + Vy, Vf as carry
		case (opcode & 0xF00F) == 0x8004:
			decodedOpcode = "Load Vx with Vx + Vy, sets Vf as carry"
			c.Flag = (int(c.Registers[(opcode&0x0F00)>>8]) + int(c.Registers[(opcode&0x00F0)>>4])) > 255
			c.Registers[(opcode&0x0F00)>>8] += c.Registers[(opcode&0x00F0)>>4]

		// Load Vx with Vx - Vy, Vf set if no carry happened
		case (opcode & 0xF00F) == 0x8005:
			decodedOpcode = "Load Vx with Vx - Vy, sets Vf if no carry"

		// SHR Vx - If the least significant bit of Vx is set, then Vf is set, otherwise is cleared.
		// Vx is this divided by 2.
		case (opcode & 0xF00F) == 0x8006:
			decodedOpcode = "SHR Vx"

		// Load Vx with Vy - Vx, Vf set if no carry happened
		case (opcode & 0xF00F) == 0x8007:
			decodedOpcode = "Load Vy with Vy - Vx, sets Vf if no carry"

		// SHL Vx - If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0.
		// Then Vx is multiplied by 2.
		case (opcode & 0xF00F) == 0x800E:
			decodedOpcode = "SHL Vx"

		// Skip next instruction if Vx != Vy
		case (opcode & 0xF00F) == 0x9000:
			decodedOpcode = "Skip next opcode if Vx != Vy"
			if c.Registers[(opcode&0x0F00)>>8] != c.Registers[(opcode&0x00F0)>>4] {
				c.PC += 2
			}

		// Load Indexer with nnn : Annn
		case (opcode & 0xF000) == 0xA000:
			decodedOpcode = "Load Indexer"
			c.Indexer = (opcode & 0x0FFF)

		// Jump to nnn + V0 : Bnnn
		case (opcode & 0xF000) == 0xB000:
			decodedOpcode = "Jump to nnn + V0"

		// Generate random byte then AND it with kk, store result in Vx : Cxkk
		case (opcode & 0xF000) == 0xC000:
			decodedOpcode = "Store random byte in Vx, ANDed with kk"
			c.Registers[(opcode&0x0F00)>>8] = byte(rand.Int()) & byte(opcode&0x00FF)

		// Draw spite
		// The interpreter reads n bytes from memory, starting at the address stored in I. These bytes are then displayed as sprites on screen at coordinates (Vx, Vy). Sprites are XORed onto the existing screen. If this causes any pixels to be erased, VF is set to 1, otherwise it is set to 0. If the sprite is positioned so part of it is outside the coordinates of the display, it wraps around to the opposite side of the screen.
		// Dxyn
		case (opcode & 0xF000) == 0xD000:
			decodedOpcode = "Draw sprite."
			updateScreen = true
			c.DrawSprite(opcode)

		// Skip next if key with value Vx is pressed : Ex9E
		case (opcode & 0xF0FF) == 0xE09E:
			decodedOpcode = "Skip next if key Vx is pressed"

		// Skip next if key with value Vx is not pressed : ExA1
		case (opcode & 0xF0FF) == 0xE0A1:
			decodedOpcode = "Skip next if key Vx is not pressed"

		// Load Vx with the value in the delay timer
		case (opcode & 0xF0FF) == 0xF007:
			decodedOpcode = "Value of delay timer is loaded in Vx"

		// Wait for key press, store pressed key in Vx : Fx0A
		case (opcode & 0xF0FF) == 0xF00A:
			decodedOpcode = "Wait for keypress"

		// Load delay timer with Vx
		case (opcode & 0xF0FF) == 0xF015:
			decodedOpcode = "Load delay timer"

		// Load sound timer with Vx
		case (opcode & 0xF0FF) == 0xF018:
			decodedOpcode = "Load sound timer"

		// Add Vx to Indexer
		case (opcode & 0xF0FF) == 0xF01E:
			decodedOpcode = "Add Vx to Indexer"

		// Set Indexer to location for digit sprite that of value in Vx
		case (opcode & 0xF0FF) == 0xF029:
			decodedOpcode = "Load I with address of digit sprite from value Vx"

		// Store BCD representation of Vx in I, I+1, and I+2
		case (opcode & 0xF0FF) == 0xF033:
			decodedOpcode = "Store BCD of Vx at address in the Indexer"

		// Store all of the V registers in memory beginning at the Indexer
		case (opcode & 0xF0FF) == 0xF055:
			decodedOpcode = "Store registers beginning at Indexer"

		// Load all of the V registers beginning at the Indexer
		case (opcode & 0xF065) == 0xF055:
			decodedOpcode = "Load registers beginning at Indexer"

		// Something's gone wrong, unknown opcode.
		// We've crashed.
		default:
			decodedOpcode = "Unknown Opcode"
			c.Crashed = true
		}

		if DEBUG_VERBOSE {
			fmt.Println(decodedOpcode)
		}

		if updateScreen {
			screen <- &c.Buffer
			updateScreen = false
		}
	}
}
