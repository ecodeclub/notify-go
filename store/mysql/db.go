package mysql

import (
	"log/slog"
	"xorm.io/xorm"
)

type DBConfig struct {
	DriverName string `toml:"driver_name"`
	Dsn        string `toml:"dsn"`
}

func NewEngine(cfg DBConfig) *xorm.Engine {
	e, err := xorm.NewEngine(cfg.DriverName, cfg.Dsn)
	if err != nil {
		slog.Error("[db]创建db失败")
	}
	return e
}
