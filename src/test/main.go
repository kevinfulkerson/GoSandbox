package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
)

var filename = "res/test.txt"

func main() {
	// Read in the file using the basic utility method
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Allocate a map to hold the contents, then decode the contents
	fileMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal(contents, fileMap)
	if err != nil {
		panic(err)
	}

	// Dump the entire string
	fmt.Printf("\nparsed contents:\n%v\n\n", fileMap)

	// Print out each of the elements read from the file
	recursivePrint("root", fileMap)

	// Search for a value in the map
	keys := []string{"second", "second", "1"}
	value, success := findValueInMap(keys, fileMap)
	if success {
		fmt.Printf("\nvalue = %s\n", value)
	}
}

// recursivePrint prints the structure of the given parsed YAML map.
// BUG(ksf): Currently makes a lot of assumptions about how the data is structured.
func recursivePrint(chain string, structure interface{}) {
	localMap, isMap := structure.(map[interface{}]interface{})
	localSlice, isSlice := structure.([]interface{})
	localInteger, isInteger := structure.(int)
	localString, isString := structure.(string)

	switch {
	case isMap:
		for key, value := range localMap {
			nextChain := chain + "_" + key.(string)
			recursivePrint(nextChain, value)
		}
	case isSlice:
		for i := 0; i < len(localSlice); i++ {
			fmt.Printf("%s_value_%d:\n%v\n", chain, i, localSlice[i])
		}
	case isInteger:
		fmt.Printf("%s_value:\n%v\n", chain, localInteger)
	case isString:
		fmt.Printf("%s_value:\n%v\n", chain, localString)
	}
}

// findValueInMap finds a value in a passed structure. It assumes that the mapToSearch is a map
// but is also lenient toward simply an array (or even a single value), abstracted as a blank
// interface.
// It returns the value as a string and a flag indicating if the operation was successful.
// BUG(ksf): Currently makes a lot of assumptions about how the data is structured.
func findValueInMap(keys []string, mapToSearch interface{}) (value string, foundValue bool) {
	// Default return value to be false
	foundValue = false

	localStructure := mapToSearch
	for i := 0; i < len(keys); i++ {
		localMap, isMap := localStructure.(map[interface{}]interface{})
		localSlice, isSlice := localStructure.([]interface{})
		localInteger, isInteger := localStructure.(int)
		localString, isString := localStructure.(string)

		switch {
		case isMap:
			localStructure = localMap[keys[i]]
			if localStructure == nil {
				foundValue = false
				return
			}
		case isSlice:
			keyAsInteger, err := strconv.ParseInt(keys[i], 10, 32)
			if err != nil {
				panic("findValueInMap(keys() - Failed to convert key as an integer")
			}

			value = fmt.Sprintf("%v", localSlice[keyAsInteger])
			foundValue = true
		case isInteger:
			value = strconv.Itoa(localInteger)
			foundValue = true
		case isString:
			value = localString
			foundValue = true
		}
	}
	return
}
