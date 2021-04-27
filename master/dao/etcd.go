package dao

import (
	"context"
	"crontab/master/api"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

type Etcd struct {
	KV     clientv3.KV
	Lease  clientv3.Lease
	JobKey string
}

func (e *Etcd) CreateJob(ctx context.Context, j *api.Job) (*api.Job, error) {
	b, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("cannnot marshal: %v", err)
	}

	resp, err := e.KV.Put(ctx, e.JobKey+j.Name, string(b), clientv3.WithPrevKV())
	if err != nil {
		return nil, fmt.Errorf("cannot put kv: %v", err)
	}

	if resp.PrevKv == nil {
		return nil, nil
	}

	// 如果是更新，那么返回旧值
	var oldJ api.Job
	err = json.Unmarshal(resp.PrevKv.Value, &oldJ)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal previous kv: %v", err)
	}

	return &oldJ, nil
}
