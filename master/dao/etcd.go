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
	JobSaveDir string
	JobKillDir string
	KV         clientv3.KV
	Lease      clientv3.Lease
	Logger     *zap.Logger
}

func (e *Etcd) CreateJob(ctx context.Context, j *model.Job) (*model.Job, error) {
	b, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("cannnot marshal: %v", err)
	}

	resp, err := e.KV.Put(ctx, e.JobSaveDir+j.Name, string(b), clientv3.WithPrevKV())
	if err != nil {
		return nil, fmt.Errorf("cannot put kv: %v", err)
	}

	if resp.PrevKv == nil {
		return nil, nil
	}

	// 如果是更新，那么返回旧值
	var oldJ model.Job
	err = json.Unmarshal(resp.PrevKv.Value, &oldJ)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal previous kv: %v", err)
	}

	return &oldJ, nil
}

func (e *Etcd) GetJobs(ctx context.Context) ([]*model.Job, error) {
	resp, err := e.KV.Get(ctx, e.JobSaveDir, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("cannot get jobs: %v", err)
	}

	var jobs []*model.Job
	for _, kv := range resp.Kvs {
		var job model.Job
		err := json.Unmarshal(kv.Value, &job)
		if err != nil {
			e.Logger.Error("cannot unmarshal job", zap.Error(err))
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (e *Etcd) DeleteJob(ctx context.Context, jobName string) (*model.Job, error) {
	resp, err := e.KV.Delete(ctx, e.JobSaveDir+jobName, clientv3.WithPrevKV())
	if err != nil {
		return nil, fmt.Errorf("cannot delete job: %v", err)
	}

	if resp.PrevKvs == nil {
		return nil, fmt.Errorf("cannot delete unexist job: %v", err)
	}

	var oldJ model.Job
	err = json.Unmarshal(resp.PrevKvs[0].Value, &oldJ)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal: %v", err)
	}

	return &oldJ, nil
}

func (e *Etcd) KillJob(ctx context.Context, name string) error {
	k := e.JobKillDir + name

	// 让 worker 监听到一次 put 操作，创建一个租约让其稍后自动过期即可
	lsResp, err := e.Lease.Grant(ctx, 1)
	if err != nil {
		return fmt.Errorf("cannot grant lease: %v", err)
	}

	_, err = e.KV.Put(ctx, k, "", clientv3.WithLease(lsResp.ID))
	if err != nil {
		return fmt.Errorf("cannot put kv: %v", err)
	}

	return nil
}
