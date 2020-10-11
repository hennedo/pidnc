package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hennedo/godnc/config"
	"github.com/hennedo/godnc/database"
	"github.com/hennedo/godnc/file"
	"github.com/hennedo/godnc/serial"
	"github.com/markbates/pkger"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var db *gorm.DB

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		FullTimestamp:             true,
	})
	flag.Int("port", 8000, "Port where the shop listens on")
	flag.String("host", "", "IP to bind to")
	flag.String("gcode-folder", "", "Folder to store gcode in.")
	flag.String("serial-port", "", "Serial Port Path")
	_ = viper.BindPFlags(flag.CommandLine)
	flag.Parse()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	config.Config = config.SConfig{
		Port:           viper.GetInt("port"),
		Host:           viper.GetString("host"),
		GCodeFolder:    strings.TrimSuffix(viper.GetString("gcode-folder"), "/"),
		SerialPort:    	viper.GetString("serial-port"),
	}

	file.InitFolder(config.Config.GCodeFolder)
	db = database.InitDatabase(config.Config.GCodeFolder)
	file.SyncFiles()
	file.RenderAll()
	p := serial.InitSerial(config.Config.SerialPort)
	defer p.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/files", ApiFilesHandler).Methods("GET")
	r.HandleFunc("/api/upload", ApiUploadHandler).Methods("POST")
	r.HandleFunc("/api/{id}/run", ApiRunHandler).Methods("GET")
	r.HandleFunc("/api/{id}/delete", ApiDeleteHandler).Methods("GET")
	r.HandleFunc("/api/{id}/lock", ApiLockHandler).Methods("GET")
	r.HandleFunc("/api/info", ApiInfoHandler).Methods("GET")
	r.PathPrefix("/svg").Handler(http.StripPrefix("/svg", http.FileServer(http.Dir(config.Config.GCodeFolder + "/svg"))))
	r.PathPrefix("").Handler(http.FileServer(pkger.Dir("/static/")))
	http.Handle("/", r)

	logrus.Info(fmt.Sprintf("Listening on %s:%d", config.Config.Host, config.Config.Port))
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port), nil))
}

func ApiInfoHandler(w http.ResponseWriter, r *http.Request) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
			continue
		}
		for _, a := range addrs {
			switch a.(type) {
			case *net.IPNet:
				ip := a.(*net.IPNet).IP.To4()
				if ip == nil || ip.IsLoopback() {
					continue
				}
				fmt.Printf("%v : %s\n", i.Name, ip)
			}

		}
	}
}

func ApiFilesHandler(w http.ResponseWriter, r *http.Request) {
	var files []database.File
	err := db.Find(&files, database.File{}).Error
	if err != nil {
		logrus.Error(err)
	}
	json.NewEncoder(w).Encode(files)
}

func ApiDeleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("wrong id"))
		return
	}
	var file database.File
	if err := db.Find(&file, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(404)
			w.Write([]byte("not found"))
		} else {
			logrus.Error(err)
			w.WriteHeader(500)
			w.Write([]byte("stuff happened"))
		}
		return
	}
	if file.Type == "locked" {
		w.WriteHeader(500)
		w.Write([]byte("file is locked"))
		return
	}
	os.Remove(fmt.Sprintf("%s/%s/%s", config.Config.GCodeFolder, file.Type, file.Name))
	os.Remove(fmt.Sprintf("%s/svg/%d.svg", config.Config.GCodeFolder, file.ID))
	db.Delete(&file)
	w.Write([]byte("file deleted"))
}

func ApiLockHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("wrong id"))
		return
	}
	var file database.File
	if err := db.Find(&file, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(404)
			w.Write([]byte("not found"))
		} else {
			logrus.Error(err)
			w.WriteHeader(500)
			w.Write([]byte("stuff happened"))
		}
		return
	}
	if file.Type == "locked" {
		w.WriteHeader(500)
		w.Write([]byte("file was already locked"))
		return
	}
	os.Rename(fmt.Sprintf("%s/%s/%s", config.Config.GCodeFolder, file.Type, file.Name), fmt.Sprintf("%s/%s/%s", config.Config.GCodeFolder, "locked", file.Name))
	file.Type = "locked"
	db.Save(&file)
	w.Write([]byte("file locked"))
}

func ApiUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	upload, handler, err := r.FormFile("gcode")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer upload.Close()
	f, err := os.OpenFile(fmt.Sprintf("%s/uploaded/%s", config.Config.GCodeFolder, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, upload)
	var dbFile = database.File{
		Type:  "uploaded",
		Name:  handler.Filename,
		Size:  handler.Size,
	}
	db.Create(&dbFile)
	file.RenderFile(dbFile)
	json.NewEncoder(w).Encode(dbFile)
}

func ApiRunHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var file database.File
	err := db.Find(&file, params["id"]).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	f, _ := os.Open(fmt.Sprintf("%s/%s/%s", config.Config.GCodeFolder, file.Type, file.Name))
	stat, _ := f.Stat()
	var pos int64 = 0
	for {
		logrus.Infof("%d %d", pos, stat.Size())
		if pos >= stat.Size() {
			break
		}
		buf := make([]byte, 10)
		f.Read(buf)
		n, err := serial.Write(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %v bytes\n", n)
		pos += 10
	}
	if file.Type == "uploaded" {
		file.Type = "machined"
		os.Rename(fmt.Sprintf("%s/uploaded/%s", config.Config.GCodeFolder, file.Name), fmt.Sprintf("%s/machined/%s", config.Config.GCodeFolder, file.Name))
		db.Save(&file)
	}
	http.Redirect(w, r, "/", 302)
}
