package instrumentation

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/wspr-ncsu/5GAC/gocode/Instrumenter/cfg"
	"github.com/wspr-ncsu/5GAC/gocode/Shared/accesspattern"
	"github.com/wspr-ncsu/5GAC/gocode/Shared/openapi"
	"golang.org/x/tools/go/ast/astutil"
)

var (
	// Keeps track of the last seen function during parsing of the AST.
	lastSeenFunc struct {
		callerFunc     *ast.FuncDecl
		calleeFuncName string
	}

	// Keeps track of if we need to make changes to the AST or not.
	makeAstChanges bool
)

func (f *GoFile) addImports(imports []string) {
	fset := accesspattern.GetFileSet()
	for _, importPkg := range imports {
		astutil.AddImport(fset, f.File, importPkg)
	}
}

// UpdateAst updates the AST if we're not improving the policy, just do a basic trial run of this and don't modify the AST in a way that will affect our future processing.
func UpdateAst(file accesspattern.GoFile, trialRun bool) {
	igf := GoFile(file)

	// err := copier.Copy(&igf.GoFile, file)
	// if err != nil {
	// 	log.Printf("Failed to copy Go File to Instrumented Go File [%v -> %v]", file, igf.GoFile)
	// }

	makeAstChanges = !trialRun

	if cfg.GetCfg().Verbosity {
		fmt.Printf("Updating AST of NF %v, File: %v\n", file.NfGuess, file.RelativePath)
	}
	astutil.Apply(igf.File, igf.findOpenAPIMethods, igf.findPrepareRequests)

	if !makeAstChanges {
		return
	}

	writeFile(igf.RelativePath, file.File)
}

func writeFile(relPath string, file *ast.File) {
	outputFile := "../_output/" + relPath

	fset := accesspattern.GetFileSet()

	cmap := ast.NewCommentMap(fset, file, file.Comments)

	file.Comments = cmap.Filter(file).Comments()

	err := os.MkdirAll(filepath.Dir(outputFile), 0o777)
	if err != nil {
		log.Fatalf("mkdir all - %v", err)
	}

	output, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := printer.Fprint(output, fset, file); err != nil {
		log.Fatal(err)
	}

	output.Close()
}

func (f *GoFile) addContextImport() string {
	fset := accesspattern.GetFileSet()

	if len(f.NfGuess) > 0 && f.File.Name.Name != "context" && f.NfGuess != "NRF" {
		context := fmt.Sprintf("%v_context", strings.ToLower(string(f.NfGuess)))
		importStmt := fmt.Sprintf("github.com/free5gc/%v/internal/context", strings.ToLower(string(f.NfGuess)))

		astutil.AddNamedImport(fset, f.File, context, importStmt)
		return context

	} else if f.NfGuess == "NRF" && f.File.Name.Name != "accesstoken" {
		astutil.AddImport(fset, f.File, fmt.Sprintf("github.com/free5gc/%v/pkg/factory", strings.ToLower(string(f.NfGuess))))
	}

	return ""
}

// validateGinObject takes the gin object variable and validates the name of it.
func validateGinObject(fields *ast.FieldList) (string, bool) {
	// gin object SHOULD be the first  and only parameter, check that.
	field := fields.List[0]
	ginObjName := fields.List[0].Names[0].Name

	if len(fields.List) != 1 {
		if ginObjName != "ctx" && ginObjName != "c" {
			return "", false
		}

		// Even if the var is "ctx", check its type is gin.Context:
		theType, ok := field.Type.(*ast.SelectorExpr)
		if !ok {
			log.Printf("Error: %v was not selector expr", field)
		}

		anotherType, ok := theType.X.(*ast.Ident)
		if !ok {
			log.Printf("Error: %v was not ast.ident", anotherType)
		}

		if theType.Sel.Name != "Context" || anotherType.Name != "gin" {
			return "", false
		}
	}

	return ginObjName, true
}

// checkScopeExists checks if a source code file has already been instrumented if returns true if it has.
func checkScopeExists(funcBody []ast.Stmt) bool {
	if len(funcBody) == 0 {
		return false
	}

	testScopesVar := funcBody[0]

	oldScope, ok := testScopesVar.(*ast.AssignStmt)
	if !ok {
		return false
	}

	oldScopeIdent, ok := oldScope.Lhs[0].(*ast.Ident)
	if !ok {
		return false
	}

	if oldScopeIdent.Name == "scopes" {
		return true
	}

	return false
}

