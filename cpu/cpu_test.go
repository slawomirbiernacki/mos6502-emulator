package cpu

import (
	"testing"
	"time"

	"github.com/slawomirbiernacki/mos6502-emulator/memory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_cpu(t *testing.T) {

	cpu := NewCpu(nil, &memory.DummyMemoryMapper{})
	err := cpu.Load("../roms/functional_test/6502_functional_test_no_decimal.bin", 0x0, 0x0400)
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
