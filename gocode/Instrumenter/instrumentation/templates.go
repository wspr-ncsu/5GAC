package instrumentation

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/wspr-ncsu/5GAC/gocode/Instrumenter/cfg"
	"github.com/wspr-ncsu/5GAC/gocode/Shared/accesspattern"
)

func MoveTemplates() {
	templates, err := os.Open("../_templates")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	files, err := templates.Readdir(0)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		cmd := exec.Command("cp", "--recursive", "-f", fmt.Sprintf("../_templates/%v", f.Name()), "../_output")
		err := cmd.Run()
		if err != nil {
			log.Fatal("Failed to move files from template folder into output")
		}
	}
}

// GetAdditionalParamsCreation returns "additional_params := url.Values{}" as an AST.
func GetAdditionalParamsCreation(token token.Token) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{
			Name: "additional_params",
			Obj: &ast.Object{
				Kind: ast.Var,
				Name: "additional_params",
			},
		}},
		Rhs: []ast.Expr{&ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "url",
				},
				Sel: &ast.Ident{
					Name: "Values",
				},
			},
		}},
		Tok: token,
	}
}

// AddAdditionalParamsToContext returns "c := context.WithValue(context.Background(), openapi.ContextOAuthAdditionalParams, additional_params)" as an AST.
func AddAdditionalParamsToContext(token token.Token) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{
			Name: "c",
			Obj: &ast.Object{
				Kind: ast.Var,
				Name: "c",
			},
		}},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "context",
					},
					Sel: &ast.Ident{
						Name: "WithValue",
					},
				},
				Args: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "Background",
							},
						},
						Args: []ast.Expr{},
					},
					&ast.Ident{
						Name: "openapi.ContextOAuthAdditionalParams",
					},
					&ast.Ident{
						Name: "additional_params",
					},
				},
			},
		},
		Tok: token,
	}
}

// getNfInstanceIDVar returns "%s_context.%s_Self().NfId", where both %s's are the NF, as an AST.
func getNfInstanceIDVar(nfGuess accesspattern.NetworkFunction, pkgName, ctxImport string) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X: &ast.CallExpr{
			Fun: GetSelfCall(nfGuess, pkgName, ctxImport),
		},
		Sel: &ast.Ident{
			Name: GetNfIDVar(nfGuess),
		},
	}
}

// getNrfUri returns "%s_context.%s_Self().NrfUri", where both %s's are the NF, as an AST.
func getNrfUri(nfGuess accesspattern.NetworkFunction, pkgName, ctxImport string) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X: &ast.CallExpr{
			Fun: GetSelfCall(nfGuess, pkgName, ctxImport),
		},
		Sel: &ast.Ident{
			Name: "NrfUri",
		},
	}
}

func getContextConfigCopy(nfGuess accesspattern.NetworkFunction) *ast.AssignStmt {
	var configVar string
	var contextVar string

	if nfGuess == "N3IWF" || nfGuess == "NSSF" {
		contextVar = fmt.Sprintf("%vContext", strings.ToLower(string(nfGuess)))
		title := string(nfGuess)[:1]
		nfNameMinusTitle := string(nfGuess)[1:]

		configVar = fmt.Sprintf("factory.%v%vConfig.Configuration", strings.ToUpper(title), strings.ToLower(nfNameMinusTitle))
	} else if nfGuess == "SMF" || nfGuess == "UDM" {
		contextVar = fmt.Sprintf("%vContext", strings.ToLower(string(nfGuess)))
		configVar = "configuration"
	} else {
		configVar = "configuration"
		contextVar = "context"
	}

	config := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent(contextVar),
				Sel: ast.NewIdent("OAuth"),
			},
		},
		Rhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent(configVar),
				Sel: ast.NewIdent("OAuth"),
			},
		},
		Tok: token.ASSIGN,
	}

	return config
}

func getEnforceOAuthVar(nfGuess accesspattern.NetworkFunction, pkgName, ctxImport string) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "strconv",
			},
			Sel: &ast.Ident{
				Name: "FormatBool",
			},
		},
		Args: []ast.Expr{
			&ast.SelectorExpr{
				X: &ast.CallExpr{
					Fun: GetSelfCall(nfGuess, pkgName, ctxImport),
				},
				Sel: &ast.Ident{
					Name: "OAuth",
				},
			},
		},
	}
}

