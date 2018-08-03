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
	"encoding/json"
	"os"
	"encoding/xml"
)

var header string = `<?xml version="1.0" encoding="UTF-8"?>
            <graphml xmlns="http://graphml.graphdrawing.org/xmlns"
            xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
            xsi:schemaLocation="http://graphml.graphdrawing.org/xmlns 
                                http://graphml.graphdrawing.org/xmlns/1.0/graphml.xsd">`

func main() {
	graphFile, err := os.Open("testdata/lndgraph.json")
	var graph = Graph{
		Id: "LND",
		NodeIds: "free",
		EdgeIds: "free",
		ParseOrder: "nodesfirst",
	}
    if err != nil {
        fmt.Printf("Error opening file %s\n", err)
    }

    jsonParser := json.NewDecoder(graphFile)
    if err = jsonParser.Decode(&graph); err != nil {
        fmt.Printf("parsing config file %s\n", err)
    }

	graph.NumNodes = len(graph.Nodes)
	graph.NumEdges = len(graph.Edges)

	directedGraph := graph.makeDirected()
	
	output, err := xml.MarshalIndent(directedGraph, "  ", "    ")
	if err != nil {
        fmt.Printf("error encoding xml %#v\n", err)
	}
	os.Stdout.Write([]byte(header))
	os.Stdout.Write(output)
    return
}