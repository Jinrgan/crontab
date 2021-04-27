package dao

import (
	"context"
	"crontab/master/api"
	etcdtesting "crontab/shared/etcd/testing"
	"github.com/coreos/etcd/clientv3"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"os"
	"testing"
)

func TestEtcd_CreateJob(t *testing.T) {
	ctx := context.Background()
	clt, err := etcdtesting.NewClient(ctx)
	if err != nil {
		t.Fatalf("cannot connect etcd: %v", err)
	}

	e := Etcd{
		KV:     clientv3.NewKV(clt),
		Lease:  clientv3.NewLease(clt),
		JobKey: "/test",
	}

	cases := []struct {
		name    string
		jobName string
		command string
		want    *api.Job
		wantErr bool
	}{
		{
			name:    "new_job",
			jobName: "job_1",
			command: "echo hello1",
		},
		{
			name:    "exist_job",
			jobName: "job_1",
			command: "echo hello2",
			want: &api.Job{
				Name:    "job_1",
				Command: "echo hello1",
			},
		},
		{
			name:    "another_new_job",
			jobName: "job_2",
			command: "echo hello3",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			got, err := e.CreateJob(ctx, &api.Job{
				Name:     cc.jobName,
				Command:  cc.command,
				CronExpr: "",
			})
			if cc.wantErr {
				if err == nil {
					t.Error("want err; got none")
				}

				return
			}
			if err != nil {
				t.Errorf("cannot create job: %v", err)
			}
			if diff := cmp.Diff(cc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("result differs; -want +got: %s", diff)
			}
		})
	}
}

func TestMain(m *testing.M) {
	os.Exit(etcdtesting.RunWithEtcdInDocker(m))
}
