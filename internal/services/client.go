package services

import (
	"fmt"
	"net"
	"time"

	"github.com/xiaotiancaipro/nextunnel-client/internal/configs"
	"github.com/xiaotiancaipro/nextunnel-client/internal/utils"
	"go.uber.org/zap"
)

type Client struct {
	Config  *configs.Client
	Proxies []configs.Proxy
	Logger  *zap.Logger
	Conn    net.Conn
}

func (c *Client) Login() error {
	payload := utils.LoginMsg{
		Id:    c.Config.Id,
		Token: c.Config.Token,
	}
	if err := utils.WriteMsg(c.Conn, utils.MsgLogin, payload); err != nil {
		c.Logger.Error(fmt.Sprintf("failed to write login msg: %v", err))
		return fmt.Errorf("failed to send LoginMsg")
	}
	return nil
}

func (c *Client) LoginResponse() (*string, error) {

	_ = c.Conn.SetDeadline(time.Now().Add(10 * time.Second))
	msgType, payload, err := utils.ReadMsg(c.Conn)
	_ = c.Conn.SetDeadline(time.Time{})
	if err != nil {
		c.Logger.Error(fmt.Sprintf("failed to read login msg: %v", err))
		return nil, fmt.Errorf("failed to read LoginResp")
	}
	if msgType != utils.MsgLoginResp {
		c.Logger.Error(fmt.Sprintf("invalid login msg type: %v", msgType))
		return nil, fmt.Errorf("expected LoginResp")
	}

	var loginResp utils.LoginRespMsg
	if err := utils.Decode(payload, &loginResp); err != nil {
		c.Logger.Error(fmt.Sprintf("failed to decode LoginResp: %v", err))
		return nil, fmt.Errorf("failed to parse LoginResp")
	}
	if loginResp.Error != "" {
		c.Logger.Error(fmt.Sprintf("login response error: %v", loginResp.Error))
		return nil, fmt.Errorf("login error")
	}

	return &loginResp.RunID, nil

}

func (c *Client) ProxiesApply(conn net.Conn) error {

	proxies := make([]utils.ProxiesApplyMsgItem, 0, len(c.Proxies))
	for _, proxy := range c.Proxies {
		proxies = append(proxies, utils.ProxiesApplyMsgItem{
			Name:       proxy.Name,
			Type:       proxy.Type,
			RemotePort: proxy.RemotePort,
		})
	}

	payload := utils.ProxiesApplyMsg{
		Proxies: proxies,
	}
	if err := utils.WriteMsg(conn, utils.MsgProxiesApply, payload); err != nil {
		c.Logger.Error(fmt.Sprintf("failed to write proxies msg: %v", err))
		return fmt.Errorf("failed to send ApplyConfigMsg")
	}

	return nil

}

func (c *Client) ProxiesApplyResponse(conn net.Conn) error {

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	msgType, payload, err := utils.ReadMsg(conn)
	_ = conn.SetDeadline(time.Time{})
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Failed to read proxies msg: %v", err))
		return fmt.Errorf("failed to read ApplyConfigResp")
	}
	if msgType != utils.MsgProxiesApplyResp {
		c.Logger.Error(fmt.Sprintf("Invalid proxies msg type: %v", msgType))
		return fmt.Errorf("expected ApplyConfigResp")
	}

	var resp utils.ProxiesApplyRespMsg
	if err := utils.Decode(payload, &resp); err != nil {
		c.Logger.Error(fmt.Sprintf("Failed to decode ApplyConfigResp: %v", err))
		return fmt.Errorf("failed to parse ApplyConfigResp")
	}
	if resp.Error != "" {
		c.Logger.Error(fmt.Sprintf("Failed to parse ApplyConfigResp: %v", resp.Error))
		return fmt.Errorf("apply config rejected by server")
	}

	for _, proxy := range c.Proxies {
		c.Logger.Info(fmt.Sprintf("Proxy applied successfully: name=%s", proxy.Name))
	}
	return nil

}
