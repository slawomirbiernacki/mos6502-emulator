package main

import "C"
import (
	"fmt"
	"io/ioutil"

	"mos6502-emulator/addressing"
	"mos6502-emulator/memory"
	"mos6502-emulator/opcode"
)

// Little Endian architecture

type Cpu struct {
	//TODO optimise data formats
	A  byte
	X  byte
	Y  byte
	S  byte //stack pointer (low 8 bits), prepend $01 to address memory (high bits)
	PC uint16
	//P  byte // status

	// flags
	/* From: https://www.nesdev.org/wiki/Status_flags
	7  bit  0
	---- ----
	NV1s DIZC
	|||| ||||
	|||| |||+- Carry
	|||| ||+-- Zero
	|||| |+--- Interrupt Disable
	|||| +---- Decimal
	||++------ No CPU effect, see: the B flag
	|+-------- Overflow
	+--------- Negative
	*/
	C byte // Carry
	Z byte // Zero
	I byte // Interrupt
	D byte // Decimal
	//B byte
	V byte // Overflow
	N byte // Negative

	memory.Memory
}

func NewCpu() Cpu {
	return Cpu{}
}

// bFlag needs to be 1 or 0
func (c *Cpu) getStatusFlags(bFlag byte) byte {
	return (c.N << 7) + (c.V << 6) + 0b00100000 + (bFlag << 4) + (c.D << 3) + (c.I << 2) + (c.Z << 1) + c.C
}

func (c *Cpu) setStatusFlags(value byte) {
	c.N = value >> 7
	c.V = value & 0b01000000 >> 6
	c.D = value & 0b00001000 >> 3
	c.I = value & 0b00000100 >> 2
	c.Z = value & 0b00000010 >> 1
	c.C = value & 0b00000001
}

//TODO how to handle IRQ and other interruptions?

func (c *Cpu) Reset() {
	//c.P = 0x20
	c.Z = 0 //Z
	c.N = 0 //N
	c.V = 0 //V
	//c.B = 0 //B
	c.D = 0
	c.I = 1
	c.C = 0 //C
	c.S = 0xFF
	//The 6502 stores addresses in low byte/hi byte format, so $FFFD contains the upper 8 bits of the
	//address and $FFFC the lower 8 bits. This line assembles the 2 bytes into a 16-bit address.
	c.PC = uint16(c.Mem[0xFFFD])<<8 | uint16(c.Mem[0xFFFC])

	c.A = 0
	c.X = 0
	c.Y = 0
}

func (c *Cpu) Load(path string, offset int, programCounter uint16) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	//offset := 0xE000
	//offset := 0x9FF0
	//offset := 0x400

	for i, b := range bytes {
		c.Mem[offset+i] = b
	}

	programCounterLo := byte(programCounter & 0x00FF)
	programCounterHi := byte(programCounter >> 8)

	c.Mem[0xFFFC] = programCounterLo
	c.Mem[0xFFFD] = programCounterHi

	c.Reset()

	return nil
}

