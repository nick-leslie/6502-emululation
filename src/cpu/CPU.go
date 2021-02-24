package cpu

import (
	"encoding/binary"
	"fmt"

	"local.com/memory"
)

//CPU is a emulation of a 6502 CPU
type CPU struct {
	ProgramCounter uint16         // the pointer to the currently acsest locations
	StackPointer   uint16         //pointer to the current locaion in the stack the second 255 bytes of memory
	X, Y, A        byte           // reisters x y and accumulator
	Flags          status         // the state flags
	Mem            *memory.Memory // memory
}
type status struct {
	negitive bool
	zero     bool
	overflow bool
	cary     bool
	breaks   bool
	decimal  bool
	interupt bool
}

const (
	//Acumulator instructions

	//LDAI loads acumulator imediant
	LDAI byte = 0xA9
	//LDAZ loads acumulator from zero page
	LDAZ byte = 0xA5

	// X  register Instructions -----------------------

	//LDXI loads value into X register in imedint mode in imident mode 2 cycles
	LDXI byte = 0xA2

	//Y register instructions--------------------------

	//LDYI loads value into Y register in imeient mode in imident mode 2 cycles
	LDYI byte = 0xA0

	//jump instructions ------------------------------------

	//JMPA Jumps to a direct location in memory takes 3 cycles
	JMPA byte = 0x4c
	//JMPI jumps to an indirecnt locaion in memory takes 5 cycles
	JMPI byte = 0x6C

	//logic instructions-------------------------------------

	//ANDI ands the acummulator and the next byte in memory in imident mode 2 cycles
	ANDI byte = 0x29

	//XORI or's the acummulator and the next byte in memory in imident mode 2 cycles
	XORI byte = 0x49

	//ORAI or's the acumulator and the next byte in memorty in imident mode 2 cyclse
	ORAI byte = 0x09
	//-------------------------------------------------------
)

//NewCPU creates a new cpu
func NewCPU() CPU {
	cpu := new(CPU)
	cpu.Mem = memory.CreateMemory() // resets memory
	cpu.ResetCPU()                  //resets cpu
	return *cpu                     //returns the initlised cpu
}

//ResetCPU resets the cpu to startup
func (cpu *CPU) ResetCPU() {
	cpu.ProgramCounter = 0xfffc // set progtram counter to start of program exicution
	cycles := 3
	cpu.JumpDirect(&cycles)
	fmt.Printf("Program counter:%x\n", cpu.ProgramCounter)
	cpu.StackPointer = 0x0100     // set stack pointer to start of stack
	cpu.X, cpu.Y, cpu.A = 0, 0, 0 // clear the registers
	cpu.Flags.zero = false
	cpu.Flags.negitive = false
	cpu.Flags.overflow = false
	cpu.Flags.cary = false
	cpu.Flags.breaks = false
	cpu.Flags.decimal = false
	cpu.Flags.interupt = false
}

//FetchByte grabs a byte from memory uses a cycle and incruments the program counter
func (cpu *CPU) FetchByte(cycle *int) byte {
	targetByte := cpu.Mem.Memory[cpu.ProgramCounter]
	cpu.ProgramCounter++
	*cycle--
	return targetByte
}

//FetchByteAtZeroPage this fetches a byte with a zero page adress
func (cpu *CPU) FetchByteAtZeroPage(cycle *int, MemLoc byte) byte {
	targetByte := cpu.Mem.Memory[MemLoc]
	*cycle--
	return targetByte
}

//Execute Runs instructions for however many cycles on the CPU
func (cpu *CPU) Execute(cycle *int) {
	//continue unill there are no more cycles
	for *cycle > 0 {
		opCode := cpu.FetchByte(cycle)
		fmt.Printf("opCode:%x \n", opCode)
		switch opCode {
		case LDAI:
			cpu.LoadAcumulatorImedient(cycle)
			break
		case LDAZ:
			cpu.LoadAcumulatorZeroPage(cycle)
			break
		case LDXI:
			cpu.LoadXImedient(cycle)
			break
		case LDYI:
			cpu.LoadYImedient(cycle)
			break
		case ANDI:
			cpu.ANDImedianet(cycle)
			break
		case XORI:
			cpu.EXORImedianant(cycle)
			break
		case ORAI:
			cpu.ORAImedianant(cycle)
			break
		case JMPA:
			cpu.JumpDirect(cycle)
			break
		case JMPI:
			cpu.JumpIndirect(cycle)
		default:
			fmt.Printf("Instruction:%x not handled\n", opCode)
			break
		}
		fmt.Printf("Acumulator:%x\n", cpu.A)
		//this is where we switch on opcodes
	}
}

