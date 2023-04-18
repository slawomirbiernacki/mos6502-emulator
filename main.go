package main

import "github.com/slawomirbiernacki/mos6502-emulator/memory"

func main() {

	cpu := NewCpu(nil, &memory.DummyMemoryMapper{})

	err := cpu.Load("roms/functional_test/6502_functional_test_no_decimal.bin", 0x0, 0x0400)
	if err != nil {
		panic(err)
	}

	// for true {
	// 	cpu.ExecuteOpcode()
	// }
	cpu.Run(10000)
	print("finished!")

}
