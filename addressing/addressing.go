// Confidential and only for the information of the intended recipient.
// Copyright 2023 Improbable Worlds Limited.

package addressing

import "fmt"

type AccesssMode int

// Based on https://llx.com/Neil/a2/opcodes.html
// Opcodes are treated as aaabbbcc
// Addressing mode values depend on the cc part, thus different groups below.
const (
	ZeroPage AccesssMode = iota
	Immediate
	Absolute
	IndirectX // (zero page,X)
	IndirectY // (zero page),Y
	ZeroPageX
	ZeroPageY
	AbsoluteY
	AbsoluteX
	Accumulator
)

func GetForCC01Code(bbb byte) AccesssMode {
	switch bbb {
	case 0b00000000:
		return IndirectX // (zero page,X)
	case 0b00000100:
		return ZeroPage
	case 0b00001000:
		return Immediate
	case 0b00001100:
		return Absolute
	case 0b00010000:
		return IndirectY // (zero page),Y
	case 0b00010100:
		return ZeroPageX
	case 0b00011000:
		return AbsoluteY
	case 0b00011100:
		return AbsoluteX
	default:
		panic(fmt.Sprintf("unrecognised addressing mode mask: %v", bbb))
	}
}

func GetForCC10Code(bbb byte) AccesssMode {
	switch bbb {
	case 0b00000000:
		return Immediate
	case 0b00000100:
		return ZeroPage
	case 0b00001000:
		return Accumulator
	case 0b00001100:
		return Absolute
	case 0b00010100:
		return ZeroPageX
	case 0b00011100:
		return AbsoluteX
	default:
		panic(fmt.Sprintf("unrecognised addressing mode mask: %v", bbb))
	}
}

//const (
//	CC01_IndirectX AccesssMode = 0b00000000 // (zero page,X)
//	CC01_ZeroPage  AccesssMode = 0b00000100
//	CC01_Immediate AccesssMode = 0b00001000
//	CC01_Absolute  AccesssMode = 0b00001100
//	CC01_IndirectY AccesssMode = 0b00010000 // (zero page),Y
//	CC01_ZeroPageX AccesssMode = 0b00010100
//	CC01_AbsoluteY AccesssMode = 0b00011000
//	CC01_AbsoluteX AccesssMode = 0b00011100
//)

//const (
//	CC10_Immediate   AccesssMode = 0b00000000
//	CC10_ZeroPage    AccesssMode = 0b00000100
//	CC10_Accumulator AccesssMode = 0b00001000
//	CC10_Absolute    AccesssMode = 0b00001100
//	CC10_ZeroPageX   AccesssMode = 0b00010100 // ZeroPageY for STX and LDX
//	CC10_AbsoluteX   AccesssMode = 0b00011100 // AbsoluteY for STX and LDX
//)

const (
	CC00_Immediate AccesssMode = 0b00000000
	CC00_ZeroPage  AccesssMode = 0b00000100
	CC00_Absolute  AccesssMode = 0b00001100
	CC00_ZeroPageX AccesssMode = 0b00010100
	CC00_AbsoluteX AccesssMode = 0b00011100
)
