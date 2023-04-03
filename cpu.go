package main

import "C"
import (
	"fmt"
	"io/ioutil"

	"mos6502-emulator/addressing"
	"mos6502-emulator/memory"
	"mos6502-emulator/opcode"
)

const (
	StackPointerHiByte = uint16(0x0100) // high byte of the stack pointer
	Mask8Bit           = 0xFF
)

type InterruptType int

const (
	InterruptTypeIRQ InterruptType = iota
	InterruptTypeNMI
)

// Cpu struct is the core of the 6502 emulator.
// Useful information:
//
// * Little Endian architecture - least significant (low) byte is always stored first.
type Cpu struct {

	// Registers
	A byte
	X byte
	Y byte

	//Stack pointer (low 8 bits), prepend $01 to address memory (high bits). Stack starts at $ff and works up towards $100.
	S  byte
	PC uint16

	// flags
	/* Based on: https://web.archive.org/web/20160406122905/http://homepage.ntlworld.com/cyborgsystems/CS_Main/6502/6502.htm#FLAGS
	Additional info bout flags behaviour: https://www.nesdev.org/wiki/Status_flags
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
	V byte // Overflow
	N byte // Negative

	memory.Memory

	interruptChannel chan InterruptType
}

func NewCpu(interruptChannel chan InterruptType) Cpu {
	return Cpu{interruptChannel: interruptChannel}
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

func (c *Cpu) Reset() {
	c.Z = 0
	c.N = 0
	c.V = 0
	c.D = 0
	c.I = 1
	c.C = 0
	c.S = 0xFF
	//The 6502 stores addresses in low byte/hi byte format, so $FFFD contains the upper 8 bits of the
	//address and $FFFC the lower 8 bits. This line assembles the 2 bytes into a 16-bit address.
	c.PC = uint16(c.Mem[0xFFFD])<<8 | uint16(c.Mem[0xFFFC])

	c.A = 0
	c.X = 0
	c.Y = 0
}

// https://www.nesdev.org/wiki/CPU_interrupts
func (c *Cpu) interrupt(interruptType InterruptType) {
	if interruptType == InterruptTypeIRQ && c.I == byte(1) {
		return
	}
	pcHi := byte(c.PC >> 8)
	pcLo := byte(c.PC & Mask8Bit)
	c.pushToStack(pcHi)
	c.pushToStack(pcLo)
	c.pushToStack(c.getStatusFlags(0))
	c.I = 1

	switch interruptType {
	case InterruptTypeIRQ:
		c.PC = uint16(c.Mem[0xFFFF])<<8 | uint16(c.Mem[0xFFFE])
	case InterruptTypeNMI:
		c.PC = uint16(c.Mem[0xFFFB])<<8 | uint16(c.Mem[0xFFFA])
	default:
		panic(fmt.Sprintf("Unhandled interrupt type: %v", interruptType))
	}
}

func (c *Cpu) Load(path string, offset int, startAddress uint16) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	for i, b := range bytes {
		c.Mem[offset+i] = b
	}

	startAddressLo := byte(startAddress & Mask8Bit)
	startAddressHi := byte(startAddress >> 8)

	c.Mem[0xFFFC] = startAddressLo
	c.Mem[0xFFFD] = startAddressHi

	c.Reset()
	return nil
}

func (c *Cpu) Cycle() {

	select {
	case interruptType := <-c.interruptChannel:
		c.interrupt(interruptType)
	default:
	}

	operation := c.Mem[c.PC]
	c.PC++

	opcodeSpec := opcode.Lookup(operation)
	memoryAccessMode := opcodeSpec.AccessMode
	switch opcodeSpec.Operation {

	case opcode.ORA:
		val, _ := c.read(memoryAccessMode)
		c.ora(val)
	case opcode.AND:
		val, _ := c.read(memoryAccessMode)
		c.and(val)
	case opcode.EOR:
		val, _ := c.read(memoryAccessMode)
		c.eor(val)
	case opcode.ADC:
		val, _ := c.read(memoryAccessMode)
		c.adc(val)
	case opcode.STA:
		address := c.nextByteToAddress(memoryAccessMode)
		c.write(address, c.A, memoryAccessMode)
	case opcode.LDA:
		val, _ := c.read(memoryAccessMode)
		c.lda(val)
	case opcode.CMP:
		val, _ := c.read(memoryAccessMode)
		c.cmp(val)
	case opcode.SBC:
		val, _ := c.read(memoryAccessMode)
		c.sbc(val)
	case opcode.ASL:
		val, address := c.read(memoryAccessMode)
		shifted := c.asl(val)
		c.write(address, shifted, memoryAccessMode)
	case opcode.ROL:
		val, address := c.read(memoryAccessMode)
		rolled := c.rol(val)
		c.write(address, rolled, memoryAccessMode)
	case opcode.LSR:
		val, address := c.read(memoryAccessMode)
		rolled := c.lsr(val)
		c.write(address, rolled, memoryAccessMode)
	case opcode.ROR:
		val, address := c.read(memoryAccessMode)
		rolled := c.ror(val)
		c.write(address, rolled, memoryAccessMode)
	case opcode.STX:
		if memoryAccessMode == addressing.ZeroPageX {
			memoryAccessMode = addressing.ZeroPageY
		}
		address := c.nextByteToAddress(memoryAccessMode)
		c.write(address, c.X, memoryAccessMode)
	case opcode.LDX:
		if memoryAccessMode == addressing.ZeroPageX {
			memoryAccessMode = addressing.ZeroPageY
		} else if memoryAccessMode == addressing.AbsoluteX {
			memoryAccessMode = addressing.AbsoluteY
		}

		val, _ := c.read(memoryAccessMode)
		c.ldx(val)
	case opcode.DEC:
		val, address := c.read(memoryAccessMode)
		val = c.dec(val)
		c.write(address, val, memoryAccessMode)
	case opcode.INC:
		val, address := c.read(memoryAccessMode)
		val = c.inc(val)
		c.write(address, val, memoryAccessMode)
	case opcode.BIT:
		val, _ := c.read(memoryAccessMode)
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
		address := c.nextByteToAddress(memoryAccessMode)
		c.write(address, c.Y, memoryAccessMode)
	case opcode.LDY:
		val, _ := c.read(memoryAccessMode)
		c.ldy(val)
	case opcode.CPY:
		val, _ := c.read(memoryAccessMode)
		c.cpy(val)
	case opcode.CPX:
		val, _ := c.read(memoryAccessMode)
		c.cpx(val)
	case opcode.BRK:
		pc := c.PC + 1
		pcHi := byte(pc >> 8)
		pcLo := byte(pc & Mask8Bit)
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
		status := c.getStatusFlags(1)
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
	if offset < 0x80 {
		return address + uint16(offset)
	} else {
		return address - (0x100 - uint16(offset))
	}
}

func (c *Cpu) read(accessMode addressing.Mode) (byte, uint16) {
	if accessMode == addressing.Accumulator {
		return c.A, 0
	} else {
		address := c.nextByteToAddress(accessMode)
		return c.Mem[address], address
	}
}

func (c *Cpu) write(address uint16, value byte, accessMode addressing.Mode) {
	if accessMode == addressing.Accumulator {
		c.A = value
	} else {
		c.Mem[address] = value
	}
}

// See https://www.pagetable.com/c64ref/6502/?tab=3
// And https://web.archive.org/web/20160406122905/http://homepage.ntlworld.com/cyborgsystems/CS_Main/6502/6502.htm#ADDR_MODE
// for addressing modes reference
func (c *Cpu) nextByteToAddress(accessMode addressing.Mode) uint16 {
	switch accessMode {
	case addressing.Immediate:
		address := c.PC
		c.PC++
		return address
	case addressing.ZeroPage:
		address := c.Mem[c.PC]
		c.PC++
		return uint16(address)
	case addressing.ZeroPageX:
		val := c.Mem[c.PC]
		c.PC++
		address := (val + c.X) & Mask8Bit
		return uint16(address)
	case addressing.ZeroPageY:
		val := c.Mem[c.PC]
		c.PC++
		address := (val + c.Y) & Mask8Bit
		return uint16(address)
	case addressing.Absolute:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		return uint16(hi)<<8 | uint16(lo)
	case addressing.AbsoluteX:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		return address + uint16(c.X)
	case addressing.AbsoluteY:
		lo := c.Mem[c.PC]
		c.PC++
		hi := c.Mem[c.PC]
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		return address + uint16(c.Y)
	case addressing.IndirectX:
		loAddr := c.Mem[c.PC]
		c.PC++
		lo := c.Mem[(loAddr+c.X)&Mask8Bit]
		hi := uint16(c.Mem[(loAddr+c.X+1)&Mask8Bit]) << 8
		return hi | uint16(lo)
	case addressing.IndirectY:
		loAddr := c.Mem[c.PC]
		c.PC++
		lo := c.Mem[loAddr]
		hi := uint16(c.Mem[(loAddr+1)&Mask8Bit]) << 8
		return (hi | uint16(lo)) + uint16(c.Y)
	default:
		panic(fmt.Sprintf("Invalid memory access mode: %v", accessMode))
	}
}

func (c *Cpu) pushToStack(value byte) {
	stackAddress := StackPointerHiByte | uint16(c.S)
	c.Mem[stackAddress] = value
	c.S = c.S - 1
}

func (c *Cpu) pullFromStack() byte {
	c.S = c.S + 1
	stackAddress := StackPointerHiByte | uint16(c.S)
	return c.Mem[stackAddress]
}