// AddAdditionalParamsValueNfInstanceID returns "additional_params.Add("nfInstanceId", %s_context.%s_Self().NfId)" where both %s's are the NF, as an AST.
func AddAdditionalParamsValueNfInstanceID(nfGuess accesspattern.NetworkFunction, pkgName, ctxImport string) *ast.ExprStmt {
	return addAdditionalParamsValueImpl("nfInstanceId", getNfInstanceIDVar(nfGuess, pkgName, ctxImport))
}

// AddAdditionalParamsValueNrfURI returns "additional_params.Add("NrfUri", %s_context.%s_Self().NrfUri)" where both %s's are the NF, as an AST.
func AddAdditionalParamsValueNrfURI(nfGuess accesspattern.NetworkFunction, pkgName, ctxImport string) *ast.ExprStmt {
	return addAdditionalParamsValueImpl("NrfUri", getNrfUri(nfGuess, pkgName, ctxImport))
}

func AddAdditionalParamsEnforceOAuth(nfGuess accesspattern.NetworkFunction, pkgName, ctxImport string) *ast.ExprStmt {
	return addAdditionalParamsValueImpl("OAuth", getEnforceOAuthVar(nfGuess, pkgName, ctxImport))
}

func addAdditionalParamsValueImpl(key string, value ast.Expr) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "additional_params",
				},
				Sel: &ast.Ident{
					Name: "Add",
				},
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%v\"", key),
				},
				value,
			},
		},
	}
}

func GetOAuthVariable(nf accesspattern.NetworkFunction, pkgName, ctxImport string) string {
	if nf == "NRF" || nf == "N3IWF" {
		title := string(nf)[:1]
		nfNameMinusTitle := string(nf)[1:]

		return fmt.Sprintf("factory.%v%vConfig.Configuration.OAuth", strings.ToUpper(title), strings.ToLower(nfNameMinusTitle))
	} else {
		return fmt.Sprintf("%v().OAuth", GetSelfCall(nf, pkgName, ctxImport).Value)
	}
}

// AugmentServerScopeChecks returns
// "	scopes := []string{"%s"} - where %s is a string array
//
//	_, oauth_err := openapi.CheckOAuth(c.Request.Header.Get("Authorization"), scopes)
//	if oauth_err != nil {
//		c.JSON(http.StatusUnauthorized, gin.H{})
//		return
//	}"
func AugmentServerScopeChecks(nf accesspattern.NetworkFunction, pkgName, ctxImport, funcName, ginName string, secRequirements openapi3.SecurityRequirements) []ast.Stmt {
	scopesVar := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{
			&ast.Ident{
				Name: "scopes",
			},
		},
	}
	scopes := &ast.CompositeLit{
		Type: &ast.ArrayType{
			Elt: &ast.Ident{
				Name: "string",
			},
		},
		Elts: []ast.Expr{},
	}

	for _, secReq := range secRequirements {
		for _, req := range secReq {
			for _, scope := range req {
				if cfg.GetCfg().Verbosity {
					log.Printf("Req: %T %v\n", scope, scope)
				}

				scopes.Elts = append(scopes.Elts, &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%v\"", scope),
				})
			}
		}
	}

	if len(scopes.Elts) == 0 {
		if cfg.GetCfg().Verbosity {
			log.Printf("Server function \"%v\" does not need to check for scopes since there is none.\n", funcName)
		}

		return []ast.Stmt{}
	}

	// Add the call to oauth procedure
	oauthCall := &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: "oauth_err",
			},
		},
		Rhs: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("openapi.VerifyOAuth(%v.Request.Header.Get(\"Authorization\"), scopes, \"\")", ginName),
			},
		},
	}

	oauthVar := GetOAuthVariable(nf, pkgName, ctxImport)

	ifStmt := &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			Op: token.LAND,
			X: &ast.BinaryExpr{
				Op: token.NEQ,
				X: &ast.Ident{
					Name: "oauth_err",
				},
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Y: &ast.BinaryExpr{
				Op: token.EQL,
				X: &ast.Ident{
					Name: oauthVar,
				},
				Y: &ast.Ident{
					Name: "true",
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: fmt.Sprintf("%v", ginName),
							},
							Sel: &ast.Ident{
								Name: "JSON",
							},
						},
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "http",
								},
								Sel: &ast.Ident{
									Name: "StatusUnauthorized",
								},
							},
							&ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "gin",
									},
									Sel: &ast.Ident{
										Name: "H",
									},
								},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"error\"",
										},
										Value: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("oauth_err"),
												Sel: ast.NewIdent("Error"),
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{},
			},
		},
	}

	scopesVar.Rhs = append(scopesVar.Rhs, scopes)

	return []ast.Stmt{scopesVar, oauthCall, ifStmt}
}

