// Confidential and only for the information of the intended recipient.
// Copyright 2023 Improbable Worlds Limited.

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
