package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase"
	"github.com/bamboo-firewall/be/cmd/server/pkg/repository"
	"github.com/bamboo-firewall/be/cmd/server/pkg/storage"
	"github.com/bamboo-firewall/be/cmd/server/route"
	"github.com/bamboo-firewall/be/config"
)

type App interface {
	Start() error
	Stop(ctx context.Context) error
}

type app struct {
	httpServer  *httpbase.Server
	policyMongo *storage.PolicyMongo
}

func NewApp(cfg config.Config) (App, error) {
	policyMongo, err := storage.NewPolicyMongo(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	router := route.RegisterHandler(repository.NewPolicyMongo(policyMongo))
	return &app{
		httpServer: httpbase.NewServer(router, httpbase.ConfigTimeout{
			ReadTimeout:       cfg.HTTPServerReadTimeout,
			ReadHeaderTimeout: cfg.HTTPServerReadHeaderTimeout,
			WriteTimeout:      cfg.HTTPServerWriteTimeout,
			IdleTimeout:       cfg.HTTPServerIdleTimeout,
		}),
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
