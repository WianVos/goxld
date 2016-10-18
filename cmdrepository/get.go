package cmdrepository

import (
	"fmt"
	"os"

	"github.com/WianVos/goxld/utils"
	"github.com/spf13/cobra"
)

var glistLong = `Return a list of templates in the system
Example:
  repository get /Environment/Dictionary1
`

func addGet() {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a certain ci",
		Long:  glistLong,
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
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	output := utils.RenderJSON(ci)

	if flagOutFile != "" {
		utils.WriteToFile(output, flagOutFile)
		return
	}

	fmt.Println(output)
}
