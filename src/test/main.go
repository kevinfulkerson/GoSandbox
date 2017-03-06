package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
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
	fmt.Printf("map value:\n%v\n\n", fileMap)

	// TODO: clean up the following and make it generic
	// Print out a value from the structure
	third := fileMap["third"]
	fmt.Printf("third:\n%v\n\n", third)

	// Use a type assertion to extract the map to use again
	thirdMap, ok := third.(map[interface{}]interface{})
	if !ok {
		// Not sure what to use here...
		panic(-1)
	}

	// Print out a value from the structure
	third_first := thirdMap["third_first"]
	fmt.Printf("third_first:\n%v\n\n", third_first)

	// Use another type assertion
	thirdFirstMap, ok := third_first.(map[interface{}]interface{})
	if !ok {
		// Not sure what to use here...
		panic(-2)
	}

	// Print out a value from the structure
	third_first_first := thirdFirstMap["third_first_first"]
	fmt.Printf("third_first_first:\n%v\n\n", third_first_first)

	// Reflect the value out of the interface, if possible
	reflectedValue := reflect.ValueOf(third_first_first)
	if reflectedValue.Kind() != reflect.Slice {
		// Not sure what to use here...
		panic(-3)
	}

	// Make a new slice with the same size as the existing slice
	thirdFirstFirstSlice := make([]interface{}, reflectedValue.Len())

	// Copy the values
	for i := 0; i < reflectedValue.Len(); i++ {
		thirdFirstFirstSlice[i] = reflectedValue.Index(i).Interface()
	}

	// Print out a value from the structure
	third_first_first_value2 := thirdFirstFirstSlice[1]
	fmt.Printf("third_first_first_value2:\n%d\n\n", third_first_first_value2)
}
