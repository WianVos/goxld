package cmdsecurity

import (
	"github.com/WianVos/goxld/utils"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var passwordLong = `(re-)set a password for a user in XL-Deploy
example: user password --password changeme user
`

func addPassword() {
	cmd := &cobra.Command{
		Use:   "password",
		Short: "(re-)set a password for a user in XL-Deploy",
		Long:  passwordLong,
		Run:   runPassword,
	}

	cmd.Flags().StringVarP(&flagPassWord, "password", "", "", "password to use for the new user")

	//add local long listing flag to the Command
	//cmd.Flags(dsfa).BoolVarP(&flagJSON, "json", "j", false, "display in json format")
	secCmd.AddCommand(cmd)

}

//runWrite executes the action define by the config write command to goxld
func runPassword(cmd *cobra.Command, args []string) {
	c := utils.GetClient()
	var username string

	if len(args) == 1 {
		username = args[0]
	}

	if flagPassWord != "" {
		err := c.Security.SetPasswordForUser(username, flagPassWord)
		if err != nil {
			jww.DEBUG.Println(err)
			jww.ERROR.Panicln("user add: unable to set password for user")
		}
	}
	jww.FEEDBACK.Printf("password set for user %s \n", username)

}
