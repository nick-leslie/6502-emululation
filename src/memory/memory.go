package memory

struct Memory type {
	maxMemory uint32
	Memory []byte
}

func CreateMemory() Memory {
	mem := new(Memory)
	mem.maxMemory = 1025 * 64
	mem.Memory = [mem.maxMemory]byte
	return Memory
}
//ManipulateMemory you pass in a map of adresses and bytes 
func (mem *Memory) ManipulateMemory(instructionSet map[uint32]byte) {
	for key, value := range instructionSet {
		mem.Memory[key] = value
	}
}