//-----------------------------------------instruction fucntions

//JumpDirect jumps to the adress spesifiyed by the two folowing bytes
func (cpu *CPU) JumpDirect(cycle *int) {
	address := cpu.grabAdress(cycle)
	fmt.Printf("Program counter:%x\n", address)
	*cycle--
	cpu.ProgramCounter = address
}

//JumpIndirect first jumps to locaions after the instruction then reads the byte and the next byte at that locaion for new adress
func (cpu *CPU) JumpIndirect(cycle *int) {
	//jumps to spesifiyed memory location
	cpu.JumpDirect(cycle)
	address := cpu.grabAdress(cycle)
	cpu.ProgramCounter = address
}
func (cpu *CPU) grabAdress(cycle *int) uint16 {
	bytes := make([]byte, 2)
	i := 0
	//remember if there is the wrong amount of cycles it will fail
	//refactored to only go twice remeber fetch byte will deduct cycles still
	for i < 2 {
		bytes[i] = cpu.FetchByte(cycle)
		i++
	}
	address := binary.LittleEndian.Uint16(bytes)
	return address
}

//--------------------------------------------------------------acumulator

//LoadAcumulatorImedient loads the byte after the adress into into the acumulaor register
func (cpu *CPU) LoadAcumulatorImedient(cycle *int) {
	value := cpu.FetchByte(cycle)
	cpu.ManipulateAcumulator(value)
}

//LoadAcumulatorZeroPage loads acummulator with an adress within the zero page
func (cpu *CPU) LoadAcumulatorZeroPage(cycle *int) {
	value := cpu.FetchByteAtZeroPage(cycle, cpu.FetchByte(cycle))
	cpu.ManipulateAcumulator(value)
}

//ManipulateAcumulator is the fucnction that manipulates the acumulator register addionaly it configers the flag
func (cpu *CPU) ManipulateAcumulator(value byte) {
	cpu.A = value
	//set flags
	//checks if A = 0
	if cpu.A == 0x00 {
		cpu.Flags.zero = true
	} else {
		cpu.Flags.zero = false
	}
	//this might be bugged becuse of numbers like 0xffff
	if (0b10000000 & cpu.A) == 0b10000000 {
		cpu.Flags.negitive = true
	} else {
		cpu.Flags.negitive = false
	}
}

//---------------------------------- x register

//LoadXImedient loads the next value into the X register
func (cpu *CPU) LoadXImedient(cycle *int) {
	value := cpu.FetchByte(cycle)
	cpu.ManipulateXRegister(value)
}

//ManipulateXRegister Changes the X Register and sets the flags
func (cpu *CPU) ManipulateXRegister(value byte) {
	cpu.X = value
	//set flags
	//checks if A = 0
	if cpu.X == 0x00 {
		cpu.Flags.zero = true
	} else {
		cpu.Flags.zero = false
	}
	//this might be bugged becuse of numbers like 0xffff
	if (0b10000000 & cpu.X) == 0b10000000 {
		cpu.Flags.negitive = true
	} else {
		cpu.Flags.negitive = false
	}
}

//--------------------------------- Y register

//LoadYImedient loads the next value into the X register
func (cpu *CPU) LoadYImedient(cycle *int) {
	value := cpu.FetchByte(cycle)
	cpu.ManipulateYRegister(value)
}

//ManipulateYRegister Changes the Y Register and sets the flags
func (cpu *CPU) ManipulateYRegister(value byte) {
	cpu.Y = value
	//set flags
	//checks if A = 0
	if cpu.Y == 0x00 {
		cpu.Flags.zero = true
	} else {
		cpu.Flags.zero = false
	}
	//this might be bugged becuse of numbers like 0xffff
	if (0b10000000 & cpu.Y) == 0b10000000 {
		cpu.Flags.negitive = true
	} else {
		cpu.Flags.negitive = false
	}
}

//---------------------------------------- And instruction set

//ANDImedianet ands in imediant mode
func (cpu *CPU) ANDImedianet(cycle *int) {
	cpu.A &= cpu.FetchByte(cycle)
}

//---------------------------------------- Exclusive OR Instruction set

//EXORImedianant does an EXOR on the accumulator in imident mode
func (cpu *CPU) EXORImedianant(cycle *int) {
	cpu.A ^= cpu.FetchByte(cycle)
}

//---------------------------------------- Inclusive OR Instruction set

//ORAImedianant Dose the inclusive or instruction in imideant mode
func (cpu *CPU) ORAImedianant(cycle *int) {
	cpu.A = cpu.A | cpu.FetchByte(cycle)
}

//---------------------------------------- Addtion stuff
