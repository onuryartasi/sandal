package v1

import (
	"os"
	"github.com/docker/docker/client"
	"log"
	"github.com/docker/docker/api/types"
	"context"

	"github.com/onuryartasi/scaler/pkg/api/v1"
	"github.com/golang/protobuf/ptypes/empty"
	"time"
	"github.com/docker/docker/api/types/container"

	"io/ioutil"
	"encoding/json"
)

var (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

type Service struct{}
type Error struct {

}
var cli *client.Client
func init() {
	var err error
	if len(os.Getenv("DOCKER_API_VERSION")) < 1{
		os.Setenv("DOCKER_API_VERSION", "1.37")
	}

	cli,err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf(ErrorColor,"Error: %s",err)
	}
}

var timeout time.Duration = 10*time.Second
func (s *Service) ContainerList(ctx context.Context,empty *empty.Empty) (*v1.Containers, error){
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatalf(ErrorColor,"Error: %s",err)
	}

	rContainers := []*v1.Container{}
	for _,container := range containers{
		rContainers = append(rContainers,&v1.Container{Names:container.Names,Id:container.ID,Image:container.Image})
	}
	return &v1.Containers{Container:rContainers},nil
}


func (s *Service) ContainerStop(ctx context.Context,containerId *v1.ContainerId) (*v1.ContainerId,error){


	err := cli.ContainerStop(ctx,containerId.GetContainerId(),&timeout)
	return containerId,err
}

func (s *Service) ContainerStart(ctx context.Context,containerId *v1.ContainerId)(*v1.ContainerId,error){

	err := cli.ContainerStart(ctx,containerId.GetContainerId(),types.ContainerStartOptions{})
	return containerId,err
}

func (s *Service) ContainerCreate(ctx context.Context,config *v1.ContainerConfig)(*v1.Container,error){
	resp,err := cli.ContainerCreate(ctx,&container.Config{Image:config.GetImage()},nil,nil,"")
	return &v1.Container{Id:resp.ID},err
}

func (s *Service) ContainerStatStream(containerId *v1.ContainerId,stream v1.ContainerService_ContainerStatStreamServer) error{
	var err error
	for {
		stats, err := cli.ContainerStats(context.Background(), containerId.GetContainerId(), false)
		if err != nil {
			log.Fatalf("Stats Stream Error %s", err)
		}
		containerStats, _ := ioutil.ReadAll(stats.Body)
		stat := Stats{}
		json.Unmarshal(containerStats, &stat)
		stream.Send(&v1.ContainerStat{Name: stat.Name})
	}
	return err
}