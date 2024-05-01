package policy

import (
	"fmt"
	"log"
	"strings"

	"github.com/wspr-ncsu/5GAC/gocode/Analyzer/cfg"
	"github.com/wspr-ncsu/5GAC/gocode/Shared/accesspattern"
)

type temporaryPatterns map[string]map[accesspattern.NetworkFunction][]accesspattern.API

func getTempScopes(nfInstances []accesspattern.AccessPattern) temporaryPatterns {
	tmpScopes := make(temporaryPatterns)

	for _, nfInstance := range nfInstances {
		// if nfInstance.RequiredScopes == "nnrf-disc" || nfInstance.RequiredScopes == "nnrf-nfm" {
		// 	continue
		// }

		if tmpScopes[nfInstance.RequiredScopes] == nil {
			tmpScopes[nfInstance.RequiredScopes] = make(map[accesspattern.NetworkFunction][]accesspattern.API)
		}

		found := false
		// Handle duplicates
		for _, API := range tmpScopes[nfInstance.RequiredScopes][nfInstance.Caller.Nf] {
			if API.HTTPType == nfInstance.HTTPType && API.URL == nfInstance.URL {
				if cfg.GetCfg().Verbosity {
					log.Printf("scope: %v\n", nfInstance.RequiredScopes)
				}

				found = true

				break
			}
		}

		if found {
			continue
		}

		API := nfInstance.API
		caller := nfInstance.Caller.Nf
		scopes := nfInstance.RequiredScopes

		tmpScopes[scopes][caller] = append(tmpScopes[scopes][caller], API)
	}

	return tmpScopes
}

// improvePolicy takes in the nfInteractions and creates a security policy based on it.
// Map URLs to calling functions - Access Pattern.
// Compare to existing security label for access token.
// First get all the NFs that contact a service.
func GenerateSecurePolicy(nfInstances []accesspattern.AccessPattern) map[string][]accesspattern.API {
	proposedScopes := make(map[string][]accesspattern.API)
	scopes := make(map[string][]accesspattern.API)

	for _, nfInstance := range nfInstances {
		scopes[nfInstance.RequiredScopes] = append(scopes[nfInstance.RequiredScopes], nfInstance.API)

		if cfg.GetCfg().Verbosity {
			log.Printf("nfInstance: %v\n", nfInstance)
		}
	}

	tmpScopes := getTempScopes(nfInstances)

	for scope, nfs := range tmpScopes {
		// We re-use an already existing scope in the spec for the NF that calls the most APIs for each service.
		// We don't need to do this, but this seems the most logical way to tame unnecessary scopes
		// and bridge the gap to already existing specs and implementation.
		if len(nfs) > 1 {
			// Get the NF that has the most APIs.
			// largest := 0

			// var largestNf accesspattern.NetworkFunction

			// for nf, apis := range nfs {
			// 	newLargest := len(apis)
			// 	if newLargest > largest {
			// 		largest = newLargest
			// 		largestNf = nf
			// 	}
			// }
			// fmt.Printf("largest nf for %v: %v\n", scope, largest_nf)

			for nf, apis := range nfs {
				// if nf == largestNf {
				// 	continue
				// }

				// if cfg.GetCfg().Verbosity {
				// 	log.Printf("%v calls %v w/ scope %v -> %v", )
				// }

				proposedScope := scope
				if scope != "nnrf-disc" && scope != "nnrf-nfm" {
					proposedScope = fmt.Sprintf("%v:%v", scope, strings.ToLower(string(nf)))
				}
				proposedScopes[proposedScope] = apis
			}
		}
	}

	return proposedScopes
}
