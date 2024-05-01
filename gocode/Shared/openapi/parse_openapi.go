// Package openapi handles parsing OpenAPI files, just the specs of the 3GPP in this project.
package openapi

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v2"
)

type (
	// Copy of OpenAPI3.T structure.
	OpenAPI openapi3.T

	// Sec matches security requirements to a HTTP path and function/package name.
	Sec struct {
		Scheme               string
		Path                 Path
		FunctionName         string
		PackageName          string
		SecurityRequirements openapi3.SecurityRequirements
	}

	// Path Path object that separates the root and actual path into different variables.
	Path struct {
		Root string
		Path string
	}

	APIs map[string]*Sec
)

// TODO: make this a map? Will it be faster?
var APIList []*Sec

// toCamel converts the input string into camel case.
func toCamel(in string) string {
	runes := []rune(in)
	out := make([]rune, 0, len(in))

	for i, r := range runes {
		if r == '_' {
			continue
		}

		if i == 0 || runes[i-1] == '_' {
			out = append(out, unicode.ToUpper(r))

			continue
		}

		out = append(out, r)
	}

	return string(out)
}

// removeCharacters removes characters that should not be put into a function name based on openapi-generator.
func removeCharacters(operationID string) string {
	tmp := strings.ReplaceAll(operationID, "{", ":")
	tmp = strings.ReplaceAll(tmp, "}", "")
	tmp = strings.ReplaceAll(tmp, "-", "/")
	tmp = strings.ReplaceAll(tmp, " ", "")
	tmp = strings.ReplaceAll(tmp, "_", "")

	return tmp
}

// getOrGenerateOperationId - Adapated from the (OpenAPI-Generator project)[https://github.com/OpenAPITools/openapi-generator/blob/37905e8bfac75aab0902f0ca0d50de2e070ab544/modules/openapi-generator/src/main/java/org/openapitools/codegen/DefaultCodegen.java#L4970].
func getOrGenerateOperationID(op *openapi3.Operation, path string, httpOp string) string {
	operationID := op.OperationID
	if len(op.OperationID) != 0 {
		operationID = removeCharacters(operationID)

		return strings.ToLower(operationID)
	}

	tmpPath := removeCharacters(path)

	parts := strings.Split(tmpPath+"/"+strings.ToLower(httpOp), "/")

	var builder []string

	if tmpPath == "/" {
		builder = append(builder, "root")
	}

	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		joinedString := strings.Join(builder, "")
		if len(joinedString) == 0 {
			builder = append(builder, string(unicode.ToUpper(rune(part[0])))+part[1:])
		} else {
			builder = append(builder, toCamel(part))
		}
	}

	joinedString := strings.Join(builder, "")

	return strings.ToLower(regexp.QuoteMeta(joinedString))
}

// AssignPackageAndPath maps a function to a package, incase there's many different functions with the same name.
func AssignPackageAndPath(funcName, path, scheme, packageName string) (*Sec, error) {
	for _, securityPolicy := range APIList {
		if (strings.EqualFold(funcName, securityPolicy.FunctionName) ||
			strings.EqualFold(funcName, "http"+securityPolicy.FunctionName)) && strings.EqualFold(path, securityPolicy.Path.Path) {
			if scheme == securityPolicy.Scheme {
				securityPolicy.PackageName = packageName

				return securityPolicy, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find a function in OpenAPI Spec")
}

// FindAPIWithFunctionName gets the security requirements based on a function name - in operationId format.
func FindAPIWithFunctionName(name string) (*Sec, error) {
	list := findAPIsWithFunctionName(name)

	if len(list) == 0 {
		return nil, fmt.Errorf("could not find an API with \"%v\" name", name)
	}

	return list[0], nil
}

func findAPIsWithFunctionName(name string) []*Sec {
	// name = strings.ToLower(name)

	var funcNameAPIList []*Sec

	for _, v := range APIList {
		if strings.EqualFold(v.FunctionName, name) || strings.EqualFold("http"+v.FunctionName, name) {
			funcNameAPIList = append(funcNameAPIList, v)
		}
	}

	return funcNameAPIList
}

func FindAPIWithRoute(scheme string, path Path) (*Sec, error) {
	for _, v := range APIList {
		if path == v.Path && strings.EqualFold(scheme, v.Scheme) {
			return v, nil
		}
	}

	return nil, fmt.Errorf("could not find an API with matching route \"%v%v\"", path.Root, path.Path)
}

// ReadSpec parses a file's OpenAPI.
func ReadSpec(file string) APIs {
	doc, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	data := OpenAPI{}
	err = yaml.Unmarshal(doc, &data)

	if err != nil {
		log.Fatalf("Err: %v, file: %v\n", err, file)
	}

	fileAPIs := make(APIs)
	globalSecurity := data.Security

	// See which APIs overwrite global security:
	for i, v := range data.Paths {
		for httpOp, details := range v.Operations() {
			var path Path
			if len(data.Servers) != 0 {
				path = Path{Path: i, Root: data.Servers[0].URL}
			} else {
				path = Path{Path: i}
			}

			tmp := strings.ReplaceAll(path.Path, "{", ":")
			tmp = strings.ReplaceAll(tmp, "}", "")

			path.Path = tmp

			api := Sec{Path: path}
			api.Scheme = httpOp
			api.FunctionName = strings.ReplaceAll(getOrGenerateOperationID(details, path.Path, httpOp), ":", "")

			if details.Security != nil {
				api.SecurityRequirements = *details.Security
			} else {
				globSec := make(openapi3.SecurityRequirements, len(globalSecurity))
				copy(globSec, globalSecurity)
				api.SecurityRequirements = append(api.SecurityRequirements, globSec...)
			}

			fileAPIs[api.FunctionName] = &api
			APIList = append(APIList, &api)

			file, err := os.OpenFile("3gpp.csv", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0o660)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			csvwriter := csv.NewWriter(file)

			err = csvwriter.Write([]string{api.FunctionName, api.Scheme, api.Path.Root + " " + api.Path.Path}) // calls Flush internally
			if err != nil {
				log.Fatal(err)
			}

			csvwriter.Flush()

			file.Close()

		}
	}

	return fileAPIs
}
