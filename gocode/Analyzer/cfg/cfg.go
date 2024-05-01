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
	flag.StringVar(&cfg.AccessPatterns, "ap", "", "Take in an access pattern file for policy analysis from CodeQL")

	cfg.ToolName = "5GAC-Analyzer"
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
