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
	"encoding/xml"
)

var _ = fmt.Printf // For debugging; delete when done.

type Data struct {
	XMLName   xml.Name `xml:"data"`
	Key string `xml:"key,attr"`
	Value interface{} `xml:",chardata"`
}

type CheckedData struct {
	XMLName   xml.Name `xml:"data"`
	Key string `xml:"key,attr"`
	Value interface{} `xml:",chardata"`
}

func (d *Data) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if (d.Value == "") {
		return nil
	}
	e.EncodeElement(CheckedData(*d), start)
	return nil
}

type Address struct {
	XMLName   xml.Name `xml:"address"`
	Network string `json:"network" xml:"network,attr"`
	Addr string `json:"addr" xml:"addr,attr"`	
}

type Node struct {
	XMLName   xml.Name `xml:"node"`
	LastUpdate int `json:"last_update" xml:"-"`
	PubKey string `json:"pub_key" xml:"id,attr"`
	Alias string `json:"alias" xml:"-"`
	Addresses []Address `xml:"-"`
	Color string `json:"color" xml:"-"`
	Attrs []Data `xml:"data"`
}

type NodePolicy struct {
	TimeLockDelta int `json:"time_lock_delta"`
	MinHtlc string `json:"min_htlc"`
	FeeBaseMsat string `json:"fee_base_msat"`
	FeeRateMilliMsat string `json:"fee_rate_milli_msat"`
	Disabled bool `json:"disabled"`
}

type Edge struct {
	XMLName   xml.Name `xml:"edge"`
	ChannelId string `json:"channel_id" xml:"id,attr"`
	Node1Pub string `json:"node1_pub" xml:"source,attr"`
	Node2Pub string `json:"node2_pub" xml:"target,attr"`

	ChanPoint string `json:"chan_point" xml:"-"`
	LastUpdate int `json:"last_update" xml:"-"`
	Capacity string `json:"capacity" xml:"-"`
	Node1Policy NodePolicy `json:"node1_policy" xml:"-"`
	Node2Policy NodePolicy `json:"node2_policy" xml:"-"`

	Attrs []Data `xml:"data"`
}

type DirectedEdge struct {
	XMLName   xml.Name `xml:"edge"`
	Id string `json:"channel_id" xml:"id,attr"`
	Source string `json:"node1_pub" xml:"source,attr"`
	Target string `json:"node2_pub" xml:"target,attr"`
	Attrs []Data `xml:"data"`
}

type Key struct {
	XMLName   xml.Name `xml:"key"`
	Id string `xml:"id,attr"`
	For string `xml:"for,attr"`
	Name string `xml:"attr.name,attr"`
	Type string `xml:"attr.type,attr"`
	Default interface{} `xml:"default"`
}

type Graph struct {
	XMLName   xml.Name `xml:"graph"`
	Keys []Key
	Id string `xml:"id,attr"`
	NumNodes int `xml:"parse.nodes,attr"`
	NumEdges int `xml:"parse.edges,attr"`
	NodeIds string `xml:"parse.nodeids,attr"`
	EdgeIds string `xml:"parse.edgeids,attr"`
	ParseOrder string `xml:"parse.order,attr"`
	EdgeDefault string `xml:"edgedefault,attr"`
	Nodes []*Node
	Edges []*Edge
}

type DirectedGraph struct {
	Graph
	Nodes []*Node
	Edges []*DirectedEdge
}

func (g *Graph) setupDataAttrs() {
	for _, node := range g.Nodes {
		node.Attrs = []Data{
			//Data{Key: "last_update", Value: node.LastUpdate},
			Data{Key: "pub_key", Value: node.PubKey},
			Data{Key: "name", Value: node.Alias},
			//Data{Key: "color", Value: node.Color},
		}
	}
	for _, edge := range g.Edges {
		edge.Attrs = []Data{
			Data{Key: "chan_point", Value: edge.ChanPoint},
			Data{Key: "last_update", Value: edge.LastUpdate},
			Data{Key: "capacity", Value: edge.Capacity},
		}
	}
}

func (g *Graph) makeDirected() (directedGraph *DirectedGraph) {
	var directedEdges []*DirectedEdge
	directedGraph = &DirectedGraph{}
	for _, edge := range g.Edges {
		outDir := &DirectedEdge{
			Id: edge.ChannelId,
			Source: edge.Node1Pub,
			Target: edge.Node2Pub,
		}
		outDir.Attrs = []Data{
			Data{Key: "chan_point", Value: edge.ChanPoint},
			Data{Key: "last_update", Value: edge.LastUpdate},
			Data{Key: "capacity", Value: edge.Capacity},
			
			Data{Key: "time_lock_delta", Value: edge.Node1Policy.TimeLockDelta},
			Data{Key: "min_htlc", Value: edge.Node1Policy.MinHtlc},
			Data{Key: "fee_base_msat", Value: edge.Node1Policy.FeeBaseMsat},
			Data{Key: "fee_rate_milli_msat", Value: edge.Node1Policy.FeeRateMilliMsat},
			Data{Key: "disabled", Value: edge.Node1Policy.Disabled},
		}
		
		inDir := &DirectedEdge{
			Id: edge.ChannelId,
			Source: edge.Node2Pub,
			Target: edge.Node1Pub,

			Attrs: []Data{
				Data{Key: "chan_point", Value: edge.ChanPoint},
				Data{Key: "last_update", Value: edge.LastUpdate},
				Data{Key: "capacity", Value: edge.Capacity},
				
				Data{Key: "time_lock_delta", Value: edge.Node2Policy.TimeLockDelta},
				Data{Key: "min_htlc", Value: edge.Node2Policy.MinHtlc},
				Data{Key: "fee_base_msat", Value: edge.Node2Policy.FeeBaseMsat},
				Data{Key: "fee_rate_milli_msat", Value: edge.Node2Policy.FeeRateMilliMsat},
				Data{Key: "disabled", Value: edge.Node2Policy.Disabled},
			},
		}
		directedEdges = append(directedEdges, inDir, outDir)
	}
	directedGraph.XMLName = g.XMLName
	directedGraph.Keys = g.Keys
	directedGraph.Id = g.Id
	directedGraph.NumNodes = g.NumNodes
	directedGraph.NumEdges = len(directedEdges)
	directedGraph.NodeIds = g.NodeIds
	directedGraph.EdgeIds = g.EdgeIds
	directedGraph.ParseOrder = g.ParseOrder
	directedGraph.EdgeDefault = g.EdgeDefault
	directedGraph.Nodes = g.Nodes
	directedGraph.Edges = directedEdges

	return
}
