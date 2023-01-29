package main

import (
	"fmt"
	"io/ioutil"
	"math/bits"
)

// Little Endian architecture

/*
Memory map:
------
|  $0000
|  Zero page
|  $00FF
------
|  $0100
|  Stack (page 1) - push(go up), pull (go down)
|  $01FF
------
|  $0200
|  Peripherals
|  $02FF
------
|  $0300
|  Free RAM
|  $DFFF(?)
------
|  $E000
|  ROM
|  $FFFF
------
*/

type Cpu struct {
	//TODO optimise data formats
	A  byte
	X  byte
	Y  byte
	S  byte //stack pointer (low 8 bits), prepend $01 to address memory (high bits)
	PC uint16
	P  byte // status
	// flags
	ZeroFlag      int
	NegativeFlag  int
	OverflowFlag  int
	BreakFlag     int
	DecimalFlag   int
	InterruptFlag int
	CarryFlag     int

	Mem [65536]byte
}

func (c *Cpu) Reset() {
	c.P = 0x20
	c.ZeroFlag = 0     //Z
	c.NegativeFlag = 0 //N
	c.OverflowFlag = 0 //V
	c.BreakFlag = 0    //B
	c.DecimalFlag = 0
	c.InterruptFlag = 0
	c.CarryFlag = 0 //C
	c.S = 0xFF
	//The 6502 stores addresses in low byte/hi byte format, so $FFFD contains the upper 8 bits of the
	//address and $FFFC the lower 8 bits. This line assembles the 2 bytes into a 16-bit address.
	c.PC = uint16(c.Mem[0xFFFD])<<8 | uint16(c.Mem[0xFFFC])

	c.A = 0
	c.X = 0
	c.Y = 0
}

func (c *Cpu) Load(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	//offset := 0xE000
	offset := 0x9FF0

	for i, b := range bytes {
		c.Mem[offset+i] = b
	}

	//c.Mem[0xFFFC] = 0x00
	//c.Mem[0xFFFD] = 0xE0
	c.Mem[0xFFFC] = 0xF0
	c.Mem[0xFFFD] = 0x9F

	return nil
}

func (c *Cpu) Cycle() {

	opcode := c.Mem[c.PC]
	c.PC++

	opType := opcode & 0b00000011
	operation := opcode >> 5
	memoryAccessMode := opcode & 0b00011100

	if opType == 00 {
		switch operation {
		case 0b000: // CLC (CLear Carry)
			c.CarryFlag = 0
		case 0b001: // SEC (SEt Carry)
			c.CarryFlag = 1
		case 0b010: // CLI (CLear Interrupt)
			c.InterruptFlag = 0
		case 0b011: // SEI (SEt Interrupt)
			c.InterruptFlag = 1
		case 0b101: // CLV (CLear oVerflow)
			c.OverflowFlag = 0
		case 0b110: // CLD (CLear Decimal)
			c.DecimalFlag = 0
		case 0b111: // SED (SEt Decimal)
			c.DecimalFlag = 0
		}
		return
	}

	switch operation {
	case 0b011: // ADC
		val := c.readMemory(memoryAccessMode)
		res, carry := bits.Add(uint(val), uint(c.A), uint(c.CarryFlag))
		if carry == 1 {
			c.CarryFlag = 1
		} else {
			c.CarryFlag = 0
		}
		c.A = byte(res)
		//TODO flags, signed values
	case 0b101: // LDA
		val := c.readMemory(memoryAccessMode)
		c.A = val
		if val == 0 {
			c.ZeroFlag = 1
		} else {
			c.ZeroFlag = 0
		}
		c.NegativeFlag = int(val >> 7)
	}
}

func (c *Cpu) readMemory(accessMode byte) byte {
	switch accessMode {
	case 0b00001000: // Immediate
		val := c.Mem[c.PC]
		c.PC++
		return val

	case 0b00000100: // Zero page
		val := c.Mem[c.PC]
		c.PC++
		return c.Mem[val]

	case 0b00010100: // Zero page, X
		val := c.Mem[c.PC]
		c.PC++
		valX := c.Mem[c.X]

		address := (val + valX) & 0b00001111
		return c.Mem[address]
	case 0b00001100: // Absolute
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi<<8) | uint16(lo)
		return c.Mem[address]
	case 0b00011100: // Absolute, X
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi<<8) | uint16(lo)
		valX := uint16(c.Mem[c.X])
		return c.Mem[address+valX]
	case 0b00011000: // Absolute, Y
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi<<8) | uint16(lo)
		valY := uint16(c.Mem[c.Y])
		return c.Mem[address+valY]
	case 0b00000000: // Indirect, X
		val := c.Mem[c.PC]
		c.PC++
		valX := c.Mem[c.X]

		address := (val + valX) & 0b00001111
		lo := c.Mem[address]
		hi := c.Mem[address+1]
		dataAddress := uint16(hi<<8) | uint16(lo)
		return c.Mem[dataAddress]
	case 0b00010000: // Indirect, Y
		val := c.Mem[c.PC]
		c.PC++
		valY := c.Mem[c.Y]

		address := (val + valY) & 0b00001111
		lo := c.Mem[address]
		hi := c.Mem[(address+1)&0b00001111]
		dataAddress := uint16(hi<<8) | uint16(lo)
		return c.Mem[dataAddress]

	default:
		panic(fmt.Sprintf("Invalid memory access mode: %v", accessMode))
	}
}

// 8-bit stack pointer (fixed at RAM address $100, so can address $100-$1ff)

//The Stack Pointer(SP)is used to keep track of the current position of the stack.
//For example the stack on the 6502 is at memory locations $1ff-$100, it starts at
//$1ff and works it's way down towards $100. The stack pointer is 8 bits wide so
//it would start out at $ff (the processor knows it really means $1ff). When a
//value is pushed onto the stack it will be put at memory location $1ff and then
//the SP will be de-incremented to it points to $1fe. When data is pulled of the
//stack, the SP is incremented, then the data is read from that memory location.

func (c *Cpu) PushToStack(value byte) {

}
