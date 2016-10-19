package cmdrepository

import (
	"fmt"
	"os"

	"github.com/WianVos/goxld/utils"
	"github.com/WianVos/xld"
	"github.com/spf13/cobra"
)

var glistLong = `List the ci's in the xld repository
Example:
  repo list /Environments
`
var flagLongList bool

func addList() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list ci's",
		Long:  glistLong,
		Run:   runList,
	}

	//add local long listing flag to the Command
	//cmd.Flags(dsfa).BoolVarP(&flagJSON, "json", "j", false, "display in json format")
	cmd.Flags().StringVarP(&flagOutFile, "outfile", "o", "", "File to use for output")
	cmd.Flags().BoolVarP(&flagLongList, "long", "l", false, "get everything")

	relCmd.AddCommand(cmd)

}

func runList(cmd *cobra.Command, args []string) {

	client := utils.GetClient()

	cis, err := client.Repository.ListCis(args[0])

	if err != nil {
		panic(fmt.Errorf("Fatal error repo list: %s \n", err))
	}

	if flagLongList == true {
		runListLong(cis)
		os.Exit(0)
	}

	output := utils.RenderJSON(cis)
	if flagOutFile != "" {
		utils.WriteToFile(output, flagOutFile)
		return
	}
	fmt.Println(output)
}

func runListLong(cis xld.CiList) {
	client := utils.GetClient()
	var ciCollection xld.Cis

	for _, c := range cis {

		ci, err := client.Repository.GetCi(c.ID)

		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		ciCollection = append(ciCollection, ci)
	}

	output := utils.RenderJSON(ciCollection)

	if flagOutFile != "" {
		utils.WriteToFile(output, flagOutFile)
		return
	}

	fmt.Println(output)
}
