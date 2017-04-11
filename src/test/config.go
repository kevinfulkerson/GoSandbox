package main

import (
	"io/ioutil"
	"os"
)

type DirectoryMap struct {
	Files []os.FileInfo
}

type DataMap struct {
	Mapping map[string]DirectoryMap
}

func NewDataMap() *DataMap {
	return &DataMap{make(map[string]DirectoryMap)}
}

// mapFiles places the files of a directory into a FileMapping struct with appropriate working directory.
// It is tolerant of nested directories, and handles them using recursion.
// It returns an error code if a directory in the hierarchy cannot be read.
func (fileMap *DataMap) MapFiles(directory string) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			if directory[len(directory)-1] != '/' {
				directory += "/"
			}

			directory += file.Name()
			fileMap.MapFiles(directory)
		} else {
			directoryMap := new(DirectoryMap)

			if directory[len(directory)-1] != '/' {
				directory += "/"
			}

			directoryMap.Files = append(directoryMap.Files, file)
			fileMap.Mapping[directory] = *directoryMap
		}
	}

	return nil
}
