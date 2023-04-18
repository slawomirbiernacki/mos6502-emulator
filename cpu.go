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

	memoryMapper     memory.MemoryMapper
	interruptChannel chan InterruptType

	PendingIRQ bool
	PendingNMI bool
}

func NewCpu(interruptChannel chan InterruptType, memoryMapper memory.MemoryMapper) Cpu {
	return Cpu{interruptChannel: interruptChannel, memoryMapper: memoryMapper}
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
	c.PC = uint16(c.memoryMapper.Read(0xFFFD))<<8 | uint16(c.memoryMapper.Read(0xFFFC))

	c.A = 0
	c.X = 0
	c.Y = 0
}

// https://www.nesdev.org/wiki/CPU_interrupts
func (c *Cpu) interrupt(interruptType InterruptType) int {
	if interruptType == InterruptTypeIRQ && c.I == byte(1) {
		return 0
	}
	pcHi := byte(c.PC >> 8)
	pcLo := byte(c.PC & Mask8Bit)
	c.pushToStack(pcHi)
	c.pushToStack(pcLo)
	c.pushToStack(c.getStatusFlags(0))
	c.I = 1

	switch interruptType {
	case InterruptTypeIRQ:
		c.PC = uint16(c.readFromMemory(0xFFFF))<<8 | uint16(c.readFromMemory(0xFFFE))
	case InterruptTypeNMI:
		c.PC = uint16(c.readFromMemory(0xFFFB))<<8 | uint16(c.readFromMemory(0xFFFA))
	default:
		panic(fmt.Sprintf("Unhandled interrupt type: %v", interruptType))
	}
	return 7
}

func (c *Cpu) Load(path string, offset int, startAddress uint16) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	for i, b := range bytes {
		c.writeToMemory(uint16(offset+i), b)
	}

	startAddressLo := byte(startAddress & Mask8Bit)
	startAddressHi := byte(startAddress >> 8)

	c.writeToMemory(0xFFFC, startAddressLo)
	c.writeToMemory(0xFFFD, startAddressHi)

	c.Reset()
	return nil
}

func (c *Cpu) Run(cycles int) int {
	cycles_executed := cycles

	for cycles_executed > 0 {
		cycles_executed -= c.ExecuteOpcode()
	}
	return cycles - cycles_executed
}

