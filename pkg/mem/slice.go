package mem

import (
	"log"
	"time"

	"go.bug.st/serial"
)

func SendSlice(port serial.Port, sl string) error {
	log.Println("trying", sl)

	port.SetReadTimeout(100 * time.Millisecond)

	tx := make([]byte, 0)
	ack := make([]byte, 0)
	r := make([]byte, 1)

	for i := 0; i < len(sl); i++ {
		b := byte(sl[i])
		tx = append(tx, b)
		if _, err := port.Write([]byte{b}); err != nil {
			log.Println("TX issue at", tx)
			return err
		}
		n, err := port.Read(r)
		if err != nil || n != 1 {
			log.Println("RX issue")
			log.Println(tx)
			log.Println(ack)
			//return err
		}
		ack = append(ack, r...)
	}
	// todo: CR ?
	return nil
}
