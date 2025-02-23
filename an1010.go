package main

import (
	"flag"
	"log"

	"github.com/dmisol/gru3880/pkg/asm"
	"github.com/dmisol/gru3880/pkg/boot"
	"github.com/dmisol/gru3880/pkg/mem"
	"go.bug.st/serial"
)

/*
	var ldr = []byte{
		// init serial 09600
		0x8E, 0x00, 0xFF, // 0000 8E00FF LDS #$FF
		0xCE, 0x10, 0x00, // 0003 CE1000 LDX #$1000 Offset for control registers.
		0x6F, 0x2C, // 0006 6F2C CLR SCCR1,X Initialise SCI for 8 data bits, 9600 baud
		0xCC, 0x30, 0x0C, // 0008 CC300C LDD #$300C
		0xA7, 0x2B, // 000B A72B STAA BAUD,X
		0xE7, 0x2D, // 000D E72D STAB SCCR2,X
		0x1C, 0x3C, 0x20, //000F 1C3C20 BSET HPRIO,X,#MDA Force Special Test

		// just print!

		// read memory with code
		0xCE, 0xB6, 0x00, // ldx 0xB600 - https://www.data-chip.ru/sheet/tipaudio-GRUNDIG/

		0xE7, 0x2F, // E72F STAB SCDR,X
		0x1F, 0x2E, 0x80, 0xFC, // 1F2E80FC WRITEC BRCLR SCSR,X,#TDRE,*

		0xE7, 0x2F, // E72F STAB SCDR,X
		0x1F, 0x2E, 0x80, 0xFC, // 1F2E80FC WRITEC BRCLR SCSR,X,#TDRE,*

		0xE7, 0x2F, // E72F STAB SCDR,X
		0x1F, 0x2E, 0x80, 0xFC, // 1F2E80FC WRITEC BRCLR SCSR,X,#TDRE,*

}
*/

func main() {
	file := flag.String("file", "bin/0322.bin", "file write")
	flag.Parse()

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(ports)

	port, err := boot.Bootload(ports[0], asm.BootAN1010)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer port.Close()

	log.Println("bootloader is fine")

	mode := &serial.Mode{
		BaudRate: 9600,
	}

	if err = port.SetMode(mode); err != nil {
		log.Fatal(err)
	}

	if err = mem.SetMode(port, "I"); err != nil {
		log.Fatal(err)
	}
	if err = mem.SendFile(port, 0xB600, 16, *file); err != nil {
		log.Fatal(err)
	}
	log.Println("ok")

}
