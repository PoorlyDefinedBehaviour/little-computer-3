package main

import "fmt"

type LittleComputer3 struct {
	memory [65536]int16
	// Registers from r0 to r7.
	registers [8]int16
	pc        int

	// Condition codes
	negative bool
	zero     bool
	positive bool
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

	OpCodeBr OpCode = 0b0000

	OpCodeJMP OpCode = 0b1100

	OpCodeTrap OpCode = 0b1111
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
	case OpCodeBr:
		return "BR"
	case OpCodeJMP:
		return "JMP"
	case OpCodeTrap:
		return "TRAP"
	default:
		panic(fmt.Sprintf("unexpected opcode: %b", opCode))
	}
}

type DecodedInstruction struct {
	Dst          byte
	Src1         byte
	RegisterMode bool
	Src2         byte
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

type BranchInstruction struct {
	N        bool
	Z        bool
	P        bool
	PcOffset int16
}

type TrapInstruction struct {
	TrapVector byte
}

func NewLittleComputer3() *LittleComputer3 {
	return &LittleComputer3{
		memory:    [65536]int16{},
		registers: [8]int16{},
		pc:        0,
		negative:  false,
		zero:      false,
		positive:  false,
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

		value := a + b

		computer.registers[inst.Dst] = value

		computer.zero = value == 0
		computer.negative = value < 0
		computer.positive = value > 0

	case OpCodeAnd:
		inst := decodeOperateInstruction(instruction)

		a := computer.registers[inst.Src1]

		var b int16
		if inst.RegisterMode {
			b = computer.registers[inst.Src2]
		} else {
			b = int16(inst.Src2)
		}

		value := a & b

		computer.registers[inst.Dst] = value

		computer.zero = value == 0
		computer.negative = value < 0
		computer.positive = value > 0

	case OpCodeNot:
		inst := decodeOperateInstruction(instruction)

		value := ^computer.registers[inst.Src1]

		computer.registers[inst.Dst] = value

		computer.zero = value == 0
		computer.negative = value < 0
		computer.positive = value > 0

	case OpCodeLd:
		inst := decodePcRelative(instruction)

		value := computer.memory[computer.pc+int(inst.PcOffset)]

		computer.registers[inst.Register] = value

		computer.zero = value == 0
		computer.negative = value < 0
		computer.positive = value > 0

	case OpCodeSt:
		inst := decodePcRelative(instruction)
		computer.memory[computer.pc+int(inst.PcOffset)] = computer.registers[inst.Register]

	case OpCodeLdi:
		inst := decodePcRelative(instruction)
		addr := computer.memory[computer.pc+int(inst.PcOffset)]

		value := computer.memory[addr]

		computer.registers[inst.Register] = value

		computer.zero = value == 0
		computer.negative = value < 0
		computer.positive = value > 0

	case OpCodeSti:
		inst := decodePcRelative(instruction)
		addr := computer.memory[computer.pc+int(inst.PcOffset)]
		computer.memory[addr] = computer.registers[inst.Register]

	case OpCodeLdr:
		inst := decodeBaseRelative(instruction)
		addr := computer.registers[inst.Base] + int16(inst.Offset)

		value := computer.memory[addr]

		computer.registers[inst.Register] = value

		computer.zero = value == 0
		computer.negative = value < 0
		computer.positive = value > 0

	case OpCodeStr:
		inst := decodeBaseRelative(instruction)
		addr := computer.registers[inst.Base] + int16(inst.Offset)
		computer.memory[addr] = computer.registers[inst.Register]

	case OpCodeLea:
		inst := decodePcRelative(instruction)

		value := int16(computer.pc + int(inst.PcOffset))

		computer.registers[inst.Register] = value

		computer.zero = value == 0
		computer.negative = value < 0
		computer.positive = value > 0

	case OpCodeBr:
		inst := decodeBranchInstruction(instruction)

		// TODO: only branch if set bit is set.
		computer.pc += int(inst.PcOffset)

	case OpCodeJMP:
		inst := decodeBaseRelative(instruction)

		computer.pc = int(computer.registers[inst.Base])

	case OpCodeTrap:
		_ = decodeTrapInstruction(instruction)

		// TODO: call service routine

		// TODO: set pc to the instruction following TRAP

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
		Dst:          byte(dst),
		Src1:         byte(src1),
		RegisterMode: registerMode == 0,
		Src2:         byte(src2),
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

func decodeBranchInstruction(instruction uint16) BranchInstruction {
	n := getBits(int(instruction), 1, 12)
	z := getBits(int(instruction), 1, 11)
	p := getBits(int(instruction), 1, 10)
	pcOffset := getBits(int(instruction), 9, 1)

	return BranchInstruction{
		N:        n == 1,
		Z:        z == 1,
		P:        p == 1,
		PcOffset: int16(pcOffset),
	}
}

func decodeTrapInstruction(instruction uint16) TrapInstruction {
	trapVector := getBits(int(instruction), 8, 1)

	return TrapInstruction{TrapVector: byte(trapVector)}
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
