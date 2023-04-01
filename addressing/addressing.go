package addressing

type Mode int

const (
	ZeroPage Mode = iota
	Immediate
	Implied
	Relative // branches
	Absolute
	Indirect  // only used in JMP
	IndirectX // (zero page,X)
	IndirectY // (zero page),Y
	ZeroPageX
	ZeroPageY
	AbsoluteY
	AbsoluteX
	Accumulator
)