func (c *Cpu) Cycle() {

	operation := c.Mem[c.PC]
	c.PC++

	opcodeSpec := opcode.Lookup(operation)
	memoryAccessMode := opcodeSpec.AccessMode
	switch opcodeSpec.Operation {

	case opcode.ORA:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.ora(val)
	case opcode.AND:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.and(val)

	case opcode.EOR:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.eor(val)

	case opcode.ADC:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.adc(val)

	case opcode.STA:
		address := c.nextOpAsAddress(memoryAccessMode)
		c.write(address, c.A, memoryAccessMode)
	case opcode.LDA:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.lda(val)
	case opcode.CMP:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.cmp(val)
	case opcode.SBC:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.sbc(val)
	case opcode.ASL:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		shifted := c.asl(val)
		c.write(address, shifted, memoryAccessMode)
	case opcode.ROL:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		rolled := c.rol(val)
		c.write(address, rolled, memoryAccessMode)

	case opcode.LSR:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		rolled := c.lsr(val)
		c.write(address, rolled, memoryAccessMode)

	case opcode.ROR:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		rolled := c.ror(val)
		c.write(address, rolled, memoryAccessMode)
	case opcode.STX:
		if memoryAccessMode == addressing.ZeroPageX {
			memoryAccessMode = addressing.ZeroPageY
		}
		address := c.nextOpAsAddress(memoryAccessMode)
		c.write(address, c.X, memoryAccessMode)
	case opcode.LDX:
		if memoryAccessMode == addressing.ZeroPageX {
			memoryAccessMode = addressing.ZeroPageY
		} else if memoryAccessMode == addressing.AbsoluteX {
			memoryAccessMode = addressing.AbsoluteY
		}

		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.ldx(val)
	case opcode.DEC:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		val = c.dec(val)
		c.write(address, val, memoryAccessMode)
	case opcode.INC:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		val = c.inc(val)
		c.write(address, val, memoryAccessMode)
	case opcode.BIT:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.bit(val)
	case opcode.JMP:
		if memoryAccessMode == addressing.Indirect {
			lo := c.Mem[c.PC]
			c.PC++
			hi := c.Mem[c.PC]
			c.PC++
			address := uint16(hi)<<8 | uint16(lo)
			final_lo := c.Mem[address]
			final_hi := c.Mem[address+1]
			jumpAddress := uint16(final_hi)<<8 | uint16(final_lo)
			c.PC = jumpAddress
		} else if memoryAccessMode == addressing.Absolute {
			lo := c.Mem[c.PC]
			c.PC++
			hi := c.Mem[c.PC]
			c.PC++
			jumpAddress := uint16(hi)<<8 | uint16(lo)
			c.PC = jumpAddress
		}
	case opcode.STY:
		address := c.nextOpAsAddress(memoryAccessMode)
		c.write(address, c.Y, memoryAccessMode)
	case opcode.LDY:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.ldy(val)
	case opcode.CPY:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.cpy(val)
	case opcode.CPX:
		address := c.nextOpAsAddress(memoryAccessMode)
		val := c.read(address, memoryAccessMode)
		c.cpx(val)
	case opcode.BRK:
		// bFlag=1
		pc := c.PC + 1
		pcHi := byte(pc >> 8)
		pcLo := byte(pc & 0xFF)
		c.pushToStack(pcHi)
		c.pushToStack(pcLo)
		flags := c.getStatusFlags(1)
		c.pushToStack(flags)
		c.I = 1
		c.PC = uint16(c.Mem[0xFFFF])<<8 | uint16(c.Mem[0xFFFE])
	case opcode.JSR:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)

		t := c.PC - 1
		tHi := byte(t >> 8)
		tLo := byte(t & 0xFF)
		c.pushToStack(tHi)
		c.pushToStack(tLo)
		c.PC = address
	case opcode.RTI:
		flags := c.pullFromStack()
		c.setStatusFlags(flags)
		lo := c.pullFromStack()
		hi := c.pullFromStack()
		address := uint16(hi)<<8 | uint16(lo)
		c.PC = address
	case opcode.RTS:
		lo := c.pullFromStack()
		hi := c.pullFromStack()
		address := uint16(hi)<<8 | uint16(lo)
		c.PC = address + 1
	case opcode.PHP:
		status := c.getStatusFlags(1) //TODO b=0 maybe?
		c.pushToStack(status)
	case opcode.PLP:
		flags := c.pullFromStack()
		c.setStatusFlags(flags)
	case opcode.PHA:
		c.pushToStack(c.A)
	case opcode.PLA:
		val := c.pullFromStack()
		c.pla(val)
	case opcode.DEY:
		c.dey()
	case opcode.TAY:
		c.tay()
	case opcode.INY:
		c.iny()
	case opcode.INX:
		c.inx()
	case opcode.CLC: // (CLear Carry)
		c.C = 0
	case opcode.SEC: // (SEt Carry)
		c.C = 1
	case opcode.CLI: // (CLear Interrupt)
		c.I = 0
	case opcode.SEI: // (SEt Interrupt)
		c.I = 1
	case opcode.CLV: // (CLear oVerflow)
		c.V = 0
	case opcode.CLD: // (CLear Decimal)
		c.D = 0
	case opcode.SED: // (SEt Decimal)
		c.D = 1
	case opcode.TYA:
		c.tya()
	case opcode.TXA:
		c.txa()
	case opcode.TXS:
		c.S = c.X
	case opcode.TAX:
		c.tax()
	case opcode.TSX:
		c.tsx()
	case opcode.DEX:
		c.dex()
	case opcode.NOP:
		return
	case opcode.BCC:
		if c.C == 0 {
			offset := c.Mem[c.PC]
			c.PC++
			c.PC = getRelativeAddress(c.PC, offset)
		} else {
			c.PC++
		}
	case opcode.BCS:
		if c.C == 1 {
			offset := c.Mem[c.PC]
			c.PC++
			c.PC = getRelativeAddress(c.PC, offset)
		} else {
			c.PC++
		}
	case opcode.BEQ:
		if c.Z == 1 {
			offset := c.Mem[c.PC]
			c.PC++
			c.PC = getRelativeAddress(c.PC, offset)
		} else {
			c.PC++
		}
	case opcode.BMI:
		if c.N == 1 {
			offset := c.Mem[c.PC]
			c.PC++
			c.PC = getRelativeAddress(c.PC, offset)
		} else {
			c.PC++
		}
	case opcode.BNE:
		if c.Z == 0 {
			offset := c.Mem[c.PC]
			c.PC++
			c.PC = getRelativeAddress(c.PC, offset)
		} else {
			c.PC++
		}
	case opcode.BPL:
		if c.N == 0 {
			offset := c.Mem[c.PC]
			c.PC++
			c.PC = getRelativeAddress(c.PC, offset)
		} else {
			c.PC++
		}
	case opcode.BVC:
		if c.V == 0 {
			offset := c.Mem[c.PC]
			c.PC++
			c.PC = getRelativeAddress(c.PC, offset)
		} else {
			c.PC++
		}
	case opcode.BVS:
		if c.V == 1 {
			offset := c.Mem[c.PC]
			c.PC++
			c.PC = getRelativeAddress(c.PC, offset)
		} else {
			c.PC++
		}
	default:
		panic(fmt.Sprintf("unknown opcode: %v", operation))
	}
}

