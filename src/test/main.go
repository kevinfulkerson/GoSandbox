package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	// TODO: Change the element printing to not use recursion
	recursivePrint("root", fileMap)
}

func recursivePrint(chain string, structure interface{}) {
	// Use a type assertion to attempt to extract some type of value
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
