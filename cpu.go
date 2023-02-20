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
	c.I = 0
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

	switch operation {
	case opcode.BRK:
		// bFlag=1
		pc := c.PC + 1
		pcHi := byte(pc >> 8)
		pcLo := byte(pc | 0xFF)
		c.pushToStack(pcHi)
		c.pushToStack(pcLo)
		flags := c.getStatusFlags(1)
		c.pushToStack(flags)
		c.PC = uint16(c.Mem[0xFFFF])<<8 | uint16(c.Mem[0xFFFE])
		return

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
		return
	case opcode.RTI:
		flags := c.pullFromStack()
		c.setStatusFlags(flags)
		lo := c.pullFromStack()
		hi := c.pullFromStack()
		address := uint16(hi)<<8 | uint16(lo)
		c.PC = address
		return
	case opcode.RTS:
		lo := c.pullFromStack()
		hi := c.pullFromStack()
		address := uint16(hi)<<8 | uint16(lo)
		c.PC = address + 1
		return
	case opcode.PHP:
		status := c.getStatusFlags(1) //TODO b=0 maybe?
		c.pushToStack(status)
		return
	case opcode.PLP:
		flags := c.pullFromStack()
		c.setStatusFlags(flags)
		return
	case opcode.PHA:
		c.pushToStack(c.A)
		return
	case opcode.PLA:
		c.A = c.pullFromStack()
		c.pla(c.A)
		return
	case opcode.DEY:
		c.dey()
		return
	case opcode.TAY:
		c.tay()
		return
	case opcode.INY:
		c.iny()
		return
	case opcode.INX:
		c.inx()
		return

	case opcode.CLC: // (CLear Carry)
		c.C = 0
		return
	case opcode.SEC: // (SEt Carry)
		c.C = 1
		return
	case opcode.CLI: // (CLear Interrupt)
		c.I = 0
		return
	case opcode.SEI: // (SEt Interrupt)
		c.I = 1
		return
	case opcode.CLV: // (CLear oVerflow)
		c.V = 0
		return
	case opcode.CLD: // (CLear Decimal)
		c.D = 0
		return
	case opcode.SED: // (SEt Decimal)
		c.D = 0
		return
	case opcode.TYA:
		c.tya()
		return
	case opcode.TXA:
		c.txa()
		return
	case opcode.TXS:
		c.S = c.X
		return
	case opcode.TAX:
		c.tax()
		return
	case opcode.TSX:
		c.tsx()
		return
	case opcode.DEX:
		c.dex()
		return
	case opcode.NOP:
		return
	}

	cc := operation & 0b00000011
	// aaabbbccc
	// cc = 01 group
	if cc == 0b01 {
		opcodePrefix := operation >> 5
		memoryAccessMode := addressing.GetForCC01Code(operation & 0b00011100)
		switch opcodePrefix {
		case opcode.ORA:
			val := c.readMemory(memoryAccessMode)
			c.ora(val)
			return
		case opcode.AND:
			val := c.readMemory(memoryAccessMode)
			c.and(val)
			return
		case opcode.EOR:
			val := c.readMemory(memoryAccessMode)
			c.eor(val)
			return
		case opcode.ADC:
			val := c.readMemory(memoryAccessMode)
			c.adc(val)
			return
		case opcode.STA: // not restricting possible addressing modes
			c.writeMemory(memoryAccessMode, c.A)
			return
		case opcode.LDA:
			val := c.readMemory(memoryAccessMode)
			c.lda(val)
			return
		case opcode.CMP:
			val := c.readMemory(memoryAccessMode)
			c.cmp(val)
			return
		case opcode.SBC:
			val := c.readMemory(memoryAccessMode)
			c.sbc(val)
			return
		}
	}

	// xxy10000
	if operation&0b00010000 == 0b00010000 {
		xx := operation >> 6
		y := operation & 0b00100000 >> 5

		switch xx {
		case 0b00: // BPL or BMI
			if c.N == y {
				offset := c.Mem[c.PC]
				c.PC++
				c.PC = getRelativeAddress(c.PC, offset)
			} else {
				c.PC++
			}
			return
		case 0b01: // BVS or BVC
			if c.V == y {
				offset := c.Mem[c.PC]
				c.PC++
				c.PC = getRelativeAddress(c.PC, offset)
			} else {
				c.PC++
			}
			return
		case 0b10: // BCS or BCC
			if c.C == y {
				offset := c.Mem[c.PC]
				c.PC++
				c.PC = getRelativeAddress(c.PC, offset)
			} else {
				c.PC++
			}
			return
		case 0b11: // BNE or BEQ
			if c.Z == y {
				offset := c.Mem[c.PC]
				c.PC++
				c.PC = getRelativeAddress(c.PC, offset)
			} else {
				c.PC++
			}

			return
		}
	}

	// cc = 10 group
	if cc == 0b10 {
		opcodePrefix := operation >> 5
		memoryAccessMode := addressing.GetForCC10Code(operation & 0b00011100)

		//	CC10_ZeroPageX   AccesssMode = 0b00010100 // ZeroPageY for STX and LDX
		//	CC10_AbsoluteX   AccesssMode = 0b00011100 // AbsoluteY for STX and LDX

		switch opcodePrefix {
		case opcode.ASL:
			val := c.readMemory(memoryAccessMode)
			shifted := c.asl(val)
			c.writeMemory(memoryAccessMode, shifted)
			return
		case opcode.ROL:
			val := c.readMemory(memoryAccessMode)
			rolled := c.rol(val)
			c.writeMemory(memoryAccessMode, rolled)
			return
		case opcode.LSR:
			val := c.readMemory(memoryAccessMode)
			rolled := c.lsr(val)
			c.writeMemory(memoryAccessMode, rolled)
			return
		case opcode.ROR:
			val := c.readMemory(memoryAccessMode)
			rolled := c.ror(val)
			c.writeMemory(memoryAccessMode, rolled)
			return
		case opcode.STX:
			if memoryAccessMode == addressing.ZeroPageX {
				memoryAccessMode = addressing.ZeroPageY
			}
			c.writeMemory(memoryAccessMode, c.X)
			return
		case opcode.LDX:
			if memoryAccessMode == addressing.ZeroPageX {
				memoryAccessMode = addressing.ZeroPageY
			} else if memoryAccessMode == addressing.AbsoluteX {
				memoryAccessMode = addressing.AbsoluteY
			}
			val := c.readMemory(memoryAccessMode)
			c.ldx(val)
			return

		case opcode.DEC:
			val := c.readMemory(memoryAccessMode)
			val = c.dec(val)
			c.writeMemory(memoryAccessMode, val)
			return
		case opcode.INC:
			val := c.readMemory(memoryAccessMode)
			val = c.inc(val)
			c.writeMemory(memoryAccessMode, val)
			return

		}
	}

	// cc = 00 group
	if cc == 0b00 {
		opcodePrefix := operation >> 5
		// Addressing same as for 10, except accumulator missing
		memoryAccessMode := addressing.GetForCC10Code(operation & 0b00011100)
		switch opcodePrefix {
		case opcode.BIT:
			val := c.readMemory(memoryAccessMode)
			c.bit(val)
			return
		case opcode.JMP:
			lo := c.Mem[c.PC]
			c.PC++
			hi := c.Mem[c.PC]
			c.PC++
			address := uint16(hi)<<8 | uint16(lo)
			final_lo := c.Mem[address]
			final_hi := c.Mem[address+1]
			jumpAddress := uint16(final_hi)<<8 | uint16(final_lo)
			c.PC = jumpAddress
			return
		case opcode.JMP_ABS:
			lo := c.Mem[c.PC]
			c.PC++
			hi := c.Mem[c.PC]
			c.PC++
			jumpAddress := uint16(hi)<<8 | uint16(lo)
			c.PC = jumpAddress
			return
		case opcode.STY:
			c.writeMemory(memoryAccessMode, c.Y)
			return
		case opcode.LDY:
			val := c.readMemory(memoryAccessMode)
			c.ldy(val)
			return
		case opcode.CPY:
			val := c.readMemory(memoryAccessMode)
			c.cpy(val)
			return
		case opcode.CPX:
			val := c.readMemory(memoryAccessMode)
			c.cpx(val)
			return
		}
	}

	panic(fmt.Sprintf("opcode not implemented %v", operation))
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

