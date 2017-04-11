package main

import (
	"fmt"
	"go-pkg-optarg"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
)

var inputFileMap *DataMap
var dataFileMap *DataMap
var recipeFileMap *DataMap

func main() {
	// Initialize maps
	inputFileMap = NewDataMap()
	dataFileMap = NewDataMap()
	recipeFileMap = NewDataMap()

	// Read arguments
	parseArguments()

	// Allocate some space to this variable, it can expand
	dataSource := make([]byte, 100)

	// Read in the data source files
	for directory, directoryMap := range dataFileMap.Mapping {
		for _, file := range directoryMap.Files {
			// For now, just check if its the correct file
			if file.Name() == "test.txt" {
				ds, err := ioutil.ReadFile(directory + file.Name())
				if err != nil {
					panic(err)
				}

				dataSource = ds
			}
		}
	}

	// Allocate some space to this variable, it can expand
	lookupValue := make([]byte, 100)

	// Read in the lookup file
	for directory, directoryMap := range inputFileMap.Mapping {
		for _, file := range directoryMap.Files {
			// For now, just check if its the correct file
			if file.Name() == "lookup.txt" {
				lv, err := ioutil.ReadFile(directory + file.Name())
				if err != nil {
					panic(err)
				}

				lookupValue = lv
			}
		}
	}

	// Allocate a map to hold the contents, then decode the contents
	fileMap := make(map[interface{}]interface{})
	err := yaml.Unmarshal(dataSource, fileMap)
	if err != nil {
		panic(err)
	}

	// Dump the entire string
	fmt.Printf("\nparsed contents:\n%v\n\n", fileMap)

	// Print out each of the elements read from the file
	recursivePrint("", fileMap)

	fmt.Print("\n\n")

	// Convert the read-in value to a list of elements to use for looking up a value
	lookupList := ParseLookupString(string(lookupValue))
	for _, val := range lookupList {
		fmt.Printf("val = %v\n", val)
	}

	// Search for a value in the map
	value, success := findValueInMap(lookupList, fileMap)
	if success {
		fmt.Printf("\nvalue = %s\n", value)
	}
}

// recursivePrint prints the structure of the given parsed YAML map.
func recursivePrint(chain string, structure interface{}) {
	localMap, isMap := structure.(map[interface{}]interface{})
	localSlice, isSlice := structure.([]interface{})

	switch {
	case isMap:
		for key, value := range localMap {
			nextChain := chain + "/" + fmt.Sprintf("%v", key)
			recursivePrint(nextChain, value)
		}
	case isSlice:
		for i := 0; i < len(localSlice); i++ {
			nextChain := chain + fmt.Sprintf("[%d]", i)
			recursivePrint(nextChain, localSlice[i])
		}
	default:
		fmt.Printf("%s:\n%v\n", chain, structure)
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

// parseArguments parses the command-line arguments for the current execution of the program and sets the
// internal state of the program to the configured values.
func parseArguments() {
	optarg.Add("i",
		"inputDirectory",
		"The path to the directory to draw input files from.",
		"./res/input/")
	optarg.Add("o",
		"outputDirectory",
		"The path to the root directory to write output files to.",
		"./res/output/")
	optarg.Add("r",
		"recipeDirectory",
		"The path to the directory to draw recipe files from.",
		"./res/recipes/")
	optarg.Add("d",
		"dataDirectory",
		"The path to the directory to draw global data files from.",
		"./res/data/")

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "i":
			err := inputFileMap.MapFiles(opt.String())
			if err != nil {
				panic(err)
			}
		case "o":

		case "r":
			err := recipeFileMap.MapFiles(opt.String())
			if err != nil {
				panic(err)
			}
		case "d":
			err := dataFileMap.MapFiles(opt.String())
			if err != nil {
				panic(err)
			}
		}
	}
}
