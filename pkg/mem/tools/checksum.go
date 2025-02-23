package tools

import "fmt"

func MkData(addr uint16, data []byte) string {
	l := 2 + len(data) + 1
	s := fmt.Sprintf("S1%02X%04X", l, addr)
	sum := l
	sum += int(addr >> 8)
	sum += int(addr)

	for i := 0; i < len(data); i++ {
		s += fmt.Sprintf("%02X", data[i])
		sum += int(data[i])
	}

	s += fmt.Sprintf("%02X", (sum^0xFF)&0xFF)

	return s
}

func MkFinal(_ uint16) string {
	return fmt.Sprintf("S9030000FC")
}
