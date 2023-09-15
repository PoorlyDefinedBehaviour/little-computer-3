package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	t.Parallel()

	computer := NewLittleComputer3()
	computer.registers[0] = 10
	computer.registers[1] = 20
	computer.executeInstruction(0b0001_000_000_0_00_001)
	assert.EqualValues(t, 30, computer.registers[0])
}

func TestDecodeOperateInstruction(t *testing.T) {
	t.Parallel()

	cases := []struct {
		instruction uint16
		expected    DecodedInstruction
	}{
		{
			instruction: 0b0001_001_010_000_001,
			expected: DecodedInstruction{
				Dst:          1,
				Src1:         2,
				RegisterMode: true,
				Src2:         1,
			},
		},
		{
			instruction: 0b0001_001_010_100_010,
			expected: DecodedInstruction{
				Dst:          1,
				Src1:         2,
				RegisterMode: false,
				Src2:         2,
			},
		},
	}

	for _, tt := range cases {
		instruction := decodeOperateInstruction(tt.instruction)
		fmt.Printf("instruction.Dst %b\n", instruction.Dst)
		fmt.Printf("instruction.Src1 %b\n", instruction.Src1)
		fmt.Printf("instruction.Src2 %b\n", instruction.Src2)
		assert.Equal(t, tt.expected, instruction)
	}
}
