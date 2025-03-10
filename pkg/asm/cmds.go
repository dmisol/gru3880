package asm

const (
	Serial_init       = "8E00FFCE10006F2CCC300CA72BE72D1C3C20"
	Serial_init_short = "8E00FFCE10006F2CCC300CA72BE72D" //"1C3C20"
	Read_n_echo       = "1F2E20FCE62F1F2E80FCE72F"
	Send_again        = "1F2E80FCE72F"

	LDY__    = "18CE"
	LDY_0    = "18CE0000"
	LDY_B600 = "18CEB600"

	INCY = "1808"

	// X нельзя трогать!
	LDX_0 = "CE0000"
	INCX  = "08"

	//LDAA_0_Y = "180000" //"18B600"
	// STAA_SCDR = "972F"

	LDAB_0_Y  = "18E600"
	STAB_SCDR = "E72F"
	STAB_0_Y  = "18E700"

	WAIT4TX = "1F2E80FC" //"132E80FC"
	WAIT4RX = "1F2E20FC"

	CPX_BYTES = "8C0020"
	LOOP_ADDR = "26EF" //

	JUMP_TO__    = "7E"
	JUMP_TO_0000 = "7E0000"
	JUMP_TO_0022 = "7E0022"

	LDAB_SERDATA = "E62F"
)
