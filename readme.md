


instructions implumented
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
instructions to be implmented
