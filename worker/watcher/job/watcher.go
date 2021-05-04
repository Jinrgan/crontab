package job

import (
	"context"
	"crontab/shared/model"
	"crontab/worker/dao/kv"
	"encoding/json"
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"

	"go.uber.org/zap"
)

type PutWatcher struct {
	Dir string
	kv.Watcher
	Logger *zap.Logger
}

func (w *PutWatcher) Watch(ctx context.Context, rev int64) (chan *model.Job, error) {
	jobCh := make(chan *model.Job)
	rawCh := w.Watcher.Watch(ctx, w.Dir, rev, mvccpb.PUT)
	go func() {
		for jobKV := range rawCh {
			var job model.Job
			err := json.Unmarshal(jobKV.Value, &job)
			if err != nil {
				w.Logger.Error("cannot unmarshal job", zap.Any("job", jobKV.Value), zap.Error(err))
				continue
			}

			jobCh <- &job
		}
		close(rawCh)
	}()
	return jobCh, nil
}

type DeleteWatcher struct {
	Dir string
	kv.Watcher
}

func (w *DeleteWatcher) Watch(ctx context.Context, rev int64) chan string {
	nameCh := make(chan string)
	rawCh := w.Watcher.Watch(ctx, w.Dir, rev, mvccpb.DELETE)
	go func() {
		for jobKV := range rawCh {
			nameCh <- strings.TrimPrefix(string(jobKV.Key), w.Dir)
		}

		close(nameCh)
	}()

	return nameCh
}
