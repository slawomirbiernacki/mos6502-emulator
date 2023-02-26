// Confidential and only for the information of the intended recipient.
// Copyright 2023 Improbable Worlds Limited.

package addressing

type AccessMode int

const (
	ZeroPage AccessMode = iota
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
