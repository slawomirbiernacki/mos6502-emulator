package opcode

import "mos6502-emulator/addressing"

//const (
//	// aaa part of opcodes, cc = 01
//
//	ORA = 0b000
//	AND = 0b001
//	EOR = 0b010
//	ADC = 0b011
//	STA = 0b100
//	LDA = 0b101
//	CMP = 0b110
//	SBC = 0b111
//
//	// aaa part of opcodes, cc = 10
//
//	ASL = 0b000
//	ROL = 0b001
//	LSR = 0b010
//	ROR = 0b011
//	STX = 0b100
//	LDX = 0b101
//	DEC = 0b110
//	INC = 0b111
//
//	// cc = 00
//	BIT     = 0b001
//	JMP_ABS = 0b010
//	JMP     = 0b011
//	STY     = 0b100
//	LDY     = 0b101
//	CPY     = 0b110
//	CPX     = 0b111
//
//	BRK = 0x00
//	JSR = 0x20
//	RTI = 0x40
//	RTS = 0x60
//
//	PHP = 0x08
//	PLP = 0x28
//	PHA = 0x48
//	PLA = 0x68
//	DEY = 0x88
//	TAY = 0xA8
//	INY = 0xC8
//	INX = 0xE8
//
//	CLC = 0x18 // CLear Carry
//	SEC = 0x38 // SEt Carry)
//	CLI = 0x58 // CLear Interrupt
//	SEI = 0x78 // SEt Interrupt
//	TYA = 0x98
//	CLV = 0xB8 // CLear oVerflow
//	CLD = 0xD8 // CLear Decimal
//	SED = 0xF8 // SEt Decimal
//
//	TXA = 0x8A
//	TXS = 0x9A
//	TAX = 0xAA
//	TSX = 0xBA
//	DEX = 0xCA
//	NOP = 0xEA
//)
//
//const (
//	// ADd with Carry
//	ADC_IMM  = 0x69
//	ADC_ZP   = 0x65
//	ADC_ZPX  = 0x75
//	ADC_ABS  = 0x6D
//	ADC_ABSX = 0x7D
//	ADC_ABSY = 0x79
//	ADC_INX  = 0x61
//	ADC_INY  = 0x71
//
//	AND_IMM  = 0x29
//	AND_ZP   = 0x25
//	AND_ZPX  = 0x35
//	AND_ABS  = 0x2D
//	AND_ABSX = 0x3D
//	AND_ABSY = 0x39
//	AND_INX  = 0x21
//	AND_INY  = 0x31
//
//	LDA_IMM  = 0xA9
//	LDA_ZP   = 0xA5
//	LDA_ZPX  = 0xB5
//	LDA_ABS  = 0xAD
//	LDA_ABSX = 0xBD
//	LDA_ABSY = 0xB9
//	LDA_INX  = 0xA1
//	LDA_INY  = 0xB1
//)

type Operation int

const (
	ORA Operation = iota
	AND
	EOR
	ADC
	STA
	LDA
	CMP
	SBC
	ASL
	ROL
	LSR
	ROR
	STX
	LDX
	DEC
	INC
	BIT
	JMP
	STY
	LDY
	CPY
	CPX
	BRK
	JSR
	RTI
	RTS
	PHP
	PLP
	PHA
	PLA
	DEY
	TAY
	INY
	INX
	CLC // CLear Carry
	SEC // SEt Carry)
	CLI // CLear Interrupt
	SEI // SEt Interrupt
	TYA
	CLV // CLear oVerflow
	CLD // CLear Decimal
	SED // SEt Decimal
	TXA
	TXS
	TAX
	TSX
	DEX
	NOP

	BCC
	BCS
	BEQ
	BMI
	BNE
	BPL
	BVC
	BVS
)

type OpcodeSpec struct {
	Operation  Operation
	AccessMode addressing.AccesssMode
}

