package utils

import (
	cfg "github.com/eris-ltd/tendermint/config"
	tmcfg "github.com/eris-ltd/tendermint/config/tendermint"
)

func init() {
	cfg.ApplyConfig(tmcfg.GetConfig(""))
}
