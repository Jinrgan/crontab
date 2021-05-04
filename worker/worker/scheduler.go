package worker

import (
	"context"
	"crontab/shared/model"
	"crontab/worker/dao"
	"fmt"
)

type PutJobWatcher interface {
	Watch(ctx context.Context, rev int64) (chan *model.Job, error)
}

type DeleteJobWatcher interface {
	Watch(ctx context.Context, rev int64) chan string
}

type KillerWatcher interface {
	Watch(ctx context.Context, rev int64) (chan string, error)
}

type Scheduler struct {
	DB               *dao.Etcd
	PutJobWatcher    PutJobWatcher
	DeleteJobWatcher DeleteJobWatcher
	KillerWatcher    KillerWatcher
}

func (s *Scheduler) Run() error {
	rev, jobs, err := s.DB.GetJobsWithRev()
	if err != nil {
		return fmt.Errorf("cannot get jobs with rev: %v", err)
	}
	for _, job := range jobs {
		fmt.Println("init job:", job)
	}

	putCh, err := s.PutJobWatcher.Watch(context.Background(), rev)
	if err != nil {
		return fmt.Errorf("fail to watch put job: %v", err)
	}

	delCh := s.DeleteJobWatcher.Watch(context.Background(), rev)

	//killCh, err := s.KillerWatcher.Watch(context.Background(), rev)
	//if err != nil {
	//	return fmt.Errorf("fail to watch killer: %v", err)
	//}

	for {
		select {
		case job := <-putCh:
			fmt.Println("put job: ", job)
		case jobName := <-delCh:
			fmt.Println("deleted job: ", jobName)
			//case <-killCh:

		}
	}
}
