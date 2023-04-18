package main

import (
	"mos6502-emulator/memory"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_cpu_functional(t *testing.T) {

	cpu := NewCpu(nil, &memory.DummyMemoryMapper{})
	err := cpu.Load("roms/functional_test/6502_functional_test_no_decimal.bin", 0x0, 0x0400)
	require.NoError(t, err)
	start := time.Now()
	for true {
		if cpu.PC == 0x336D {
			return // success!
		}
		now := time.Now()
		timeTaken := now.Sub(start)
		if timeTaken > 10*time.Second {
			assert.FailNow(t, "Test hit a trap, sth went wrong 🪦💀🪦")
		}
		cpu.ExecuteOpcode()
	}
}

func Test_cpu_interruptions(t *testing.T) {

	// interruptChannel := make(chan InterruptType, 100)
	memporyMapper := InterruptTestMemoryMapper{}
	cpu := NewCpu(nil, &memporyMapper)

	memporyMapper.cpu = &cpu

	// bytes, err := utils.ReadMemoryFromGzipFile("roms/6502_interrupt_test.bin.gz")
	// if err != nil {
	// 	panic(err)
	// }

	err := cpu.Load("roms/6502_interrupt_test.bin", 0x0, 0x0800)
	// err := cpu.Load("roms/6502_interrupt_test.bin", 0x0, 0x0a93)
	memporyMapper.interruptsOn = true
	require.NoError(t, err)
	start := time.Now()
	for true {
		if cpu.PC == 0x0af5 {
			return // success!
		}
		now := time.Now()
		timeTaken := now.Sub(start)
		if timeTaken > 10*time.Second {
			// assert.FailNow(t, "Test hit a trap, sth went wrong 🪦💀🪦")
		}
		cpu.ExecuteOpcode()
	}
}

type InterruptTestMemoryMapper struct {
	Mem              [65536]byte
	InterruptChannel chan InterruptType
	interruptsOn     bool
	cpu              *Cpu
}

func (m *InterruptTestMemoryMapper) Read(address uint16) byte {
	return m.Mem[address]
}

func (m *InterruptTestMemoryMapper) Write(address uint16, value byte) {

	if m.interruptsOn && address == 0xbffc {
		oldValue := m.Read(address)
		m.TriggerInterrupt(oldValue, value)
		m.Mem[address] = value
		return
	}

	m.Mem[address] = value
}

func (m *InterruptTestMemoryMapper) TriggerInterrupt(oldValue uint8, value uint8) {
	oldInterrupt := (oldValue & 0x1) == 0x1
	oldNMI := (oldValue & 0x2) == 0x2

	interrupt := (value & 0x1) == 0x1
	NMI := (value & 0x2) == 0x2

	if oldInterrupt != interrupt {
		// m.InterruptChannel <- InterruptTypeIRQ
		m.cpu.PendingIRQ = interrupt
	}

	if (oldNMI != NMI) && NMI {
		// m.InterruptChannel <- InterruptTypeNMI
		m.cpu.PendingNMI = NMI
	}
}
