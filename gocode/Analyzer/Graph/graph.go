package graph

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/wspr-ncsu/5GAC/gocode/Instrumenter/cfg"
	"github.com/wspr-ncsu/5GAC/gocode/Shared/accesspattern"
)

func createGraph() *gographviz.Graph {
	g, _ := gographviz.ParseString(`digraph G { ratio="fill"; size="15,11.7!"; margin=4; rankdir=LR}`)
	graph := gographviz.NewGraph()

	if err := gographviz.Analyse(g, graph); err != nil {
		panic(err)
	}

	// This is a directed graph
	err := graph.SetDir(true)
	if err != nil {
		log.Fatalf("Error occurred setting graph to directed. - %v\n", err)
	}

	err = graph.SetStrict(true)
	if err != nil {
		log.Fatalf("Error occurred setting graph to strict. - %v\n", err)
	}

	return graph
}

func sanitizeScopes(scope string) string {
	scopes := strings.ReplaceAll(scope, "-", "")

	return strings.ReplaceAll(scopes, ":", "")
}

func createNodes(graph *gographviz.Graph, nfInstance accesspattern.AccessPattern) (string, string, string, string) {
	sanitizedScopes := sanitizeScopes(nfInstance.RequiredScopes)
	strCaller := string(nfInstance.Caller.Nf)
	strCallee := string(nfInstance.Callee.Nf)

	if !graph.IsNode(strCaller) {
		err := graph.AddNode("G", strCaller, map[string]string{"size": "\"10,10\"", "ratio": "fill"})
		if err != nil {
			log.Fatalf("Error creating node %v\n", strCaller)
		}
	}

	url := fmt.Sprintf("%v%v", nfInstance.URL.Root, nfInstance.URL.Path)
	node := fmt.Sprintf("\"%v %v\"", nfInstance.HTTPType, url)

	if graph.IsNode(node) && cfg.GetCfg().Verbosity {
		log.Printf("Warning: Node %v already exists\n", node)
	}

	return sanitizedScopes, strCaller, strCallee, url
}

func createSubGraphs(graph *gographviz.Graph, scopes, callee string) string {
	cluster := fmt.Sprintf("cluster_%v", callee)

	if !graph.IsSubGraph(cluster) {
		err := graph.AddSubGraph("G", cluster, map[string]string{"label": callee})
		if err != nil {
			log.Fatalf("(1) Error creating subgraph %v\n", callee)
		}
	}

	clusterScopes := fmt.Sprintf("cluster_%v", scopes)

	if !graph.IsSubGraph(clusterScopes) {
		err := graph.AddSubGraph(cluster, clusterScopes, map[string]string{"label": scopes})
		if err != nil {
			log.Fatalf("(2) Error creating subgraph %v:%v\n", callee, scopes)
		}
	}

	return clusterScopes
}

// generateGraph create a dot graph representing the NF interactions. Converted to PDF outside of this program.
func GenerateGraph(nfinstances []accesspattern.AccessPattern, name string) {
	// Set some defaults for the graph
	graph := createGraph()

	// For each NF, get its scopes and set an edge to and from it's interactions.
	for _, nfInstance := range nfinstances {
		sanitizedScopes, strCaller, strCallee, url := createNodes(graph, nfInstance)
		cluster := createSubGraphs(graph, sanitizedScopes, strCallee)

		err := graph.AddNode(cluster, "\""+nfInstance.HTTPType+" "+url+"\"", nil)
		if err != nil {
			log.Fatalf("Error (%v) adding node %v", err, "\""+nfInstance.HTTPType+" "+url+"\"")
		}

		err = graph.AddEdge(strCaller, "\""+nfInstance.HTTPType+" "+url+"\"", true, nil)
		if err != nil {
			log.Fatalf("Error creating edge %v %v - %v\n", strCaller, "\""+nfInstance.HTTPType+" "+url+"\"", err)
		}
	}

	writeGraph(graph, name)
}

func writeGraph(graph *gographviz.Graph, name string) {
	// Write to the file
	file, err := os.Create(fmt.Sprintf("../_output/%v.dot", name))
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.WriteString(graph.String())
	if err != nil {
		file.Close()
		log.Fatalf("Error writing graph to file. %v", err)
	}

	file.Close()
}
