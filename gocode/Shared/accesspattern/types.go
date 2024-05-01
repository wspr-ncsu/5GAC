package accesspattern

import (
	"go/ast"

	"github.com/wspr-ncsu/5GAC/gocode/Shared/openapi"
)

type (
	// fileInfo Keeps track of file paths.
	fileInfo struct {
		FullPath     string
		name         string
		RelativePath string
	}

	API struct {
		HTTPType string       `yaml:"HTTPType"`
		URL      openapi.Path `yaml:"URL"`
	}

	NetworkFunction string

	NetworkFunctionCall struct {
		// API
		// ApiImplementations
		Nf           NetworkFunction `yaml:"NetworkFunction"`
		FunctionName string          `yaml:"FuncName"`
	}

	ServerFunction string

	NFCalls map[NetworkFunction][]ServerFunction

	AccessPattern struct {
		API            `yaml:"API"`
		Caller         NetworkFunctionCall `yaml:"Caller"`
		Callee         NetworkFunctionCall `yaml:"Callee"`
		RequiredScopes string              `yaml:"Scope"`
	}

	GoFile struct {
		fileInfo
		File *ast.File
		// The NF this file's source code is in.
		NfGuess NetworkFunction
		// What Go imports are used in this file
		Imports []string
		BaseURL string
	}

	// Store all data and spec file parsing in RAM.
	FileTracker struct {
		Specs map[fileInfo]openapi.APIs
		Code  map[fileInfo]GoFile
	}
)
