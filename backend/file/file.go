package file

import (
	"fmt"
	"os"
)

var path string

func InitFolder(p string) {
	path = p
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			panic("gcode path does not exist")
		}
	}
	if !info.IsDir() {
		panic("gcode path needs to be a folder")
	}
	initSubFolder("uploaded")
	initSubFolder("machined")
	initSubFolder("locked")
	initSubFolder("svg")
}

func initSubFolder(subfolder string) {
	info, err := os.Stat(fmt.Sprintf("%s/%s", path, subfolder))
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(fmt.Sprintf("%s/%s", path, subfolder), os.ModePerm)
		}
	} else if !info.IsDir() {
		panic(subfolder + " in gcode path needs to be a folder")
	}
}
