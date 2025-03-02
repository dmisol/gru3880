package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dmisol/gru3880/pkg/asm"
	"github.com/dmisol/gru3880/pkg/boot"
	"go.bug.st/serial"
)

var (
	startFrom = "B600"

	need2save = map[int]bool{
		0:      true,
		0x1000: true,
		0xB600: true,
		0xB700: true,
	}

	STAx = map[byte]bool{
		0xA7: true,
		0x97: true,
		0xE7: true,
	}
	EE = map[byte]bool{
		0x35: true, // PPROG
		0x3B: true, //BPROT
	}
)

const single = false

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
	//file := flag.String("file", "bin/0322.bin", "file write")
	//flag.Parse()

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(ports)

	port, err := boot.Bootload(ports[0], asm.BootDumpFrom(startFrom))
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

	/*
		if err = mem.SetMode(port, "I"); err != nil {
			log.Fatal(err)
		}
	*/
	ack_once(port)
	read_forever(port)

	/*
		if err = mem.SendFile(port, 0xB600, 16, "code1234.bin"); err != nil {
			log.Fatal(err)
		}
	*/
	/*
	   b0, _ := hex.DecodeString("FFFFFFFFFFFFFFFFFF")
	   b1, _ := hex.DecodeString("03FC5922DD7819E643")

	   sl := tools.MkData(0xB600, b0)

	   log.Println("data slice:", sl)

	   	if err = mem.SendSlice(port, sl); err != nil {
	   		log.Fatal(err)
	   	}

	   sl = tools.MkData(0xB600, b1)

	   log.Println("data slice:", sl)

	   	if err = mem.SendSlice(port, sl); err != nil {
	   		log.Fatal(err)
	   	}

	   	if err = mem.SendSlice(port, "S9030000FC"); err != nil {
	   		log.Fatal(err)
	   	}
	*/
	log.Println("ok")

}

func read_forever(port serial.Port) {
	r := make([]byte, 1)
	addr, _ := strconv.ParseInt(startFrom, 16, 32)
	slice := make([]byte, 0)

	for {
		n, err := port.Read(r)
		if err != nil {
			log.Println("rd error: ", err, n)
			return
		}
		// log.Println("GOT", n, err, r)

		slice = append(slice, r[0])
		if len(slice)%0x100 == 0 && len(slice) >= 0x100 {
			fmt.Printf("%04X\n", addr)

			// single file for ghidra
			if single {
				if addr >= 0xFF00 {
					os.WriteFile("dump.bin", slice, os.FileMode(0644))
					return
				}
				addr += 0x0100
				continue
			}

			if addr >= 0xFF00 {
				os.WriteFile("0xD000.bin", slice, os.FileMode(0644))
				break
			}
			if addr >= 0xD000 {
				addr += 0x0100
				continue
			}
			if need2save[int(addr)] {
				if addr == 0x1000 {
					os.WriteFile(fmt.Sprintf("0x%04X.bin", addr), slice[:0x40], os.FileMode(0644))
				} else {
					os.WriteFile(fmt.Sprintf("0x%04X.bin", addr), slice, os.FileMode(0644))
				}
			} //fmt.Println(hex.EncodeToString(slice))
			slice = make([]byte, 0)
			addr += 0x0100

		}
	}
	found0x10 := false
	for i := 0; i < len(slice)-2; i++ {
		b := slice[i]
		if found0x10 {
			switch b {
			case 0x3B:
				fmt.Printf("?PPROG %03X %s\n", i-1, hex.EncodeToString(slice[i-5:i+5]))
			case 0x35:
				fmt.Printf("?BPROT %03X %s\n", i-1, hex.EncodeToString(slice[i-5:i+5]))
			}
		}
		found0x10 = STAx[b]
	}
}

func ack_once(port serial.Port) {
	time.Sleep(1 * time.Second)

	port.ResetInputBuffer()
	port.ResetOutputBuffer()

	r := make([]byte, 1)

	port.Write([]byte{0x81})
	n, err := port.Read(r)
	log.Println(n, err, r)
	/*
	   port.Write([]byte{0x82})
	   n, err = port.Read(r)
	   log.Println(n, err, r)

	   port.Write([]byte{0x83})
	   n, err = port.Read(r)
	   log.Println(n, err, r)
	*/
}
