package config

type SConfig struct {
	Port            int
	Host            string
	GCodeFolder     string
	SerialPort		string
}

var Config SConfig
