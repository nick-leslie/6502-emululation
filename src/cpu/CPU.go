package cpu

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
	ProgramCounter uint16 // the pointer to the currently acsest locations
	StackPointer   uint16 //pointer to the current locaion in the stack the second 255 bytes of memory
	X, Y, A        byte   // reisters x y and accumulator
	State          status // the state flags
}

//NewCPU creates a new cpu
func NewCPU() CPU {
	cpu := new(CPU)
	cpu.ResetCPU() //resets cpu
	return *cpu    //returns the initlised cpu
}

//ResetCPU resets the cpu to startup
func (cpu CPU) ResetCPU() {
	cpu.ProgramCounter = 0xfffc   // set progtram counter to start of program exicution
	cpu.StackPointer = 0x0100     // set stack pointer to start of stack
	cpu.X, cpu.Y, cpu.A = 0, 0, 0 // clear the registers
}
