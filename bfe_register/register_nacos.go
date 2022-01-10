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
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type RegisterNacos struct {
	registerInfo    RegisterInfo
	productPoolsAPI *ProductPoolsAPI
}

func (register *RegisterNacos) SetRegisterInfo(registerInfo RegisterInfo) {
	register.registerInfo = registerInfo
}
func (register *RegisterNacos) SetProductPoolsAPI(productPoolsAPI *ProductPoolsAPI) {
	register.productPoolsAPI = productPoolsAPI
}
func (register *RegisterNacos) Init() error {
	registerInfo := register.registerInfo
	sc := make([]constant.ServerConfig, len(registerInfo.Address))
	for addressIndex, address := range registerInfo.Address {
		sc[addressIndex] = constant.ServerConfig{
			IpAddr: address.IpAddr,
			Port:   address.Port,
		}
	}

	for _, spaceInfo := range registerInfo.SpaceInfos {
		cc := constant.ClientConfig{
			NamespaceId:         spaceInfo.SpaceName, //namespace id
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogDir:              "./log",
			CacheDir:            "./cache",
			RotateTime:          "1h",
			MaxAge:              3,
			LogLevel:            "debug",
		}
		client, err := clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig:  &cc,
				ServerConfigs: sc,
			},
		)
		if err != nil {
			return err
		}
		if spaceInfo.ServierInfos == nil {
			var num uint32 = 0
			for {
				param := vo.GetAllServiceInfoParam{
					NameSpace: spaceInfo.SpaceName,
					PageNo:    num,
					PageSize:  50,
				}
				services, _ := client.GetAllServicesInfo(param)
				for _, service := range services.Doms {
					param := &vo.SubscribeParam{
						ServiceName: service,
						SubscribeCallback: func(services []model.SubscribeService, err error) {
							productPool := CreateProductPool(spaceInfo.PoolName, service, services)
							//spaceInfo.PoolName
							register.productPoolsAPI.registerPools(service, &productPool)
						},
					}
					client.Subscribe(param)
				}
				if services.Count < 50 {
					break
				}
				num += num
			}
		} else {
			for _, serviceInfo := range spaceInfo.ServierInfos {
				param := &vo.SubscribeParam{
					ServiceName: serviceInfo.ServiceName,
					SubscribeCallback: func(services []model.SubscribeService, err error) {
						//spaceInfo.PoolName
						var poolName string
						if serviceInfo.PoolName == "" {
							poolName = spaceInfo.PoolName
						} else {
							poolName = serviceInfo.PoolName
						}
						productPool := CreateProductPool(poolName, serviceInfo.ServiceName, services)
						register.productPoolsAPI.registerPools(serviceInfo.ServiceName, &productPool)
					},
				}
				client.Subscribe(param)
			}
		}
	}
	return nil
}

func CreateProductPool(poolName string, servicesName string, services []model.SubscribeService) ProductPool {
	productPool := ProductPool{}
	productPool.Name = servicesName + "." + poolName
	instances := make([]Instance, len(services))
	for index, service := range services {
		instances[index] = CreateInstance(service)
	}

	productPool.Instances = instances
	return productPool
}

func CreateInstance(service model.SubscribeService) Instance {
	instance := Instance{
		Ip:       service.Ip,
		Ports:    map[string]uint64{"Default": service.Port},
		Weigth:   service.Weight,
		Hostname: service.ServiceName,
		Tags:     service.Metadata,
	}

	return instance
}
