package main

import (
	"crontab/master/dao"
	"crontab/master/master"
	"crontab/shared/server"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/coreos/etcd/clientv3"
	"go.uber.org/zap"
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

	var conf config
	err = server.GetConfig("master/config.json",
		func(m func([]byte, interface{}) error, b []byte) error {
			err := m(b, &conf)
			if err != nil {
				return fmt.Errorf("cannot unmarshal: %v", err)
			}
			return nil
		})
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
			JobSaveDir: server.JobSaveDir,
			JobKillDir: server.JobKillDir,
			KV:         clientv3.NewKV(clt),
			Lease:      clientv3.NewLease(clt),
			Logger:     logger,
		},
		Logger: logger,
	}
	h.Register()
	http.Handle("/", http.FileServer(http.Dir("./master/front")))

	s := &http.Server{
		Handler:      http.DefaultServeMux,
		ReadTimeout:  time.Duration(conf.HandlerReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(conf.HandlerWriteTimeout) * time.Millisecond,
	}
	logger.Info("server started", zap.String("name", "user"), zap.String("Addr", conf.HandlerPort))
	logger.Sugar().Fatal(s.Serve(lis))
}
