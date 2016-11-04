package cmdsecurity

import "github.com/spf13/cobra"

// flag variables

var flagOutFile string

//
var secCmd = &cobra.Command{
	Use:   "user",
	Short: "user allows handeling users in xl-deploy",
}

//GetCommands grab and return commands in this package
func GetUserCommands() *cobra.Command {

	//collect the commands in the package
	addUser()
	return secCmd
}
