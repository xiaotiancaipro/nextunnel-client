package internal

import (
	"fmt"

	"github.com/xiaotiancaipro/nextunnel-client/internal/configs"
	"github.com/xiaotiancaipro/nextunnel-client/internal/services"
	"go.uber.org/zap"
)

type App struct {
	logger        *zap.Logger
	tlsService    *services.Tls
	serverService *services.Server
	clientService *services.Client
}

func NewApp(config *configs.Configs, logger *zap.Logger) *App {
	tls := services.Tls{
		Config: config.Tls,
		Logger: logger,
	}
	server := services.Server{
		Config: config.Server,
		Logger: logger,
	}
	client := services.Client{
		Config:  config.Client,
		Proxies: config.Proxies,
		Logger:  logger,
	}
	return &App{
		logger:        logger,
		tlsService:    &tls,
		serverService: &server,
		clientService: &client,
	}
}

func (a *App) Start() error {

	c, err := a.tlsService.Init()
	if err != nil {
		return err
	}

	conn, err := a.serverService.DialServer(c)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to connect to server: %s", err))
		return fmt.Errorf("failed to connect to server")
	}
	a.clientService.Conn = conn

	if err = a.clientService.Login(); err != nil {
		_ = conn.Close()
		a.logger.Error(fmt.Sprintf("Failed to login: %s", err))
		return fmt.Errorf("failed to login")
	}

	runIdP, err := a.clientService.LoginResponse()
	if err != nil {
		_ = conn.Close()
		a.logger.Error(fmt.Sprintf("Failed to login: %s", err))
		return fmt.Errorf("failed to login")
	}
	runId := *runIdP
	a.logger.Info(fmt.Sprintf("Running with id: %s", runId))

	// TODO

	return nil

}
