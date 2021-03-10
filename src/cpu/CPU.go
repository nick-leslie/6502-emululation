package cpu

import (
	"encoding/binary"
	"fmt"

	"local.com/memory"
)

//CPU is a emulation of a 6502 CPU
type CPU struct {
	ProgramCounter uint16         // the pointer to the currently acsest locations
	StackPointer   byte           //pointer to the current locaion in the stack the second 255 bytes of memory
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

	//LDAI loads acumulator imediant 2 cycles
	LDAI byte = 0xA9
	//LDAZ loads acumulator from zero page 3 cycles
	LDAZ byte = 0xA5

	//TXA copys X to accumulator
	TXA byte = 0x8A

	//TYA copys Y to accumulator
	TYA byte = 0x98

	//STAZ stores the acumulator in some place in the zero page 3 cycles
	STAZ byte = 0x85

	// X  register Instructions -----------------------

	//LDXI loads value into X register in imedint mode in imident mode 2 cycles
	LDXI byte = 0xA2

	//STXZ stores the X value in the acumulator
	STXZ byte = 0x86

	//TAX copys accumulator to X
	TAX byte = 0xAA
	//Y register instructions--------------------------

	//LDYI loads value into Y register in imeient mode in imident mode 2 cycles
	LDYI byte = 0xA0

	//TAY copys Y to the acumulator
	TAY byte = 0xA8

	//STYZ stores the y in a zero page
	STYZ byte = 0x84
	//jump instructions ------------------------------------

	//JMPA Jumps to a direct location in memory takes 3 cycles
	JMPA byte = 0x4c
	//JMPI jumps to an indirecnt locaion in memory takes 5 cycles
	JMPI byte = 0x6C

	//logic instructions-------------------------------------

	//ANDI ands the acummulator and the next byte in memory in imident mode 2 cycles
	ANDI byte = 0x29

	//XOR or's the acummulator and the next byte in memory in imident mode 2 cycles
	XOR byte = 0x49

	//ORA or's the acumulator and the next byte in memorty in imident mode 2 cyclse
	ORA byte = 0x09
	//------------------------------------------------------- stack operation

	//TXS transfers X to stack 2 cycles
	TXS byte = 0x9A
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
	cpu.StackPointer = 0x00       // start of the stack is at 0x0100 end is 0x01ff set stack pointer to start of stack stack assumes start is always within 100
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
		case STXZ:
			cpu.storeInZeroPage(cycle, cpu.X)
			break
		case STAZ:
			cpu.storeInZeroPage(cycle, cpu.A)
			break
		case TAX:
			cpu.CopyToRegister(cycle, &cpu.X, cpu.A)
			break
		case TAY:
			cpu.CopyToRegister(cycle, &cpu.Y, cpu.A)
			break
		case TXA:
			cpu.CopyToRegister(cycle, &cpu.A, cpu.X)
			break
		case TYA:
			cpu.CopyToRegister(cycle, &cpu.A, cpu.Y)
			break
		case LDYI:
			cpu.LoadYImedient(cycle)
			break
		case STYZ:
			cpu.storeInZeroPage(cycle, cpu.Y)
			break
		case ANDI:
			cpu.AND(cycle)
			break
		case XOR:
			cpu.XOR(cycle)
			break
		case ORA:
			cpu.ORA(cycle)
			break
		case JMPA:
			cpu.JumpDirect(cycle)
			break
		case JMPI:
			cpu.JumpIndirect(cycle)
			break
		case TXS:
			cpu.TransferXtoStackPointer(cycle)
			break
		default:
			fmt.Printf("Instruction:%x not handled\n", opCode)
			break
		}
		fmt.Printf("Acumulator:%x\n", cpu.A)
		//this is where we switch on opcodes
	}
}

//-----------------------------------------instruction fucntions

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

//--------------------------------------------------------------Genaral registers

//ManipulateRegister Changes the Target register to the value passed in
func (cpu *CPU) ManipulateRegister(targetRegister *byte, value byte) {
	*targetRegister = value
	//set flags
	//checks if A = 0
	if *targetRegister == 0x00 {
		cpu.Flags.zero = true
	} else {
		cpu.Flags.zero = false
	}
	//this might be bugged becuse of numbers like 0xffff
	if (0b10000000 & *targetRegister) == 0b10000000 {
		cpu.Flags.negitive = true
	} else {
		cpu.Flags.negitive = false
	}
}

//CopyToRegister copys the passed in value to the targetRegister and decrements a cycle
func (cpu *CPU) CopyToRegister(cycle *int, target *byte, register byte) {
	cpu.ManipulateRegister(target, register)
	*cycle--
}

func (cpu *CPU) storeInZeroPage(cycle *int, value byte) {
	ZeroPageadress := cpu.FetchByte(cycle)
	cpu.Mem.Memory[ZeroPageadress] = value
	*cycle--
}

//--------------------------------------------------------------acumulator

//LoadAcumulatorImedient loads the byte after the adress into into the acumulaor register
func (cpu *CPU) LoadAcumulatorImedient(cycle *int) {
	value := cpu.FetchByte(cycle)
	cpu.ManipulateRegister(&cpu.A, value)
}

//LoadAcumulatorZeroPage loads acummulator with an adress within the zero page
func (cpu *CPU) LoadAcumulatorZeroPage(cycle *int) {
	value := cpu.FetchByteAtZeroPage(cycle, cpu.FetchByte(cycle))
	cpu.ManipulateRegister(&cpu.A, value)
}

//---------------------------------- x register

//LoadXImedient loads the next value into the X register
func (cpu *CPU) LoadXImedient(cycle *int) {
	value := cpu.FetchByte(cycle)
	cpu.ManipulateRegister(&cpu.X, value)
}

//--------------------------------- Y register

//LoadYImedient loads the next value into the X register
func (cpu *CPU) LoadYImedient(cycle *int) {
	value := cpu.FetchByte(cycle)
	cpu.ManipulateRegister(&cpu.Y, value)
}

//---------------------------------------- And instruction set

//AND ands in imediant mode
func (cpu *CPU) AND(cycle *int) {
	cpu.A &= cpu.FetchByte(cycle)
	fmt.Printf("Acumulator:%x\n", cpu.A)
}

//---------------------------------------- Exclusive OR Instruction set

//XOR does an XOR on the accumulator in imident mode
func (cpu *CPU) XOR(cycle *int) {
	cpu.A ^= cpu.FetchByte(cycle)
	fmt.Printf("Acumulator:%x\n", cpu.A)
}

//---------------------------------------- Inclusive OR Instruction set

//ORA Dose the inclusive or instruction in imideant mode
func (cpu *CPU) ORA(cycle *int) {
	cpu.A = cpu.A | cpu.FetchByte(cycle)
	fmt.Printf("Acumulator:%x\n", cpu.A)
}

//---------------------------------------- Addtion stuff

//---------------------------------------- stack stuff
//TODO on jsr to stack look at ben eater video

//TransferXtoStackPointer ,oves X into The stack Pointer
func (cpu *CPU) TransferXtoStackPointer(cycle *int) {
	cpu.StackPointer = cpu.X
	*cycle--
}

//---------------------------------------- basic jump

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
