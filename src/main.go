package main

import (
	"fmt"

	"local.com/cpu"
)

func main() {
	CPU := cpu.NewCPU()
	//this Inline program sets the starter adress to 0x1311 then it loads 0xfd into the accumulator it jumps to 01av and loads 0x11 into the accumulator
	CPU.Mem.Memory[0xfffc] = 0x11     // set least signifecent bit of adress
	CPU.Mem.Memory[0xfffd] = 0x13     // set most significent bit of adress
	CPU.Mem.Memory[0x1311] = cpu.LDAZ // Load accumulator zero page 3 cycles
	CPU.Mem.Memory[0x1312] = 0xfd     // zero page adress
	CPU.Mem.Memory[0x00fd] = 0x4c     // acumulator value value
	CPU.Mem.Memory[0x1313] = cpu.JMPA // jump 3 cycles
	CPU.Mem.Memory[0x1314] = 0xab     // set least signifecent bit of jump adress
	CPU.Mem.Memory[0x1315] = 0x01     // set most signifecent bit of jump adress
	CPU.Mem.Memory[0x01ab] = cpu.LDAI // loads acumulator with next byte 2 cycles
	CPU.Mem.Memory[0x01ac] = 0x11     // set accumulator
	CPU.ResetCPU()
	cycles := 8
	CPU.Execute(&cycles)
	CPU.PrintMemory()
	//fmt.Printf("%x\n", 0x0f&0xf0)
}

//this is all to show for presintaions
func littleENdiantest() {
	fmt.Println("gamer")
	fmt.Printf("%x %b\n", 0x001f<<2, 0x001f<<2)
	value := 0x11
	mask := 0x2
	final := 0x0000
	var test uint16
	var test2 uint16
	test = 0x00ff
	test2 = 0x001f
	fmt.Printf("Value:%b\nmask:%b\nfinal:%b\n", value, mask, final)
	fmt.Printf("%x %b\n", value&mask, value&mask)
	fmt.Printf("%x %b\n", final|value, final|value)
	fmt.Printf("%b %b %x\n", test, test2, (test<<8)|test2)
}
