// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

package main

import (

	//external libraries

	"fmt"

	"github.com/WianVos/goxld/cmdconfig"
	"github.com/WianVos/goxld/cmdrepository"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config environmental variables.

//Host holds the address of the xldeploy host to connect to
var Host string

//Context holds the context root of the xld server
var Context string

//User to use when authenticating to the xld server
var User string

//Password to use when authentication to the xld server
var Password string

//Port the xld server is listening on
var Port string

//Scheme identifeis the http scheme to use in communication with the xld server
var Scheme string

var goxld = &cobra.Command{
	Use:   "goxld",
	Short: "goxld provides a command line interface to work with XL-Deploy",
}

func init() {
	goxld.AddCommand(cmdconfig.GetCommands())
	goxld.AddCommand(cmdrepository.GetCommands())
	goxld.PersistentFlags().StringVarP(&Host, "host", "x", "blah", "XL-Deploy hostname")
	goxld.PersistentFlags().StringVarP(&Context, "context", "c", "/deployit", "XL-Deploy context")
	goxld.PersistentFlags().StringVarP(&User, "user", "u", "", "XL-Deploy username")
	goxld.PersistentFlags().StringVarP(&Password, "password", "p", "", "XL-Deploy password")
	goxld.PersistentFlags().StringVarP(&Port, "port", "P", "80", "portnumber to reach XL-Deploymk on")
	goxld.PersistentFlags().StringVarP(&Scheme, "scheme", "s", "http", "http scheme to user")
	viper.BindPFlag("port", goxld.PersistentFlags().Lookup("port"))
	viper.BindPFlag("host", goxld.PersistentFlags().Lookup("host"))
	viper.BindPFlag("context", goxld.PersistentFlags().Lookup("context"))
	viper.BindPFlag("user", goxld.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", goxld.PersistentFlags().Lookup("password"))
	viper.BindPFlag("scheme", goxld.PersistentFlags().Lookup("scheme"))

}
func main() {
	// initialze config
	initializeConfig()

	goxld.Execute()
}

//initialize the viper config
func initializeConfig() {
	// get input from config files

	// configfile name is goxld
	viper.SetConfigName("goxld")

	// add the filepaths that will be used
	viper.AddConfigPath("/etc/goxld/")
	viper.AddConfigPath("$HOME/.goxld")
	viper.AddConfigPath(".")

	// Handle errors reading the config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}