// checkCallToOpenAPI checks if the SelectorExpr is a call to an openAPI client function.
func checkCallToOpenAPI(rhsFunc *ast.SelectorExpr) bool {
	switch rhsFuncClass := rhsFunc.X.(type) {
	case *ast.SelectorExpr:
		switch ident := rhsFuncClass.X.(type) {
		case *ast.Ident:
			// assumption: The variable referencing the openAPI client requesting thing is not called "client". TODO
			if !strings.Contains(strings.ToLower(ident.Name), "client") {
				log.Printf("Error - didn't make the expected \"client\" variable name.\n")
				log.Fatal("This is meant to prevent a match on when functions of the same name match one in the OpenAPI implementation.\n")
			}
		case *ast.SelectorExpr:
			break
		default:
			return false
		}
	default:
		if cfg.GetCfg().Verbosity {
			log.Printf("outter: %T %v lastseenfunc: %v\n", rhsFuncClass, rhsFuncClass, lastSeenFunc)
		}

		return false
	}

	return true
}

// findOpenAPIMethods does a lot more than what the name of the function would suggest. It parses metadata for server side functions and augments all of the server-side checks.
// and, it does so for the client-side functions. If there is weird inconsistencies in the AST, it may post process and augment the AST at a different point.
func (f *GoFile) findOpenAPIMethods(cursor *astutil.Cursor) bool {
	if cursor == nil {
		return false
	}

	serverFunctions := accesspattern.GetNFCalls(f.NfGuess)

	node := cursor.Node()

	switch t := node.(type) {
	case *ast.FuncDecl:
		lastSeenFunc.callerFunc = t

		// Handles augmenting AUSF, NSSF, SMF.
		// We can skip doing this for NRF. If we need to read NRFs value we can read directly from the configuration structure, not the context one.
		if f.NfGuess != "NRF" && (f.File.Name.Name == "context" || f.File.Name.Name == "util") && strings.EqualFold(fmt.Sprintf("Init%vContext", f.NfGuess), t.Name.Name) {
			if cfg.GetCfg().Verbosity {
				log.Printf("Found context Init%vContext function.\n", f.NfGuess)
			}
			returnStmt := t.Body.List[len(t.Body.List)-1]

			config := getContextConfigCopy(f.NfGuess)

			if stmt, ok := returnStmt.(*ast.ReturnStmt); ok {
				t.Body.List = append(t.Body.List[:len(t.Body.List)-1], config, stmt)
			} else {
				t.Body.List = append(t.Body.List, config)
			}
		}

		// if f.NfGuess == "NRF" && f.File.Name.Name == "accesstoken" && t.Name.Name == accesspattern.GetAccessTokenMethod() {
		// 	fmt.Printf("Found %v\n", t.Name.Name)

		// 	var insertSpot int = -1
		// 	for i, stmt := range t.Body.List {
		// 		if assignment, ok := stmt.(*ast.AssignStmt); ok {
		// 			if len(assignment.Lhs) == 1 {
		// 				if ident, ok := assignment.Lhs[0].(*ast.Ident); ok {
		// 					if ident.Name == "req" {
		// 						insertSpot = i
		// 						break
		// 					}
		// 				}
		// 			}
		// 		}
		// 	}
		// 	var list []ast.Stmt

		// 	list = append(list, &ast.IfStmt{
		// 		Init: &ast.AssignStmt{
		// 			Tok: token.DEFINE,
		// 			Lhs: []ast.Expr{
		// 				ast.NewIdent("err"),
		// 			},
		// 			Rhs: []ast.Expr{
		// 				&ast.CallExpr{
		// 					Fun: &ast.SelectorExpr{
		// 						X: &ast.SelectorExpr{
		// 							X:   ast.NewIdent("c"),
		// 							Sel: ast.NewIdent("Request"),
		// 						},
		// 						Sel: ast.NewIdent("ParseForm"),
		// 					},
		// 				},
		// 			},
		// 		},
		// 		Cond: &ast.BinaryExpr{
		// 			Op: token.NEQ,
		// 			X:  ast.NewIdent("err"),
		// 			Y:  ast.NewIdent("nil"),
		// 		},
		// 		Body: &ast.BlockStmt{
		// 			List: []ast.Stmt{
		// 				&ast.AssignStmt{
		// 					Tok: token.DEFINE,
		// 					Lhs: []ast.Expr{
		// 						ast.NewIdent("problemDetail"),
		// 					},
		// 					Rhs: []ast.Expr{
		// 						&ast.BinaryExpr{
		// 							Op: token.ADD,
		// 							X: &ast.BasicLit{
		// 								Value: "\"[Request Body] \"",
		// 								Kind:  token.STRING,
		// 							},
		// 							Y: &ast.CallExpr{
		// 								Fun: &ast.SelectorExpr{
		// 									X:   ast.NewIdent("err"),
		// 									Sel: ast.NewIdent("Error"),
		// 								},
		// 							},
		// 						},
		// 					},
		// 				},
		// 				&ast.AssignStmt{
		// 					Tok: token.DEFINE,
		// 					Lhs: []ast.Expr{
		// 						ast.NewIdent("rsp"),
		// 					},
		// 					Rhs: []ast.Expr{
		// 						ast.NewIdent("models.ProblemDetails{Status: http.StatusInternalServerError, Cause:	\"SYSTEM_FAILURE\", Detail:	err.Error(),}"),
		// 					},
		// 				},
		// 				&ast.ExprStmt{
		// 					X: &ast.CallExpr{
		// 						Fun: &ast.SelectorExpr{
		// 							X: &ast.SelectorExpr{
		// 								X:   ast.NewIdent("logger"),
		// 								Sel: ast.NewIdent("AccessTokenLog"),
		// 							},
		// 							Sel: ast.NewIdent("Warnln"),
		// 						},
		// 						Args: []ast.Expr{
		// 							ast.NewIdent("problemDetail"),
		// 						},
		// 					},
		// 				},

		// 				&ast.ExprStmt{
		// 					X: &ast.CallExpr{
		// 						Fun: &ast.SelectorExpr{
		// 							X:   ast.NewIdent("c"),
		// 							Sel: ast.NewIdent("JSON"),
		// 						},
		// 						Args: []ast.Expr{
		// 							&ast.SelectorExpr{
		// 								X:   ast.NewIdent("http"),
		// 								Sel: ast.NewIdent("StatusBadRequest"),
		// 							},
		// 							&ast.Ident{Name: "rsp"},
		// 						},
		// 					},
		// 				},
		// 				&ast.ReturnStmt{},
		// 			},
		// 		},
		// 	})
		// 	// list = append(list)

		// 	list = append(list, &ast.AssignStmt{
		// 		Lhs: []ast.Expr{
		// 			ast.NewIdent("accessTokenReq.NfInstanceId"),
		// 		},
		// 		Rhs: []ast.Expr{
		// 			ast.NewIdent("c.Request.FormValue(\"nfInstanceId\")"),
		// 		},
		// 		Tok: token.ASSIGN,
		// 	})
		// 	list = append(list, &ast.AssignStmt{
		// 		Lhs: []ast.Expr{
		// 			ast.NewIdent("accessTokenReq.GrantType"),
		// 		},
		// 		Rhs: []ast.Expr{
		// 			ast.NewIdent("c.Request.FormValue(\"grant_type\")"),
		// 		},
		// 		Tok: token.ASSIGN,
		// 	})
		// 	list = append(list, &ast.AssignStmt{
		// 		Lhs: []ast.Expr{
		// 			ast.NewIdent("accessTokenReq.Scope"),
		// 		},
		// 		Rhs: []ast.Expr{
		// 			ast.NewIdent("c.Request.FormValue(\"scope\")"),
		// 		},
		// 		Tok: token.ASSIGN,
		// 	})
		// 	list = append(list, &ast.AssignStmt{
		// 		Lhs: []ast.Expr{
		// 			ast.NewIdent("accessTokenReq.NfType"),
		// 		},
		// 		Rhs: []ast.Expr{
		// 			ast.NewIdent("models.NfType(c.Request.FormValue(\"nfType\"))"),
		// 		},
		// 		Tok: token.ASSIGN,
		// 	})
		// 	list = append(list, &ast.AssignStmt{
		// 		Lhs: []ast.Expr{
		// 			ast.NewIdent("accessTokenReq.TargetNfType"),
		// 		},
		// 		Rhs: []ast.Expr{
		// 			ast.NewIdent("models.NfType(c.Request.FormValue(\"targetNfType\"))"),
		// 		},
		// 		Tok: token.ASSIGN,
		// 	})
		// 	list = append(list, &ast.AssignStmt{
		// 		Lhs: []ast.Expr{
		// 			ast.NewIdent("accessTokenReq.TargetNfInstanceId"),
		// 		},
		// 		Rhs: []ast.Expr{
		// 			ast.NewIdent("c.Request.FormValue(\"targetNfInstanceId\")"),
		// 		},
		// 		Tok: token.ASSIGN,
		// 	})
		// 	if insertSpot != -1 {
		// 		t.Body.List = append(t.Body.List[:insertSpot], append(
		// 			list, t.Body.List[insertSpot:]...)...)
		// 	}
		// }

		for _, funcName := range serverFunctions {
			strFuncName := string(funcName)
			if strFuncName != t.Name.Name {
				continue
			}
			// Check params to make sure this is actually a gin function.

			ginObjName, valid := validateGinObject(t.Type.Params)

			if !valid {
				continue
			}

			if strFuncName == "HTTPRegisterNFInstance" {
				if cfg.GetCfg().Verbosity {
					log.Printf("Skipping \"%v\" - This API should not enfroce OAuth.\n", strFuncName)
				}
				continue
			}

			// Make sure we're not inserting a scope check when one already exists:

			bodyExist := checkScopeExists(t.Body.List)

			if bodyExist {
				// fmt.Printf("Already instrumented %v\n", funcName)
				continue
			}

			secReq, err := openapi.FindAPIWithFunctionName(strFuncName)
			if err != nil {
				if cfg.GetCfg().Verbosity {
					log.Printf("Did not find server function %v even though it was in ServerFunctions slice.\n", strFuncName)
				}

				continue
			}

			if cfg.GetCfg().Verbosity {
				log.Printf("----Found server function %v, scopes: %v\n", strFuncName, secReq)
			}

			if !makeAstChanges {
				continue
			}

			ctxImport := f.addContextImport()
			pkgName := f.File.Name.Name

			if args := AugmentServerScopeChecks(f.NfGuess, pkgName, ctxImport, strFuncName, ginObjName, secReq.SecurityRequirements); args == nil {
				continue
			} else {
				f.addContextImport()
				f.addImports([]string{"github.com/free5gc/openapi", "net/http"}) //"strconv"})

				// Put the arguments at the beginning of the function call. This authorization needs to happen before anything else.
				t.Body.List = append(args, t.Body.List...)
			}
		}

		return true

	case *ast.AssignStmt:
		for _, rhs := range t.Rhs {
			r, ok := rhs.(*ast.CallExpr)
			if !ok {
				if cfg.GetCfg().Verbosity {
					log.Printf("Error: %v was not callexpr", rhs)
				}

				continue
			}

			rhsFunc, ok := r.Fun.(*ast.SelectorExpr)
			if !ok {
				if cfg.GetCfg().Verbosity {
					log.Printf("Error: %v was not *ast.SelectorExpr", r.Fun)
				}

				continue
			}

			lastSeenFunc.calleeFuncName = rhsFunc.Sel.Name

			_, err := openapi.FindAPIWithFunctionName(rhsFunc.Sel.Name)
			if err != nil {
				return true
			}

			if !checkCallToOpenAPI(rhsFunc) {
				return true
			}

			if makeAstChanges {
				// Change the first argument to our new variable "c" in the openapi call
				// client.DefaultApi.PolicyDataBdtDataGet(openapi.CreateContext(pcf_context.PCF_Self().EnforceOAuth, pcf_context.PCF_Self().NfId, pcf_context.PCF_Self().NrfUri))

				ctxImport := fmt.Sprintf("%v_context", strings.ToLower(string(f.NfGuess)))
				selfCall := GetSelfCall(f.NfGuess, f.File.Name.Name, ctxImport).Value
				oauthVar := GetOAuthVariable(f.NfGuess, f.File.Name.Name, ctxImport)

				r.Args[0] = ast.NewIdent(fmt.Sprintf("openapi.CreateContext(%v, %v().%v, %v().NrfUri, \"%v\")", oauthVar, selfCall, GetNfIDVar(f.NfGuess), selfCall, f.NfGuess))
				f.addContextImport()
				f.addImports([]string{"github.com/free5gc/openapi", "github.com/free5gc/openapi/models"})

				found := false

				ast.Inspect(f.File, func(n ast.Node) bool {
					if found {
						return false
					}

					switch x := n.(type) {
					case *ast.SelectorExpr:
						if ident, ok := x.X.(*ast.Ident); ok {
							if ident.Name == "context" && (x.Sel.Name == "Background" || x.Sel.Name == "TODO" || x.Sel.Name == "WithValue") {
								found = true
							}
						}
					}

					return true
				})

				if !found {
					if f.File.Name.Name != "association" {
						astutil.DeleteImport(accesspattern.GetFileSet(), f.File, "context")
					}
				}
			}
		}

	case *ast.GenDecl:
		if (f.File.Name.Name == "context" || f.File.Name.Name == "factory") && t.Tok == token.TYPE {
			for _, spec := range t.Specs {

				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || (typeSpec.Name.Name != fmt.Sprintf("%vContext", string(f.NfGuess)) && typeSpec.Name.Name != "Configuration") {
					if cfg.GetCfg().Verbosity {
						log.Printf("Skipping %v - not Configuration struct", spec)
					}

					continue
				}

				typeStruct, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					if cfg.GetCfg().Verbosity {
						log.Printf("Found configuration struct, but its not a StructType. %v", typeSpec)
					}

					continue
				}

				typeStruct.Fields.List = append(typeStruct.Fields.List, GetEnforceOAuthConfigVar())
			}
		}
		if f.File.Name.Name == "openapi" && t.Tok == token.CONST {
			for _, spec := range t.Specs {
				if valSpec, ok := spec.(*ast.ValueSpec); ok {
					if len(valSpec.Names) == 1 && valSpec.Names[0].Name == "ContextAccessToken" {
						t.Specs = append(t.Specs, &ast.ValueSpec{
							Names: []*ast.Ident{
								ast.NewIdent("ContextOAuthAdditionalParams"),
							},
							Values: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("ContextKey"),
									Args: []ast.Expr{
										&ast.BasicLit{
											Value: "\"additional_params\"",
										},
									},
								},
							},
						})
					}
				}
			}
		}
	}
	return true
}

