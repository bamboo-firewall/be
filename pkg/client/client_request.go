package client

import "github.com/bamboo-firewall/be/pkg/http"

type apiServer struct {
	client *http.Client
}

func NewAPIServer(address string) *apiServer {
	return &apiServer{client: http.NewClient(address)}
}
