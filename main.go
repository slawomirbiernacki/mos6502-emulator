package main

func main() {

	cpu := Cpu{}

	//cpu.Load("roms/lda-clc.bin")
	err := cpu.Load("roms/nestest.nes")
	if err != nil {
		panic(err)
	}
	cpu.Reset()
	cpu.Cycle()
	//
	//print("A")
	//test:= 0xA9
	//
	//masked:= (test & 0b00011100) >>2
	//
	//fmt.Printf("%b\n",  0b10000000>>7)
}
