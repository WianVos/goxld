package utils

import (
	"encoding/json"
	"errors"
	"fmt"
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

//HandleErr handles an error by panicing like a little bitch
func HandleErr(err error) {
	if err != nil {
		panic(fmt.Errorf("Fatal error dict merge: %s \n", err))
	}
}

//GetCi retrieve a CI form xld and check it for type
func GetCi(n, t string) xld.Ci {

	var err error

	client := GetClient()
	c, err := client.Repository.GetCi(n)

	HandleErr(err)

	// check to see if the type is what we expect
	if c.Type != t {
		HandleErr(errors.New("requested ci is not of the correct type"))
	}

	return c
}

//NewCiObject returns a ci object
func NewCiObject(n, t string, p map[string]interface{}) xld.Ci {
	c := GetClient()
	nc, err := c.Repository.CreateCi(n, t, p)
	HandleErr(err)

	return nc
}

//VerifyType checks if a ci is a certain type
func VerifyType(c xld.Ci, t string) bool {
	if c.Type == t {
		return true
	}

	return false
}
