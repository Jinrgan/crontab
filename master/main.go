package main

import (
	"crontab/master/dao"
	"crontab/master/master"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type config struct {
	HandlerPort         string   `json:"handler_port"`
	HandlerReadTimeout  int      `json:"handler_read_timeout"`
	HandlerWriteTimeout int      `json:"handler_write_timeout"`
	EtcdEndPoints       []string `json:"etcd_end_points"`
	EtcdDialTimeout     int      `json:"etcd_dial_timeout"`
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

	// etcd
	clt, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.EtcdEndPoints,
		DialTimeout: time.Duration(conf.EtcdDialTimeout) * time.Millisecond,
	})
	if err != nil {
		logger.Fatal("cannot connect etcd", zap.Error(err))
	}

	// create http service
	lis, err := net.Listen("tcp", conf.HandlerPort)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	h := master.Handler{
		DB: &dao.Etcd{
			KV:     clientv3.NewKV(clt),
			Lease:  clientv3.NewLease(clt),
			JobKey: "/cron/job",
		},
		Logger: logger,
	}

	w := master.Wrapper{Logger: logger}
	http.HandleFunc("/job/save", w.WrapErr(h.SaveJob))

	s := &http.Server{
		Handler:      http.DefaultServeMux,
		ReadTimeout:  time.Duration(conf.HandlerReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(conf.HandlerWriteTimeout) * time.Millisecond,
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
