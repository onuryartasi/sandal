package v1

import (
	"os"
	"github.com/docker/docker/client"
	"log"
	"github.com/docker/docker/api/types"
	"context"

	"github.com/onuryartasi/scaler/pkg/api/v1"
	"github.com/golang/protobuf/ptypes/empty"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

type Service struct{}
var cli *client.Client
func init() {
	var err error
	os.Setenv("DOCKER_API_VERSION", "1.37")
	cli,err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf(ErrorColor,"Error: %s",err)
	}
}

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