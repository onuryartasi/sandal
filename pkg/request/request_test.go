package request

import (
	"testing"
	"net/http"
)

func TestRequest_DockerSocket(t *testing.T){
	client := NewClient()
	response, err := client.Get("http://v1.37/_ping")
	if err != nil {
		t.Logf("Docker Socket Error %v",err)
	}

	if response.StatusCode != http.StatusOK{
		t.Logf("Docker socket not running, %v",err)
	}
}