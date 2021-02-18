package cpu

import (
	"encoding/binary"
	"fmt"

	"local.com/memory"
)

type status struct {
	negitive bool
	zero     bool
	overflow bool
	cary     bool
	breaks   bool
	decimal  bool
	interupt bool
}

//CPU is a emulation of a 6502 CPU
type CPU struct {
	ProgramCounter uint16         // the pointer to the currently acsest locations
	StackPointer   uint16         //pointer to the current locaion in the stack the second 255 bytes of memory
	X, Y, A        byte           // reisters x y and accumulator
	Flags          status         // the state flags
	mem            *memory.Memory // memory
}

const (
	//Acumulator instructions
	//LDAI loads acumulator imediant
	LDAI byte = 0xA9

	//jump instructions
	JMPA byte = 0x4c
	JMPI byte = 0x6C
)

//NewCPU creates a new cpu
func NewCPU() CPU {
	cpu := new(CPU)
	cpu.mem = memory.CreateMemory() // resets memory
	cpu.ResetCPU()                  //resets cpu
	return *cpu                     //returns the initlised cpu
}

//FetchByte grabs a byte from memory uses a cycle and incruments the program counter
func (cpu *CPU) FetchByte(cycle *int) byte {
	targetByte := cpu.mem.Memory[cpu.ProgramCounter]
	cpu.ProgramCounter++
	*cycle--
	return targetByte
}

//ResetCPU resets the cpu to startup
func (cpu *CPU) ResetCPU() {
	cpu.ProgramCounter = 0xfffc // set progtram counter to start of program exicution
	cycles := 3
	cpu.JumpDirect(&cycles)
	fmt.Printf("Program counter:%x\n", cpu.ProgramCounter)
	cpu.StackPointer = 0x0100     // set stack pointer to start of stack
	cpu.X, cpu.Y, cpu.A = 0, 0, 0 // clear the registers
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
		case JMPA:
			cpu.JumpDirect(cycle)
			break
		case JMPI:
			cpu.JumpIndirect(cycle)
		default:
			fmt.Println("Instruction not handled")
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

//LoadAcumulatorImedient loads a byte into the acumulaor register
func (cpu *CPU) LoadAcumulatorImedient(cycle *int) {
	value := cpu.FetchByte(cycle)
	cpu.A = value
	//set flags
	//checks if A = 0
	if cpu.A == 0 {
		cpu.Flags.zero = true
	} else {
		cpu.Flags.zero = true
	}
	if (0b10000000 & cpu.A) == 0b10000000 {
		cpu.Flags.negitive = true
	} else {
		cpu.Flags.negitive = false
	}
}
