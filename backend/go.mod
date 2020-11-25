module github.com/hennedo/godnc

go 1.14

replace github.com/256dpi/gcode => github.com/hennedo/gcode v0.4.0

require (
	github.com/256dpi/gcode v0.2.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/markbates/pkger v0.17.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	go.bug.st/serial v1.1.0
	gorm.io/driver/sqlite v1.1.3
	gorm.io/gorm v1.20.2
)
