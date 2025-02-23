package mem

import (
	"io"
	"log"
	"os"

	"github.com/dmisol/gru3880/pkg/mem/tools"
	"go.bug.st/serial"
)

func SendFile(port serial.Port, addr uint16, sz int, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	offs := 0

	x := make([]byte, sz)
	for {
		n, err := f.Read(x)
		if err == io.EOF {
			log.Println("file DONE")
			break
		}
		if err != nil {
			return err
		}
		if n != sz {
			log.Println("got only", n)
			x = x[:n]
		}

		sl := tools.MkData(addr+uint16(offs), x)
		offs += n

		if err = SendSlice(port, sl); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return SendSlice(port, "S9030000FC")
}