func getRelativeAddress(address uint16, offset byte) uint16 {
	//lo := byte(address & 0xFF)
	//hi := address >> 8
	//lo = lo + offset
	//return hi<<8 | uint16(lo)

	if offset < 0x80 {
		return address + uint16(offset)
	} else {
		return address - (0x100 - uint16(offset))
	}
}

func (c *Cpu) nextOpAsAddress(accessMode addressing.AccesssMode) (address uint16) {
	switch accessMode {
	case addressing.Accumulator:
		return 0
	case addressing.Immediate:
		address = c.PC
		c.PC++
		return
	case addressing.ZeroPage:
		val := c.Mem[c.PC]
		c.PC++
		return uint16(val)

	case addressing.ZeroPageX:
		val := c.Mem[c.PC]
		c.PC++
		address := (val + c.X) & 0x00FF
		return uint16(address)
	case addressing.ZeroPageY:
		val := c.Mem[c.PC]
		c.PC++

		address := (val + c.Y) & 0x00FF
		return uint16(address)
	case addressing.Absolute:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		return address
	case addressing.AbsoluteX:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		valX := uint16(c.X)
		address = (address + valX) & 0xFFFF
		return address
	case addressing.AbsoluteY:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		valY := uint16(c.Y)
		address = (address + valY) & 0xFFFF
		return address
	case addressing.IndirectX:
		loAddr := c.Mem[c.PC]
		c.PC++
		lo := c.Mem[(loAddr+c.X)&0xFF]
		hi := uint16(c.Mem[(loAddr+c.X+1)&0xFF]) << 8
		address := hi | uint16(lo)
		return address
	case addressing.IndirectY:
		loAddr := c.Mem[c.PC]
		c.PC++
		lo := c.Mem[loAddr]
		hi := uint16(c.Mem[(loAddr+1)&0xFF]) << 8
		address := (hi | uint16(lo) + uint16(c.Y)) & 0xFFFF
		return address

	default:
		panic(fmt.Sprintf("Invalid memory access mode: %v", accessMode))
	}
}

func (c *Cpu) write(address uint16, value byte, accessMode addressing.AccesssMode) {
	if accessMode == addressing.Accumulator {
		c.A = value
	} else {
		c.Mem[address] = value
	}
}

func (c *Cpu) read(address uint16, accessMode addressing.AccesssMode) byte {
	if accessMode == addressing.Accumulator {
		return c.A
	} else {
		return c.Mem[address]
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

func (c *Cpu) pushToStack(value byte) {
	stackPrefix := uint16(0x0100)
	stackAddress := stackPrefix | uint16(c.S)
	c.Mem[stackAddress] = value
	c.S = c.S - 1
}

func (c *Cpu) pullFromStack() byte {
	stackPrefix := uint16(0x0100)
	c.S = c.S + 1
	stackAddress := stackPrefix | uint16(c.S)
	return c.Mem[stackAddress]
}
