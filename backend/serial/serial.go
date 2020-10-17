package serial

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

var port serial.Port
var readChannel chan bool

func InitSerial(path string) (serial.Port, error) {
	if port != nil {
		err := port.Close()
		if err != nil {
			logrus.Error(err)
		}
		port = nil
	}
	if readChannel != nil {
		readChannel = make(chan bool)
	}
	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 7,
		Parity:   serial.EvenParity,
		StopBits: serial.OneStopBit,
	}
	var err error
	port, err = serial.Open(path, mode)
	if err != nil {
		port = nil
		return nil, err
	}
	//go ReadSerial()

	return port, nil
}

func Write(data []byte) (int, error) {
	return port.Write(data)
}

func ReadSerial() {
	buff := make([]byte, 100)
	for {
		select {
			case <- readChannel:
				return
			default:
				n, err := port.Read(buff)
				if err != nil {
					logrus.Error(err)
					break
				}
				fmt.Printf("%v", string(buff[:n]))
		}
	}
}

func Close() {
	if port != nil {
		port.Close()
	}
}

func GetPortsList() ([]string, error) {
	return serial.GetPortsList()
}