// findPrepareRequests finds the calls to "PrepareRequest" - which is in openapi, and right before were we need to insert OAuth stuff.
func (f *GoFile) findPrepareRequests(cursor *astutil.Cursor) bool {
	if cursor == nil {
		return false
	}

	node := cursor.Node()

	assignmentNode, ok := node.(*ast.AssignStmt)
	if !ok {
		return true
	}

	for _, callNode := range assignmentNode.Rhs {
		call, err := callNode.(*ast.CallExpr)
		if !err {
			return true
		}

		selector, err := call.Fun.(*ast.SelectorExpr)
		if !err {
			return true
		}

		funcName := lastSeenFunc.callerFunc.Name.Name

		if selector.Sel.Name == "PrepareRequest" {

			if funcName == "RegisterNFInstance" {
				if cfg.GetCfg().Verbosity {
					log.Printf("Not inserting client template to \"%v\" - this API should not enforce OAuth.", funcName)
				}
				return true
			}

			api, err := openapi.FindAPIWithFunctionName(funcName)
			if err != nil {
				return true
			}

			if len(api.SecurityRequirements) == 0 {
				if cfg.GetCfg().Verbosity {
					log.Printf("%v doesn't have any security requirements... %v - %v\n\n", funcName, api, f)
				}
				// don't process this function anymore, since there's no scopes we don't need to request and OAuth token
				return true
			}

			var scopesNeeded []string

			// TODO: Allow multiple scope definitions
			for _, secReq := range api.SecurityRequirements {
				for _, scopes := range secReq {
					scopesNeeded = scopes
				}
			}

			if !makeAstChanges {
				return true
			}

			if len(scopesNeeded) > 0 {
				match := accesspattern.FindNfFromScope.FindStringSubmatch(scopesNeeded[0])

				var calleeGuess string

				if len(match) >= 1 {
					calleeGuess = strings.ToUpper(match[1])
				}
				template := GetClientTemplate(lastSeenFunc.callerFunc, scopesNeeded, calleeGuess)
				// outerLoop:
				for _, node := range template {
					cursor.InsertBefore(node)
				}
				f.addImports([]string{"fmt", "net/url", "crypto/tls", "strconv", "golang.org/x/net/http2", "golang.org/x/oauth2", "golang.org/x/oauth2/clientcredentials"})
			}
		}
	}

	return true
}
