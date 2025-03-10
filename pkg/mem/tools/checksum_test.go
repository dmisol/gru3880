package tools

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestChkSum(t *testing.T) {
	d, _ := hex.DecodeString("BDB653CC0800FD1018FD101AFD101CFD101EBDB665CC0002DD007F00067F0010")
	s := MkData(0xB600, d)
	if s != "S123B600BDB653CC0800FD1018FD101AFD101CFD101EBDB665CC0002DD007F00067F001055" {
		t.Fatal()
	}
	fmt.Println([]byte(s))
	/*
		d, _ = hex.DecodeString("A70D383239FEB680DFC7B6B68297C986AAB710048604B7100986C4B710280E3972")
		s = MkData(0xB660, d)
		if s != "S123B660A70D383239FEB680DFC7B6B68297C986AAB710048604B7100986C4B710280E3972\n" {
			t.Fatal()
		}
	*/
	d, _ = hex.DecodeString("0B29BDC02A18386A3B6F3B39180926FC39")
	s = MkData(0xC01E, d)
	if s != "S114C01E0B29BDC02A18386A3B6F3B39180926FC39DE" {
		fmt.Println(s)
		fmt.Println("S114C01E0B29BDC02A18386A3B6F3B39180926FC39DE")
		t.Fatal()
	}

	MkFinal(0)
}
