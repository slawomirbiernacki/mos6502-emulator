package main

func main() {

	cpu := Cpu{}

	//cpu.Load("roms/lda-clc.bin")
	//err := cpu.Load("roms/6502_functional_test.bin", 0x0, 0x0400)
	err := cpu.Load("roms/6502_functional_test.bin", 0x0, 0x16e2)
	if err != nil {
		panic(err)
	}

	for true {
		cpu.Cycle()
	}

}
