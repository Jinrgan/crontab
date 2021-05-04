package kv

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

type JobKV struct {
	Key   []byte
	Value []byte
}

type Watcher struct {
	Watcher clientv3.Watcher
}

func (w *Watcher) Watch(ctx context.Context, dir string, rev int64, et mvccpb.Event_EventType) chan *JobKV {
	kvCh := make(chan *JobKV)
	wCh := w.Watcher.Watch(ctx, dir, clientv3.WithRev(rev), clientv3.WithPrefix())
	go func() {
		for resp := range wCh {
			for _, event := range resp.Events {
				if event.Type == et {
					kvCh <- &JobKV{
						Key:   event.Kv.Key,
						Value: event.Kv.Value,
					}
				}
			}
		}
		close(kvCh)
	}()

	return kvCh
}
