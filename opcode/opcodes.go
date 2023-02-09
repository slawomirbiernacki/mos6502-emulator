package opcode

const (
	// aaa part of opcodes, cc = 01

	ORA = 0b000
	AND = 0b001
	EOR = 0b010
	ADC = 0b011
	STA = 0b100
	LDA = 0b101
	CMP = 0b110
	SBC = 0b111

	// aaa part of opcodes, cc = 10

	ASL = 0b000
	ROL = 0b001
	LSR = 0b010
	ROR = 0b011
	STX = 0b100
	LDX = 0b101
	DEC = 0b110
	INC = 0b111

	// cc = 00
	BIT     = 0b001
	JMP_ABS = 0b010
	JMP     = 0b011
	STY     = 0b100
	LDY     = 0b101
	CPY     = 0b110
	CPX     = 0b111

	BRK = 0x00
	JSR = 0x20
	RTI = 0x40
	RTS = 0x60

	PHP = 0x08
	PLP = 0x28
	PHA = 0x48
	PLA = 0x68
	DEY = 0x88
	TAY = 0xA8
	INY = 0xC8
	INX = 0xE8

	CLC = 0x18 // CLear Carry
	SEC = 0x38 // SEt Carry)
	CLI = 0x58 // CLear Interrupt
	SEI = 0x78 // SEt Interrupt
	TYA = 0x98
	CLV = 0xB8 // CLear oVerflow
	CLD = 0xD8 // CLear Decimal
	SED = 0xF8 // SEt Decimal

	TXA = 0x8A
	TXS = 0x9A
	TAX = 0xAA
	TSX = 0xBA
	DEX = 0xCA
	NOP = 0xEA
)

const (
	// ADd with Carry
	ADC_IMM  = 0x69
	ADC_ZP   = 0x65
	ADC_ZPX  = 0x75
	ADC_ABS  = 0x6D
	ADC_ABSX = 0x7D
	ADC_ABSY = 0x79
	ADC_INX  = 0x61
	ADC_INY  = 0x71

	AND_IMM  = 0x29
	AND_ZP   = 0x25
	AND_ZPX  = 0x35
	AND_ABS  = 0x2D
	AND_ABSX = 0x3D
	AND_ABSY = 0x39
	AND_INX  = 0x21
	AND_INY  = 0x31

	LDA_IMM  = 0xA9
	LDA_ZP   = 0xA5
	LDA_ZPX  = 0xB5
	LDA_ABS  = 0xAD
	LDA_ABSX = 0xBD
	LDA_ABSY = 0xB9
	LDA_INX  = 0xA1
	LDA_INY  = 0xB1
)
