package modbus

import (
	"fmt"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/goburrow/modbus"
)

type Mode string

const (
	ModeTcp Mode = "tcp"
	ModeRtu Mode = "rtu"
)

type handler interface {
	modbus.ClientHandler
	Connect() error
	Close() error
}

type MbClient struct {
	modbus.Client
	handler
}

func NewClient(cfg SlaveConfig) (*MbClient, error) {
	var cli MbClient
	switch Mode(cfg.Mode) {
	case ModeTcp:
		// Modbus TCP
		h := modbus.NewTCPClientHandler(cfg.Address[6:])
		h.SlaveId = cfg.ID
		h.Timeout = cfg.Timeout
		h.IdleTimeout = cfg.IdleTimeout
		cli.handler = h
	case ModeRtu:
		// Modbus RTU
		h := modbus.NewRTUClientHandler(cfg.Address)
		h.BaudRate = cfg.BaudRate
		h.DataBits = cfg.DataBits
		h.Parity = cfg.Parity
		h.StopBits = cfg.StopBits
		h.SlaveId = cfg.ID
		h.Timeout = cfg.Timeout
		h.IdleTimeout = cfg.IdleTimeout
		cli.handler = h
	default:
		return nil, errors.Errorf("method not supported")
	}
	return &cli, nil
}

func (m *MbClient) Reconnect() error {
	if err := m.Close(); err != nil {
		return err
	}
	return m.Connect()
}

func (m *MbClient) Connect() error {
	err := m.handler.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %s", err.Error())
	}
	m.Client = modbus.NewClient(m.handler)
	return nil
}

func (m *MbClient) Close() error {
	err := m.handler.Close()
	if err != nil {
		return fmt.Errorf("failed to close client: %s", err.Error())
	}
	return nil
}
