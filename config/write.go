package cmdconfig

import (
	"fmt"

	"github.com/WianVos/goxld/utils"
	"github.com/spf13/cobra"
)

var writeLong = `write a cli config file template to the filesystem.
This config file does not contain the proper settings yet , but acts as a template
Example:
  goxld config write /tmp/goxld.json
`

func addWrite() {
	cmd := &cobra.Command{
		Use:   "write",
		Short: "write a config file template",
		Long:  writeLong,
		Run:   runWrite,
	}

	//add local long listing flag to the Command
	//cmd.Flags(dsfa).BoolVarP(&flagJSON, "json", "j", false, "display in json format")
	relCmd.AddCommand(cmd)

}

func runWrite(cmd *cobra.Command, args []string) {
	utils.WriteToFile(getFormattedTemplate(), args[0])
}

func getFormattedTemplate() string {
	config := utils.GetConfig()

	return fmt.Sprintf(`{
  "host": %s,
  "port": %s,
  "context": %s,
  "user": %s,
  "password": %s,
  "scheme": %s
  }`, config.Host,
		config.Port,
		config.Context,
		config.User,
		config.Password,
		config.Scheme)
}
