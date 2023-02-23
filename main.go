package main

func main() {

	cpu := Cpu{}

	//cpu.Load("roms/lda-clc.bin")
	err := cpu.Load("roms/ca65/6502_functional_test_no_decimal.bin", 0x0, 0x0400)
	//err := cpu.Load("roms/ca65/6502_functional_test_no_decimal.bin", 0x0, 0x3308)
	//err := cpu.Load("roms/ca65/6502_functional_test_no_decimal.bin", 0x0, 0x3373)
	if err != nil {
		panic(err)
	}

	for true {
		cpu.Cycle()
	}

}
