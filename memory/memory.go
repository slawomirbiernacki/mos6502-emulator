package memory

type Memory struct {
	// ($0 - $FFFF in hex)
	RAM [65536]byte
}

type AccesssMode int

const (
	Immediate AccesssMode = iota
	ZeroPage
	ZeroPageX
	Absolute
	AbsoluteX
	AbsoluteY
	IndirectX
	IndirectY
)
