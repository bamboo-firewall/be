package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bamboo-firewall/be/cmd/server/route"
	"github.com/bamboo-firewall/be/config"
	"github.com/bamboo-firewall/be/pkg/httpbase"
	"github.com/bamboo-firewall/be/pkg/repository"
	"github.com/bamboo-firewall/be/pkg/storage"
	"github.com/bamboo-firewall/be/pkg/validator"
)

type App interface {
	Start() error
	Stop(ctx context.Context) error
}

type app struct {
	httpServer *httpbase.Server
	policyDB   *storage.PolicyDB
}

func NewApp(cfg config.Config) (App, error) {
	policy, err := storage.NewPolicyDB(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	validator.Init()

	router := route.RegisterHandler(repository.NewPolicy(policy))
	return &app{
		httpServer: httpbase.NewServer(
			fmt.Sprintf("%s:%s", cfg.HTTPServerHost, cfg.HTTPServerPort),
			router,
			httpbase.ConfigTimeout{
				ReadTimeout:       cfg.HTTPServerReadTimeout,
				ReadHeaderTimeout: cfg.HTTPServerReadHeaderTimeout,
				WriteTimeout:      cfg.HTTPServerWriteTimeout,
				IdleTimeout:       cfg.HTTPServerIdleTimeout,
			},
		),
		policyDB: policy,
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
	if err := a.policyDB.Stop(ctx); err != nil {
		return err
	}
	return nil
}
