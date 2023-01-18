package main

// Little Endian architecture

type Cpu struct {
	A  byte
	X  byte
	Y  byte
	S  byte // 1 + 7 bits     , stack pointer
	PC uint16
	P  byte // status
	// flags
	ZeroFlag      int
	SignFlag      int
	OverflowFlag  int
	BreakFlag     int
	DecimalFlag   int
	InterruptFlag int
	CarryFlag     int

	Mem [65536]byte
	Clk int
}

func (c *Cpu) Reset() {
	c.P = 0x20
	c.ZeroFlag = 0
	c.SignFlag = 0
	c.OverflowFlag = 0
	c.BreakFlag = 0
	c.DecimalFlag = 0
	c.InterruptFlag = 0
	c.CarryFlag = 0
	c.S = 0xFF
	//The 6502 stores addresses in low byte/hi byte format, so $FFFD contains the upper 8 bits of the
	//address and $FFFC the lower 8 bits. This line assembles the 2 bytes into a 16-bit address.
	c.PC = uint16(c.Mem[0xFFFD])<<8 | uint16(c.Mem[0xFFFC])

	c.A = 0
	c.X = 0
	c.Y = 0
}

func (c *Cpu) Cycle() {

	switch c.PC {

	}
}

func GetOpCode()

// 8-bit stack pointer (fixed at RAM address $100, so can address $100-$1ff)

//The Stack Pointer(SP)is used to keep track of the current position of the stack.
//For example the stack on the 6502 is at memory locations $1ff-$100, it starts at
//$1ff and works it's way down towards $100. The stack pointer is 8 bits wide so
//it would start out at $ff (the processor knows it really means $1ff). When a
//value is pushed onto the stack it will be put at memory location $1ff and then
//the SP will be de-incremented to it points to $1fe. When data is pulled of the
//stack, the SP is incremented, then the data is read from that memory location.

func (c *Cpu) PushToStack(value byte) {

}
