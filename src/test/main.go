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

	fmt.Print("\n\n")

	// Search for a value in the map
	keys := []string{"third", "first", "first", "0", "five", "2"}
	value, success := findValueInMap(keys, fileMap)
	if success {
		fmt.Printf("first value = %s\n", value)
	}

	// Search for another value in the map
	keys = []string{"second", "second", "1"}
	value, success = findValueInMap(keys, fileMap)
	if success {
		fmt.Printf("second value = %s\n", value)
	}
}

// recursivePrint prints the structure of the given parsed YAML map.
func recursivePrint(chain string, structure interface{}) {
	localMap, isMap := structure.(map[interface{}]interface{})
	localSlice, isSlice := structure.([]interface{})

	switch {
	case isMap:
		for key, value := range localMap {
			nextChain := chain + "_" + fmt.Sprintf("%v", key)
			recursivePrint(nextChain, value)
		}
	case isSlice:
		for i := 0; i < len(localSlice); i++ {
			nextChain := chain + fmt.Sprintf("[%d]", i)
			recursivePrint(nextChain, localSlice[i])
		}
	default:
		fmt.Printf("%s_bool:\n%v\n", chain, structure)
	}
}

// findValueInMap finds a value in a passed structure. It assumes that the mapToSearch is a map
// but is also lenient toward simply an array (or even a single value), abstracted as a blank
// interface.
// It returns the value as a string and a flag indicating if the operation was successful.
func findValueInMap(keys []string, mapToSearch interface{}) (value string, foundValue bool) {
	// Default return value to be false
	foundValue = false

	localStructure := mapToSearch
	for i := 0; i < len(keys); i++ {
		localMap, isMap := localStructure.(map[interface{}]interface{})
		localSlice, isSlice := localStructure.([]interface{})

		switch {
		case isMap:
			localStructure = localMap[keys[i]]
			if localStructure == nil {
				return
			} else {
				// If the value is not a map or slice, then decrement i to reinterpret the value
				// using the same key
				localMap, isMap = localStructure.(map[interface{}]interface{})
				localSlice, isSlice = localStructure.([]interface{})
				if !isMap && !isSlice {
					i--
				}
			}
		case isSlice:
			keyAsInteger, err := strconv.Atoi(keys[i])
			if err != nil {
				return
			}

			if keyAsInteger <= len(localSlice) {
				localStructure = localSlice[keyAsInteger]

				// If the value is not a map or slice, then decrement i to reinterpret the value
				// using the same key
				localMap, isMap = localStructure.(map[interface{}]interface{})
				localSlice, isSlice = localStructure.([]interface{})
				if !isMap && !isSlice {
					i--
				}
			} else {
				return
			}
		default:
			value = fmt.Sprintf("%v", localStructure)
			foundValue = true
		}
	}
	return
}
