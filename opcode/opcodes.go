package opcode

import "mos6502-emulator/addressing"

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
	AccessMode addressing.Mode
	Cycles     int
}

var mapping = map[byte]OpcodeSpec{

	0x09: {Operation: ORA, AccessMode: addressing.Immediate, Cycles: 2},
	0x05: {Operation: ORA, AccessMode: addressing.ZeroPage, Cycles: 2},
	0x15: {Operation: ORA, AccessMode: addressing.ZeroPageX, Cycles: 3},
	0x0D: {Operation: ORA, AccessMode: addressing.Absolute, Cycles: 4},
	0x1D: {Operation: ORA, AccessMode: addressing.AbsoluteX, Cycles: 4},
	0x19: {Operation: ORA, AccessMode: addressing.AbsoluteY, Cycles: 4},
	0x01: {Operation: ORA, AccessMode: addressing.IndirectX, Cycles: 6},
	0x11: {Operation: ORA, AccessMode: addressing.IndirectY, Cycles: 5},

	0x29: {Operation: AND, AccessMode: addressing.Immediate, Cycles: 2},
	0x25: {Operation: AND, AccessMode: addressing.ZeroPage, Cycles: 2},
	0x35: {Operation: AND, AccessMode: addressing.ZeroPageX, Cycles: 3},
	0x2D: {Operation: AND, AccessMode: addressing.Absolute, Cycles: 4},
	0x3D: {Operation: AND, AccessMode: addressing.AbsoluteX, Cycles: 4},
	0x39: {Operation: AND, AccessMode: addressing.AbsoluteY, Cycles: 4},
	0x21: {Operation: AND, AccessMode: addressing.IndirectX, Cycles: 6},
	0x31: {Operation: AND, AccessMode: addressing.IndirectY, Cycles: 5},

	0x49: {Operation: EOR, AccessMode: addressing.Immediate, Cycles: 2},
	0x45: {Operation: EOR, AccessMode: addressing.ZeroPage, Cycles: 3},
	0x55: {Operation: EOR, AccessMode: addressing.ZeroPageX, Cycles: 4},
	0x4D: {Operation: EOR, AccessMode: addressing.Absolute, Cycles: 4},
	0x5D: {Operation: EOR, AccessMode: addressing.AbsoluteX, Cycles: 4},
	0x59: {Operation: EOR, AccessMode: addressing.AbsoluteY, Cycles: 4},
	0x41: {Operation: EOR, AccessMode: addressing.IndirectX, Cycles: 6},
	0x51: {Operation: EOR, AccessMode: addressing.IndirectY, Cycles: 5},

	0x69: {Operation: ADC, AccessMode: addressing.Immediate, Cycles: 2},
	0x65: {Operation: ADC, AccessMode: addressing.ZeroPage, Cycles: 3},
	0x75: {Operation: ADC, AccessMode: addressing.ZeroPageX, Cycles: 4},
	0x6D: {Operation: ADC, AccessMode: addressing.Absolute, Cycles: 4},
	0x7D: {Operation: ADC, AccessMode: addressing.AbsoluteX, Cycles: 4},
	0x79: {Operation: ADC, AccessMode: addressing.AbsoluteY, Cycles: 4},
	0x61: {Operation: ADC, AccessMode: addressing.IndirectX, Cycles: 6},
	0x71: {Operation: ADC, AccessMode: addressing.IndirectY, Cycles: 5},

	0x85: {Operation: STA, AccessMode: addressing.ZeroPage, Cycles: 3},
	0x95: {Operation: STA, AccessMode: addressing.ZeroPageX, Cycles: 4},
	0x8D: {Operation: STA, AccessMode: addressing.Absolute, Cycles: 4},
	0x9D: {Operation: STA, AccessMode: addressing.AbsoluteX, Cycles: 5},
	0x99: {Operation: STA, AccessMode: addressing.AbsoluteY, Cycles: 5},
	0x81: {Operation: STA, AccessMode: addressing.IndirectX, Cycles: 6},
	0x91: {Operation: STA, AccessMode: addressing.IndirectY, Cycles: 6},

	0xA9: {Operation: LDA, AccessMode: addressing.Immediate, Cycles: 2},
	0xA5: {Operation: LDA, AccessMode: addressing.ZeroPage, Cycles: 3},
	0xB5: {Operation: LDA, AccessMode: addressing.ZeroPageX, Cycles: 4},
	0xAD: {Operation: LDA, AccessMode: addressing.Absolute, Cycles: 4},
	0xBD: {Operation: LDA, AccessMode: addressing.AbsoluteX, Cycles: 4},
	0xB9: {Operation: LDA, AccessMode: addressing.AbsoluteY, Cycles: 4},
	0xA1: {Operation: LDA, AccessMode: addressing.IndirectX, Cycles: 6},
	0xB1: {Operation: LDA, AccessMode: addressing.IndirectY, Cycles: 5},

	0xC9: {Operation: CMP, AccessMode: addressing.Immediate, Cycles: 2},
	0xC5: {Operation: CMP, AccessMode: addressing.ZeroPage, Cycles: 3},
	0xD5: {Operation: CMP, AccessMode: addressing.ZeroPageX, Cycles: 4},
	0xCD: {Operation: CMP, AccessMode: addressing.Absolute, Cycles: 4},
	0xDD: {Operation: CMP, AccessMode: addressing.AbsoluteX, Cycles: 4},
	0xD9: {Operation: CMP, AccessMode: addressing.AbsoluteY, Cycles: 4},
	0xC1: {Operation: CMP, AccessMode: addressing.IndirectX, Cycles: 6},
	0xD1: {Operation: CMP, AccessMode: addressing.IndirectY, Cycles: 5},

	0xE9: {Operation: SBC, AccessMode: addressing.Immediate, Cycles: 2},
	0xE5: {Operation: SBC, AccessMode: addressing.ZeroPage, Cycles: 3},
	0xF5: {Operation: SBC, AccessMode: addressing.ZeroPageX, Cycles: 4},
	0xED: {Operation: SBC, AccessMode: addressing.Absolute, Cycles: 4},
	0xFD: {Operation: SBC, AccessMode: addressing.AbsoluteX, Cycles: 4},
	0xF9: {Operation: SBC, AccessMode: addressing.AbsoluteY, Cycles: 4},
	0xE1: {Operation: SBC, AccessMode: addressing.IndirectX, Cycles: 6},
	0xF1: {Operation: SBC, AccessMode: addressing.IndirectY, Cycles: 5},

	0x0A: {Operation: ASL, AccessMode: addressing.Accumulator, Cycles: 2},
	0x06: {Operation: ASL, AccessMode: addressing.ZeroPage, Cycles: 5},
	0x16: {Operation: ASL, AccessMode: addressing.ZeroPageX, Cycles: 6},
	0x0E: {Operation: ASL, AccessMode: addressing.Absolute, Cycles: 6},
	0x1E: {Operation: ASL, AccessMode: addressing.AbsoluteX, Cycles: 7},

	0x2A: {Operation: ROL, AccessMode: addressing.Accumulator, Cycles: 2},
	0x26: {Operation: ROL, AccessMode: addressing.ZeroPage, Cycles: 5},
	0x36: {Operation: ROL, AccessMode: addressing.ZeroPageX, Cycles: 6},
	0x2E: {Operation: ROL, AccessMode: addressing.Absolute, Cycles: 6},
	0x3E: {Operation: ROL, AccessMode: addressing.AbsoluteX, Cycles: 7},

	0x4A: {Operation: LSR, AccessMode: addressing.Accumulator, Cycles: 2},
	0x46: {Operation: LSR, AccessMode: addressing.ZeroPage, Cycles: 5},
	0x56: {Operation: LSR, AccessMode: addressing.ZeroPageX, Cycles: 6},
	0x4E: {Operation: LSR, AccessMode: addressing.Absolute, Cycles: 6},
	0x5E: {Operation: LSR, AccessMode: addressing.AbsoluteX, Cycles: 7},

	0x6A: {Operation: ROR, AccessMode: addressing.Accumulator, Cycles: 2},
	0x66: {Operation: ROR, AccessMode: addressing.ZeroPage, Cycles: 5},
	0x76: {Operation: ROR, AccessMode: addressing.ZeroPageX, Cycles: 6},
	0x6E: {Operation: ROR, AccessMode: addressing.Absolute, Cycles: 6},
	0x7E: {Operation: ROR, AccessMode: addressing.AbsoluteX, Cycles: 7},

	0x86: {Operation: STX, AccessMode: addressing.ZeroPage, Cycles: 3},
	0x96: {Operation: STX, AccessMode: addressing.ZeroPageY, Cycles: 4},
	0x8E: {Operation: STX, AccessMode: addressing.Absolute, Cycles: 4},

	0xA2: {Operation: LDX, AccessMode: addressing.Immediate, Cycles: 2},
	0xA6: {Operation: LDX, AccessMode: addressing.ZeroPage, Cycles: 3},
	0xB6: {Operation: LDX, AccessMode: addressing.ZeroPageY, Cycles: 4},
	0xAE: {Operation: LDX, AccessMode: addressing.Absolute, Cycles: 4},
	0xBE: {Operation: LDX, AccessMode: addressing.AbsoluteY, Cycles: 4},

	0xC6: {Operation: DEC, AccessMode: addressing.ZeroPage, Cycles: 5},
	0xD6: {Operation: DEC, AccessMode: addressing.ZeroPageX, Cycles: 6},
	0xCE: {Operation: DEC, AccessMode: addressing.Absolute, Cycles: 6},
	0xDE: {Operation: DEC, AccessMode: addressing.AbsoluteX, Cycles: 7},

	0xE6: {Operation: INC, AccessMode: addressing.ZeroPage, Cycles: 5},
	0xF6: {Operation: INC, AccessMode: addressing.ZeroPageX, Cycles: 6},
	0xEE: {Operation: INC, AccessMode: addressing.Absolute, Cycles: 6},
	0xFE: {Operation: INC, AccessMode: addressing.AbsoluteX, Cycles: 7},

	0x24: {Operation: BIT, AccessMode: addressing.ZeroPage, Cycles: 3},
	0x2C: {Operation: BIT, AccessMode: addressing.Absolute, Cycles: 4},

	0x4C: {Operation: JMP, AccessMode: addressing.Absolute, Cycles: 3},
	0x6C: {Operation: JMP, AccessMode: addressing.Indirect, Cycles: 5},

	0x84: {Operation: STY, AccessMode: addressing.ZeroPage, Cycles: 3},
	0x94: {Operation: STY, AccessMode: addressing.ZeroPageX, Cycles: 4},
	0x8C: {Operation: STY, AccessMode: addressing.Absolute, Cycles: 4},

	0xA0: {Operation: LDY, AccessMode: addressing.Immediate, Cycles: 2},
	0xA4: {Operation: LDY, AccessMode: addressing.ZeroPage, Cycles: 3},
	0xB4: {Operation: LDY, AccessMode: addressing.ZeroPageX, Cycles: 4},
	0xAC: {Operation: LDY, AccessMode: addressing.Absolute, Cycles: 4},
	0xBC: {Operation: LDY, AccessMode: addressing.AbsoluteX, Cycles: 4},

	0xC0: {Operation: CPY, AccessMode: addressing.Immediate, Cycles: 2},
	0xC4: {Operation: CPY, AccessMode: addressing.ZeroPage, Cycles: 3},
	0xCC: {Operation: CPY, AccessMode: addressing.Absolute, Cycles: 4},

	0xE0: {Operation: CPX, AccessMode: addressing.Immediate, Cycles: 2},
	0xE4: {Operation: CPX, AccessMode: addressing.ZeroPage, Cycles: 3},
	0xEC: {Operation: CPX, AccessMode: addressing.Absolute, Cycles: 4},

	0x00: {Operation: BRK, AccessMode: addressing.Implied, Cycles: 7},
	0x20: {Operation: JSR, AccessMode: addressing.Absolute, Cycles: 6},
	0x40: {Operation: RTI, AccessMode: addressing.Implied, Cycles: 6},
	0x60: {Operation: RTS, AccessMode: addressing.Implied, Cycles: 6},
	0x08: {Operation: PHP, AccessMode: addressing.Implied, Cycles: 3},
	0x28: {Operation: PLP, AccessMode: addressing.Implied, Cycles: 4},
	0x48: {Operation: PHA, AccessMode: addressing.Implied, Cycles: 3},
	0x68: {Operation: PLA, AccessMode: addressing.Implied, Cycles: 4},
	0x88: {Operation: DEY, AccessMode: addressing.Implied, Cycles: 2},
	0xA8: {Operation: TAY, AccessMode: addressing.Implied, Cycles: 2},
	0xC8: {Operation: INY, AccessMode: addressing.Implied, Cycles: 2},
	0xE8: {Operation: INX, AccessMode: addressing.Implied, Cycles: 2},
	0x18: {Operation: CLC, AccessMode: addressing.Implied, Cycles: 2},
	0x38: {Operation: SEC, AccessMode: addressing.Implied, Cycles: 2},
	0x58: {Operation: CLI, AccessMode: addressing.Implied, Cycles: 2},
	0x78: {Operation: SEI, AccessMode: addressing.Implied, Cycles: 2},
	0x98: {Operation: TYA, AccessMode: addressing.Implied, Cycles: 2},
	0xB8: {Operation: CLV, AccessMode: addressing.Implied, Cycles: 2},
	0xD8: {Operation: CLD, AccessMode: addressing.Implied, Cycles: 2},
	0xF8: {Operation: SED, AccessMode: addressing.Implied, Cycles: 2},
	0x8A: {Operation: TXA, AccessMode: addressing.Implied, Cycles: 2},
	0x9A: {Operation: TXS, AccessMode: addressing.Implied, Cycles: 2},
	0xAA: {Operation: TAX, AccessMode: addressing.Implied, Cycles: 2},
	0xBA: {Operation: TSX, AccessMode: addressing.Implied, Cycles: 2},
	0xCA: {Operation: DEX, AccessMode: addressing.Implied, Cycles: 2},
	0xEA: {Operation: NOP, AccessMode: addressing.Implied, Cycles: 2},

	// branches
	0x90: {Operation: BCC, AccessMode: addressing.Relative, Cycles: 2},
	0xB0: {Operation: BCS, AccessMode: addressing.Relative, Cycles: 2},
	0xF0: {Operation: BEQ, AccessMode: addressing.Relative, Cycles: 2},
	0x30: {Operation: BMI, AccessMode: addressing.Relative, Cycles: 2},
	0xD0: {Operation: BNE, AccessMode: addressing.Relative, Cycles: 2},
	0x10: {Operation: BPL, AccessMode: addressing.Relative, Cycles: 2},
	0x50: {Operation: BVC, AccessMode: addressing.Relative, Cycles: 2},
	0x70: {Operation: BVS, AccessMode: addressing.Relative, Cycles: 2},
}

func Lookup(opcode byte) OpcodeSpec {
	return mapping[opcode]
}
