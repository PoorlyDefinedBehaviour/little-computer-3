package main

import "fmt"

type LittleComputer3 struct {
	memory [65536]int16
	// Registers from r0 to r7.
	registers [8]int16
	pc        int
}

type OpCode uint8

const (
	OpCodeAdd OpCode = 0b0001
	OpCodeAnd OpCode = 0b0101
	OpCodeNot OpCode = 0b1001

	// Load -- read data from memory to register
	// PC-relative mode
	OpCodeLd OpCode = 0b0010
	// Base + offset mode
	OpCodeLdr OpCode = 0b0110
	// Indirect mode
	OpCodeLdi OpCode = 0b1010

	// Store -- write data from register to memory
	// PC-relative mode
	OpCodeSt OpCode = 0b0011
	// Base + offset mode
	OpCodeStr OpCode = 0b0111
	// Indirect mode
	OpCodeSti OpCode = 0b1011

	// Load effective address - computer address, save in register
	// Immediate mode
	OpCodeLea OpCode = 0b1110
)

func (opCode OpCode) String() string {
	switch opCode {
	case OpCodeAdd:
		return "ADD"
	case OpCodeAnd:
		return "AND"
	case OpCodeNot:
		return "NOT"
	case OpCodeLd:
		return "LD"
	case OpCodeSt:
		return "ST"
	case OpCodeLdi:
		return "LDI"
	case OpCodeSti:
		return "STI"
	case OpCodeLdr:
		return "LDR"
	case OpCodeStr:
		return "STR"
	case OpCodeLea:
		return "LEA"
	default:
		panic(fmt.Sprintf("unexpected opcode: %b", opCode))
	}
}

type DecodedInstruction struct {
	Dst          uint16
	Src1         uint16
	RegisterMode bool
	Src2         uint16
}

type PcRelativeInstruction struct {
	Register uint16
	PcOffset int16
}

type BaseRelativeInstruction struct {
	Register uint8
	Base     uint8
	Offset   int8
}

func NewLittleComputer3() *LittleComputer3 {
	return &LittleComputer3{
		memory:    [65536]int16{},
		registers: [8]int16{},
		pc:        0,
	}
}

func (computer *LittleComputer3) executeInstruction(instruction uint16) {
	opcode := getOpcode(instruction)

	switch opcode {
	case OpCodeAdd:
		inst := decodeOperateInstruction(instruction)

		a := computer.registers[inst.Src1]

		var b int16
		if inst.RegisterMode {
			b = computer.registers[inst.Src2]
		} else {
			b = int16(inst.Src2)
		}

		computer.registers[inst.Dst] = a + b

	case OpCodeAnd:
		inst := decodeOperateInstruction(instruction)

		a := computer.registers[inst.Src1]

		var b int16
		if inst.RegisterMode {
			b = computer.registers[inst.Src2]
		} else {
			b = int16(inst.Src2)
		}

		computer.registers[inst.Dst] = a & b

	case OpCodeNot:
		inst := decodeOperateInstruction(instruction)

		computer.registers[inst.Dst] = ^computer.registers[inst.Src1]

	case OpCodeLd:
		inst := decodePcRelative(instruction)
		computer.registers[inst.Register] = computer.memory[computer.pc+int(inst.PcOffset)]

	case OpCodeSt:
		inst := decodePcRelative(instruction)
		computer.memory[computer.pc+int(inst.PcOffset)] = computer.registers[inst.Register]

	case OpCodeLdi:
		inst := decodePcRelative(instruction)
		addr := computer.memory[computer.pc+int(inst.PcOffset)]
		computer.registers[inst.Register] = computer.memory[addr]

	case OpCodeSti:
		inst := decodePcRelative(instruction)
		addr := computer.memory[computer.pc+int(inst.PcOffset)]
		computer.memory[addr] = computer.registers[inst.Register]

	case OpCodeLdr:
		inst := decodeBaseRelative(instruction)
		addr := computer.registers[inst.Base] + int16(inst.Offset)
		computer.registers[inst.Register] = computer.memory[addr]

	case OpCodeStr:
		inst := decodeBaseRelative(instruction)
		addr := computer.registers[inst.Base] + int16(inst.Offset)
		computer.memory[addr] = computer.registers[inst.Register]

	case OpCodeLea:
		inst := decodePcRelative(instruction)
		computer.registers[inst.Register] = int16(computer.pc + int(inst.PcOffset))

	default:
		panic(fmt.Sprintf("unexpected instruction: %b", instruction))
	}
}

func getOpcode(instruction uint16) OpCode {
	return OpCode(getBits(int(instruction), 4, 13))
}

func decodeOperateInstruction(instruction uint16) DecodedInstruction {
	dst := getBits(int(instruction), 3, 10)
	src1 := getBits(int(instruction), 3, 7)
	registerMode := getBits(int(instruction), 1, 6)

	var src2 int
	if registerMode == 0 {
		src2 = getBits(int(instruction), 3, 1)
	} else {
		src2 = getBits(int(instruction), 5, 1)
	}

	return DecodedInstruction{
		Dst:          uint16(dst),
		Src1:         uint16(src1),
		RegisterMode: registerMode == 0,
		Src2:         uint16(src2),
	}
}

func decodePcRelative(instruction uint16) PcRelativeInstruction {
	dst := getBits(int(instruction), 3, 10)
	pcOffset := getBits(int(instruction), 9, 1)

	return PcRelativeInstruction{
		Register: uint16(dst),
		PcOffset: int16(pcOffset),
	}
}

func decodeBaseRelative(instruction uint16) BaseRelativeInstruction {
	dst := getBits(int(instruction), 3, 10)
	base := getBits(int(instruction), 3, 7)
	offset := getBits(int(instruction), 6, 1)

	return BaseRelativeInstruction{
		Register: uint8(dst),
		Base:     uint8(base),
		Offset:   int8(offset),
	}
}

func getBits(n, k, p int) int {
	return (((1 << k) - 1) & (n >> (p - 1)))
}

func main() {
	// computer := NewLittleComputer3()
	// computer.executeInstruction()
	// n := 0b1001000
	// fmt.Printf("%b\n", getBits(n, 4, 1))
}