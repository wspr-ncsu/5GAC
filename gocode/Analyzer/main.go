package main

import (
	"log"
	"os"

	graph "github.com/wspr-ncsu/5GAC/gocode/Analyzer/Graph"
	policy "github.com/wspr-ncsu/5GAC/gocode/Analyzer/Policy"
	"github.com/wspr-ncsu/5GAC/gocode/Analyzer/cfg"
	"github.com/wspr-ncsu/5GAC/gocode/Shared/accesspattern"
)

// main entry point into the program, where everything begins.
func main() {
	// Command line args
	c := cfg.RegisterCommandLineArgs()

	if c.Verbosity {
		log.Printf("%v [%v] started.\n", c.ToolName, c.Version)
	}

	// The first step, is to parse all specification files and free5gc code files
	files := accesspattern.ParseAllFiles()

	if c.Verbosity {
		log.Printf("%v\n", files)
	}

	_, err := os.Stat(c.AccessPatterns)

	if len(c.AccessPatterns) > 0 && err == nil {
		accesspattern.CompileRegexp()
		accesspattern.ParseJSON(c.AccessPatterns)
	}

	accessPatterns := accesspattern.GetAccessPatterns()
	graph.GenerateGraph(accessPatterns, "3gpp")

	newPolicy := policy.GenerateSecurePolicy(accessPatterns)

	accesspattern.UpdateSecurity(newPolicy)

	accessPatterns = accesspattern.GetAccessPatterns()
	graph.GenerateGraph(accessPatterns, "improved")
	policy.WriteAuthorizationServerPolicy(newPolicy)
	policy.WriteAccessPatterns(accessPatterns)
	policy.WriteSimpleAccessControlPolicy(accessPatterns)
}
