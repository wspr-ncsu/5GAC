package accesspattern

import (
	"log"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/wspr-ncsu/5GAC/gocode/Analyzer/cfg"
	"github.com/wspr-ncsu/5GAC/gocode/Shared/openapi"
)

// updateSecurity takes in the API list from the specifications and the functional interactions
// based on code analysis to determine a more least privileged policy.
func UpdateSecurity(policy interface{}) {
	APIs := openapi.APIList
	accesspatterns := GetAccessPatterns()

	proposedPolicy, ok := policy.(map[string][]API)
	if !ok {
		ap, ok := policy.([]AccessPattern)
		if !ok {
			log.Fatalf("Error: Cannot convert policy into readable format %v\n", policy)
		}

		proposedPolicy = make(map[string][]API)

		for _, accessPattern := range ap {
			scope := accessPattern.RequiredScopes
			if proposedPolicy[scope] == nil {
				proposedPolicy[scope] = make([]API, 0, 3)
			}
			proposedPolicy[scope] = append(proposedPolicy[scope], accessPattern.API)
		}
	}

	for newScope, proposedApis := range proposedPolicy {
		for _, API := range proposedApis {
			for _, origAPI := range APIs {
				if origAPI.Scheme == API.HTTPType && origAPI.Path == API.URL {
					if cfg.GetCfg().Verbosity {
						log.Printf("Found matching API for %v: %v\n", origAPI, API)
					}

					newReq := make(openapi3.SecurityRequirement)
					newReq["oAuth2ClientCredentials"] = make([]string, 0, 1)
					newReq["oAuth2ClientCredentials"] = append(newReq["oAuth2ClientCredentials"], newScope)

					origAPI.SecurityRequirements = make(openapi3.SecurityRequirements, 1, 2) // handle optional OAuth
					origAPI.SecurityRequirements = append(origAPI.SecurityRequirements, newReq)
				}
			}

			for i, ap := range accesspatterns {
				if ap.API.HTTPType == API.HTTPType && ap.API.URL == API.URL {
					accesspatterns[i].RequiredScopes = newScope
				}
			}
		}
	}
	// TODO: Write updated security policy to OpenAPI spec directly. Instrumenter should consume this.
	openapi.APIList = APIs
}
