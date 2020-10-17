package config

type SConfig struct {
	Port            int `json:"port"`
	Host            string `json:"host"`
	GCodeFolder     string `json:"gCodeFolder"`
	SerialPort		string `json:"serialPort"`
}

var Config SConfig
