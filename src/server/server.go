package server

import (
	"zerologix-homework/bootstrap"
	errUtil "zerologix-homework/src/pkg/util/error"
	"zerologix-homework/src/server/router"
)

func Run() error {
	cfg, err := bootstrap.Get()
	if err != nil {
		return errUtil.NewError(err)
	}

	serverRouter := router.EntryRouter(cfg)
	serverAddr := cfg.Server.Addr()
	if err := serverRouter.Run(serverAddr); err != nil {
		return err
	}
	return nil
}
