package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var filename = "res/test.txt"
var fields = []string{"third", "third_first", "third_first_first", "third_first_first_value2"}

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
	fmt.Printf("parsed contents:\n%v\n\n", fileMap)

	// Print out a value from the structure
	third := fileMap[fields[0]]
	fmt.Printf("%s:\n%v\n\n", fields[0], third)

	currentStructure := fileMap[fields[0]]
	for i := 1; i < len(fields); i++ {
		// Use a type assertion to attempt to extract some type of value
		localMap, isMap := currentStructure.(map[interface{}]interface{})
		localSlice, isSlice := currentStructure.([]interface{})

		// TODO: consider other cases here?
		switch {
		case isMap:
			// This was a map of blank interfaces
			currentStructure = localMap[fields[i]]
			fmt.Printf("%s:\n%v\n\n", fields[i], currentStructure)
		case isSlice:
			// This was a slice of blank interfaces
			value := localSlice[1]
			fmt.Printf("value:\n%d\n\n", value)
		}
	}
}
