package test

import "zerologix-homework/bootstrap"

type IDbCfgRepo interface {
	Get(cfg *bootstrap.Config) (repoCfg bootstrap.Db)
	Set(cfg *bootstrap.Config, testName string)
}

type IRedisCfgRepo interface {
	Get(cfg *bootstrap.Config) (repoCfg bootstrap.Db, keyRoot string)
	Set(cfg *bootstrap.Config, testName string)
}
