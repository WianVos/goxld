package cmdrepository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/WianVos/goxld/utils"
	"github.com/WianVos/xld"
	"github.com/spf13/cobra"
)

var gCreateLong = `Return a ci
Example:
  repo create Environment/Dictionary1
`

var flagInFile string

func addCreate() {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a ci from filebased input",
		Long:  gCreateLong,
		Run:   runCreate,
	}

	//add local long listing flag to the Command
	//cmd.Flags(dsfa).BoolVarP(&flagJSON, "json", "j", false, "display in json format")
	//cmd.Flags().StringVarP(&flagInFile, "infile", "i", "", "File to use for input")
	relCmd.AddCommand(cmd)

}

func runCreate(cmd *cobra.Command, args []string) {

	// declaring variables
	var ci xld.Ci

	if len(args) < 1 {
		utils.HandleErr(errors.New("insufficient number of arguments given"))
	}

	in := utils.ReadFromFile(args[0])

	if err := json.NewDecoder(bytes.NewReader(in)).Decode(&ci); err != nil {
		utils.HandleErr(err)
	}
	// create the ci from the in string

	client := utils.GetClient()

	_, err := client.Repository.SaveCi(ci)
	utils.HandleErr(err)

	output := utils.RenderJSON(ci)

	fmt.Println(output)
}
