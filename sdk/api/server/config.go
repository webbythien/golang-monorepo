package server

import (
	"fmt"
	"net"
)

// Listen represents a network end point address.
type Listen struct {
	Host string `json:"host" mapstructure:"host" yaml:"host"`
	Port int    `json:"port" mapstructure:"port" yaml:"port"`
}

func (l *Listen) CreateListener() (net.Listener, error) {
	lis, err := net.Listen("tcp", l.String())
	if err != nil {
		return nil, fmt.Errorf("failed to listen %s: %w", l.String(), err)
	}
	return lis, nil
}

func (l *Listen) String() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

// Config hold http/grpc server config
type Config struct {
	HTTP Listen `json:"http" mapstructure:"http" yaml:"http"`
}
