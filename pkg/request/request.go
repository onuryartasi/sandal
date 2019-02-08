package request

import (
	"net/http"
	"net"
	"context"
)



func NewClient() http.Client{
	client := http.Client{
		Transport:&http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix","/var/run/docker.sock")
			},
		},
	}
	return client
}
