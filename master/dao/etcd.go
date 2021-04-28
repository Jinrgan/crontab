package dao

import (
	"context"
	"crontab/master/api"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"go.uber.org/zap"
)

type Etcd struct {
	KV     clientv3.KV
	Lease  clientv3.Lease
	JobKey string
	logger *zap.Logger
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

func (e *Etcd) GetJobs(ctx context.Context) ([]*api.Job, error) {
	resp, err := e.KV.Get(ctx, e.JobKey, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("cannot get jobs: %v", err)
	}

	var jobs []*api.Job
	for _, kv := range resp.Kvs {
		var job api.Job
		err := json.Unmarshal(kv.Value, &job)
		if err != nil {
			e.logger.Error("cannot unmarshal job", zap.Error(err))
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (e *Etcd) DeleteJob(ctx context.Context, jobName string) (*api.Job, error) {
	resp, err := e.KV.Delete(ctx, e.JobKey+jobName, clientv3.WithPrevKV())
	if err != nil {
		return nil, fmt.Errorf("cannot delete job: %v", err)
	}

	if resp.PrevKvs == nil {
		return nil, nil
	}

	var oldJ api.Job
	err = json.Unmarshal(resp.PrevKvs[0].Value, &oldJ)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal: %v", err)
	}

	return &oldJ, nil
}
