package v1

import (
	"context"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/golang/protobuf/ptypes/empty"
	v1 "github.com/onuryartasi/scaler/pkg/api/v1"

	"encoding/json"
	"fmt"
	"strconv"

	"github.com/onuryartasi/scaler/pkg/metric"
	p2 "github.com/onuryartasi/scaler/pkg/types"
)

var (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

type Service struct{}
type Error struct{}

var projects []p2.Project
var cli *client.Client

func init() {
	var err error
	if len(os.Getenv("DOCKER_API_VERSION")) < 1 {
		os.Setenv("DOCKER_API_VERSION", "1.38")
	}

	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf(ErrorColor, "Error: %s", err)
	}
}

var timeout time.Duration = 10 * time.Second

func (s *Service) ContainerList(ctx context.Context, empty *empty.Empty) (*v1.Containers, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatalf(ErrorColor, "Error: %s", err)
	}

	rContainers := []*v1.Container{}
	for _, container := range containers {
		rContainers = append(rContainers, &v1.Container{Names: container.Names, Id: container.ID, Image: container.Image})
	}
	return &v1.Containers{Container: rContainers}, nil
}

func (s *Service) ContainerStop(ctx context.Context, containerId *v1.ContainerId) (*v1.ContainerId, error) {

	err := cli.ContainerStop(ctx, containerId.GetContainerId(), &timeout)
	return containerId, err
}

func (s *Service) ContainerStart(ctx context.Context, containerId *v1.ContainerId) (*v1.ContainerId, error) {

	err := cli.ContainerStart(ctx, containerId.GetContainerId(), types.ContainerStartOptions{})
	return containerId, err
}

func (s *Service) ContainerCreate(ctx context.Context, config *v1.ContainerConfig) (*v1.Container, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{Image: config.GetImage()}, nil, nil, "")
	return &v1.Container{Id: resp.ID}, err
}
func (s *Service) ContainerRemove(ctx context.Context, containerID *v1.ContainerId) (*v1.ContainerId, error) {
	err := cli.ContainerRemove(ctx, containerID.GetContainerId(), types.ContainerRemoveOptions{})
	if err != nil {
		log.Fatalf("Container Remove Error: %v", err)
	}
	return containerID, err
}

func (s *Service) ContainerStatStream(containerId *v1.ContainerId, stream v1.ContainerService_ContainerStatStreamServer) error {
	var err error
	for {
		response, err := cli.ContainerStats(context.Background(), containerId.GetContainerId(), true)
		if err != nil {
			log.Fatalf("Stats Stream Error %s", err)
		}
		stat := p2.Metric{}
		json.NewDecoder(response.Body).Decode(&stat)

		stream.Send(&v1.ContainerStat{Name: stat.Name})
	}
	return err
}

func (s *Service) CreateProject(ctx context.Context, project *v1.Project) (*v1.ProjectInfo, error) {
	image := project.GetImage()
	max := int(project.GetMax())
	min := int(project.GetMin())

	tmp := p2.Project{Max: max, Min: min, Image: image, Name: project.GetName()}

	var containers []string
	for i := 0; i < min; i++ {

		resp, err := cli.ContainerCreate(context.Background(), &container.Config{Image: image}, nil, nil, fmt.Sprintf("%s%s", project.GetName(), strconv.Itoa(i+1)))

		if err != nil {
			log.Printf("[CREATE_PROJECT] Creating Container error: %v", err)

		}
		err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
		if err != nil {
			log.Printf("Contaner Not starting, %s", err)
		}
		containers = append(containers, string(resp.ID))
	}

	tmp.Containers = containers
	go metric.ProjectMetric(tmp)
	projects = append(projects, tmp)
	return &v1.ProjectInfo{ContainerId: containers}, nil

}

func GetProjects() *[]p2.Project {
	return &projects
}

//func SetProjects()  Maybe later need it.