func (c *Cpu) ExecuteOpcode() int {
	cycles := 0
	// select {
	// case interruptType := <-c.interruptChannel:
	// 	cycles += c.interrupt(interruptType) // TODO how do I account for interrupts in cycles, like that?
	// 	return cycles
	// default:
	// }

	if c.PendingIRQ {
		cycles += c.interrupt(InterruptTypeIRQ)
		c.PendingIRQ = false
		return cycles
	}

	if c.PendingNMI {
		cycles += c.interrupt(InterruptTypeNMI)
		c.PendingNMI = false
		return cycles
	}

	operation := c.readFromMemory(c.PC)
	c.PC++

	opcodeSpec := opcode.Lookup(operation)
	memoryAccessMode := opcodeSpec.AccessMode
	cycles += opcodeSpec.Cycles
	switch opcodeSpec.Operation {

	case opcode.ORA:
		val, _, pageCrossed := c.readNext(memoryAccessMode)
		c.ora(val)
		cycles += pageCrossed
	case opcode.AND:
		val, _, pageCrossed := c.readNext(memoryAccessMode)
		c.and(val)
		cycles += pageCrossed
	case opcode.EOR:
		val, _, pageCrossed := c.readNext(memoryAccessMode)
		c.eor(val)
		cycles += pageCrossed
	case opcode.ADC:
		val, _, pageCrossed := c.readNext(memoryAccessMode)
		c.adc(val)
		cycles += pageCrossed
	case opcode.STA:
		address, _ := c.nextByteToAddress(memoryAccessMode)
		c.write(address, c.A, memoryAccessMode)
	case opcode.LDA:
		val, _, pageCrossed := c.readNext(memoryAccessMode)
		c.lda(val)
		cycles += pageCrossed
	case opcode.CMP:
		val, _, pageCrossed := c.readNext(memoryAccessMode)
		c.cmp(val)
		cycles += pageCrossed
	case opcode.SBC:
		val, _, pageCrossed := c.readNext(memoryAccessMode)
		c.sbc(val)
		cycles += pageCrossed
	case opcode.ASL:
		val, address, _ := c.readNext(memoryAccessMode)
		shifted := c.asl(val)
		c.write(address, shifted, memoryAccessMode)
	case opcode.ROL:
		val, address, _ := c.readNext(memoryAccessMode)
		rolled := c.rol(val)
		c.write(address, rolled, memoryAccessMode)
	case opcode.LSR:
		val, address, _ := c.readNext(memoryAccessMode)
		rolled := c.lsr(val)
		c.write(address, rolled, memoryAccessMode)
	case opcode.ROR:
		val, address, _ := c.readNext(memoryAccessMode)
		rolled := c.ror(val)
		c.write(address, rolled, memoryAccessMode)
	case opcode.STX:
		if memoryAccessMode == addressing.ZeroPageX {
			memoryAccessMode = addressing.ZeroPageY
		}
		address, _ := c.nextByteToAddress(memoryAccessMode)
		c.write(address, c.X, memoryAccessMode)
	case opcode.LDX:
		if memoryAccessMode == addressing.ZeroPageX {
			memoryAccessMode = addressing.ZeroPageY
		} else if memoryAccessMode == addressing.AbsoluteX {
			memoryAccessMode = addressing.AbsoluteY
		}

		val, _, pageCrossed := c.readNext(memoryAccessMode)
		c.ldx(val)
		cycles += pageCrossed
	case opcode.DEC:
		val, address, _ := c.readNext(memoryAccessMode)
		val = c.dec(val)
		c.write(address, val, memoryAccessMode)
	case opcode.INC:
		val, address, _ := c.readNext(memoryAccessMode)
		val = c.inc(val)
		c.write(address, val, memoryAccessMode)
	case opcode.BIT:
		val, _, _ := c.readNext(memoryAccessMode)
		c.bit(val)
	case opcode.JMP:
		if memoryAccessMode == addressing.Indirect {
			lo := c.readFromMemory(c.PC)
			c.PC++
			hi := c.readFromMemory(c.PC)
			c.PC++
			address := uint16(hi)<<8 | uint16(lo)
			final_lo := c.readFromMemory(address)
			final_hi := c.readFromMemory(address + 1)
			jumpAddress := uint16(final_hi)<<8 | uint16(final_lo)
			c.PC = jumpAddress
		} else if memoryAccessMode == addressing.Absolute {
			lo := c.readFromMemory(c.PC)
			c.PC++
			hi := c.readFromMemory(c.PC)
			c.PC++
			jumpAddress := uint16(hi)<<8 | uint16(lo)
			c.PC = jumpAddress
		}
	case opcode.STY:
		address, _ := c.nextByteToAddress(memoryAccessMode)
		c.write(address, c.Y, memoryAccessMode)
	case opcode.LDY:
		val, _, pageCrossed := c.readNext(memoryAccessMode)
		c.ldy(val)
		cycles += pageCrossed
	case opcode.CPY:
		val, _, _ := c.readNext(memoryAccessMode)
		c.cpy(val)
	case opcode.CPX:
		val, _, _ := c.readNext(memoryAccessMode)
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
		c.PC = uint16(c.readFromMemory(0xFFFF))<<8 | uint16(c.readFromMemory(0xFFFE))
	case opcode.JSR:
		lo := c.readFromMemory(c.PC)
		c.PC++
		hi := c.readFromMemory(c.PC)
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
		// do nothing
	case opcode.BCC:
		if c.C == 0 {
			cycles += c.takeBranch()
		} else {
			c.PC++
		}
	case opcode.BCS:
		if c.C == 1 {
			cycles += c.takeBranch()
		} else {
			c.PC++
		}
	case opcode.BEQ:
		if c.Z == 1 {
			cycles += c.takeBranch()
		} else {
			c.PC++
		}
	case opcode.BMI:
		if c.N == 1 {
			cycles += c.takeBranch()
		} else {
			c.PC++
		}
	case opcode.BNE:
		if c.Z == 0 {
			cycles += c.takeBranch()
		} else {
			c.PC++
		}
	case opcode.BPL:
		if c.N == 0 {
			cycles += c.takeBranch()
		} else {
			c.PC++
		}
	case opcode.BVC:
		if c.V == 0 {
			cycles += c.takeBranch()
		} else {
			c.PC++
		}
	case opcode.BVS:
		if c.V == 1 {
			cycles += c.takeBranch()
		} else {
			c.PC++
		}
	default:
		panic(fmt.Sprintf("unknown opcode: %v", operation))
	}
	return cycles
}

