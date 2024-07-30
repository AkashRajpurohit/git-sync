package client

import "github.com/AkashRajpurohit/git-sync/pkg/config"

type Client interface {
	Sync(cfg config.Config) error
}
