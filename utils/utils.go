package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/WianVos/xld"
	jww "github.com/spf13/jwalterweatherman"
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

//ReadFromFile reads the contents of a file to a byte array
func ReadFromFile(f string) []byte {
	contents, err := ioutil.ReadFile(f)
	HandleErr(err)
	return contents
}

//ReadFromFileToString reads contents from file and returns a string
func ReadFromFileToString(f string) string {
	return string(ReadFromFile(f))
}

//GetClient returns a xld.Client object
func GetClient() *xld.Client {
	//get the much needed config for the xlr client
	config := GetConfig()

	// instantiate the xlr client
	client := xld.NewClient(config)
	if client.VerifyConnection() == false {
		// !HandleErr(errors.New("unable to connect to XL-Deploy"))
		jww.FATAL.Fatalf("unable to connect XL-Deploy at: %s \n", client.Config.Host)
	}
	jww.FEEDBACK.Printf("connected to XL-Deploy at: %s, as %s \n", client.Config.Host, client.Config.User)
	return client

}

//HandleErr handles an error by panicing like a little bitch
func HandleErr(err error) {
	if err != nil {
		fmt.Printf("Goxld: fatal error : %s \n", err)
		os.Exit(2)
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
