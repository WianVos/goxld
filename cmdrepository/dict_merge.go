package cmdrepository

import (
	"fmt"
	"os"
	"time"

	"github.com/WianVos/goxld/utils"
	"github.com/WianVos/xld"
	"github.com/spf13/cobra"
)

var dmlistLong = `Compare two dictionaries
Example:
  repository dict_compare /Environments/Dictionary1 /Environments/Dictionary2
`

type conflict struct {
	original interface{}
	conflict interface{}
}

var flagPersist bool
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
	cmd.Flags().StringVarP(&flagDictID, "id", "", "", "ID of the new dictionary")
	cmd.Flags().StringVarP(&flagEnvID, "env", "", "", "ID of the environment to merge dictionaries for (will merge all dictionaries for that environment into one)")

	relCmd.AddCommand(cmd)

}

func runDictMerge(cmd *cobra.Command, args []string) {

	var ci xld.Ci

	if flagEnvID != "" {
		ci = runDictMergeEnv(flagEnvID)
	} else {

		if len(args) > 2 {
			fmt.Println("too many arguments")
			os.Exit(2)
		}

		ci = mergeDictionary(args[0], args[1], flagDictID)
	}

	output := utils.RenderJSON(ci)

	if flagOutFile != "" {
		utils.WriteToFile(output, flagOutFile)
	} else {
		fmt.Println(output)
	}

	if flagPersist == true {
		client := utils.GetClient()
		client.Repository.CreateCi(ci.ID, "udm.Dictionary", ci.Properties)
		fmt.Println("dictionary persisted")
	}
}

func runDictMergeEnv(i string) xld.Ci {

	newProps := make(map[string]interface{})
	// var cis xld.Cis

	client := utils.GetClient()
	envCi, err := client.Repository.GetCi(i)
	if err != nil {
		panic(fmt.Errorf("Fatal error dict merge: %s \n", err))
	}

	dictionaries := envCi.Properties["dictionaries"].([]interface{})

	s, dictionaries := dictionaries[0].(string), dictionaries[1:]

	fmt.Println(s)
	fmt.Println(dictionaries)

	newname := i + "mergeDict"

	newProps["entries"] = map[string]string{}
	newProps["encryptedEntries"] = map[string]string{}
	newProps["restrictToApplications"] = []string{}
	newProps["restrictToContainers"] = []string{}

	ci, err := client.Repository.NewCi(newname, "udm.Dictionary", newProps)

	if err != nil {
		panic(fmt.Errorf("Fatal error dict merge: %s \n", err))
	}

	for _, d := range dictionaries {
		fmt.Println(d)
	}
	return ci
}

//mergeDictionary
// takes care or the mergin of two dictionaries
// s represents the source dictionary
// m represents the merge candidate
// i is the new dict name

func mergeDictionary(s, m, i string) xld.Ci {

	var newName string
	client := utils.GetClient()

	// get the source dictionary
	source, err := client.Repository.GetCi(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// get the dictionary
	merge, err := client.Repository.GetCi(m)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// create a name for the new dictionary
	if i == "" {
		newName = source.Path() + "/merged_" + source.Name() + "_" + merge.Name()
	} else {
		newName = i
	}

	a := getDictEntries(source, false)
	b := getDictEntries(merge, false)
	ae := getDictEntries(source, true)
	be := getDictEntries(merge, true)

	mu, conflicts := mergeEntries(a, b, true)
	me, conflictsEncrypted := mergeEntries(ae, be, true)

	// find duplicates between the m and me .. (those cause problems as well)
	//ec := mutualKeys(m, me)

	ci := xld.Ci{}
	ci.Type = "udm.Dictionary"

	//add merge info to dictionary
	mu["merge_date"] = time.Now()
	mu["merge_base"] = source.ID
	mu["merge_target"] = merge.ID
	mu["merge_conflicts"] = conflictsEncrypted

	// add merge conflicts to the dictionary in a non intrusive way
	for k, v := range conflicts {
		mu[k] = v.original
		mu[k+"_cnflct"] = v.conflict
	}

	ci.ID = newName
	ci.Properties = make(map[string]interface{})
	ci.Properties["entries"] = make(map[string]interface{})
	ci.Properties["entries"] = mu
	ci.Properties["encryptedEntries"] = make(map[string]interface{})
	ci.Properties["encryptedEntries"] = me
	ci.Properties["restrictToApplications"] = mergeRestrictions(source.Properties["restrictToApplications"], merge.Properties["restrictToApplications"])
	ci.Properties["restrictToContainers"] = mergeRestrictions(source.Properties["restrictToContainers"], merge.Properties["restrictToContainers"])

	return ci

}

// merges
func mergeRestrictions(a, b interface{}) []string {
	var newRestrict []string
	aa := a.([]interface{})

	for _, x := range b.([]interface{}) {
		aa = append(aa, x)
	}

	for _, r := range aa {
		if mapContains(newRestrict, r.(string)) == false {
			newRestrict = append(newRestrict, r.(string))

		}
	}
	return newRestrict
}

func mergeEntries(a, b map[string]interface{}, i bool) (map[string]interface{}, map[string]conflict) {

	newDict := make(map[string]interface{})

	conflicts := make(map[string]conflict)

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
		conflicts[k] = conflict{
			original: a[k],
			conflict: b[k],
		}
	}
	return newDict, conflicts
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

func getDictEntries(d xld.Ci, e bool) map[string]interface{} {

	if isDict(d) == false {
		fmt.Println("argument not a dictionary")
		os.Exit(2)
	}

	// determine if we need the regular or encrypted entries
	f := "entries"

	if e == true {
		f = "encryptedEntries"
	}

	return d.Properties[f].(map[string]interface{})

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

// func (d *xld.Ci) dictMerge(m xld.Ci) {
//
// 	a := getDictEntries(d, false)
// 	b := getDictEntries(m, false)
// 	ae := getDictEntries(d, true)
// 	be := getDictEntries(m, true)
//
// 	mu, conflicts := mergeEntries(a, b, true)
// 	me, conflictsEncrypted := mergeEntries(ae, be, true)
//
// 	//add merge info to dictionary
// 	mu["merge_date"] = time.Now()
// 	mu["merge_base"] = source.ID
// 	mu["merge_target"] = merge.ID
// 	mu["merge_conflicts"] = conflictsEncrypted
//
// 	// add merge conflicts to the dictionary in a non intrusive way
// 	for k, v := range conflicts {
// 		mu[k] = v.original
// 		mu[k+"_cnflct"] = v.conflict
// 	}
//
// 	if val, ok := d.Properties["entries"]; ok != false {
// 		d.Properties["entries"] = make(map[string]interface{})
// 	}
//
// 	if val, ok := d.Properties["encryptedEntries"]; ok != false {
// 		d.Properties["encryptedEntries"] = make(map[string]interface{})
// 	}
//
// 	d.Properties = make(map[string]interface{})
//
// 	d.Properties["entries"] = mu
// 	d.Properties["encryptedEntries"] = me
// 	d.Properties["restrictToApplications"] = mergeRestrictions(d.Properties["restrictToApplications"], m.Properties["restrictToApplications"])
// 	d.Properties["restrictToContainers"] = mergeRestrictions(d.Properties["restrictToContainers"], m.Properties["restrictToContainers"])
//
// }
