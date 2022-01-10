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

package bfe_register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ProductPoolsAPI struct {
	cfg BfeRegisterConfig
}

type ProductPool struct {
	Name      string
	Instances []Instance
}

type Instance struct {
	Hostname string
	Ip       string
	Weigth   float64
	Ports    map[string]uint64
	Tags     map[string]string
}

type ReturnData struct {
	ErrNum int
	ErrMsg string
}

func (productPoolsAPI *ProductPoolsAPI) registerPools(product_name string, productPool *ProductPool) {
	body, err := productPoolsAPI.UpdatePools(product_name, productPool)
	if err != nil {
		return
	}
	var data ReturnData
	json.Unmarshal(body, &data)
	if data.ErrNum != 200 {
		productPoolsAPI.InstancePools(product_name, productPool)
	}
}

func (productPoolsAPI *ProductPoolsAPI) InstancePools(product_name string, productPool *ProductPool) ([]byte, error) {
	return productPoolsAPI.Send("POST", "http://"+productPoolsAPI.cfg.APIService+"/open-api/v1/products/"+product_name+"/instance-pools", productPool)
}

func (productPoolsAPI *ProductPoolsAPI) UpdatePools(product_name string, productPool *ProductPool) ([]byte, error) {
	return productPoolsAPI.Send("PATCH", "http://"+productPoolsAPI.cfg.APIService+"/open-api/v1/products/"+product_name+"/instance-pools/"+productPool.Name, productPool)
}

func (productPoolsAPI *ProductPoolsAPI) One(product_name string, productPool *ProductPool) {
	productPoolsAPI.Send("GET", "http://"+productPoolsAPI.cfg.APIService+"/open-api/v1/products/"+product_name+"/instance-pools/"+productPool.Name, nil)

}

func (productPoolsAPI *ProductPoolsAPI) Send(method string, url string, data interface{}) ([]byte, error) {
	client := &http.Client{}
	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(bytesData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", productPoolsAPI.cfg.Authorization)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	str := string(body[:])
	fmt.Print(str)
	return body, err
}
