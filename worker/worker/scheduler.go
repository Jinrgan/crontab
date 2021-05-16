package worker

import (
	"context"
	"crontab/shared/model"
	"crontab/worker/dao"
	"fmt"

	"go.uber.org/zap"
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

type JobPlan struct {
	executor *Executor
	cancelFn context.CancelFunc
	doneChan chan struct{}
}

type Scheduler struct {
	DB               *dao.Etcd
	JobPlanTable     map[string]*JobPlan
	LogCh            chan string
	PutJobWatcher    PutJobWatcher
	DeleteJobWatcher DeleteJobWatcher
	KillerWatcher    KillerWatcher
	Logger           *zap.Logger
}

func (s *Scheduler) Run() error {
	s.JobPlanTable = make(map[string]*JobPlan)

	rev, jobs, err := s.DB.GetJobsWithRev()
	if err != nil {
		return fmt.Errorf("cannot get jobs with rev: %v", err)
	}
	for _, job := range jobs {
		s.Logger.Info("init job", zap.Any("job info", job))

		exe, err := NewExecutor(job.Command, job.CronExpr)
		if err != nil {
			return fmt.Errorf("cannot create executor: %v", err)
		}
		s.JobPlanTable[job.Name] = &JobPlan{
			executor: exe,
		}

		ctx, cancelFunc := context.WithCancel(context.Background())
		s.JobPlanTable[job.Name].cancelFn = cancelFunc
		go exe.Execute(ctx)

		go func(name string) {
			for log := range s.JobPlanTable[name].executor.ResCh {
				s.Logger.Info("execute result", zap.String("Job", name), zap.Any("log", log))
			}
		}(job.Name)
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
			s.Logger.Info("put job", zap.Any("job info", job))
			// 判断是否是 create
			if _, ok := s.JobPlanTable[job.Name]; !ok {
				_, err := NewExecutor(job.Command, job.CronExpr)
				if err != nil {
					return fmt.Errorf("cannot create executor: %v", err)
				}

				//exe.Execute()
				//s.JobPlanTable[job.Name].executor = exe

			}
		case jobName := <-delCh:
			fmt.Println("deleted job: ", jobName)
			// done chan ang close chan
			//case <-killCh:

		}
	}
}