// getSelfCall returns the function name of the functio nthat returns the configuration of a NF.
// For some reason, in free5gc, in the AUSf this is called "GetSelf" and in other NFs,
// this is just called "%s_Self" where %s is the name of the NF.
func GetSelfCall(nfGuess accesspattern.NetworkFunction, pkgName, ctxImport string) *ast.BasicLit {
	selfCallString := fmt.Sprintf("%v_Self", nfGuess)

	if pkgName == "context" {
		// The smf_context (or another NF context) is calling this function. We don't need to add the context class to the parent of this call
		ctxImport = ""
	}

	// TODO: Is this even necessary anymore or are all NFs using GetSelf?
	if nfGuess == "AUSF" || nfGuess == "AMF" || nfGuess == "NSSF" || nfGuess == "PCF" || nfGuess == "SMF" || nfGuess == "UDR" {
		selfCallString = "GetSelf"
	}

	// Yep, they really just don't capitalize this one...
	if nfGuess == "UDM" {
		selfCallString = "Getself"
	}

	var getSelfCallAST *ast.BasicLit

	if len(ctxImport) != 0 {
		getSelfCallAST = &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("%v.%v", ctxImport, selfCallString),
		}
	} else {
		getSelfCallAST = &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("%v", selfCallString),
		}
	}

	//"github.com/free5gc/nrf/pkg/factory"

	return getSelfCallAST
}

// getNfIDVar  returns the name of the NfInstanceID variable. For some reason, in free5gc it has a different name in the SMF than other NFs.
func GetNfIDVar(nfGuess accesspattern.NetworkFunction) string {
	if nfGuess == "SMF" {
		return "NfInstanceID"
	}

	return "NfId"
}

func checkErrors(returnParameters int, err string) *ast.IfStmt {
	// Check that caller function actually returns 3 values. If it doesn't, only return two values.
	var result *ast.IfStmt
	if returnParameters != 3 {
		result = &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				Op: token.NEQ,
				X: &ast.Ident{
					Name: "err",
				},
				Y: ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: "nil",
							},
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "fmt",
									},
									Sel: &ast.Ident{
										Name: "Errorf",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: err,
									},
								},
							},
						},
					},
				},
			},
		}
	} else {
		// 	if !ok {
		//		return localVarReturnValue, nil, fmt.Errorf("OAuth parameters are invalid")
		//	}

		result = &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				Op: token.NEQ,
				X: &ast.Ident{
					Name: "err",
				},
				Y: ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: "localVarReturnValue",
							},
							&ast.Ident{
								Name: "nil",
							},
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "fmt",
									},
									Sel: &ast.Ident{
										Name: "Errorf",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: err,
									},
								},
							},
						},
					},
				},
			},
		}
	}

	return result
}

