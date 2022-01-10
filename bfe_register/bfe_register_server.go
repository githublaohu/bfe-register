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
	"fmt"
)

func StartUp(cfg BfeRegisterConfig) error {
	var err error
	cfg.Authorization = "Session " + cfg.Authorization
	productPoolsAPI := ProductPoolsAPI{cfg: cfg}
	for _, register := range cfg.Register {
		var registerServier RegisterNacos
		switch {
		case register.Name == "nacos":
			registerServier = RegisterNacos{}
		default:
			fmt.Printf("差\n")
			continue
		}
		registerServier.SetRegisterInfo(register)
		registerServier.SetProductPoolsAPI(&productPoolsAPI)
		err := registerServier.Init()
		if err != nil {
			return err
		}
	}

	return err
}
