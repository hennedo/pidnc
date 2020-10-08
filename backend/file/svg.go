package file

import (
	"fmt"
	"github.com/256dpi/gcode"
	"github.com/hennedo/godnc/config"
	"github.com/hennedo/godnc/database"
	"github.com/sirupsen/logrus"
	"os"
)

func RenderAll() {
	var db = database.GetDatabase()
	var files []database.File
	err := db.Find(&files, database.File{}).Error
	if err != nil {
		logrus.Error("error getting files from database", err)
	}
	for _, file := range files {
		RenderFile(file)
	}
}

func RenderId(id int) {
	var db = database.GetDatabase()
	var file database.File
	err := db.Find(&file, id).Error
	if err != nil {
		logrus.Error("error getting file from database ", err)
	}
	RenderFile(file)
}

func RenderFile(file database.File) {
	f, err := os.Open(fmt.Sprintf("%s/%s/%s", config.Config.GCodeFolder, file.Type, file.Name))
	if err != nil {
		logrus.Error("error opening gcode ", err)
		return
	}
	defer f.Close()
	gf, pe := gcode.ParseFile(f)
	if pe != nil {
		logrus.Error("error creating svg ", pe)
	}
	svg := gcode.ConvertToSVG(gf)
	svgFile, serr := os.Create(fmt.Sprintf("%s/svg/%d.svg", config.Config.GCodeFolder, file.ID))
	if serr != nil {
		logrus.Error("error opening svg file ", serr)
		return
	}
	svgFile.Write([]byte(svg))
}
