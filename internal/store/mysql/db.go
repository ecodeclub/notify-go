package mysql

import (
	"github.com/ecodeclub/notify-go/internal/pkg/logger"
	"xorm.io/xorm"
)

type DBConfig struct {
	DriverName string `toml:"driver_name"`
	Dsn        string `toml:"dsn"`
}

func NewEngine(cfg DBConfig) *xorm.Engine {
	e, err := xorm.NewEngine(cfg.DriverName, cfg.Dsn)
	if err != nil {
		logger.Fatal("[db]创建db失败")
	}
	return e
}
