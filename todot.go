/* Copyright (c) 2018 Kulpreet Singh
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */
package main

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"encoding/json"
	"os"
)

func export(graph Graph){
	graphName := "LightningNetwork"
	graphCanvas := gographviz.NewGraph()
	graphCanvas.SetName(graphName)
	graphCanvas.SetDir(true)

	// For each node within the graph, we'll add a new vertex to the graph.
	for _, node := range graph.Nodes {
		nodeID := fmt.Sprintf(`"%v"`, node.PubKey)
		var attrs map[string]string
		graphCanvas.AddNode(graphName, nodeID, attrs)
	}

	for _, edge := range graph.Edges {
		src := fmt.Sprintf(`"%v"`, edge.Node1Pub)
		dest := fmt.Sprintf(`"%v"`, edge.Node2Pub)
		fmt.Println(src)
		fmt.Println(dest)

		fmt.Printf("tld: %d\n", edge.Node1Policy.TimeLockDelta)
		fmt.Printf("feebase: %v\n", edge.Node1Policy.FeeBaseMsat)
		if err := graphCanvas.AddEdge(src, dest, true, map[string]string{
			"weight":   fmt.Sprintf("%v", edge.Capacity),
			"time_lock_delta": fmt.Sprintf("%v", edge.Node1Policy.TimeLockDelta),
			"min_htlc": edge.Node1Policy.MinHtlc,
			"fee_base_msat": fmt.Sprintf("%v", edge.Node1Policy.FeeBaseMsat),
			"fee_rate_milli_msat": edge.Node1Policy.FeeRateMilliMsat,
			"disabled": fmt.Sprintf("%v", edge.Node1Policy.Disabled),
		}); err != nil {
			fmt.Println(err)
		}

		// if err := graphCanvas.AddEdge(dest, src, true, map[string]string{
		// 	"weight":   fmt.Sprintf("%v", edge.Capacity),
		// 	"time_lock_delta": fmt.Sprintf("%d", edge.Node2Policy.TimeLockDelta),
		// 	"min_htlc": edge.Node2Policy.MinHtlc,
		// 	"fee_base_msat": edge.Node2Policy.FeeBaseMsat,
		// 	"fee_rate_milli_msat": edge.Node2Policy.FeeRateMilliMsat,
		// 	"disabled": fmt.Sprintf("%v", edge.Node2Policy.Disabled),
		// }); err != nil {
		// 	fmt.Println(err)
		// }
	}
	//graphDotString := graphCanvas.String()
	//os.Stdout.Write([]byte(graphDotString))
}

func ToDot() {
	graphFile, err := os.Open("testdata/lndgraph.json")
	var graph = Graph{}
    if err != nil {
        fmt.Printf("Error opening file %s\n", err)
    }

    jsonParser := json.NewDecoder(graphFile)
    if err = jsonParser.Decode(&graph); err != nil {
        fmt.Printf("parsing config file %s\n", err)
    }
	export(graph)
    return
}
