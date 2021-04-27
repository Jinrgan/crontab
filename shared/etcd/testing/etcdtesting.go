package etcdtesting

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"testing"
)

const (
	image         = "etcd:3.3.9"
	containerPort = "2379/tcp"
)

var etcdEndPoint string

const defaultEtcdEndPoint = ":2379"

//RunWithEtcdInDocker runs the tests with
// a etcd instance in a docker container.
func RunWithEtcdInDocker(m *testing.M) int {
	c, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	resp, err := c.ContainerCreate(ctx, &container.Config{
		ExposedPorts: nat.PortSet{
			containerPort: {},
		},
		Image: image,
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0",
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	containerID := resp.ID
	defer func() {
		err := c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			panic(err)
		}
	}()

	err = c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	ispRes, err := c.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)
	}

	hostPort := ispRes.NetworkSettings.Ports[containerPort][0]
	etcdEndPoint = fmt.Sprintf("%s:%s", hostPort.HostIP, hostPort.HostPort)

	return m.Run()
}

//NewClient creates a client connected to the etcd instance in docker.
func NewClient(c context.Context) (*clientv3.Client, error) {
	if etcdEndPoint == "" {
		return nil, fmt.Errorf("etcd end point not set, Please run RunWithEtcdInDocker in TestMain")
	}

	return clientv3.New(clientv3.Config{
		Endpoints: []string{etcdEndPoint},
	})
}

//NewDefaultClient creates a client connected to localhost:2379
func NewDefaultClient(c context.Context) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints: []string{defaultEtcdEndPoint},
	})
}
