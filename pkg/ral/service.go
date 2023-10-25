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

package ral

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Service struct {
	Resources []Resource `toml:"Resources"`
}

type Resource struct {
	Name         string      `toml:"Name"`
	ConnTimeOut  int         `toml:"ConnTimeOut"`
	WriteTimeOut int         `toml:"WriteTimeOut"`
	ReadTimeOut  int         `toml:"ReadTimeOut"`
	Retry        int         `toml:"Retry"`
	Protocol     string      `toml:"Protocol"`
	Converter    string      `toml:"Converter"`
	Interface    []Interface `toml:"Interface"`
}

type Interface struct {
	Name   string `toml:"Name"`
	URL    string `toml:"Url"`
	Method string `toml:"Method"`
	Host   string `toml:"Host"`
	Port   string `toml:"Port"`
}

func NewService(file string) Service {
	service := new(Service)
	if _, err := toml.DecodeFile(file, service); err != nil {
		log.Printf("[ral] 初始化失败 %v.", err)
	}
	return *service
}
