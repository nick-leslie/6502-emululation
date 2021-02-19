package cpu

import (
	"testing"
)

//genral tests----------------------------

//tests the CPU Reset sequence to make shure it can enter a program correctly
func TestReset(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint32]byte)
	Program[0xfffc] = 0xaa
	Program[0xfffd] = 0xbb
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	if cpu.ProgramCounter != 0xbbaa {
		t.Errorf("ProgramCounter:%x wanted:%x", cpu.ProgramCounter, 0xbbaa)
	}
}

//-------------------------------------------test Instructions

//accumulator
func TestLoadAcumulatorImedient(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint32]byte)
	Program[0xfffc] = 0xaa
	Program[0xfffd] = 0xbb
	Program[0xbbaa] = 0xA9
	Program[0xbbab] = 0xCC
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 2
	cpu.Execute(&cycles)
	//add tests for flags
	if cpu.A != 0xCC {
		t.Errorf("Acumulator:%x wanted:%x\nProgram counter:%x", cpu.A, 0xCC, cpu.ProgramCounter)
	} else {
		if cpu.Flags.zero == true {
			t.Errorf("Zero flag incorect Acumulator:%x wanted:%x\nProgram counter:%x", cpu.A, 0xCC, cpu.ProgramCounter)
		}
	}
}

//jump
func TestJumpImedient(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint32]byte)
	Program[0xfffc] = 0xaa
	Program[0xfffd] = 0xbb
	Program[0xbbaa] = 0x4c
	Program[0xbbab] = 0x11
	Program[0xbbac] = 0x11
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 3
	cpu.Execute(&cycles)
	if cpu.ProgramCounter != 0x1111 {
		t.Errorf("ProgramCounter:%x wanted:%x", cpu.ProgramCounter, 0x1111)
	}
}
func TestJumpIndirect(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint32]byte)
	Program[0xfffc] = 0xaa
	Program[0xfffd] = 0xbb
	Program[0xbbaa] = 0x6c
	Program[0xbbab] = 0x11
	Program[0xbbac] = 0x11
	Program[0x1111] = 0x22
	Program[0x1112] = 0x22
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 5
	cpu.Execute(&cycles)
	if cpu.ProgramCounter != 0x2222 {
		t.Errorf("ProgramCounter:%x wanted:%x", cpu.ProgramCounter, 0x2222)
	}
}
