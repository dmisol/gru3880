package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

func main() {
	file := flag.String("file", "0xD000.bin", "file to read from")
	seq := flag.String("seq", "A735", "sequence")
	flag.Parse()

	b, err := os.ReadFile(*file)
	if err != nil {
		fmt.Println(err)
		return
	}

	h, err := hex.DecodeString(*seq)
	if err != nil {
		fmt.Println(err)
		return
	}
	for s := 0; s < len(b)-2; s++ {
		if b[s] == h[0] {
			fine := true
			for i := 0; i < len(h); i++ {
				if h[i] != b[s+i] {
					fine = false
					break
				}
			}
			if fine {
				fmt.Printf("%03X %s %s %s\n", s,
					hex.EncodeToString(b[s-6:s]),
					hex.EncodeToString(b[s:s+len(h)]),
					hex.EncodeToString(b[s+len(h):s+len(h)+6]))
			}
		}
	}
}