func (c *Cpu) readMemory(accessMode addressing.AccesssMode) byte {
	switch accessMode {
	case addressing.Accumulator:
		return c.A
	case addressing.Immediate:
		val := c.Mem[c.PC]
		c.PC++
		return val

	case addressing.ZeroPage:
		val := c.Mem[c.PC]
		c.PC++
		return c.Mem[val]

	case addressing.ZeroPageX:
		val := c.Mem[c.PC]
		c.PC++
		valX := c.Mem[c.X]

		address := (val + valX) & 0b00001111
		return c.Mem[address]
	case addressing.ZeroPageY:
		val := c.Mem[c.PC]
		c.PC++
		valY := c.Mem[c.Y]

		address := (val + valY) & 0b00001111
		return c.Mem[address]
	case addressing.Absolute: //TODO wrapping?
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		return c.Mem[address]
	case addressing.AbsoluteX:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		valX := uint16(c.Mem[c.X])
		return c.Mem[address+valX]
	case addressing.AbsoluteY:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		valY := uint16(c.Mem[c.Y])
		return c.Mem[address+valY]
	case addressing.IndirectX:
		val := c.Mem[c.PC]
		c.PC++
		valX := c.Mem[c.X]

		address := (val + valX) & 0b00001111
		lo := c.Mem[address]
		hi := c.Mem[address+1]
		dataAddress := uint16(hi)<<8 | uint16(lo)
		return c.Mem[dataAddress]
	case addressing.IndirectY:
		val := c.Mem[c.PC]
		c.PC++
		valY := c.Mem[c.Y]

		address := (val + valY) & 0b00001111
		lo := c.Mem[address]
		hi := c.Mem[(address+1)&0b00001111]
		dataAddress := uint16(hi)<<8 | uint16(lo)
		return c.Mem[dataAddress]

	default:
		panic(fmt.Sprintf("Invalid memory access mode: %v", accessMode))
	}
}

