package cmdrepository

import (
	"fmt"

	"github.com/WianVos/goxld/utils"
	"github.com/WianVos/xld"
	"github.com/spf13/cobra"
)

var dmlistLong = `Compare two dictionaries
Example:
  repository dict_compare /Environments/Dictionary1 /Environments/Dictionary2
`

type conflicts []conflict

type conflict struct {
	Key          string      `json:"key"`
	Original     interface{} `json:"original"`
	ConflictDict string      `json:"conflict dictionary"`
	Conflict     interface{} `json:"conflict"`
}

var flagPersist bool
var flagForgetConflicts bool
var flagDictID string
var flagEnvID string

func addDictMerge() {
	cmd := &cobra.Command{
		Use:   "dict_merge",
		Short: "merges the entries for two dictionaries",
		Long:  dmlistLong,
		Run:   runDictMerge,
	}

	//add local long listing flag to the Command
	cmd.Flags().StringVarP(&flagOutFile, "outfile", "o", "", "File to use for output")
	cmd.Flags().BoolVarP(&flagPersist, "persist", "", false, "persist the merged dictionary to xld")
	cmd.Flags().BoolVarP(&flagForgetConflicts, "fc", "", false, "Forget Conflicts: do not write the conflict information to the new dictionary")
	cmd.Flags().StringVarP(&flagDictID, "id", "", "", "ID of the new dictionary")
	cmd.Flags().StringVarP(&flagEnvID, "env", "", "", "ID of the environment to merge dictionaries for (will merge all dictionaries for that environment into one)")

	relCmd.AddCommand(cmd)

}

func runDictMerge(cmd *cobra.Command, args []string) {

	var dicts []string
	var cis xld.Cis
	var nc xld.Ci
	var conf conflicts
	var name string

	// get collection of dict ci's to merge
	// the list of dictionaries to merge should be in args
	// if --env flag is set get the list of dicts from the environment specified
	if flagEnvID != "" {
		dicts = getDictListFromEnv(flagEnvID)
	} else {
		dicts = args
	}

	cis = getDictCis(dicts)

	// create the new ci object representing the result of our merge operation

	// get the new name
	if flagDictID == "" {
		name = createNewDictName(cis)
	} else {
		name = flagDictID
	}

	// fully initialize the new xld.Ci object so that we can merge in the properties from the other dictionaries
	nc = getNewDict(name)

	// TODO: should we be merging dictionaries with restrictions ?? (maybe not if we are working on an environment)

	// loop over the ci collection of ci's
	for _, c := range cis {
		newConf, err := mergeDictProperties(&nc, c)
		utils.HandleErr(err)
		for _, cn := range newConf {
			conf = append(conf, cn)

			//fmt.Println(utils.RenderJSON(conf))

		}
		// merge each dictionarie with the newly created one
		// for range
		//  we should be working on the actual object not a copy

		//  mergeDictCIS(&t, ms) conflicts, err

	}

	// adding the conflicting key value pairs to the dictionary (prefixing the key with the originating dictionary)
	if flagForgetConflicts == false {
		addConflictsToCi(&nc, conf)
	}

	// checking properties for null fields ..
	output := utils.RenderJSON(nc)

	if flagOutFile != "" {
		utils.WriteToFile(output, flagOutFile)
		return
	}

	if flagPersist == true {
		c := utils.GetClient()
		_, err := c.Repository.SaveCi(nc)
		utils.HandleErr(err)
	}

	fmt.Println(output)
	// handle the new dictionary
	// present on screen
	// write to file
	// persist in xld

}

func addConflictsToCi(cc *xld.Ci, cnflcts conflicts) {
	nc := *cc
	entries := nc.Properties["entries"].(map[string]interface{})

	for _, cnf := range cnflcts {
		key := cnf.ConflictDict + ":" + cnf.Key
		entries[key] = cnf.Conflict.(string)
	}
	nc.Properties["entries"] = entries
	*cc = nc
}

