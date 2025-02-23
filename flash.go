package main

import (
	"log"

	"github.com/dmisol/gru3880/pkg/asm"
	"github.com/dmisol/gru3880/pkg/boot"
	"go.bug.st/serial"
)

func main() {

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(ports)

	port, err := boot.Bootload(ports[0], asm.Write())
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

	r := make([]byte, 1)
	for {
		n, err := port.Read(r)
		if err != nil {
			log.Println("rd error: ", err, n)
			return
		}
	}
}
