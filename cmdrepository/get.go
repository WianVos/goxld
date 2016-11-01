package cmdrepository

import (
	"fmt"

	"github.com/WianVos/goxld/utils"
	"github.com/spf13/cobra"
)

var gGetLong = `Return a ci
Example:
  repo get /Environment/Dictionary1
`

func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a certain ci",
		Long:  gGetLong,
		Run:   runGet,
	}

	//add local long listing flag to the Command
	//cmd.Flags(dsfa).BoolVarP(&flagJSON, "json", "j", false, "display in json format")
	cmd.Flags().StringVarP(&flagOutFile, "outfile", "o", "", "File to use for output")
	relCmd.AddCommand(cmd)

}

func runGet(cmd *cobra.Command, args []string) {

	client := utils.GetClient()

	ci, err := client.Repository.GetCi(args[0])
	utils.HandleErr(err)

	output := utils.RenderJSON(ci)

	if flagOutFile != "" {
		utils.WriteToFile(output, flagOutFile)
		return
	}

	fmt.Println(output)
}
