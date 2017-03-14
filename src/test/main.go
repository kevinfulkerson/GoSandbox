package main

import (
	"fmt"
	"go-pkg-optarg"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
)

type FileMapping struct {
	files            []os.FileInfo
	workingDirectory string
}

var inputFileMap FileMapping
var dataFileMap FileMapping

func main() {
	// Read arguments
	parseArguments()

	// Read in the data source file
	dataSource, err := ioutil.ReadFile(dataFileMap.getFileLocation(0))
	if err != nil {
		panic(err)
	}

	// Read in the lookup file
	lookupValue, err := ioutil.ReadFile(inputFileMap.getFileLocation(0))
	if err != nil {
		panic(err)
	}

	// Allocate a map to hold the contents, then decode the contents
	fileMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal(dataSource, fileMap)
	if err != nil {
		panic(err)
	}

	// Dump the entire string
	fmt.Printf("\nparsed contents:\n%v\n\n", fileMap)

	// Print out each of the elements read from the file
	recursivePrint("", fileMap)

	fmt.Print("\n\n")

	// Convert the read-in value to a list of elements to use for looking up a value
	lookupList := parseLookupString(string(lookupValue))
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

// parseLookupString parses the provided string for elements to insert into a lookup array.
// It returns the ordered values as a slice.
// BUG(ksf) Does not use a proper interface, so it may have some strange interactions
func parseLookupString(lookupString string) (lookupValues []string) {
	tempString := ""
	elementStarted := false
	for i := 0; i < len(lookupString); i++ {
		switch lookupString[i] {
		case '@':
			// In most cases, this will be the start of an element id indication, so
			// look for the rest of the indicator. The indicator is currently:
			// @id
			addToString := true
			if i+2 <= len(lookupString) {
				if lookupString[i+1] == 'i' && lookupString[i+2] == 'd' {
					addToString = false
					elementStarted = true
					i += 2
				}
			}

			if addToString {
				tempString += string(lookupString[i])
			}
		case '/':
			if elementStarted {
				// Check if we have anything to insert, and if so, insert it.
				if tempString != "" {
					lookupValues = append(lookupValues, tempString)
					tempString = ""
				}
			} else {
				tempString += string(lookupString[i])
			}
		case '[':
			if elementStarted {
				// This is an array indication, so we need to insert the index
				// after its parent. To do this, insert the parent first.
				lookupValues = append(lookupValues, tempString)

				// Set the temp value to the empty string so the next '/' token
				// we encounter (after the ending array indicator) won't cause an
				// insert of a blank value in the array.
				tempString = ""

				// The parent is in correctly, so begin inserting the index value.
				str := ""
				i++
				for lookupString[i] != ']' {
					str += string(lookupString[i])
					i++
				}

				lookupValues = append(lookupValues, str)
			} else {
				tempString += string(lookupString[i])
			}
		default:
			if lookupString[i] != '\n' && lookupString[i] != '\r' {
				tempString += string(lookupString[i])
			}
		}
	}

	return append(lookupValues, tempString)
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
			files, err := ioutil.ReadDir(opt.String())
			if err != nil {
				panic(err)
			}

			err = inputFileMap.mapFiles(files, opt.String())
			if err != nil {
				panic(err)
			}
		case "o":
		case "r":
		case "d":
			files, err := ioutil.ReadDir(opt.String())
			if err != nil {
				panic(err)
			}

			err = dataFileMap.mapFiles(files, opt.String())
			if err != nil {
				panic(err)
			}
		}
	}
}

// mapFiles places the files of a directory into a FileMapping struct with appropriate working directory.
// It is tolerant of nested directories, and handles them using recursion.
// It returns an error code if a directory in the hierarchy cannot be read.
func (fileMap *FileMapping) mapFiles(files []os.FileInfo, directory string) error {
	for _, file := range files {
		if file.IsDir() {
			innerFiles, err := ioutil.ReadDir(file.Name())
			if err != nil {
				return err
			}
			fileMap.mapFiles(innerFiles, directory+file.Name())
		} else {
			fileMap.files = append(fileMap.files, file)

			if directory[len(directory)-1] != '/' {
				directory += "/"
			}

			fileMap.workingDirectory = directory
		}
	}

	return nil
}

// getFileLocation generates a string for the file at the specified index in the map that can be used to locate that
// file.
func (fileMap *FileMapping) getFileLocation(fileIndex int) string {
	return fileMap.workingDirectory + fileMap.files[fileIndex].Name()
}
