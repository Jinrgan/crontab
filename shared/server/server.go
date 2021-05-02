package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

const (
	JobSaveDir = "/cron/job/"
	JobKillDir = "/cron/killer/"
)

func GetConfig(param string, unmarshal func(func([]byte, interface{}) error, []byte) error) error {
	var file string
	flag.StringVar(&file, "config", param, "specify config file")
	flag.Parse()

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("cannot read file: %v", err)
	}

	return unmarshal(json.Unmarshal, bytes)
}
