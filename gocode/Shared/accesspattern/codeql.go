package accesspattern

import (
	"encoding/json"
	"log"
	"os"

	"github.com/wspr-ncsu/5GAC/gocode/Shared/openapi"
)

type (
	codeQLAP struct {
		SelectStmt Selected `json:"#select"`
	}

	codeQLAP_OAI struct {
		SelectStmt Selected_OAI `json:"#select"`
	}

	Selected struct {
		Columns []ColumnsType `json:"columns"`
		Tuples  [][]string    `json:"tuples"`
	}

	Selected_OAI struct {
		Columns []ColumnsType `json:"columns"`
		Tuples  interface{}   `json:"tuples"`
	}

	ColumnsType struct {
		Name string `json:"name"`
		Kind string `json:"kind"`
	}

	// Results struct {
	// 	httpType            string
	// 	controlFlowSource   string
	// 	SourceNF            string
	// 	ServiceCallFunction string
	// 	TargetNF            string
	// 	TargetHandler       string
	// 	TargetPath          string
	// 	TargetRootPath      string
	// }
)

func ParseJSON(file string) []AccessPattern {
	doc, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	data := codeQLAP{}
	err = json.Unmarshal(doc, &data)
	if err != nil {
		log.Printf("Err: %v, file: %v\n", err, file)

		data := codeQLAP_OAI{}
		err = json.Unmarshal(doc, &data)
		if err != nil {
			log.Printf("OAI Err: %v, file: %v\n", err, file)
		}
	}

	accessPatterns = make([]AccessPattern, 0)

	// log.Printf("The data: %v\n", data)

	for _, outer := range data.SelectStmt.Tuples {
		p := AccessPattern{}
		p.HTTPType = outer[0]
		p.Caller = NetworkFunctionCall{
			Nf:           NetworkFunction(outer[2]),
			FunctionName: outer[3],
		}
		p.Callee = NetworkFunctionCall{
			Nf:           NetworkFunction(outer[4]),
			FunctionName: "Unknown",
		}
		p.URL = openapi.Path{
			Path: outer[5],
			Root: outer[6] + "/",
		}

		match := FindScopeFromService.FindStringSubmatch(p.URL.Root)

		if len(match) >= 1 {
			p.RequiredScopes = match[1]
		}

		accessPatterns = append(accessPatterns, p)
	}

	return accessPatterns
}
