package v1

import (
	"context"
	"github.com/bamboo-firewall/be/cmd/server/pkg/storage"
	"github.com/bamboo-firewall/be/v1/route"

	"github.com/bamboo-firewall/be/cmd/server/pkg/http"
)

type App interface {
	Start()
	Stop(ctx context.Context) error
}

type app struct {
	httpServer  *http.Server
	policyMongo *storage.PolicyMongo
}

func NewApp() (App, error) {
	// get from env or argument
	configMongo := storage.ConfigMongo{}
	policyMongo, err := storage.NewPolicyMongo(configMongo)
	if err != nil {
		return nil, err
	}

	router := route.RegisterHandler()
	return &app{
		httpServer:  http.NewServer(router),
		policyMongo: policyMongo,
	}, nil
}

func (a *app) Start() {
	// Handle ErrServerClosed
	if err := a.httpServer.Start(); err != nil {
		panic(err)
	}
}

func (a *app) Stop(ctx context.Context) error {
	if err := a.httpServer.Stop(ctx); err != nil {
		return err
	}
	if err := a.policyMongo.Stop(ctx); err != nil {
		return err
	}
	return nil
}