func mergeDictProperties(cc *xld.Ci, c xld.Ci) (conflicts, error) {
	nc := *cc
	var outputCf conflicts
	var cf conflicts

	out := make(map[string]interface{})

	//dictPropNames := [4]string{"entries", "encryptedEntries", "restrictToApplications", "restrictToContainers"}

	for _, p := range [2]string{"entries", "encryptedEntries"} {

		m := c.Properties[p].(map[string]interface{})
		t := nc.Properties[p].(map[string]interface{})

		out, cf = mergeEntries(t, m)

		for _, con := range cf {
			con.ConflictDict = c.ID
			outputCf = append(outputCf, con)
		}
		// fmt.Println("mdp")
		// fmt.Printf("%+v\n", cf)

		nc.Properties[p] = out

	}

	for _, p := range [2]string{"restrictToApplications", "restrictToContainers"} {
		m := c.Properties[p]
		t := nc.Properties[p]

		nc.Properties[p] = mergeRestrictions(t, m)
	}
	*cc = nc
	return outputCf, nil
}

func mergeEntries(a, b map[string]interface{}) (map[string]interface{}, conflicts) {

	newDict := make(map[string]interface{})

	var cf conflicts

	abKeyDiff := missingKeys(a, b)
	baKeyDiff := missingKeys(b, a)
	abMutual := mutualKeys(a, b)
	abVallDiff := valDiff(a, b)

	// now we need to make a map[string]interface{} with the mutualKeys without value differens and the different keys.

	// the mutual keys that have value difference are conflicting and thus pose a problem
	for _, k := range abKeyDiff {
		newDict[k] = a[k]
	}
	for _, k := range baKeyDiff {
		newDict[k] = b[k]
	}
	for _, k := range abMutual {
		if mapContains(abVallDiff, k) == false {
			newDict[k] = a[k]
		}
	}

	for _, k := range abVallDiff {
		cf = append(cf, conflict{Key: k, Original: a[k], Conflict: b[k]})
	}
	return newDict, cf
}

// merges array type properties
func mergeRestrictions(a, b interface{}) []string {
	var newRestrict []string
	aa := a.([]string)

	for _, x := range b.([]interface{}) {
		aa = append(aa, x.(string))
	}

	for _, r := range aa {
		if mapContains(newRestrict, r) == false {
			newRestrict = append(newRestrict, r)

		}
	}
	return newRestrict
}

//getDictListFromEnv retrieves the list of used dictionaries from and environment
func getDictListFromEnv(e string) []string {
	var dictionaries []string
	envCi := utils.GetCi(e, "udm.Environment")

	ds := envCi.Properties["dictionaries"].([]interface{})

	for _, v := range ds {
		dictionaries = append(dictionaries, v.(string))
	}
	return dictionaries
}

//getDictCis gets a collection of udm.Dictionary cis from xld
func getDictCis(d []string) xld.Cis {
	var cis xld.Cis

	for _, dn := range d {
		ci := utils.GetCi(dn, "udm.Dictionary")
		cis = append(cis, ci)
	}

	return cis
}

func createNewDictName(d xld.Cis) string {

	var nn string

	// merge all the names from dictionaries into one name ... this can be renamed later
	for _, c := range d {
		nn = nn + "_" + c.Name()
	}

	nn = d[0].Path() + "/" + nn

	return nn

}

func getNewDict(n string) xld.Ci {

	newProps := make(map[string]interface{})

	ci := utils.NewCiObject(n, "udm.Dictionary", newProps)
	ci.Properties = make(map[string]interface{})
	ci.Properties["entries"] = make(map[string]interface{})
	ci.Properties["encryptedEntries"] = make(map[string]interface{})
	ci.Properties["restrictToApplications"] = []string{}
	ci.Properties["restrictToContainers"] = []string{}

	return ci
}

//valDiff
func valDiff(l, s map[string]interface{}) []string {
	var vd []string
	mut := mutualKeys(l, s)

	for _, k := range mut {
		if l[k] != s[k] {
			vd = append(vd, k)
		}
	}
	return vd
}

func mutualKeys(l, s map[string]interface{}) []string {

	var mutualKeys []string

	for k := range l {
		if hasKey(k, s) {
			mutualKeys = append(mutualKeys, k)
		}
	}
	return mutualKeys
}

func missingKeys(l, s map[string]interface{}) []string {

	var missingKeys []string

	for k := range l {
		if hasKey(k, s) != true {
			missingKeys = append(missingKeys, k)
		}
	}
	return missingKeys
}

func hasKey(key string, m map[string]interface{}) bool {
	for k := range m {
		if k == key {
			return true
		}
	}
	return false
}

func isDict(c xld.Ci) bool {
	if c.Type == "udm.Dictionary" {
		return true
	}

	return false
}

func mapContains(m []string, key string) bool {
	for _, k := range m {
		if k == key {
			return true
		}
	}
	return false
}
