package main

import (
	"mos6502-emulator/memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCpu_adc1(t *testing.T) {

	cpu := NewCpu(nil, &memory.DummyMemoryMapper{})
	cpu.Reset()

	cpu.A = 0b11111111
	cpu.C = 0
	cpu.adc(1)
	assert.Equal(t, byte(0), cpu.A)
	assert.Equal(t, byte(1), cpu.C)
	assert.Equal(t, byte(0), cpu.V)
	assert.Equal(t, byte(1), cpu.Z)
	assert.Equal(t, byte(0), cpu.N)
}
