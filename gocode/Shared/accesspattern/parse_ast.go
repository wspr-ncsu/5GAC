package accesspattern

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/wspr-ncsu/5GAC/gocode/Shared/openapi"
)

var (
	fset     *token.FileSet
	nfCalls  NFCalls
	atMethod string
)

// parseGoAst parses the golang file into an AST and returns the result.
func parseGoAst(file string) *ast.File {
	if fset == nil {
		fset = token.NewFileSet()
	}

	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err.Error())
	}

	return f
}

func GetFileSet() *token.FileSet {
	return fset
}

func ResetFileSet() {
	fset = nil
}

func GetNFCalls(nf NetworkFunction) []ServerFunction {
	return nfCalls[nf]
}

// Manages the GoFunction structure and parses the file's AST.
func (f *GoFile) Parse(file string) {
	f.FullPath = file
	f.File = parseGoAst(file)
	f.Inspect()
}

// Inspect inspects the AST and updates some metadata before we really get into processing it.
func (f *GoFile) Inspect() {
	ast.Inspect(f.File, f.parse)
}

// Make global variables to manage our regular expressions.
var (
	importRegexp, findNfGuess             *regexp.Regexp
	FindNfFromScope, FindScopeFromService *regexp.Regexp
)

// CompileRegexp compiles some regular expressions used throughout the program. Typically used to match certain NFs
// and packages together.
func CompileRegexp() {
	if importRegexp == nil {
		importRegexp = regexp.MustCompile("\"github\\.com\\/free5gc\\/(?P<NF>.*)\\/internal/context\"")
	}

	if findNfGuess == nil {
		findNfGuess = regexp.MustCompile(`../_code\/(?P<NF>[^\/|\.]+)[\/|\.].`)
	}

	if FindNfFromScope == nil {
		FindNfFromScope = regexp.MustCompile(`^n(?P<NF>.*?)\-.*$`)
	}

	if FindScopeFromService == nil {
		FindScopeFromService = regexp.MustCompile(`.+\/n(?P<Scope>.*?)\/.*$`)
	}
}

// matchContextImport takes an import definition and checks if it's the NF context import
// (which contains configuration for the NF)
// remember this import for later use.
// e.g: to get NRF IP address for OAuth requests.
func (f *GoFile) matchContextImport(path string, importName *ast.Ident) bool {
	if importName != nil {
		idx := strings.Index(importName.Name, "_")
		if idx != -1 {
			f.NfGuess = NetworkFunction(importName.Name[:idx])
		} else {
			f.NfGuess = NetworkFunction(importName.Name)
		}
	}

	match := importRegexp.FindStringSubmatch(path)

	if len(match) <= 1 { // At least two  must elements match
		log.Fatalf("Import in wrong format? %v", match)
	}

	f.NfGuess = NetworkFunction(strings.ToUpper(match[1]))

	return true
}

// updateNfMetadata figures out which NF this file belongs to, since we process all the NFs at once.
func (f *GoFile) updateNfMetadata(importSpec *ast.ImportSpec) bool {
	path := importSpec.Path.Value
	f.Imports = append(f.Imports, path)

	if len(f.NfGuess) == 0 {
		match := findNfGuess.FindStringSubmatch(f.FullPath)

		if len(match) >= 1 {
			f.NfGuess = NetworkFunction(strings.ToUpper(match[1]))
		}
	}

	if !importRegexp.Match([]byte(path)) {
		return false
	}

	return f.matchContextImport(path, importSpec.Name)
}

func GetAccessTokenMethod() string {
	return atMethod
}

