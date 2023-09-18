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

func TestSimpleLea(t *testing.T) {
	t.Parallel()

	computer := NewLittleComputer3()

	computer.executeInstruction(0b1110001111111101)
	computer.executeInstruction(0b0001010001101110)
	computer.executeInstruction(0b0011010111111011)
	computer.executeInstruction(0b0101010010100000)
	computer.executeInstruction(0b0001010010100101)
	computer.executeInstruction(0b0111010001001110)
	computer.executeInstruction(0b1010011111110111)
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

func TestDecodePcRelative(t *testing.T) {
	t.Parallel()

	computer := NewLittleComputer3()
	computer.executeInstruction(0b1110011110101111)

	// TODO: sign extend the binary numbers, they are in 2's complement
	// cases := []struct {
	// 	description string
	// 	instruction uint64
	// 	expected    BaseRelativeInstruction
	// }{
	// 	{
	// 		description: "should decode LEA",
	// 		instruction: ,
	// 	}
	// }

	// for _, tt := range cases {
	// 	decodedInstruction := decodeBaseRelative(uint16(tt.instruction))

	// 	assert.Equal(t, tt.expected, decodedInstruction)
	// }
}
