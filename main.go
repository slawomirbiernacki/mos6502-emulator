package main

func main() {

	cpu := NewCpu()

	err := cpu.Load("roms/functional_test/6502_functional_test_no_decimal.bin", 0x0, 0x0400)
	if err != nil {
		panic(err)
	}

	for true {
		cpu.Cycle()
	}

}
