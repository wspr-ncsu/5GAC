package policy

import (
	"log"
	"os"

	"github.com/wspr-ncsu/5GAC/gocode/Shared/accesspattern"
	"gopkg.in/yaml.v2"
)

// writeNfScopes writes the output of the network function interactions to a single YAML file,
// to later be consumed by the authorization server.
func WriteAuthorizationServerPolicy(nfScopes map[string][]accesspattern.API) {
	authPolicy, err := yaml.Marshal(nfScopes)
	if err != nil {
		log.Fatal(err)
	}

	fm := os.FileMode(0o600)
	err = os.WriteFile("../_output/AuthorizationPolicy.yaml", authPolicy, fm)

	if err != nil {
		log.Fatal(err)
	}
}

func WriteAccessPatterns(nfScopes []accesspattern.AccessPattern) {
	authPolicy, err := yaml.Marshal(nfScopes)
	if err != nil {
		log.Fatal(err)
	}

	fm := os.FileMode(0o600)
	err = os.WriteFile("../_output/AccessPatterns.yaml", authPolicy, fm)

	if err != nil {
		log.Fatal(err)
	}
}

func nfInAcp(acp []accesspattern.NetworkFunction, nf accesspattern.NetworkFunction) bool {
	for _, rule := range acp {
		if rule == nf {
			return true
		}
	}

	return false
}

func WriteSimpleAccessControlPolicy(nfScopes []accesspattern.AccessPattern) {
	acp := make(map[string][]accesspattern.NetworkFunction)
	for _, nfScope := range nfScopes {
		if !nfInAcp(acp[nfScope.RequiredScopes], nfScope.Caller.Nf) {
			acp[nfScope.RequiredScopes] = append(acp[nfScope.RequiredScopes], nfScope.Caller.Nf)
		}
	}

	accessControlPolicy, err := yaml.Marshal(acp)
	if err != nil {
		log.Fatal(err)
	}

	fm := os.FileMode(0o600)
	err = os.WriteFile("../_output/acp.yaml", accessControlPolicy, fm)

	if err != nil {
		log.Fatal(err)
	}
}
