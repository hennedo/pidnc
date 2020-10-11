package serial

import (
	"fmt"
	"go.bug.st/serial"
	"log"
)

var port serial.Port

func InitSerial(path string) serial.Port {
	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 7,
		Parity:   serial.EvenParity,
		StopBits: serial.OneStopBit,
	}
	var err error
	port, err = serial.Open(path, mode)
	if err != nil {
		log.Fatal(err)
	}
	go ReadSerial()

	return port
}

func Write(data []byte) (int, error) {
	return port.Write(data)
}

func ReadSerial() {
	buff := make([]byte, 100)
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		fmt.Printf("%v", string(buff[:n]))
	}
}
