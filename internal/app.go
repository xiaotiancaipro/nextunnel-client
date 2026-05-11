package internal

import (
	"fmt"
	"net"

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

	conn, err := a.serverConn()
	if err != nil {
		return err
	}
	a.logger.Info("Successfully connected to the server")

	runIdP, err := a.clientLogin()
	if err != nil {
		_ = conn.Close()
		return err
	}
	a.logger.Info(fmt.Sprintf("Running with id: %s", *runIdP))

	if err = a.clientProxiesApply(); err != nil {
		_ = conn.Close()
		return err
	}
	a.logger.Info("Client proxies configuration application successful")

	// TODO

	return nil

}

func (a *App) serverConn() (net.Conn, error) {

	c, err := a.tlsService.Init()
	if err != nil {
		return nil, err
	}

	conn, err := a.serverService.DialServer(c)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to connect to server: %s", err))
		return nil, fmt.Errorf("failed to connect to server")
	}

	a.clientService.Conn = conn
	return conn, nil

}

func (a *App) clientLogin() (*string, error) {

	if err := a.clientService.Login(); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to login: %s", err))
		return nil, fmt.Errorf("failed to login")
	}

	runIdP, err := a.clientService.LoginResponse()
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to login: %s", err))
		return nil, fmt.Errorf("failed to login")
	}

	return runIdP, nil

}

func (a *App) clientProxiesApply() error {
	if err := a.clientService.ProxiesApply(); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to apply proxies: %s", err))
		return fmt.Errorf("failed to apply proxies")
	}
	if err := a.clientService.ProxiesApplyResponse(); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to apply proxies: %s", err))
		return fmt.Errorf("failed to apply proxies")
	}
	return nil
}
