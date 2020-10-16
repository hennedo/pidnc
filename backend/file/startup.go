package file

import (
	"archive/zip"
	"fmt"
	"github.com/hennedo/godnc/config"
	"github.com/hennedo/godnc/database"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func SyncFiles() {
	// check if all database files are there
	var db = database.GetDatabase()
	var files []database.File
	err := db.Find(&files, database.File{}).Error
	if err != nil {
		logrus.Error("error getting files", err)
	}
	for _, file := range files {
		_, err := os.Stat(fmt.Sprintf("%s/%s/%s", path, file.Type, file.Name))
		if os.IsNotExist(err) {
			db.Delete(&database.File{}, file.ID)
		}
	}
	// check if new files are there
	_ = filepath.Walk(path + "/uploaded", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Fatal("error warlking uploaded dir", err)
		}
		if info.IsDir() {
			return nil
		}
		var file database.File
		if err := db.Where(&database.File{Type: "uploaded", Name: info.Name()}).First(&file).Error; err != nil {
			file.Name = info.Name()
			file.Type = "uploaded"
			file.Size = info.Size()
			db.Create(&file)
		}
		return nil
	})

	_ = filepath.Walk(path + "/machined", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Fatal("error walking machined dir", err)
		}
		if info.IsDir() {
			return nil
		}
		var file database.File
		if err := db.Where(&database.File{Type: "machined", Name: info.Name()}).First(&file).Error; err != nil {
			file.Name = info.Name()
			file.Type = "machined"
			file.Size = info.Size()
			db.Create(&file)
		}
		return nil
	})

	_ = filepath.Walk(path + "/locked", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Fatal("error walking locked dir", err)
		}
		if info.IsDir() {
			return nil
		}
		var file database.File
		if err := db.Where(&database.File{Type: "locked", Name: info.Name()}).First(&file).Error; err != nil {
			file.Name = info.Name()
			file.Type = "locked"
			file.Size = info.Size()
			db.Create(&file)
		}
		return nil
	})
}

func Backup(w io.Writer) {
	path := config.Config.GCodeFolder
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	_ = filepath.Walk(path + "/uploaded", zipFile(zipWriter, "uploaded"))
	_ = filepath.Walk(path + "/machined", zipFile(zipWriter, "machined"))
	_ = filepath.Walk(path + "/locked", zipFile(zipWriter, "locked"))
}

func zipFile(zipWriter *zip.Writer, subdir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Fatal("error walking locked dir", err)
		}
		if info.IsDir() {
			return nil
		}
		f, _ := os.Open(path)
		defer f.Close()
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			logrus.Error(err)
		}
		header.Name = fmt.Sprintf("%s/%s", subdir, info.Name())
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			logrus.Error(err)
		}
		_, err = io.Copy(writer, f)
		if err != nil {
			logrus.Error(err)
		}
		return nil
	}
}