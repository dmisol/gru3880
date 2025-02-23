package boot

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

const nbl = 512

func Bootload(dev string, code string) (serial.Port, error) {
	//ldr, _ := hex.DecodeString(str)
	ldr, err := hex.DecodeString(code)
	if err != nil {
		return nil, err
	}

	mode := &serial.Mode{
		BaudRate: 1200,
	}

	port, err := serial.Open(dev, mode)
	if err != nil {
		log.Fatal(err)
	}

	port.ResetInputBuffer()
	port.ResetOutputBuffer()

	t := make([]byte, 1)
	r := make([]byte, 1)

	// sync rate
	if _, err := port.Write([]byte{0xFF}); err != nil {
		log.Println("wr initial 0xFF error: ", err)
		return nil, err
	}

	port.SetReadTimeout(time.Second)
	// read ack
	n, err := port.Read(r)
	if err != nil || n != 1 {
		log.Println("initial rd error: ", err, n)
		return nil, fmt.Errorf("initial ACK failed")
	}

	w := 0
	for {

		if len(ldr) > w {
			t[0] = ldr[w]
		} else {
			t[0] = 0xFF
		}
		if n, err := port.Write(t); err != nil || n != 1 {
			log.Printf("wr[%03X] %02X %d %v", w, t[0], n, err)
			return nil, err
		}

		n, err := port.Read(r)
		if err != nil {
			log.Println("rd error: ", err, n)
			break
		}
		if n != 1 {
			log.Printf("no ack at %03X", w)
		}

		if t[0] != r[0] {
			log.Printf("mismatch: rd %02X", r[0])
			return nil, fmt.Errorf("ack mismatch: %03X %02X", w, r[0])
		}
		w++
		if w >= nbl {
			log.Println("boolloader sent")
			return port, nil
		}
	}
	if w < len(ldr)-2 {
		log.Println("last byte not ack'd")
	}
	return port, nil
}
