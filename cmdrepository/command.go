package cmdrepository

import "github.com/spf13/cobra"

// flag variables

var flagOutFile string

// maskCmd represents the parent for all mask cli commands.
var relCmd = &cobra.Command{
	Use:   "repo",
	Short: "repo provides and interface to work with the xldeploy repository",
}

//GetCommands grab and return commands in this package
func GetCommands() *cobra.Command {

	//collect the commands in the package
	addGet()
	addCreate()
	addDictMerge()
	addList()
	return relCmd
}
