package cpu

import (
	"fmt"
	"testing"
)

//genral tests----------------------------

//tests the CPU Reset sequence to make shure it can enter a program correctly
func TestReset(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
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
func TestLDAI(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	Program[0xfffc] = 0xaa
	Program[0xfffd] = 0xbb
	Program[0xbbaa] = LDAI
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
func TestLDAZ(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	Program[0xfffc] = 0xaa
	Program[0xfffd] = 0xbb
	Program[0xbbaa] = LDAZ
	Program[0xbbab] = 0xCC
	Program[0x00cc] = 0x55
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 2
	cpu.Execute(&cycles)
	//add tests for flags
	if cpu.A != 0x55 {
		t.Errorf("Acumulator:%x wanted:%x\nProgram counter:%x", cpu.A, 0x55, cpu.ProgramCounter)
	} else {
		if cpu.Flags.zero == true {
			t.Errorf("Zero flag incorect Acumulator:%x wanted:%x\nProgram counter:%x", cpu.A, 0x55, cpu.ProgramCounter)
		}
	}
}

func TestLDXI(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	Program[0xfffc] = 0xaa
	Program[0xfffd] = 0xbb
	Program[0xbbaa] = LDXI
	Program[0xbbab] = 0xCC
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 2
	cpu.Execute(&cycles)
	if cpu.X != 0xCC {
		t.Errorf("X register:%x wanted:%x\nProgram counter:%x", cpu.X, 0xCC, cpu.ProgramCounter)
	}
}
func TestLDYI(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	Program[0xfffc] = 0xaa
	Program[0xfffd] = 0xbb
	Program[0xbbaa] = LDYI
	Program[0xbbab] = 0xCC
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 2
	cpu.Execute(&cycles)
	if cpu.Y != 0xCC {
		t.Errorf("Y register:%x wanted:%x\nProgram counter:%x", cpu.Y, 0xCC, cpu.ProgramCounter)
	}
}

//jump
func TestJumpImedient(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	Program[0xfffc] = 0xaa
	Program[0xfffd] = 0xbb
	Program[0xbbaa] = JMPA
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
	Program := make(map[uint16]byte)
	Program[0xfffc] = 0xaa
	Program[0xfffd] = 0xbb
	Program[0xbbaa] = JMPI
	Program[0xbbab] = 0x11 // location adress 2
	Program[0xbbac] = 0x11 // location adress 1
	Program[0x1111] = 0x22 // program
	Program[0x1112] = 0x22
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 5
	cpu.Execute(&cycles)
	if cpu.ProgramCounter != 0x2222 {
		t.Errorf("ProgramCounter:%x wanted:%x", cpu.ProgramCounter, 0x2222)
	}
}

//add tests for  and the FLag tests
func TestTAX(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDAI
	Program[0xbbab] = 0xCC
	Program[0xbbac] = TAX
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 4
	cpu.Execute(&cycles)
	if cpu.X != 0xCC {
		t.Errorf("A register:%x wanted:%x\nProgram counter:%x", cpu.X, 0xCC, cpu.ProgramCounter)
	}
}
func TestTAY(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDAI
	Program[0xbbab] = 0xCC
	Program[0xbbac] = TAY
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 4
	cpu.Execute(&cycles)
	if cpu.Y != 0xCC {
		t.Errorf("Y register:%x wanted:%x\nProgram counter:%x", cpu.Y, 0xCC, cpu.ProgramCounter)
	}
}

func TestAND(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDAI
	Program[0xbbab] = 0xf0
	Program[0xbbac] = ANDI
	Program[0xbbad] = 0x0f
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 4
	cpu.Execute(&cycles)
	if cpu.A != 0xf0&0x0f {
		t.Errorf("A register:%x wanted:%x\nProgram counter:%x", cpu.A, 0xf0&0x0f, cpu.ProgramCounter)
	}
}

func TestXOR(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDAI
	Program[0xbbab] = 0xff
	Program[0xbbac] = XOR
	Program[0xbbad] = 0x0f
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 4
	cpu.Execute(&cycles)
	if cpu.A != 0xff^0x0f {
		t.Errorf("A register:%x wanted:%x\nProgram counter:%x", cpu.A, 0xff^0x0f, cpu.ProgramCounter)
	}
}

func TestORA(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDAI
	Program[0xbbab] = 0xf0
	Program[0xbbac] = ORA
	Program[0xbbad] = 0x0f
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 4
	cpu.Execute(&cycles)
	if cpu.A != 0xf0|0x0f {
		t.Errorf("A register:%x wanted:%x\nProgram counter:%x", cpu.A, 0xf0|0x0f, cpu.ProgramCounter)
	}
}

//add tests for TYA TXA

func TestTYA(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDYI
	Program[0xbbab] = 0xCC
	Program[0xbbac] = TYA
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 4
	cpu.Execute(&cycles)
	if cpu.A != 0xCC {
		t.Errorf("A register:%x wanted:%x\nProgram counter:%x", cpu.A, 0xCC, cpu.ProgramCounter)
	}
}
func TestTXA(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDXI
	Program[0xbbab] = 0xCC
	Program[0xbbac] = TXA
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 4
	cpu.Execute(&cycles)
	if cpu.A != 0xCC {
		t.Errorf("A register:%x wanted:%x\nProgram counter:%x", cpu.A, 0xCC, cpu.ProgramCounter)
	}
}

func TestSTAZ(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDAI // 2 cycles
	Program[0xbbab] = 0xf0 // 1 cycle
	Program[0xbbac] = STAZ // 3 cycles
	Program[0xbbad] = 0x0f
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 5
	cpu.Execute(&cycles)
	if cpu.A != cpu.Mem.Memory[0x0f] {
		t.Errorf("A register:%x Memory Value:%x\nProgram counter:%x", cpu.A, cpu.Mem.Memory[0x0f], cpu.ProgramCounter)
	}
}
func TestSTXZ(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDXI // 2 cycles
	Program[0xbbab] = 0xf0 // 1 cycle
	Program[0xbbac] = STXZ // 3 cycles
	Program[0xbbad] = 0x0f
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 5
	cpu.Execute(&cycles)
	if cpu.X != cpu.Mem.Memory[0x0f] {
		t.Errorf("A register:%x Memory Value:%x\nProgram counter:%x", cpu.A, cpu.Mem.Memory[0x0f], cpu.ProgramCounter)
	}
}
func TestSTYZ(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDYI // 2 cycles
	Program[0xbbab] = 0xf0 // 1 cycle
	Program[0xbbac] = STYZ // 3 cycles
	Program[0xbbad] = 0x0f
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 5
	cpu.Execute(&cycles)
	if cpu.Y != cpu.Mem.Memory[0x0f] {
		t.Errorf("A register:%x Memory Value:%x\nProgram counter:%x", cpu.A, cpu.Mem.Memory[0x0f], cpu.ProgramCounter)
	}
}

func TestStackAddion(t *testing.T) {
	if 0x01ff != 0x0100+0xff {
		t.Errorf("%x", 0x0100+0xff)
	}
}

//ADD TEST FOR and PHA

func TestTXS(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDXI
	Program[0xbbab] = 0xff
	Program[0xbbac] = TXS
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 4
	cpu.Execute(&cycles)
	if cpu.StackPointer != 0xff {
		t.Errorf("X register:%x Stack Pointer Value:%x Program counter:%x", cpu.X, cpu.StackPointer, cpu.ProgramCounter)
	}
}

func TestTSX(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDXI
	Program[0xbbab] = 0xff
	Program[0xbbac] = TXS
	Program[0xbbad] = TSX
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 4
	cpu.Execute(&cycles)
	if cpu.X != 0xFF {
		t.Errorf("X register:%x Stack Pointer Value:%x Program counter:%x", cpu.X, cpu.StackPointer, cpu.ProgramCounter)
	}
}
func TestPHA(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	Program[0xbbaa] = LDXI //2
	Program[0xbbab] = 0xff
	Program[0xbbac] = TXS  // 2
	Program[0xbbad] = LDAI // 2
	Program[0xbbae] = 0x11
	Program[0xbbaf] = PHA // 1
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 7
	cpu.Execute(&cycles)
	if cpu.Mem.Memory[0x01ff] != 0x11 {
		t.Errorf("X register:%x Stack Pointer Value:%x Acumulator:%x Program counter:%x MemoryLocation:%x", cpu.X, cpu.StackPointer, cpu.A, cpu.ProgramCounter, cpu.Mem.Memory[0x01FF])
	}
}

//add test of jump to sub routine
func TestJumpToSubRoutine(t *testing.T) {
	cpu := NewCPU()
	Program := make(map[uint16]byte)
	cpu.Mem.SetStartAdress(0xbbaa)
	//initlise stack
	Program[0xbbaa] = LDXI
	Program[0xbbab] = 0xff
	Program[0xbbac] = TXS

	Program[0xbbad] = JSR
	Program[0xbbae] = 0x22
	Program[0xbbaf] = 0x11
	cpu.Mem.ManipulateMemory(Program)
	cpu.ResetCPU()
	cycles := 10
	cpu.Execute(&cycles)
	if cpu.ProgramCounter != 0x1122 {
		t.Errorf("Wrong program counter Curent:%x, target:0x1122", cpu.ProgramCounter)
	}
	fmt.Println(cpu.StackPointer)
	checkPoint := cpu.Mem.ReturnAdressInStack(cpu.StackPointer) + 1
	fmt.Println(checkPoint)
	if cpu.Mem.Memory[checkPoint] != 0xbb {
		t.Errorf("Wrong most significant Byte current:%x, target:0xbb", cpu.Mem.Memory[checkPoint])
	}
	checkPoint++
	if cpu.Mem.Memory[checkPoint] != 0xaf {
		t.Errorf("Wrong least significant Byte current:%x, target:0xaf", cpu.Mem.Memory[checkPoint])
	}
	fmt.Println(checkPoint)
	cpu.printStack(cpu.Mem.Memory)
	//new adress = 0x1122
	//return adress should be bbab but you when use return from sub routine it should add 1
}

// bit tests
func TestDevBYTwo(t *testing.T) {
	ogVal := 12
	devVal := ogVal >> 1
	fmt.Println(ogVal)
	fmt.Println(devVal)
	if devVal != ogVal/2 {
		t.Errorf("test failed")
	}
}
func TestMultBYTwo(t *testing.T) {
	ogVal := 12
	devVal := ogVal << 1
	fmt.Println(ogVal)
	fmt.Println(devVal)
	if devVal != ogVal*2 {
		t.Errorf("test failed")
	}
}