// GetClientTemplate returns the client template insert at all API endpoints for requesting OAuth.
// Example:
/*
	scopes := []string{"namf-comm",}
	additional_params, ok := ctx.Value(openapi.ContextOAuthAdditionalParams).(url.Values)
	if !ok {
		return localVarReturnValue, nil, fmt.Errorf("OAuth parameters are invalid")
	}
	tokenUrl := fmt.Sprintf("%v/oauth2/token", additional_params["NrfUri"][0])
	additional_params.Del("NrfUri")
	clientID := "YOUR_CLIENT_ID"
	ClientSecret := "YOUR_CLIENT_SECRET"
	conf := &clientcredentials.Config{ClientID: clientID, ClientSecret: ClientSecret, Scopes: scopes, TokenURL: tokenUrl, AuthStyle: oauth2.AuthStyleInParams, EndpointParams: additional_params}
	http_client := &http.Client{Transport: &http2.Transport{AllowHTTP: true, DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
		return net.Dial(network, addr)
	}}}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, http_client)
	token, err := conf.Token(ctx)
	if err != nil {
		return localVarReturnValue, nil, err
	}
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token.AccessToken)
*/
func GetClientTemplate(callerFunc *ast.FuncDecl, scopes []string, callee string) []ast.Stmt {
	list := []ast.Stmt{}

	var arrayStr string
	for _, str := range scopes {
		if len(arrayStr) != 0 {
			arrayStr = fmt.Sprintf("%v \"%v\",", arrayStr, str)
		} else {
			arrayStr = fmt.Sprintf("\"%v\",", str)
		}
	}

	list = append(list,
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "scopes",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.ArrayType{
						Elt: &ast.Ident{
							Name: "string",
						},
					},
					Elts: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: arrayStr,
						},
					},
				},
			},
		},

		// additional_params, ok := ctx.Value(openapi.ContextOAuthAdditionalParams).(url.Values)
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "additional_params",
				},
				&ast.Ident{
					Name: "ok",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.TypeAssertExpr{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "ctx",
							},
							Sel: &ast.Ident{
								Name: "Value",
							},
						},
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "openapi",
								},
								Sel: &ast.Ident{
									Name: "ContextOAuthAdditionalParams",
								},
							},
						},
					},
					Type: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "url",
						},
						Sel: &ast.Ident{
							Name: "Values",
						},
					},
				},
			},
		},
	)

	if len(callerFunc.Type.Results.List) != 3 {
		list = append(list, &ast.IfStmt{
			Cond: &ast.UnaryExpr{
				Op: token.NOT,
				X: &ast.Ident{
					Name: "ok",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: "nil",
							},
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "fmt",
									},
									Sel: &ast.Ident{
										Name: "Errorf",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"OAuth parameters are invalid\"",
									},
								},
							},
						},
					},
				},
			},
		})
	} else {
		// 	if !ok {
		//		return localVarReturnValue, nil, fmt.Errorf("OAuth parameters are invalid")
		//	}

		list = append(list, &ast.IfStmt{
			Cond: &ast.UnaryExpr{
				Op: token.NOT,
				X: &ast.Ident{
					Name: "ok",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: "localVarReturnValue",
							},
							&ast.Ident{
								Name: "nil",
							},
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "fmt",
									},
									Sel: &ast.Ident{
										Name: "Errorf",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"OAuth parameters are invalid\"",
									},
								},
							},
						},
					},
				},
			},
		})
	}

	// list = append(list, checkErrors(len(callerFunc.Type.Results.List), "\"OAuth parameters are invalid\""))

	list = append(list,
		// oauth, ok := strconv.ParseBool(additional_params["OAuth"])
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "oauth",
				},
				&ast.Ident{
					Name: "err",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("strconv"),
						Sel: ast.NewIdent("ParseBool"),
					},
					Args: []ast.Expr{
						&ast.IndexExpr{
							X: &ast.IndexExpr{
								X: ast.NewIdent("additional_params"),
								Index: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"OAuth\"",
								},
							},
							Index: &ast.BasicLit{
								Kind:  token.INT,
								Value: "0",
							},
						},
					},
				},
			},
		},
	)

	// Check for errors again
	list = append(list, checkErrors(len(callerFunc.Type.Results.List), "err.Error()"))

	var oauthChecks []ast.Stmt

	oauthChecks = append(oauthChecks,

		// tokenUrl := fmt.Sprintf("%v/oauth2/token", additional_params["NrfUri"][0])
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "tokenUrl",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: "fmt.Sprintf(\"%v/oauth2/token\", additional_params[\"NrfUri\"][0])",
				},
			},
		},

		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "additional_params",
					},
					Sel: &ast.Ident{
						Name: "Del",
					},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"NrfUri\"",
					},
				},
			},
		},

		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "additional_params",
					},
					Sel: &ast.Ident{
						Name: "Del",
					},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"OAuth\"",
					},
				},
			},
		},

		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "additional_params",
					},
					Sel: &ast.Ident{
						Name: "Add",
					},
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"targetNfType\"",
					},
					ast.NewIdent(fmt.Sprintf("\"%v\"", callee)),
				},
			},
		},

		// clientID := "YOUR_CLIENT_ID"
		// &ast.AssignStmt{
		// 	Lhs: []ast.Expr{
		// 		&ast.Ident{
		// 			Name: "clientID",
		// 		},
		// 	},
		// 	Tok: token.DEFINE,
		// 	Rhs: []ast.Expr{
		// 		&ast.BasicLit{
		// 			Kind:  token.STRING,
		// 			Value: "\"YOUR_CLIENT_ID\"",
		// 		},
		// 	},
		// },

		// // ClientSecret := "YOUR_CLIENT_SECRET"
		// &ast.AssignStmt{
		// 	Lhs: []ast.Expr{
		// 		&ast.Ident{
		// 			Name: "ClientSecret",
		// 		},
		// 	},
		// 	Tok: token.DEFINE,
		// 	Rhs: []ast.Expr{
		// 		&ast.BasicLit{
		// 			Kind:  token.STRING,
		// 			Value: "\"YOUR_CLIENT_SECRET\"",
		// 		},
		// 	},
		// },

		// 	conf := &clientcredentials.Config{
		// 		ClientID:       clientID,
		// 		ClientSecret:   ClientSecret,
		// 		Scopes:         scopes,
		// 		TokenURL:       tokenUrl,
		// 		AuthStyle:      oauth2.AuthStyleInParams,
		// 		EndpointParams: additional_params,
		// 	}

		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "conf",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "clientcredentials",
							},
							Sel: &ast.Ident{
								Name: "Config",
							},
						},
						Elts: []ast.Expr{
							// &ast.KeyValueExpr{
							// 	Key: &ast.Ident{
							// 		Name: "ClientID",
							// 	},
							// 	Value: &ast.Ident{
							// 		Name: "clientID",
							// 	},
							// },
							// &ast.KeyValueExpr{
							// 	Key: &ast.Ident{
							// 		Name: "ClientSecret",
							// 	},
							// 	Value: &ast.Ident{
							// 		Name: "ClientSecret",
							// 	},
							// },
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "Scopes",
								},
								Value: &ast.Ident{
									Name: "scopes",
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "TokenURL",
								},
								Value: &ast.Ident{
									Name: "tokenUrl",
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "AuthStyle",
								},
								Value: &ast.Ident{
									Name: "oauth2.AuthStyleInParams",
								},
							},
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "EndpointParams",
								},
								Value: &ast.Ident{
									Name: "additional_params",
								},
							},
						},
					},
				},
			},
		},

		// Needed to request the OAuth token over an insecure channel, the default in free5gc
		// YOU WOULD NEED TO REMOVE THIS IN A PRODUCTION ENVIRONMENT!!!!!!!!!!!
		// TODO: command line to remove this from the template.
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "http_client",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "http",
							},
							Sel: &ast.Ident{
								Name: "Client",
							},
						},
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: &ast.Ident{
									Name: "Transport",
								},
								Value: &ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "http2",
											},
											Sel: &ast.Ident{
												Name: "Transport",
											},
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "TLSClientConfig",
												},
												Value: &ast.Ident{
													Name: "&tls.Config{InsecureSkipVerify: true}",
												},
											},
											// 	},
											// },
											// &ast.KeyValueExpr{
											// 	Key: &ast.Ident{
											// 		Name: "DialTLS",
											// 	},
											// 	Value: &ast.FuncLit{
											// 		Type: &ast.FuncType{
											// 			Params: &ast.FieldList{
											// 				List: []*ast.Field{
											// 					{
											// 						Names: []*ast.Ident{
											// 							{
											// 								Name: "network",
											// 							},
											// 							{
											// 								Name: "addr",
											// 							},
											// 						},
											// 						Type: &ast.Ident{
											// 							Name: "string",
											// 						},
											// 					},
											// 					{
											// 						Names: []*ast.Ident{
											// 							{
											// 								Name: "cfg",
											// 							},
											// 						},
											// 						Type: &ast.Ident{
											// 							Name: "*tls.Config",
											// 						},
											// 						// &ast.Ident{
											// 						// 	Name: "string",
											// 						// },
											// 					},
											// 				},
											// 			},
											// 			Results: &ast.FieldList{
											// 				List: []*ast.Field{
											// 					{
											// 						Type: &ast.SelectorExpr{
											// 							X: &ast.Ident{
											// 								Name: "net",
											// 							},
											// 							Sel: &ast.Ident{
											// 								Name: "Conn",
											// 							},
											// 						},
											// 					},
											// 					{
											// 						Type: &ast.Ident{
											// 							Name: "error",
											// 						},
											// 					},
											// 				},
											// 			},
											// 		},
											// 		Body: &ast.BlockStmt{
											// 			List: []ast.Stmt{
											// 				&ast.ReturnStmt{
											// 					Results: []ast.Expr{
											// 						&ast.CallExpr{
											// 							Fun: &ast.SelectorExpr{
											// 								X: &ast.Ident{
											// 									Name: "net",
											// 								},
											// 								Sel: &ast.Ident{
											// 									Name: "Dial",
											// 								},
											// 							},
											// 							Args: []ast.Expr{
											// 								&ast.Ident{
											// 									Name: "network",
											// 								},
											// 								&ast.Ident{
											// 									Name: "addr",
											// 								},
											// 							},
											// 						},
											// 					},
											// 				},
											// 			},
											// 		},
											// 	},
											// },
										},
									},
								},
							},
						},
					},
				},
			},
		},

		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "ctx",
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "context",
						},
						Sel: &ast.Ident{
							Name: "WithValue",
						},
					},
					Args: []ast.Expr{
						&ast.Ident{
							Name: "ctx",
						},
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "oauth2",
							},
							Sel: &ast.Ident{
								Name: "HTTPClient",
							},
						},
						&ast.Ident{
							Name: "http_client",
						},
					},
				},
			},
		},

		// token, err := conf.Token(ctx)
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "token",
				},
				&ast.Ident{
					Name: "err",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.Ident{
					Name: "conf.Token(ctx)",
				},
			},
		},
		checkErrors(len(callerFunc.Type.Results.List), "err.Error()"),
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "openapi",
					},
					Sel: &ast.Ident{
						Name: "PutCachedAT",
					},
				},
				Args: []ast.Expr{
					&ast.IndexExpr{
						X: &ast.Ident{
							Name: "scopes",
						},
						Index: &ast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
					},
					&ast.Ident{
						Name: "token",
					},
				},
			},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "ctx",
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "context",
						},
						Sel: &ast.Ident{
							Name: "WithValue",
						},
					},
					Args: []ast.Expr{
						&ast.Ident{
							Name: "ctx",
						},
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "openapi",
							},
							Sel: &ast.Ident{
								Name: "ContextAccessToken",
							},
						},
						&ast.SelectorExpr{
							X: &ast.Ident{
								Name: "token",
							},
							Sel: &ast.Ident{
								Name: "AccessToken",
							},
						},
					},
				},
			},
		},
	)

	oauthChecks = append(oauthChecks)

	// token := openapi.GetCachedAT(scopes[0])
	var cacheCheck []ast.Stmt

	cacheCheck = append(cacheCheck,
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "token",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "openapi",
						},
						Sel: &ast.Ident{
							Name: "GetCachedAT",
						},
					},
					Args: []ast.Expr{
						&ast.IndexExpr{
							X: &ast.Ident{
								Name: "scopes",
							},
							Index: &ast.BasicLit{
								Kind:  token.INT,
								Value: "0",
							},
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				Op: token.EQL,
				X: &ast.Ident{
					Name: "token",
				},
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: oauthChecks,
			},
			Else: &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "ctx",
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "WithValue",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "ctx",
							},
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "openapi",
								},
								Sel: &ast.Ident{
									Name: "ContextAccessToken",
								},
							},
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "token",
								},
								Sel: &ast.Ident{
									Name: "AccessToken",
								},
							},
						},
					},
				},
			},
		},
	)

	list = append(list, &ast.IfStmt{
		Cond: ast.NewIdent("oauth"),
		Body: &ast.BlockStmt{
			List: cacheCheck,
		},
	})

	if cfg.GetCfg().Verbosity {
		printer.Fprint(os.Stdout, token.NewFileSet(), list)
	}

	return list
}

func GetEnforceOAuthConfigVar() *ast.Field {
	newField := ast.Field{
		Names: []*ast.Ident{
			{
				Name: "OAuth",
			},
		},
		Type: &ast.Ident{
			Name: "bool",
		},
		Tag: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "`yaml:\"OAuth,omitempty\"`",
		},
	}

	return &newField
}
