package main

import (
	"crontab/shared/server"
	"crontab/worker/dao"
	"crontab/worker/dao/kv"
	"crontab/worker/watcher/job"
	"crontab/worker/worker"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"go.uber.org/zap"
)

// 程序配置
type config struct {
	EtcdEndPoints   []string `json:"etcd_end_points"`   // etcd 的集群列表：配置多个，避免单点故障
	EtcdDialTimeout int      `json:"etcd_dial_timeout"` // etcd 的连接超时：单位亳秒
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	var conf config
	err = server.GetConfig("worker/config.json",
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

	kvWch := kv.Watcher{Watcher: clientv3.NewWatcher(clt)}
	sch := &worker.Scheduler{
		DB: &dao.Etcd{
			Dir:    server.JobSaveDir,
			KV:     clientv3.NewKV(clt),
			Lease:  clientv3.NewLease(clt),
			Logger: logger,
		},
		PutJobWatcher: &job.PutWatcher{
			Dir:     server.JobSaveDir,
			Watcher: kvWch,
			Logger:  logger,
		},
		DeleteJobWatcher: &job.DeleteWatcher{
			Dir:     server.JobSaveDir,
			Watcher: kvWch,
		},
	}

	logger.Sugar().Error(sch.Run())
}