func (c *Cpu) takeBranch() int {
	offset := c.readFromMemory(c.PC)
	c.PC++
	relativeAddress, pageCrossed := getRelativeAddress(c.PC, offset)
	c.PC = relativeAddress
	return pageCrossed + 1 // +1 for taking the branch
}

// return address and whether page bounduary has been crossed
func getRelativeAddress(address uint16, offset byte) (uint16, int) {
	var resultAddress uint16
	if offset < 0x80 {
		resultAddress = address + uint16(offset)
	} else {
		resultAddress = address - (0x100 - uint16(offset))
	}
	return resultAddress, hiByteDiffers(address, resultAddress)
}

// Returns value, address and if page bounduary was crossed
func (c *Cpu) readNext(accessMode addressing.Mode) (byte, uint16, int) {
	if accessMode == addressing.Accumulator {
		return c.A, 0, 0
	} else {
		address, pageCrossed := c.nextByteToAddress(accessMode)
		return c.readFromMemory(address), address, pageCrossed
	}
}

func (c *Cpu) readFromMemory(address uint16) byte {
	return c.memoryMapper.Read(address)
}

func (c *Cpu) write(address uint16, value byte, accessMode addressing.Mode) {
	if accessMode == addressing.Accumulator {
		c.A = value
	} else {
		c.writeToMemory(address, value)
	}
}

func (c *Cpu) writeToMemory(address uint16, value byte) {
	c.memoryMapper.Write(address, value)
}

// See https://www.pagetable.com/c64ref/6502/?tab=3
// And https://web.archive.org/web/20160406122905/http://homepage.ntlworld.com/cyborgsystems/CS_Main/6502/6502.htm#ADDR_MODE
// for addressing modes reference
// second int returned indicates whether page has been crossed for access modes that affect timing based on that
func (c *Cpu) nextByteToAddress(accessMode addressing.Mode) (uint16, int) {
	switch accessMode {
	case addressing.Immediate:
		address := c.PC
		c.PC++
		return address, 0
	case addressing.ZeroPage:
		address := c.readFromMemory(c.PC)
		c.PC++
		return uint16(address), 0
	case addressing.ZeroPageX:
		val := c.readFromMemory(c.PC)
		c.PC++
		address := (val + c.X) & Mask8Bit
		return uint16(address), 0
	case addressing.ZeroPageY:
		val := c.readFromMemory(c.PC)
		c.PC++
		address := (val + c.Y) & Mask8Bit
		return uint16(address), 0
	case addressing.Absolute:
		lo := c.readFromMemory(c.PC)
		c.PC++
		hi := c.readFromMemory(c.PC)
		c.PC++
		return uint16(hi)<<8 | uint16(lo), 0
	case addressing.AbsoluteX:
		lo := c.readFromMemory(c.PC)
		c.PC++
		hi := c.readFromMemory(c.PC)
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		result := address + uint16(c.X)
		return result, hiByteDiffers(result, address)
	case addressing.AbsoluteY:
		lo := c.readFromMemory(c.PC)
		c.PC++
		hi := c.readFromMemory(c.PC)
		c.PC++
		address := uint16(hi)<<8 | uint16(lo)
		result := address + uint16(c.Y)
		return result, hiByteDiffers(result, address)
	case addressing.IndirectX:
		loAddr := c.readFromMemory(c.PC)
		c.PC++
		lo := c.readFromMemory(uint16((loAddr + c.X) & Mask8Bit))
		hi := uint16(c.readFromMemory(uint16((loAddr+c.X+1)&Mask8Bit))) << 8
		return hi | uint16(lo), 0
	case addressing.IndirectY:
		loAddr := c.readFromMemory(c.PC)
		c.PC++
		lo := c.readFromMemory(uint16(loAddr))
		hi := uint16(c.readFromMemory(uint16((loAddr+1)&Mask8Bit))) << 8
		address := hi | uint16(lo)
		result := address + uint16(c.Y)
		return result, hiByteDiffers(result, address)
	default:
		panic(fmt.Sprintf("Invalid memory access mode: %v", accessMode))
	}
}

func hiByteDiffers(a, b uint16) int {
	if (a >> 4) == (b >> 0) {
		return 1
	} else {
		return 0
	}
}

func (c *Cpu) pushToStack(value byte) {
	stackAddress := StackPointerHiByte | uint16(c.S)
	c.writeToMemory(stackAddress, value)
	c.S = c.S - 1
}

func (c *Cpu) pullFromStack() byte {
	c.S = c.S + 1
	stackAddress := StackPointerHiByte | uint16(c.S)
	return c.readFromMemory(stackAddress)
}
