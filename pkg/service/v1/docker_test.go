package v1

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	v1 "github.com/onuryartasi/scaler/pkg/api/v1"

	"testing"
)

func TestService_ContainerList(t *testing.T) {
	s := Service{}

	req := &empty.Empty{}
	_, err := s.ContainerList(context.Background(), req)
	if err != nil {
		t.Errorf("Error ContainerList")
	}

}

func TestService_ContainerCreate(t *testing.T) {
	s := Service{}
	req := &v1.ContainerConfig{Image: "hello-world"}
	_, err := s.ContainerCreate(context.Background(), req)
	if err != nil {
		t.Errorf("Container Create Error : %v", err)
	}

}

func TestService_ContainerStart(t *testing.T) {
	s := Service{}
	req := &v1.ContainerConfig{Image: "hello-world"}
	resp, err := s.ContainerCreate(context.Background(), req)
	if err != nil {
		t.Errorf("Container Create Error : %v", err)
	}

	startReq := &v1.ContainerId{ContainerId: resp.GetId()}
	_, err = s.ContainerStart(context.Background(), startReq)
	if err != nil {
		t.Errorf("Container Start Error : %v", err)
	}
}

func TestService_ContainerStop(t *testing.T) {
	s := Service{}
	req := &v1.ContainerConfig{Image: "hello-world"}
	resp, err := s.ContainerCreate(context.Background(), req)
	if err != nil {
		t.Errorf("Container Create Error : %v", err)
	}

	stopReq := &v1.ContainerId{ContainerId: resp.GetId()}
	_, err = s.ContainerStop(context.Background(), stopReq)
	if err != nil {
		t.Errorf("Container Stop Error : %v", err)
	}
}

func TestService_CreateProejct(t *testing.T) {
	s := Service{}
	req := &v1.Project{Image: "hello-world", Max: 1, Min: 1, Name: "test"}
	resp, err := s.CreateProject(context.Background(), req)
	if err != nil {
		t.Errorf("Create Project Error : %v", err)
	}
	if len(resp.GetContainerId()) != 1 {
		t.Error("Container count dont match expected")
	}

}
