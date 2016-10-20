package cmdconfig

import "github.com/spf13/cobra"

// flag variables

var flagOutFile string

//
var relCmd = &cobra.Command{
	Use:   "config",
	Short: "config allows for interaction with the cli's config files",
}

//GetCommands grab and return commands in this package
func GetCommands() *cobra.Command {

	//collect the commands in the package
	addWrite()
	return relCmd
}
