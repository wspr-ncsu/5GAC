package accesspattern

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// getFileList takes a path, and an extension of files in the path
// and returns a file list that is in that path and has that extension.
func getFileList(pathParam string, ext string) []fileInfo {
	fileList := []fileInfo{}

	// Get current working directory, incase this is a relative path.
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Walk all the files in the directory
	err = filepath.Walk(pathParam, func(path string, file os.FileInfo, err error) error {
		if file == nil {
			log.Fatalf("%s - Directory doesn't exist?", path)
		}

		if file != nil && file.IsDir() {
			return nil
		}

		// Only consider *.go files
		if !strings.HasSuffix(path, "."+ext) {
			return nil
		}

		pathToScan := filepath.Join(cwd, pathParam)
		joined := filepath.Join(cwd, path)

		relPath, relerr := filepath.Rel(pathToScan, joined)
		if relerr != nil {
			log.Fatal(err)
		}

		// fmt.Printf("Adding %v - %v\n", path, relPath)

		fileList = append(fileList, fileInfo{FullPath: path, name: file.Name(), RelativePath: relPath})

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return fileList
}

// getDataFiles return all the files in the data directory that are Go files.
func getDataFiles() []fileInfo {
	return getFileList("../_code", "go")
}

// getOpenAPIFiles return all the files in the specs directory that are YAML files.
func getOpenAPIFiles() []fileInfo {
	return getFileList("../_specs/v15.2.0", "yaml")
}

// Get both the specs and data files.
func getAllFiles() ([]fileInfo, []fileInfo) {
	return getOpenAPIFiles(), getDataFiles()
}
