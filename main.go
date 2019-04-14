package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/danielingegneri/megatec-ups/megatec"
	"github.com/jacobsa/go-serial/serial"
	"log"
)

func main() {
	var (
		dev       = flag.String("dev", "/dev/ttyUSB0", "Serial device that has a Megatec UPS connected to it")
		baud      = flag.Uint("baud", 2400, "Baud rate in bits")
		dataBits  = flag.Uint("dbits", 8, "Data bits")
		stopBits  = flag.Uint("sbits", 1, "Stop bits")
		parityInt = flag.Uint("parity", 0, "Parity: 0, 1 or 2")
	)
	flag.Parse()
	//ctx := context.Background()

	var parity serial.ParityMode
	switch *parityInt {
	case 0:
		parity = serial.PARITY_NONE
	case 1:
		parity = serial.PARITY_ODD
	case 2:
		parity = serial.PARITY_EVEN
	default:
		log.Fatal("Parity must be 0, 1 or 2")
	}

	ups := megatec.NewUPS(*dev, *baud, *dataBits, *stopBits, parity)
	defer ups.Close()

	query, err := ups.Query()
	if err != nil {
		log.Fatal(err)
	}
	jsonString, err := json.Marshal(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonString))
}
