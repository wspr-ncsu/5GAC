package main

import (
	"log"
	"os"

	"github.com/wspr-ncsu/5GAC/gocode/Instrumenter/cfg"
	"github.com/wspr-ncsu/5GAC/gocode/Instrumenter/instrumentation"
	"github.com/wspr-ncsu/5GAC/gocode/Shared/accesspattern"
	"gopkg.in/yaml.v2"
)

// main entry point into the program, where everything begins.
func main() {
	// Command line args
	c := cfg.RegisterCommandLineArgs()

	if c.Verbosity {
		log.Printf("%v [%v] started.\n", c.ToolName, c.Version)
	}

	// The first step, is to parse all specification files and free5gc code files.
	ft := accesspattern.ParseAllFiles()

	if len(c.Policy) == 0 && len(c.AccessPatterns) == 0 {
		// log.Fatal("Error: Invalid authorization policy.")
	} else if _, err := os.Stat(c.Policy); err == nil {
		file, err := os.ReadFile(c.Policy)
		if err != nil {
			log.Fatalf("failed to read file %v\n", c.Policy)
		}
		ap := make(map[string][]accesspattern.API)

		err2 := yaml.Unmarshal(file, &ap)
		if err2 != nil {
			log.Fatalf("failed to read YAML in %v\n", c.Policy)
		}
		if c.Verbosity {
			log.Printf("%v\n", ap)
		}

		accesspattern.UpdateSecurity(ap)
	} else if _, err := os.Stat(c.AccessPatterns); err == nil {
		file, err := os.ReadFile(c.AccessPatterns)
		if err != nil {
			log.Fatalf("failed to read file %v\n", c.AccessPatterns)
		}

		accessPatterns := make([]accesspattern.AccessPattern, 0, 100)

		err2 := yaml.Unmarshal(file, &accessPatterns)
		if err2 != nil {
			log.Fatalf("Failed to read YAML access patterns in %v", c.AccessPatterns)
		}

		accesspattern.SetAccessPatterns(accessPatterns)
		accesspattern.UpdateSecurity(accessPatterns)

		if c.Verbosity {
			log.Printf("%v\n", accessPatterns)
		}
	}

	for _, code := range ft.Code {
		instrumentation.UpdateAst(code, false)
	}
	instrumentation.MoveTemplates()
}
