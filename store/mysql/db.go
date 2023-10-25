// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
