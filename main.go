package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.bug.st/serial"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)


type Config struct {
	Port            int
	Host            string
	GCodeFolder     string
	SerialPort		string
}

type GCode struct {
	Path 		string
	Name		string
	MD5 		string
}

var config Config
var tmpl *template.Template
var port serial.Port

func compileTemplates() {
	tmpl = template.New("")
	pkger.Walk("/templates", func(path string, info os.FileInfo, _ error) error {
		if !strings.HasSuffix(path, ".gohtml") {
			return nil
		}
		f, _ := pkger.Open(path)
		sl, _ := ioutil.ReadAll(f)
		_, err := tmpl.New(info.Name()).Parse(string(sl))
		if err != nil {
			logrus.Error(err)
		}
		return nil
	})
}

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

	config = Config{
		Port:           viper.GetInt("port"),
		Host:           viper.GetString("host"),
		GCodeFolder:    viper.GetString("gcode-folder"),
		SerialPort:    	viper.GetString("serial-port"),
	}
	compileTemplates()

	mode := &serial.Mode{
		BaudRate: 115200,
	}
	var err error
	port, err = serial.Open(config.SerialPort, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()
	go ReadSerial()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/upload", UploadHandler).Methods("POST")
	r.HandleFunc("/{id}/run", RunHandler).Methods("GET")
	//r.HandleFunc("/{id}/preview", RefundConfirmHandler).Methods("GET")
	//r.HandleFunc("/{id}/delete", RefundConfirmHandler).Methods("GET")
	r.PathPrefix("").Handler(http.FileServer(pkger.Dir("/static/")))
	http.Handle("/", r)

	logrus.Info(fmt.Sprintf("Listening on %s:%d", config.Host, config.Port))
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), nil))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var gcodes []GCode
	err := filepath.Walk(config.GCodeFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Info("foo")
			logrus.Fatal(err)
		}
		if info.IsDir() {
			return nil
		}
		gcodes = append(gcodes, GCode{
			Path: path,
			Name: info.Name(),
			MD5: fmt.Sprintf("%x", md5.Sum([]byte(info.Name()))),
		})
		return nil
	})
	err = tmpl.ExecuteTemplate(w, "index.gohtml", map[string]interface{}{
		"Files": gcodes,
	})
	if err != nil {
		logrus.Error(err)
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("gcode")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile(config.GCodeFolder+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	http.Redirect(w, r, "/", 302)
}

func RunHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	f, _ := os.Open(config.GCodeFolder + params["id"])
	stat, _ := f.Stat()
	var pos int64 = 0
	for {
		logrus.Infof("%d %d", pos, stat.Size())
		if pos >= stat.Size() {
			break
		}
		buf := make([]byte, 10)
		f.Read(buf)
		n, err := port.Write(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent %v bytes\n", n)
		pos += 10
	}
	http.Redirect(w, r, "/", 302)
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