// parseRoutes takes a "routes" variable from golang,
// parse the variables out of it: mapping HTTP paths to functions.
func (f *GoFile) parseRoutes(specs []ast.Spec) {
	packageName := f.File.Name.Name

	for _, spec := range specs {
		switch st := spec.(type) {
		case *ast.ValueSpec:
			for _, value := range st.Values {
				switch vType := value.(type) {
				case *ast.CompositeLit:
					switch routerType := vType.Type.(type) {
					case *ast.Ident:
						// Make sure this is a "Routes" variable, otherwise we don't care...
						if routerType.Name != "Routes" || routerType.Obj == nil || routerType.Obj.Name != "Routes" || routerType.Obj.Kind != ast.Typ {
							continue
						}

						for _, item := range vType.Elts {
							switch iType := item.(type) {
							case *ast.CompositeLit:
								// 3 strings, followed by a function reference in this structure.
								var funcName, path string

								v, ok := iType.Elts[3].(*ast.Ident)
								if !ok {
									log.Fatalf("Error occurred converting %v to a string.", iType.Elts[3])
								}

								funcName = v.Name

								vv, ok := iType.Elts[2].(*ast.BasicLit)
								if !ok {
									log.Fatalf("Error occurred converting %v to a string", iType.Elts[2])
								}

								path = vv.Value[1:]

								var theScheme string

								switch scheme := iType.Elts[1].(type) {
								case *ast.CallExpr:
									v, ok := scheme.Args[0].(*ast.BasicLit)
									if !ok {
										log.Fatalf("Error occurred converting %v to a string", scheme.Args[0])
									}

									theScheme = v.Value[1:]
									theScheme = theScheme[:len(theScheme)-1]

								case *ast.BasicLit:
									theScheme = scheme.Value[1:]
									theScheme = theScheme[:len(theScheme)-1]
								case *ast.SelectorExpr:
									// http.MethodGet, http.MethodPost
									v, ok := scheme.X.(*ast.Ident)
									if !ok {
										log.Fatalf("Error occurred converting %v to a string", scheme.X)
									}

									theScheme = v.Name + "." + scheme.Sel.Name

									switch theScheme {
									case "http.MethodGet":
										theScheme = http.MethodGet
									case "http.MethodPost":
										theScheme = http.MethodPost
									case "http.MethodDelete":
										theScheme = http.MethodDelete
									case "http.MethodOptions":
										theScheme = http.MethodOptions
									case "http.MethodPatch":
										theScheme = http.MethodPatch
									case "http.MethodPut":
										theScheme = http.MethodPut
									default:
										log.Fatalf("Unknown HTTP type: %T %v", theScheme, theScheme)
									}
								}

								// fmt.Printf("%v %v %v %v\n", funcName, path, theScheme, pack)

								if funcName == "Index" {
									continue
								}

								pathCutQuotes := path[:len(path)-1]

								fullPath := openapi.Path{
									Root: f.BaseURL,
									Path: pathCutQuotes,
								}

								if len(f.BaseURL) == 0 {
									// Couldn't parse URL from earlier prog analysis.
									// This happens in cases where the root URL may be determined by a variable. e.g: in the NRF
									// group := engine.Group(factory.NrfNfmResUriPrefix)
									// This factory.NrfNf... determines the base URL. In other services it looks like:
									// group := engine.Group("/nnssf-nsselection/v1")
									// Therefore, in other services we can just parse this string.
									// in the NRF's case we need to track this variable to its definition.
									// Hard code a few relative paths if we know already their base URL.
									baseURLs := map[string]string{
										"/nf-instances":               "{apiRoot}/nnrf-disc/v1",
										"/nf-instances/:nfInstanceID": "{apiRoot}/nnrf-nfm/v1",
									}

									f.BaseURL = baseURLs[pathCutQuotes]
									fullPath.Root = baseURLs[pathCutQuotes]

								}

								if packageName == "accesstoken" && fullPath.Path == "/oauth2/token" {
									atMethod = funcName
								}

								_, _ = openapi.AssignPackageAndPath(funcName, path[:len(path)-1], strings.ToUpper(theScheme), packageName)
								// if err != nil {
								// 	log.Errorf(err)
								// }

								for _, itemVal := range iType.Elts {
									switch funcIdx := itemVal.(type) {
									case *ast.Ident:
										networkFunction := f.NfGuess
										// nfCall := createOrFindNFCall(f.NfGuess, scheme, path)
										if nfCalls == nil {
											nfCalls = make(NFCalls)
										}

										api, err := openapi.FindAPIWithRoute(strings.ToUpper(theScheme), fullPath)

										if err == nil {
											updateAccessPatterns(api, networkFunction, &funcIdx.Name)
										}

										nfCalls[networkFunction] = append(nfCalls[networkFunction], ServerFunction(funcIdx.Name))
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func updateAccessPatterns(api *openapi.Sec, networkFunction NetworkFunction, serverFuncName *string) {
	apAPI := API{
		URL:      api.Path,
		HTTPType: api.Scheme,
	}

	if serverFuncName != nil {
		updateServerAccessPattern(apAPI, networkFunction, *serverFuncName)
	}
}

func updateServerAccessPattern(api API, networkFunction NetworkFunction, serverFuncName string) {
	for i, ap := range accessPatterns {
		if ap.API == api && ap.Callee.Nf == networkFunction {
			accessPatterns[i].Callee.FunctionName = serverFuncName
		}
	}
}

var callerFunc *ast.FuncDecl

// parse parses an AST, based on the Inspect method.
func (f *GoFile) parse(n ast.Node) bool {
	switch node := n.(type) {
	case *ast.ImportSpec: // Find the context for the NF in free5gc.
		// fmt.Printf("Looking at import %v\n", t)
		f.updateNfMetadata(node)

	case *ast.GenDecl:
		// We only care about a global variable gendecl
		if node.Tok != token.VAR {
			return true
		}

		f.parseRoutes(node.Specs)

	case *ast.FuncDecl:
		callerFunc = node

		if node.Name.Name == "AddService" {
			route := node.Body.List[0]
			switch r := route.(type) {
			case *ast.AssignStmt:
				switch rhs := r.Rhs[0].(type) {
				case *ast.CallExpr:
					switch rhsFunc := rhs.Fun.(type) {
					case *ast.SelectorExpr:
						class := ""
						funcName := rhsFunc.Sel.Name

						switch classIdent := rhsFunc.X.(type) {
						case *ast.Ident:
							class = classIdent.Name
						}

						if class == "engine" && funcName == "Group" {
							switch strArg := rhs.Args[0].(type) {
							case *ast.BasicLit:
								if strArg.Kind == token.STRING {
									f.BaseURL = fmt.Sprintf("{apiRoot}%v", strings.TrimSuffix(strArg.Value[1:], "\""))
								}
							}
						}
					}
				}
			}
		}

	case *ast.AssignStmt:
		for _, rhs := range node.Rhs {
			switch r := rhs.(type) {
			case *ast.CallExpr:
				switch rhsFunc := r.Fun.(type) {
				case *ast.SelectorExpr:
					api, err := openapi.FindAPIWithFunctionName(rhsFunc.Sel.Name)
					if err != nil {
						// fmt.Printf("Found %v not in OpenAPI\n", rhs_func.Sel.Name)
						// return true
						funcName := callerFunc.Name.Name

						if rhsFunc.Sel.Name == "PrepareRequest" {
							api, err := openapi.FindAPIWithFunctionName(funcName)
							if err != nil {
								return true
							}

							if len(api.SecurityRequirements) == 0 {
								return true
							}

							for _, secReq := range api.SecurityRequirements {
								for _, scopes := range secReq {
									for _, scope := range scopes {
										addAccessPattern(f.NfGuess, scope, api.Scheme, api.Path, funcName)
									}
								}
							}
						}

						return true
					}

					if !checkCallToOpenAPI(rhsFunc) {
						return true
					}

					for _, securityReq := range api.SecurityRequirements {
						for _, val := range securityReq {
							for _, scope := range val {
								// fmt.Printf("Scope: %v, NF: %v\n", scope, f.NfGuess)
								addAccessPattern(f.NfGuess, scope, api.Scheme, api.Path, rhsFunc.Sel.Name)
							}
						}
					}
				}
			}
		}
	}

	return true
}

var accessPatterns []AccessPattern

func GetAccessPatterns() []AccessPattern {
	return accessPatterns
}

func SetAccessPatterns(ap []AccessPattern) {
	accessPatterns = ap
}

func addAccessPattern(NfGuess NetworkFunction, scope, scheme string, route openapi.Path, callerFunc string) bool {
	if NfGuess == NetworkFunction("OPENAPI") {
		return false
	}

	match := FindNfFromScope.FindStringSubmatch(scope)

	var calleeGuess string

	if len(match) >= 1 {
		calleeGuess = strings.ToUpper(match[1])
	}

	caller := NetworkFunctionCall{
		Nf:           NfGuess,
		FunctionName: callerFunc,
	}

	api := API{
		URL:      route,
		HTTPType: scheme,
	}

	calleeFuncName := "Unknown"
	calleeNfGuess := NetworkFunction(strings.ToUpper(calleeGuess))

	for _, serverfunc := range nfCalls[calleeNfGuess] {
		serverAPI, err := openapi.FindAPIWithFunctionName(string(serverfunc))
		if err != nil {
			continue
		}

		servAPI := API{
			URL:      serverAPI.Path,
			HTTPType: serverAPI.Scheme,
		}

		if servAPI == api {
			calleeFuncName = string(serverfunc)
		}
	}

	// Hardcode this fix.
	// TODO: Fix this by better analysis.
	// Free5GC manually renamed this specific function to be different than the operationId.
	// We can still extract this information by looking at the Routes variable tho. TODO.
	// hardcodedFunctionNames := map[API]string{
	// 	{
	// 		HTTPType: "GET",
	// 		URL: openapi.Path{
	// 			Root: "/nnssf-nsselection/v1",
	// 			Path: "/network-slice-information",
	// 		},
	// 	}: "HTTPNetworkSliceInformationDocument",
	// 	{
	// 		HTTPType: "PUT",
	// 		URL: openapi.Path{
	// 			Root: "/nudm-uecm/v1",
	// 			Path: "/:ueId/registrations/amf-3gpp-access",
	// 		},
	// 	}: "HTTPRegistrationAmf3gppAccess",
	// 	{
	// 		HTTPType: "PUT",
	// 		URL: openapi.Path{
	// 			Root: "/nudm-uecm/v1",
	// 			Path: "/:ueId/registrations/amf-non-3gpp-access",
	// 		},
	// 	}: "HTTPRegistrationAmfNon3gppAccess",
	// }

	// if funcName, ok := hardcodedFunctionNames[api]; ok {
	// 	calleeFuncName = funcName
	// }

	callee := NetworkFunctionCall{
		Nf:           calleeNfGuess,
		FunctionName: calleeFuncName,
	}

	if len(accessPatterns) == 0 {
		accessPatterns = make([]AccessPattern, 0)
	}

	// Check for an already existing access pattern that matches this.
	for _, ap := range accessPatterns {
		if ap.API == api && ap.Callee == callee && ap.Caller == caller && ap.RequiredScopes == scope {
			// log.Printf("Access pattern already exists %v\n", ap)

			return false
		}
	}

	if callee.FunctionName == "Unknown" {
		// Try to guess server function name.
		callee.FunctionName = "HTTP" + caller.FunctionName
	}

	accessPatterns = append(accessPatterns, AccessPattern{
		API:            api,
		Caller:         caller,
		Callee:         callee,
		RequiredScopes: scope,
	})

	return true
}

// checkCallToOpenAPI checks if the SelectorExpr is a call to an openAPI client function.
func checkCallToOpenAPI(rhsFunc *ast.SelectorExpr) bool {
	switch rhsFuncClass := rhsFunc.X.(type) {
	case *ast.SelectorExpr:
		switch ident := rhsFuncClass.X.(type) {
		case *ast.Ident:
			if !strings.Contains(strings.ToLower(ident.Name), "client") {
				log.Printf("Error - didn't make the expected \"client\" variable name.\n")
				log.Printf("This is meant to prevent a match on when functions of the same name match one in the OpenAPI implementation.\n")

				return true // assumption: The variable referencing the openAPI client requesting thing is not called "client". TODO
			}

		case *ast.SelectorExpr:
			break
		default:
			return false
		}
	default:
		return false
	}

	return true
}