func (c *Cpu) writeMemory(accessMode addressing.AccesssMode, valueToWrite byte) {
	switch accessMode {
	case addressing.Accumulator:
		c.A = valueToWrite
	case addressing.ZeroPage:
		address := c.Mem[c.PC]
		c.PC++
		c.Mem[address] = valueToWrite

	case addressing.ZeroPageX:
		val := c.Mem[c.PC]
		c.PC++
		valX := c.Mem[c.X]

		address := (val + valX) & 0b00001111
		c.Mem[address] = valueToWrite
	case addressing.ZeroPageY:
		val := c.Mem[c.PC]
		c.PC++
		valY := c.Mem[c.Y]

		address := (val + valY) & 0b00001111
		c.Mem[address] = valueToWrite
	case addressing.Absolute:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		c.Mem[address] = valueToWrite
	case addressing.AbsoluteX:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		valX := uint16(c.Mem[c.X])
		c.Mem[address+valX] = valueToWrite
	case addressing.AbsoluteY:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		valY := uint16(c.Mem[c.Y])
		c.Mem[address+valY] = valueToWrite
	case addressing.IndirectX:
		val := c.Mem[c.PC]
		c.PC++
		valX := c.Mem[c.X]

		address := (val + valX) & 0b00001111
		lo := c.Mem[address]
		hi := c.Mem[address+1]
		dataAddress := uint16(hi)<<8 | uint16(lo)
		c.Mem[dataAddress] = valueToWrite
	case addressing.IndirectY:
		val := c.Mem[c.PC]
		c.PC++
		valY := c.Mem[c.Y]

		address := (val + valY) & 0b00001111
		lo := c.Mem[address]
		hi := c.Mem[(address+1)&0b00001111]
		dataAddress := uint16(hi)<<8 | uint16(lo)
		c.Mem[dataAddress] = valueToWrite

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
