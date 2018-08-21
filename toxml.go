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
	"github.com/jessevdk/go-flags"
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
var footer string = `</graphml>`

// Flags.
var opts = struct {
	Filename  string   `short:"f" long:"filename" description:"File to read"`
	Directed  bool   `short:"d" long:"directed" description:"Convert to directed"`
}{
	Filename: "",
	Directed: false,
}

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	if (opts.Filename == "") {
		fmt.Println("No file name specified")
		os.Exit(1)
	}
	graphFile, err := os.Open(opts.Filename)
	if (err != nil){
		fmt.Printf("Error reading file %s: %s\n", opts.Filename, err)
		os.Exit(1)		
	}
	keys := []Key{
		Key{Id: "pub_key", For: "node", Name: "pub_key", Type: "string"},
		Key{Id: "name", For: "node", Name: "name", Type: "string"},

		Key{Id: "chan_point", For: "edge", Name: "chan_point", Type: "string"},
		Key{Id: "last_update", For: "edge", Name: "last_update", Type: "string"},
		Key{Id: "capacity", For: "edge", Name: "capacity", Type: "int"},
		Key{Id: "time_lock_delta", For: "edge", Name: "time_lock_delta", Type: "int"},
		Key{Id: "min_htlc", For: "edge", Name: "min_htlc", Type: "int", Default: 0},
		Key{Id: "fee_base_msat", For: "edge", Name: "fee_base_msat", Type: "int"},
		Key{Id: "fee_rate_milli_msat", For: "edge", Name: "fee_rate_milli_msat", Type: "int"},
		Key{Id: "disabled", For: "edge", Name: "disabled", Type: "string"},
	}
	var graph = Graph{
		Id: "LND",
		NodeIds: "free",
		EdgeIds: "free",
		ParseOrder: "nodesfirst",
		EdgeDefault: "undirected",
	}
	if (opts.Directed) {
		graph.EdgeDefault = "directed"
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

	keysOutput, err := xml.MarshalIndent(keys, "  ", "    ")
	if err != nil {
        fmt.Printf("error encoding xml %#v\n", err)
	}

	var output []byte
	var directedGraph *DirectedGraph

	graph.setupDataAttrs()
	
	if (opts.Directed) {
		directedGraph = graph.makeDirected()
		output, err = xml.MarshalIndent(directedGraph, "  ", "    ")
		if err != nil {
			fmt.Printf("error encoding xml %#v\n", err)
		}
	} else {
		output, err = xml.MarshalIndent(graph, "  ", "    ")
		if err != nil {
			fmt.Printf("error encoding xml %#v\n", err)
		}
	}
	
	os.Stdout.Write([]byte(header))
	os.Stdout.Write(keysOutput)
	os.Stdout.Write(output)
	os.Stdout.Write([]byte(footer))
    return
}
