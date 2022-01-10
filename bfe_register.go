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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/log/log4go"
	"github.com/bfenetworks/bfe-register/bfe_register"
	"gopkg.in/yaml.v3"
)

var (
	help        = flag.Bool("h", false, "to show help")
	confRoot    = flag.String("c", "./conf", "root path of configuration")
	logPath     = flag.String("l", "./log", "dir path of log")
	stdOut      = flag.Bool("s", false, "to show log in stdout")
	showVersion = flag.Bool("v", false, "to show version of bfe")
	showVerbose = flag.Bool("V", false, "to show verbose information about bfe")
)

var version string
var commit string

func main() {
	var err error
	var config bfe_register.BfeRegisterConfig
	var logSwitch string

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}
	if *showVersion {
		fmt.Printf("bfe version: %s\n", version)
		return
	}
	if *showVerbose {
		fmt.Printf("bfe version: %s\n", version)
		fmt.Printf("go version: %s\n", runtime.Version())
		fmt.Printf("git commit: %s\n", commit)
		return
	}

	// initialize log
	log4go.SetLogBufferLength(10000)
	log4go.SetLogWithBlocking(false)
	log4go.SetLogFormat(log4go.FORMAT_DEFAULT_WITH_PID)
	log4go.SetSrcLineForBinLog(false)

	err = log.Init("bfe_register", logSwitch, *logPath, *stdOut, "midnight", 7)
	if err != nil {
		fmt.Printf("bfe_register: err in log.Init():%s\n", err.Error())
		AbnormalExit()
	}

	log.Logger.Info("bfe_register[version:%s] start", version)

	confPath := path.Join(*confRoot, "bfe_register.yaml")
	buffer, err := ioutil.ReadFile(confPath)
	err = yaml.Unmarshal(buffer, &config)
	//config, err = bfe_register.BfeRegisterConfigLoad(confPath, *confRoot)
	if err != nil {
		log.Logger.Error("main(): in BfeRegisterConfigLoad():%s", err.Error())
		AbnormalExit()
	}
	err = bfe_register.StartUp(config)
	if err != nil {
		log.Logger.Error("main(): in BfeRegisterServer():%s", err.Error())
		AbnormalExit()
	}
	// waiting for logger finish jobs
	time.Sleep(1 * time.Second)
	log.Logger.Close()
}

func AbnormalExit() {
	// waiting for logger finish jobs
	log.Logger.Close()
	// exit
	os.Exit(1)
}
