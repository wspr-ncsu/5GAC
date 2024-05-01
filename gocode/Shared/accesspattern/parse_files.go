package accesspattern

import (
	"log"
	"sync"

	"github.com/wspr-ncsu/5GAC/gocode/Shared/openapi"
)

var (
	// Mutex for the cached files, so we don't overwrite important data.
	cachedFilesMutex map[string]*sync.Mutex

	// Keep track of the files we parse, so we don't need to parse them more than once.
	cachedFiles FileTracker
)

// parseSpecs read the specifications directory and parses the OpenAPI in them.
func parseSpecs(files []fileInfo) {
	// If we already parsed the specs, just return that.
	if cachedFiles.Specs == nil {
		cachedFiles.Specs = make(map[fileInfo]openapi.APIs)
	} else {
		log.Fatal("Skipping parse specs: already parsed")
	}

	// cachedFiles.specs = make(map[fileInfo]openapi.APIs)

	waitgroup := sync.WaitGroup{}
	cachedFilesMutex["specs"] = &sync.Mutex{}

	for _, file := range files {
		waitgroup.Add(1)

		go func(file fileInfo) {
			apis := openapi.ReadSpec(file.FullPath)

			cachedFilesMutex["specs"].Lock()
			cachedFiles.Specs[file] = apis
			cachedFilesMutex["specs"].Unlock()

			waitgroup.Done()
		}(file)

		waitgroup.Wait()
	}
}

// ParseDataFiles parses the files in the data directory, assuming they're all golang files.
func ParseDataFiles(files []fileInfo) error {
	// Just return the cached files if we already parsed them before.
	if cachedFiles.Code == nil {
		cachedFiles.Code = make(map[fileInfo]GoFile)
	} else {
		// always re-parse. TODO
		log.Fatal("Skipping parse data files: already parsed")
	}

	// cachedFiles.dataFiles = make(map[fileInfo]my_ast.GoFunction)

	waitgroup := sync.WaitGroup{}
	cachedFilesMutex["data"] = &sync.Mutex{}

	for _, file := range files {
		waitgroup.Add(1)

		go func(file fileInfo) {
			goFile := GoFile{}
			goFile.Parse(file.FullPath)

			goFile.RelativePath = file.RelativePath

			cachedFilesMutex["data"].Lock()
			cachedFiles.Code[file] = goFile
			cachedFilesMutex["data"].Unlock()

			waitgroup.Done()
		}(file)

		waitgroup.Wait()
	}

	return nil
}

func initialize() ([]fileInfo, []fileInfo) {
	CompileRegexp()
	ResetFileSet()

	cachedFilesMutex = make(map[string]*sync.Mutex)

	return getAllFiles()
}

// parseAllFiles parses the specification and data files.
func ParseAllFiles() FileTracker {
	specs, dataFiles := initialize()

	// MUST parse specs before data files!
	// This is because the parsing of data files needs to map security requirements in the specs.
	parseSpecs(specs)

	err := ParseDataFiles(dataFiles)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}

	return cachedFiles
}
