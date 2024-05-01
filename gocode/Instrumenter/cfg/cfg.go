// package cfg keeps track of all the configuration for the program.
package cfg

import (
	"flag"
	"log"
)

// cfg keeps track of the main configuration.
var cfg Config

// RegisterCommandLineArgs parses the command line values sent into the program.
func RegisterCommandLineArgs() Config {
	flag.BoolVar(&cfg.Verbosity, "v", false, "Verbose output")
	flag.StringVar(&cfg.Policy, "policy", "", "Access Control Policy File Location")
	flag.StringVar(&cfg.AccessPatterns, "accesspatterns", "", "Access Patterns File Location. Policy is derived from this. Incompatible with Policy option.")
	cfg.ToolName = "5GAC-Instrumenter"
	cfg.Version.init(1, 0, 0)

	flag.Parse()

	if cfg.Verbosity {
		log.Printf("Verbosity: %v\n", cfg.Verbosity)
	}

	return cfg
}

// GetCfg gets the configuration structure.
func GetCfg() Config {
	return cfg
}
