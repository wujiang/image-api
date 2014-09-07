package main

import (
	"encoding/json"
	"log"
	"os"
)

var (
	config_file = "config.json"
	cfg         = Configuration{}
)

type Configuration struct {
	TempDir      string
	LogTo        string
	JpegQuality  int
	SharpenSigma float64
}

func ReadConfig(filename string) (err error) {
	if filename == "" {
		filename = config_file
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)

	return
}
