package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/WianVos/xld"
	"github.com/spf13/viper"
)

//GetConfig prepares a xld.Config object
func GetConfig() *xld.Config {

	return &xld.Config{
		User:     viper.GetString("user"),
		Password: viper.GetString("password"),
		Host:     viper.GetString("host"),
		Port:     viper.GetString("port"),
		Context:  viper.GetString("context"),
		Scheme:   viper.GetString("scheme"),
	}

}

//RenderJSON function to render output as json
// returns a string object with json formated output
func RenderJSON(l interface{}) string {

	b, err := json.MarshalIndent(l, "", " ")
	if err != nil {
		panic(err)
	}
	s := string(b)

	return s
}

//WriteToFile writes any string output to file
func WriteToFile(s string, f string) {
	d1 := []byte(s + "\n")
	err := ioutil.WriteFile(f, d1, 0644)
	if err != nil {
		panic(err)
	}
}

//GetClient returns a xld.Client object
func GetClient() *xld.Client {
	//get the much needed config for the xlr client
	config := GetConfig()

	// instantiate the xlr client
	client := xld.NewClient(config)

	return client

}
