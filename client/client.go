package client

import (
	"fmt"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
)

type Client struct {
	logger     logger.Logger
	config     Config
	auth       *domain.Auth
	updateAuth func() error
}

func NewClient(
	logger logger.Logger,
	config Config,
	auth *domain.Auth,
	updateAuth func() error,
) Client {
	output := Client{
		logger:     logger,
		config:     config,
		auth:       auth,
		updateAuth: updateAuth,
	}
	return output
}

func (client Client) getDomain() string {
	return fmt.Sprintf("https://%v.janeapp.com", client.auth.Domain)
}
