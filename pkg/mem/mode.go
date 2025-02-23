package mem

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

func SetMode(port serial.Port, mode string) error {
	time.Sleep(1 * time.Second)

	port.ResetInputBuffer()
	port.ResetOutputBuffer()

	for i := 0; i < 3; i++ {
		t := []byte{byte(mode[0]), 0x0D}
		if _, err := port.Write(t); err != nil {
			log.Println("mode write error: ", i, err)
			continue
		}
		r := make([]byte, 2)

		n, err := port.Read(r)
		if err != nil {
			log.Println("mode ack error: ", i, err)
			continue
		}
		log.Println("rd1", r)
		if n == 1 {
			x := make([]byte, 1)
			if _, err = port.Read(x); err != nil {
				log.Println("mode ack2 error: ", i, err)
				continue
			}
			log.Println("rd2", x)
			r[1] = x[0]
		}
		if bytes.Equal(t, r) {
			return nil
		}
		log.Println("not equal:", i, t, r)
	}
	return fmt.Errorf("can't set mode")
}
