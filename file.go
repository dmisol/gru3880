package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dmisol/gru3880/pkg/asm"
	"github.com/dmisol/gru3880/pkg/boot"
	"go.bug.st/serial"
)

func main() {
	file := flag.String("file", "bin/0322n.bin", "file to flash")
	erase := flag.Bool("erase", true, "wipe flash")
	addr := flag.String("addr", "B600", "starting address")
	mask := flag.String("mask", "F0", "write protection mask")
	flag.Parse()

	fmt.Println(*file, *erase, *addr, *mask)

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(ports)

	port, err := boot.Bootload(ports[0], asm.EraseAndFlash(*mask, *addr, *erase))
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

	// time.Sleep(time.Second)
	if *erase {
		read10(port)
	}

	log.Println(*file)

	b, err := os.ReadFile(*file)
	log.Println("wrining", len(b))
	for i := 0; i < len(b); i++ {
		if n, err := port.Write(b[i : i+1]); err != nil || n != 1 {
			log.Printf("wr[%03X] %02X %d %v", 0xB600+i, b[i:i+1], n, err)
			return
		}
		read10(port)
	}
}

func read10(port serial.Port) {
	r := make([]byte, 1)
	slice := make([]byte, 0)

	for i := 0; i < 10; i++ {
		n, err := port.Read(r)
		if err != nil || n != 1 {
			log.Println("rd error: ", err, slice)
			return
		}
		slice = append(slice, r[0])
	}
	log.Println(slice)
}
