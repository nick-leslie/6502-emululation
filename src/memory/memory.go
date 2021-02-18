package memory

type Memory struct {
	maxMemory uint16
	Memory    []byte
}

func CreateMemory() *Memory {
	mem := new(Memory)
	var maxuint16 uint16
	maxuint16 = 0xffff
	mem.maxMemory = maxuint16
	mem.Memory = make([]byte, mem.maxMemory)
	for i := 0; i < len(mem.Memory); i++ {
		mem.Memory[i] = 0
	}
	//this is an inline program sets the adress apon reset
	mem.Memory[0xfffc] = 0x11
	mem.Memory[0xfffd] = 0x13
	mem.Memory[0x1311] = 0xA9
	mem.Memory[0x1312] = 0xfd
	//-------------------------
	return mem
}

//ManipulateMemory you pass in a map of adresses and bytes
func (mem *Memory) ManipulateMemory(instructionSet map[uint32]byte) {
	for key, value := range instructionSet {
		mem.Memory[key] = value
	}
}
