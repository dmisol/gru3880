package asm

const (
	PPROG = "103B"

	LDAA_  = "86" // + byte
	STAA__ = "B7" // + addr
	LDAB_  = "C6" // + byte
	STAB__ = "F7" // + addr
	CLR__  = "7F"
)

/*

LDAA <byte>

LDAB #$02	; EELAT=1
STAB $103B	;

STAA <addr>

LDAB #$03	; EELAT=1, EPGM=1
STAB $103B	;

<send 10 serial @9600, approx 1ms each>

CLR $10

*/

func Write() string {
	s := WriteEnable("00")
	s += Serial_init
	/*
		s += WriteByte("B600", "03", 10)
		s += WriteByte("B601", "FC", 10)
		s += WriteByte("B602", "59", 10)

		s += WriteByte("B603", "22", 10)
		s += WriteByte("B604", "DD", 10)
		s += WriteByte("B605", "78", 10)
	*/
	/*
		s += WriteByte("B606", "10", 10)
		s += WriteByte("B607", "EF", 10)
		s += WriteByte("B608", "4A", 10)
	*/

	s += WriteByte("B609", "09", 10)
	s += WriteByte("B60A", "00", 10)
	s += WriteByte("B60B", "00", 10)

	s += WriteByte("B670", "00", 10)
	s += WriteByte("B671", "00", 10)

	s += WriteByte("B672", "41", 10)
	s += WriteByte("B673", "41", 10)
	s += WriteByte("B674", "41", 10)
	s += WriteByte("B675", "41", 10)

	s += WriteByte("B676", "00", 10)
	s += WriteByte("B677", "00", 10)

	s += WAIT4RX // forever
	s += WAIT4RX
	s += WAIT4RX

	return s
}

func EraseTail() string {
	s := WriteEnable("FE")
	s += Serial_init

	s += LDAB_ + "0E"
	s += STAB__ + PPROG
	s += STAB__ + "B7F0"
	s += LDAB_ + "0F"
	s += STAB__ + PPROG

	for i := 0; i < 10; i++ {
		s += ms1("00")
	}

	s += CLR__ + PPROG

	s += WAIT4RX // forever
	s += WAIT4RX
	s += WAIT4RX

	return s
}

func WriteEnable(block string) string {
	return LDAB_ + block + STAB__ + "1035"
}

func ms1(b string) string {
	return WAIT4TX + LDAB_ + b + STAB_SCDR
}

func WriteByte(addr string, b string, n int) string {

	s := LDAA_ + b +
		LDAB_ + "02" + STAB__ + PPROG +
		STAA__ + addr +
		LDAB_ + "03" + STAB__ + PPROG

	for i := 0; i < n; i++ {
		s += ms1(b)
	}

	s += CLR__ + PPROG
	return s
}
