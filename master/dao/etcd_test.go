package dao

import (
	"context"
	etcdtesting "crontab/shared/etcd/testing"
	"crontab/shared/model"
	"encoding/json"
	"os"
	"testing"

	"github.com/coreos/etcd/clientv3"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestEtcd_CreateJob(t *testing.T) {
	ctx := context.Background()
	clt, err := etcdtesting.NewClient(ctx)
	if err != nil {
		t.Fatalf("cannot connect etcd: %v", err)
	}

	e := Etcd{
		JobSaveDir: "/test/",
		KV:         clientv3.NewKV(clt),
		Lease:      clientv3.NewLease(clt),
	}

	cases := []struct {
		name    string
		jobName string
		command string
		want    *model.Job
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
			want: &model.Job{
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
			got, err := e.CreateJob(ctx, &model.Job{
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

func TestEtcd_DeleteJob(t *testing.T) {
	ctx := context.Background()
	clt, err := etcdtesting.NewClient(ctx)
	if err != nil {
		t.Fatalf("cannot connect etcd: %v", err)
	}

	e := Etcd{
		JobSaveDir: "/test/",
		KV:         clientv3.NewKV(clt),
		Lease:      clientv3.NewLease(clt),
	}

	jobs := []*model.Job{
		{
			Name:    "job_1",
			Command: "echo hello1",
		},
		{
			Name:    "job_2",
			Command: "echo hello2",
		},
	}

	for _, j := range jobs {
		b, err := json.Marshal(j)
		if err != nil {
			t.Fatalf("cannot marshal job: %v", err)
		}

		_, err = e.KV.Put(ctx, e.JobSaveDir+j.Name, string(b))
		if err != nil {
			t.Fatalf("cannot put job: %v", err)
		}
	}

	cases := []struct {
		name    string
		jobName string
		want    *model.Job
		wantErr bool
	}{
		{
			name:    "exist_job",
			jobName: "job_1",
			want: &model.Job{
				Name:    "job_1",
				Command: "echo hello1",
			},
		},
		{
			name:    "another_exist_job",
			jobName: "job_2",
			want: &model.Job{
				Name:    "job_2",
				Command: "echo hello2",
			},
		},
		{
			name:    "not_exist_job",
			jobName: "bad_job",
			wantErr: true,
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			got, err := e.DeleteJob(ctx, cc.jobName)
			if cc.wantErr {
				if err == nil {
					t.Error("want err; got none")
				}

				return
			}
			if err != nil {
				t.Errorf("cannot delete job: %v", err)
			}

			if diff := cmp.Diff(cc.want, got); diff != "" {
				t.Errorf("result differs; -want +got: %s", diff)
			}
		})
	}
}

func TestEtcd_GetJobs(t *testing.T) {
	ctx := context.Background()
	clt, err := etcdtesting.NewClient(ctx)
	if err != nil {
		t.Fatalf("cannot connect etcd: %v", err)
	}

	e := Etcd{
		JobSaveDir: "/test/",
		KV:         clientv3.NewKV(clt),
		Lease:      clientv3.NewLease(clt),
	}

	jobs := []*model.Job{
		{
			Name:    "job_1",
			Command: "echo hello1",
		},
		{
			Name:    "job_2",
			Command: "echo hello2",
		},
	}

	for _, j := range jobs {
		b, err := json.Marshal(j)
		if err != nil {
			t.Fatalf("cannot marshal job: %v", err)
		}

		_, err = e.KV.Put(ctx, e.JobSaveDir+j.Name, string(b))
		if err != nil {
			t.Fatalf("cannot put job: %v", err)
		}
	}

	got, err := e.GetJobs(ctx)
	if err != nil {
		t.Errorf("cannot get jobs: %v", err)
	}

	if diff := cmp.Diff(jobs, got, protocmp.Transform()); diff != "" {
		t.Errorf("result differs; -want +got: %s", diff)
	}
}

func TestMain(m *testing.M) {
	os.Exit(etcdtesting.RunWithEtcdInDocker(m))
}
