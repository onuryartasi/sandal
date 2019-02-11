package v1

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

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
