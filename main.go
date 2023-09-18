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

type Register byte

const (
	R0 Register = 0
	R1 Register = 1
	R2 Register = 2
	R3 Register = 3
	R4 Register = 4
	R5 Register = 5
	R6 Register = 6
	R7 Register = 7
)

func (register Register) String() string {
	switch register {
	case R0:
		return "R0"
	case R1:
		return "R1"
	case R2:
		return "R2"
	case R3:
		return "R3"
	case R4:
		return "R4"
	case R5:
		return "R5"
	case R6:
		return "R6"
	case R7:
		return "R7"
	default:
		return fmt.Sprintf("value: %d", register)
	}
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
	Dst          Register
	Src1         Register
	RegisterMode bool
	Src2         Register
}

type PcRelativeInstruction struct {
	Register Register
	PcOffset int16
}

type BaseRelativeInstruction struct {
	Register Register
	Base     Register
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
		fmt.Printf("\n\naaaaaaa ADD inst %+v\n\n", inst)

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
		fmt.Printf("\n\naaaaaaa ST inst %+v\n\n", inst)
		computer.memory[computer.pc+int(inst.PcOffset)] = computer.registers[inst.Register]

	case OpCodeLdi:
		inst := decodePcRelative(instruction)
		fmt.Printf("\n\naaaaaaa LDI inst %+v\n\n", inst)
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
		fmt.Printf("\n\naaaaaaa STR inst %+v\n\n", inst)
		addr := computer.registers[inst.Base] + int16(inst.Offset)
		computer.memory[addr] = computer.registers[inst.Register]

	case OpCodeLea:
		inst := decodePcRelative(instruction)
		fmt.Printf("\n\naaaaaaa LEA inst %+v\n\n", inst)

		fmt.Printf("\n\naaaaaaa computer.pc %+v\n\n", computer.pc)
		value := int16(computer.pc + int(inst.PcOffset))
		fmt.Printf("\n\naaaaaaa value %+v %b\n\n", value, value)

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
		Dst:          Register(dst),
		Src1:         Register(src1),
		RegisterMode: registerMode == 0,
		Src2:         Register(src2),
	}
}

func decodePcRelative(instruction uint16) PcRelativeInstruction {
	dst := getBits(int(instruction), 3, 10)
	fmt.Printf("\n\naaaaaaa instruction %b\n\n", instruction)
	pcOffset := getBits(int(instruction), 9, 1)
	fmt.Printf("\n\naaaaaaa pcOffset %+v %b\n\n", pcOffset, pcOffset)

	return PcRelativeInstruction{
		Register: Register(dst),
		PcOffset: int16(pcOffset),
	}
}

func decodeBaseRelative(instruction uint16) BaseRelativeInstruction {
	dst := getBits(int(instruction), 3, 10)
	base := getBits(int(instruction), 3, 7)
	offset := getBits(int(instruction), 6, 1)

	return BaseRelativeInstruction{
		Register: Register(dst),
		Base:     Register(base),
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
