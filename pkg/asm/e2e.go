package asm

import "fmt"

func eraseAndFlash_head(mask string, addr string, erase bool) string {
	s := WriteEnable(mask)
	s += Serial_init

	if erase {
		s += LDAB_ + "06" + STAB__ + PPROG +
			STAB__ + addr +
			LDAB_ + "07" + STAB__ + PPROG

		for i := 0; i < 10; i++ {
			s += ms1("FF")
		}
	}

	s += CLR__ + PPROG
	s += LDY__ + addr

	return s
}

func EraseAndFlash(mask string, addr string, erase bool) string {
	s := eraseAndFlash_head(mask, addr, erase)
	jumpAddr := len(s) / 2
	jaHex := fmt.Sprintf("%04X", jumpAddr)

	s += WAIT4RX + LDAB_SERDATA // data to be written in B

	s += LDAA_ + "02" + STAA__ + PPROG +
		STAB_0_Y +
		LDAA_ + "03" + STAA__ + PPROG

	for i := 0; i < 10; i++ {
		s += ms1b()
	}

	s += INCY + JUMP_TO__ + jaHex
	return s
}

func ms1b() string {
	return WAIT4TX + STAB_SCDR
}
