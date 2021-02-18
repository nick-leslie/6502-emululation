package main

import (
	"fmt"

	"local.com/cpu"
)

func main() {
	cpu := cpu.NewCPU()
	cpu.ResetCPU()
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
