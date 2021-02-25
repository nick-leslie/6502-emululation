package memory

import (
	"encoding/binary"
	"fmt"
)

//Memory holds bytes
type Memory struct {
	maxMemory uint16
	Memory    []byte
}

//CreateMemory creates empty memory
func CreateMemory() *Memory {
	mem := new(Memory)
	var maxuint16 uint16
	maxuint16 = 0xffff
	mem.maxMemory = maxuint16
	mem.Memory = make([]byte, mem.maxMemory)
	for i := 0; i < len(mem.Memory); i++ {
		mem.Memory[i] = 0
	}
	return mem
}

//ManipulateMemory you pass in a map of adresses and bytes
func (mem *Memory) ManipulateMemory(instructionSet map[uint16]byte) {
	for key, value := range instructionSet {
		mem.Memory[key] = value
		fmt.Printf("%x\n", mem.Memory[key])
	}
	fmt.Printf("called manipulate memory\n")
}

//SetStartAdress sets the start adress
func (mem *Memory) SetStartAdress(address uint16) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(address))
	mem.Memory[0xfffc] = b[0]
	mem.Memory[0xfffd] = b[1]
}
