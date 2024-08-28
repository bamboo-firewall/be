package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase"
	"github.com/bamboo-firewall/be/cmd/server/pkg/repository"
	"github.com/bamboo-firewall/be/cmd/server/pkg/storage"
	"github.com/bamboo-firewall/be/cmd/server/route"
)

type App interface {
	Start() error
	Stop(ctx context.Context) error
}

type app struct {
	httpServer  *httpbase.Server
	policyMongo *storage.PolicyMongo
}

func NewApp() (App, error) {
	// get from env or argument
	policyMongo, err := storage.NewPolicyMongo("mongodb://admin:password@localhost:27017/?w=majority&socketTimeoutMS=3000")
	if err != nil {
		return nil, err
	}

	router := route.RegisterHandler(repository.NewPolicyMongo(policyMongo))
	return &app{
		httpServer:  httpbase.NewServer(router),
		policyMongo: policyMongo,
	}, nil
}

func (a *app) Start() error {
	if err := a.httpServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
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
