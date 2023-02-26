// Confidential and only for the information of the intended recipient.
// Copyright 2023 Improbable Worlds Limited.

package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_cpu(t *testing.T) {

	cpu := NewCpu()
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
		cpu.Cycle()
	}
}
