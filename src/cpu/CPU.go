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

	JSR byte = 0x20
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

	//TSX transfers Stack pointer to X
	TSX byte = 0xBA

	//PHA pushes coppy of acumulator to stack
	PHA byte = 0x48
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
			cpu.A &= cpu.FetchByte(cycle)
			break
		case XOR:
			cpu.A ^= cpu.FetchByte(cycle)
			break
		case ORA:
			cpu.A = cpu.A | cpu.FetchByte(cycle)
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
		case TSX:
			cpu.X = cpu.StackPointer
			*cycle--
			break
		case PHA:
			cpu.pushToStack(cpu.A, cycle)
			*cycle--
		case JSR:
			cpu.JumpToSubRoutine(cycle)
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

//adds value and decremnts the stack pointer
func (cpu *CPU) pushToStack(value byte, cycle *int) {
	instructionSet := make(map[uint16]byte)
	bytes := []byte{cpu.StackPointer, 0x01}
	adress := binary.LittleEndian.Uint16(bytes)
	instructionSet[adress] = value
	cpu.Mem.ManipulateMemory(instructionSet)
	cpu.StackPointer--
	*cycle--
}

//adds a little endian adress to the stack subtracts 2 from cycles and stack
func (cpu *CPU) pushAdressToStack(adress uint16, cycle *int) {
	//break uint into bytes
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(adress))

	//maniplate stack
	instructionSet := make(map[uint16]byte)
	bytes := []byte{cpu.StackPointer, 0x01}
	stackAdress := binary.LittleEndian.Uint16(bytes)
	instructionSet[stackAdress] = b[0]
	cpu.StackPointer--
	bytes[0] = cpu.StackPointer
	stackAdress = binary.LittleEndian.Uint16(bytes)
	instructionSet[stackAdress] = b[1]
	cpu.Mem.ManipulateMemory(instructionSet)
	cpu.StackPointer--
	*cycle -= 2
}

// pull from stack  removes value and increces the stack pointer
//--------------------------------------------------------------Genaral registers INstructions

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

//---------------------------------------- Addtion stuff

//---------------------------------------- stack stuff
//TODO on jsr to stack look at ben eater video

//see if
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
func (cpu *CPU) jumpDirectWithAdress(adress uint16, cycle *int) {
	*cycle--
	cpu.ProgramCounter = adress
}

//JumpIndirect first jumps to locaions after the instruction then reads the byte and the next byte at that locaion for new adress
func (cpu *CPU) JumpIndirect(cycle *int) {
	//jumps to spesifiyed memory location
	cpu.JumpDirect(cycle)
	address := cpu.grabAdress(cycle)
	cpu.ProgramCounter = address
}

// ---------------------------------------------- sub routines
//fetch the instruction // 1
//fetch the adress // 2
//fetch the Next Instruction and subtract 1 store that in the stack  // 1
//jump direct //2

func (cpu *CPU) JumpToSubRoutine(cycle *int) {
	adress := cpu.grabAdress(cycle) // cycle -2
	returnAdress := cpu.ProgramCounter - 1
	cpu.pushAdressToStack(returnAdress, cycle) // cycle -2
	cpu.jumpDirectWithAdress(adress, cycle)    // -1 cycle
}

// debug-------------------------
func (cpu *CPU) PrintStack() {
	row := 0
	for i := len(cpu.Mem.Memory) - 1; i > 0; i-- {
		if i <= 0x01ff && i >= 0x0100 {
			fmt.Printf(" %x ", cpu.Mem.Memory[i])
			row++
			if row > 8 {
				fmt.Printf(" %x - %x\n", i+8, i)
				row = 0
			}
		}
	}
}
func (cpu *CPU) PrintMemory() {
	row := 0
	for i := len(cpu.Mem.Memory) - 1; i > 0; i-- {
		fmt.Printf(" %x ", cpu.Mem.Memory[i])
		row++
		if row > 8 {
			fmt.Printf(" %x - %x\n", i+8, i)
			row = 0
		}
	}
}
