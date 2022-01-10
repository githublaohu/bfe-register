// Copyright (c) 2019 The BFE Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// cluster framework for bfe

package bfe_register

import (
	gcfg "gopkg.in/gcfg.v1"
)

type Address struct {
	IpAddr string `yaml:"ipAddr"`
	Port   uint64 `yaml:"port"`
}

type ServierInfo struct {
	ServiceName string `yaml:"serviceName"`
	PoolName    string `yaml:"poolName"`
}

type SpaceInfo struct {
	SpaceName    string        `yaml:"spaceName"`
	PoolName     string        `yaml:"poolName"`
	ServierInfos []ServierInfo `yaml:"servierInfos"`
}

type RegisterInfo struct {
	Name       string      `yaml:"name"`
	Address    []Address   `yaml:"address"`
	SpaceInfos []SpaceInfo `yaml:"spaceInfos"`
}

type BfeRegisterConfig struct {
	Register      []RegisterInfo `yaml:"register"`
	APIService    string         `yaml:"APIService"`
	Authorization string         `yaml:"authorization"`
}

func SetDefaultConf(conf *BfeRegisterConfig) {

}

// BfeConfigLoad loads config from config file.
// NOTICE: some value will be modified when not set or out of range!!
func BfeRegisterConfigLoad(filePath string, confRoot string) (BfeRegisterConfig, error) {
	var cfg BfeRegisterConfig
	var err error

	SetDefaultConf(&cfg)

	// read config from file
	err = gcfg.ReadFileInto(&cfg, filePath)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
