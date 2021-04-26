package main

import (
	"crontab/master/master"
	"encoding/json"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type config struct {
	Port         string `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	var file string
	flag.StringVar(&file, "config", "master/config.json", "specify config file")
	flag.Parse()

	conf, err := getConfig(file)
	if err != nil {
		logger.Fatal("cannot get config", zap.Error(err))
	}

	// create http service
	lis, err := net.Listen("tcp", conf.Port)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	http.HandleFunc("job/save", master.HandleJobSave)

	s := &http.Server{
		Handler:      http.DefaultServeMux,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
	}

	err = s.Serve(lis)
	if err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}

func getConfig(f string) (*config, error) {
	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %v", err)
	}

	var conf config
	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal: %v", err)
	}

	return &conf, nil
}
