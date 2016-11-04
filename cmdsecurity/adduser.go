package cmdsecurity

import (
	"errors"
	"fmt"

	"github.com/WianVos/goxld/utils"
	"github.com/spf13/cobra"
)

var writeLong = `add a user in the XL-Deploy repository
example: user add --password changeme -a newuser
`

var flagPassWord string
var flagAdmin bool

func addUser() {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add or manipulate a user in XL-Deploy",
		Long:  writeLong,
		Run:   runUser,
	}

	cmd.Flags().StringVarP(&flagPassWord, "password", "", "", "password to use for the new user")
	cmd.Flags().BoolVarP(&flagAdmin, "admin", "a", false, "user is admin")

	//add local long listing flag to the Command
	//cmd.Flags(dsfa).BoolVarP(&flagJSON, "json", "j", false, "display in json format")
	secCmd.AddCommand(cmd)

}

//runWrite executes the action define by the config write command to goxld
func runUser(cmd *cobra.Command, args []string) {
	c := utils.GetClient()
	var username string

	if len(args) == 1 {
		username = args[0]
	}

	_, err := c.Security.CreateUser(username, flagAdmin)
	if err != nil {
		utils.HandleErr(errors.New("user add: unable to create user"))
	}

	if flagPassWord != "" {
		err := c.Security.SetPasswordForUser(username, flagPassWord)
		if err != nil {
			fmt.Println(err)
			utils.HandleErr(errors.New("user add: unable to set password for user"))
		}
	}
	fmt.Println("Goxld: user add: user created")
}
