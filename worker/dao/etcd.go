package dao

import (
	"context"
	"crontab/shared/model"
	"encoding/json"
	"fmt"

	"github.com/coreos/etcd/clientv3"
	"go.uber.org/zap"
)

type Etcd struct {
	Dir    string
	KV     clientv3.KV
	Lease  clientv3.Lease
	Logger *zap.Logger
}

func (e *Etcd) GetJobsWithRev() (int64, []*model.Job, error) {
	resp, err := e.KV.Get(context.Background(), e.Dir, clientv3.WithPrefix())
	if err != nil {
		return 0, nil, fmt.Errorf("cannot get kv: %v", err)
	}

	var jobs []*model.Job
	for _, pair := range resp.Kvs {
		var job model.Job
		err := json.Unmarshal(pair.Value, &job)
		if err != nil {
			e.Logger.Error("cannot unmarshal job", zap.Any("job", pair.Value), zap.Error(err))
		}
		jobs = append(jobs, &job)
	}

	return resp.Header.Revision + 1, jobs, nil
}