var mapping = map[byte]OpcodeSpec{

	0x09: {Operation: ORA, AccessMode: addressing.Immediate},
	0x05: {Operation: ORA, AccessMode: addressing.ZeroPage},
	0x15: {Operation: ORA, AccessMode: addressing.ZeroPageX},
	0x0D: {Operation: ORA, AccessMode: addressing.Absolute},
	0x1D: {Operation: ORA, AccessMode: addressing.AbsoluteX},
	0x19: {Operation: ORA, AccessMode: addressing.AbsoluteY},
	0x01: {Operation: ORA, AccessMode: addressing.IndirectX},
	0x11: {Operation: ORA, AccessMode: addressing.IndirectY},

	0x29: {Operation: AND, AccessMode: addressing.Immediate},
	0x25: {Operation: AND, AccessMode: addressing.ZeroPage},
	0x35: {Operation: AND, AccessMode: addressing.ZeroPageX},
	0x2D: {Operation: AND, AccessMode: addressing.Absolute},
	0x3D: {Operation: AND, AccessMode: addressing.AbsoluteX},
	0x39: {Operation: AND, AccessMode: addressing.AbsoluteY},
	0x21: {Operation: AND, AccessMode: addressing.IndirectX},
	0x31: {Operation: AND, AccessMode: addressing.IndirectY},

	0x49: {Operation: EOR, AccessMode: addressing.Immediate},
	0x45: {Operation: EOR, AccessMode: addressing.ZeroPage},
	0x55: {Operation: EOR, AccessMode: addressing.ZeroPageX},
	0x4D: {Operation: EOR, AccessMode: addressing.Absolute},
	0x5D: {Operation: EOR, AccessMode: addressing.AbsoluteX},
	0x59: {Operation: EOR, AccessMode: addressing.AbsoluteY},
	0x41: {Operation: EOR, AccessMode: addressing.IndirectX},
	0x51: {Operation: EOR, AccessMode: addressing.IndirectY},

	0x69: {Operation: ADC, AccessMode: addressing.Immediate},
	0x65: {Operation: ADC, AccessMode: addressing.ZeroPage},
	0x75: {Operation: ADC, AccessMode: addressing.ZeroPageX},
	0x6D: {Operation: ADC, AccessMode: addressing.Absolute},
	0x7D: {Operation: ADC, AccessMode: addressing.AbsoluteX},
	0x79: {Operation: ADC, AccessMode: addressing.AbsoluteY},
	0x61: {Operation: ADC, AccessMode: addressing.IndirectX},
	0x71: {Operation: ADC, AccessMode: addressing.IndirectY},

	0x85: {Operation: STA, AccessMode: addressing.ZeroPage},
	0x95: {Operation: STA, AccessMode: addressing.ZeroPageX},
	0x8D: {Operation: STA, AccessMode: addressing.Absolute},
	0x9D: {Operation: STA, AccessMode: addressing.AbsoluteX},
	0x99: {Operation: STA, AccessMode: addressing.AbsoluteY},
	0x81: {Operation: STA, AccessMode: addressing.IndirectX},
	0x91: {Operation: STA, AccessMode: addressing.IndirectY},

	0xA9: {Operation: LDA, AccessMode: addressing.Immediate},
	0xA5: {Operation: LDA, AccessMode: addressing.ZeroPage},
	0xB5: {Operation: LDA, AccessMode: addressing.ZeroPageX},
	0xAD: {Operation: LDA, AccessMode: addressing.Absolute},
	0xBD: {Operation: LDA, AccessMode: addressing.AbsoluteX},
	0xB9: {Operation: LDA, AccessMode: addressing.AbsoluteY},
	0xA1: {Operation: LDA, AccessMode: addressing.IndirectX},
	0xB1: {Operation: LDA, AccessMode: addressing.IndirectY},

	0xC9: {Operation: CMP, AccessMode: addressing.Immediate},
	0xC5: {Operation: CMP, AccessMode: addressing.ZeroPage},
	0xD5: {Operation: CMP, AccessMode: addressing.ZeroPageX},
	0xCD: {Operation: CMP, AccessMode: addressing.Absolute},
	0xDD: {Operation: CMP, AccessMode: addressing.AbsoluteX},
	0xD9: {Operation: CMP, AccessMode: addressing.AbsoluteY},
	0xC1: {Operation: CMP, AccessMode: addressing.IndirectX},
	0xD1: {Operation: CMP, AccessMode: addressing.IndirectY},

	0xE9: {Operation: SBC, AccessMode: addressing.Immediate},
	0xE5: {Operation: SBC, AccessMode: addressing.ZeroPage},
	0xF5: {Operation: SBC, AccessMode: addressing.ZeroPageX},
	0xED: {Operation: SBC, AccessMode: addressing.Absolute},
	0xFD: {Operation: SBC, AccessMode: addressing.AbsoluteX},
	0xF9: {Operation: SBC, AccessMode: addressing.AbsoluteY},
	0xE1: {Operation: SBC, AccessMode: addressing.IndirectX},
	0xF1: {Operation: SBC, AccessMode: addressing.IndirectY},

	0x0A: {Operation: ASL, AccessMode: addressing.Accumulator},
	0x06: {Operation: ASL, AccessMode: addressing.ZeroPage},
	0x16: {Operation: ASL, AccessMode: addressing.ZeroPageX},
	0x0E: {Operation: ASL, AccessMode: addressing.Absolute},
	0x1E: {Operation: ASL, AccessMode: addressing.AbsoluteX},

	0x2A: {Operation: ROL, AccessMode: addressing.Accumulator},
	0x26: {Operation: ROL, AccessMode: addressing.ZeroPage},
	0x36: {Operation: ROL, AccessMode: addressing.ZeroPageX},
	0x2E: {Operation: ROL, AccessMode: addressing.Absolute},
	0x3E: {Operation: ROL, AccessMode: addressing.AbsoluteX},

	0x4A: {Operation: LSR, AccessMode: addressing.Accumulator},
	0x46: {Operation: LSR, AccessMode: addressing.ZeroPage},
	0x56: {Operation: LSR, AccessMode: addressing.ZeroPageX},
	0x4E: {Operation: LSR, AccessMode: addressing.Absolute},
	0x5E: {Operation: LSR, AccessMode: addressing.AbsoluteX},

	0x6A: {Operation: ROR, AccessMode: addressing.Accumulator},
	0x66: {Operation: ROR, AccessMode: addressing.ZeroPage},
	0x76: {Operation: ROR, AccessMode: addressing.ZeroPageX},
	0x6E: {Operation: ROR, AccessMode: addressing.Absolute},
	0x7E: {Operation: ROR, AccessMode: addressing.AbsoluteX},

	0x86: {Operation: STX, AccessMode: addressing.ZeroPage},
	0x96: {Operation: STX, AccessMode: addressing.ZeroPageY},
	0x8E: {Operation: STX, AccessMode: addressing.Absolute},

	0xA2: {Operation: LDX, AccessMode: addressing.Immediate},
	0xA6: {Operation: LDX, AccessMode: addressing.ZeroPage},
	0xB6: {Operation: LDX, AccessMode: addressing.ZeroPageY},
	0xAE: {Operation: LDX, AccessMode: addressing.Absolute},
	0xBE: {Operation: LDX, AccessMode: addressing.AbsoluteY},

	0xC6: {Operation: DEC, AccessMode: addressing.ZeroPage},
	0xD6: {Operation: DEC, AccessMode: addressing.ZeroPageX},
	0xCE: {Operation: DEC, AccessMode: addressing.Absolute},
	0xDE: {Operation: DEC, AccessMode: addressing.AbsoluteX},

	0xE6: {Operation: INC, AccessMode: addressing.ZeroPage},
	0xF6: {Operation: INC, AccessMode: addressing.ZeroPageX},
	0xEE: {Operation: INC, AccessMode: addressing.Absolute},
	0xFE: {Operation: INC, AccessMode: addressing.AbsoluteX},

	0x24: {Operation: BIT, AccessMode: addressing.ZeroPage},
	0x2C: {Operation: BIT, AccessMode: addressing.Absolute},

	0x4C: {Operation: JMP, AccessMode: addressing.Absolute},
	0x6C: {Operation: JMP, AccessMode: addressing.Indirect},

	0x84: {Operation: STY, AccessMode: addressing.ZeroPage},
	0x94: {Operation: STY, AccessMode: addressing.ZeroPageX},
	0x8C: {Operation: STY, AccessMode: addressing.Absolute},

	0xA0: {Operation: LDY, AccessMode: addressing.Immediate},
	0xA4: {Operation: LDY, AccessMode: addressing.ZeroPage},
	0xB4: {Operation: LDY, AccessMode: addressing.ZeroPageX},
	0xAC: {Operation: LDY, AccessMode: addressing.Absolute},
	0xBC: {Operation: LDY, AccessMode: addressing.AbsoluteX},

	0xC0: {Operation: CPY, AccessMode: addressing.Immediate},
	0xC4: {Operation: CPY, AccessMode: addressing.ZeroPage},
	0xCC: {Operation: CPY, AccessMode: addressing.Absolute},

	0xE0: {Operation: CPX, AccessMode: addressing.Immediate},
	0xE4: {Operation: CPX, AccessMode: addressing.ZeroPage},
	0xEC: {Operation: CPX, AccessMode: addressing.Absolute},

	0x00: {Operation: BRK, AccessMode: addressing.Implied},
	0x20: {Operation: JSR, AccessMode: addressing.Absolute},
	0x40: {Operation: RTI, AccessMode: addressing.Implied},
	0x60: {Operation: RTS, AccessMode: addressing.Implied},
	0x08: {Operation: PHP, AccessMode: addressing.Implied},
	0x28: {Operation: PLP, AccessMode: addressing.Implied},
	0x48: {Operation: PHA, AccessMode: addressing.Implied},
	0x68: {Operation: PLA, AccessMode: addressing.Implied},
	0x88: {Operation: DEY, AccessMode: addressing.Implied},
	0xA8: {Operation: TAY, AccessMode: addressing.Implied},
	0xC8: {Operation: INY, AccessMode: addressing.Implied},
	0xE8: {Operation: INX, AccessMode: addressing.Implied},
	0x18: {Operation: CLC, AccessMode: addressing.Implied},
	0x38: {Operation: SEC, AccessMode: addressing.Implied},
	0x58: {Operation: CLI, AccessMode: addressing.Implied},
	0x78: {Operation: SEI, AccessMode: addressing.Implied},
	0x98: {Operation: TYA, AccessMode: addressing.Implied},
	0xB8: {Operation: CLV, AccessMode: addressing.Implied},
	0xD8: {Operation: CLD, AccessMode: addressing.Implied},
	0xF8: {Operation: SED, AccessMode: addressing.Implied},
	0x8A: {Operation: TXA, AccessMode: addressing.Implied},
	0x9A: {Operation: TXS, AccessMode: addressing.Implied},
	0xAA: {Operation: TAX, AccessMode: addressing.Implied},
	0xBA: {Operation: TSX, AccessMode: addressing.Implied},
	0xCA: {Operation: DEX, AccessMode: addressing.Implied},
	0xEA: {Operation: NOP, AccessMode: addressing.Implied},

	// branches
	0x90: {Operation: BCC, AccessMode: addressing.Relative},
	0xB0: {Operation: BCS, AccessMode: addressing.Relative},
	0xF0: {Operation: BEQ, AccessMode: addressing.Relative},
	0x30: {Operation: BMI, AccessMode: addressing.Relative},
	0xD0: {Operation: BNE, AccessMode: addressing.Relative},
	0x10: {Operation: BPL, AccessMode: addressing.Relative},
	0x50: {Operation: BVC, AccessMode: addressing.Relative},
	0x70: {Operation: BVS, AccessMode: addressing.Relative},
}

func Lookup(opcode byte) OpcodeSpec {
	return mapping[opcode]
}
