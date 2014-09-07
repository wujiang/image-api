package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := flag.String("p", "8001", "port")
	configFile := flag.String("c", "", "Path to the configuration file")
	flag.Parse()

	err := ReadConfig(*configFile)
	if err != nil {
		glog.Fatal(err)
	}

	if _, err := os.Stat(cfg.TempDir); err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(cfg.TempDir, 0755); err != nil {
			glog.Fatal(err)
		}
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	r := mux.NewRouter()
	r.HandleFunc("/api/resize", ResizerHandler)
	http.Handle("/", r)

	logFile, err := os.OpenFile(cfg.LogTo, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		logFile = os.Stdout
	} else {
		defer logFile.Close()
	}

	loggingHandler := handlers.LoggingHandler(logFile, r)

	server := &http.Server{
		Addr:         ":" + *port,
		Handler:      loggingHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Serves at localhost" + server.Addr)
	glog.Fatal(server.ListenAndServe())
}
