package asm

var BootDump = Serial_init + Read_n_echo +
	LDY_0 +
	// 0x0022
	WAIT4TX +
	LDAB_0_Y + STAB_SCDR +
	INCY +

	JUMP_TO__ + "0022"

func BootDumpFrom(addr string) string {
	return Serial_init + Read_n_echo +
		LDY__ + addr +

		WAIT4TX +
		LDAB_0_Y + STAB_SCDR +

		INCY +
		JUMP_TO__ + "0022"
}
