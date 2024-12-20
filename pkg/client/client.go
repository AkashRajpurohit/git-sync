package client

import (
	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/token"
)

type Client interface {
	Sync(config config.Config) error
	GetTokenManager() *token.Manager
}
