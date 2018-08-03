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
	"encoding/xml"
)

type Address struct {
	XMLName   xml.Name `xml:"address"`
	Network string `json:"network" xml:"network,attr"`
	Addr string `json:"addr" xml:"addr,attr"`	
}

type Node struct {
	XMLName   xml.Name `xml:"node"`
	LastUpdate int `json:"last_update" xml:"last_update,attr"`
	PubKey string `json:"pub_key" xml:"pub_key,attr"`
	Alias string `json:"alias" xml:"alias,attr"`
	Addresses []Address
	Color string `json:"color" xml:"color,attr"` 
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
	ChannelId string `json:"channel_id" xml:"channel_id,attr"`
	ChanPoint string `json:"chan_point" xml:"channel_point,attr"`
	LastUpdate int `json:"last_update" xml:"last_update,attr"`
	Node1Pub string `json:"node1_pub" xml:"node1_pub,attr"`
	Node2Pub string `json:"node2_pub" xml:"node2_pub,attr"`
	Capacity string `json:"capacity" xml:"capacity,attr"`
	Node1Policy NodePolicy `json:"node1_policy" xml:"node1_policy"`
	Node2Policy NodePolicy `json:"node2_policy" xml:"node2_policy"`
}

type DirectedEdge struct {
	XMLName   xml.Name `xml:"edge"`
	ChannelId string `json:"channel_id" xml:"channel_id,attr"`
	ChanPoint string `json:"chan_point" xml:"channel_point,attr"`
	LastUpdate int `json:"last_update" xml:"last_update,attr"`
	OutPub string `json:"node1_pub" xml:"out_pub,attr"`
	InPub string `json:"node2_pub" xml:"in_pub,attr"`
	Capacity string `json:"capacity" xml:"capacity,attr"`

	// NodePolicy flattened into directed edge
	TimeLockDelta int `json:"time_lock_delta,attr" xml:"timelock_delta,attr"`
	MinHtlc string `json:"min_htlc,attr" xml:"min_htlc,attr"`
	FeeBaseMsat string `json:"fee_base_msat,attr" xml:"fee_base_msat,attr"`
	FeeRateMilliMsat string `json:"fee_rate_milli_msat,attr" xml:"fee_rate_milli_msat,attr"`
	Disabled bool `json:"disabled,attr" xml:"disabled,attr"`
}

type Graph struct {
	XMLName   xml.Name `xml:"graph"`
	Id string `xml:"id,attr"`
	NumNodes int `xml:"parse.nodes,attr"`
	NumEdges int `xml:"parse.edges,attr"`
	NodeIds string `xml:"parse.nodeids,attr"`
	EdgeIds string `xml:"parse.edgeids,attr"`
	ParseOrder string `xml:"parse.order,attr"`
	Nodes []Node
	Edges []Edge
}

type DirectedGraph struct {
	XMLName   xml.Name `xml:"graph"`
	Id string `xml:"id,attr"`
	NumNodes int `xml:"parse.nodes,attr"`
	NumEdges int `xml:"parse.edges,attr"`
	NodeIds string `xml:"parse.nodeids,attr"`
	EdgeIds string `xml:"parse.edgeids,attr"`
	ParseOrder string `xml:"parse.order,attr"`
	Nodes []Node
	Edges []DirectedEdge
}

func (g *Graph) makeDirected() (directedGraph *DirectedGraph) {
	var directedEdges []DirectedEdge
	directedGraph = &DirectedGraph{}
	for _, edge := range g.Edges {
		inDir := DirectedEdge{
			ChannelId: edge.ChannelId,
			ChanPoint: edge.ChanPoint,
			LastUpdate: edge.LastUpdate,
			OutPub: edge.Node1Pub,
			InPub: edge.Node2Pub,
			Capacity: edge.Capacity,

			TimeLockDelta: edge.Node1Policy.TimeLockDelta,
			MinHtlc: edge.Node1Policy.MinHtlc,
			FeeBaseMsat: edge.Node1Policy.FeeBaseMsat,
			FeeRateMilliMsat: edge.Node1Policy.FeeRateMilliMsat,
			Disabled: edge.Node1Policy.Disabled,			
		}
		outDir := DirectedEdge{
			ChannelId: edge.ChannelId,
			ChanPoint: edge.ChanPoint,
			LastUpdate: edge.LastUpdate,
			InPub: edge.Node1Pub,
			OutPub: edge.Node2Pub,
			Capacity: edge.Capacity,

			TimeLockDelta: edge.Node2Policy.TimeLockDelta,
			MinHtlc: edge.Node2Policy.MinHtlc,
			FeeBaseMsat: edge.Node2Policy.FeeBaseMsat,
			FeeRateMilliMsat: edge.Node2Policy.FeeRateMilliMsat,
			Disabled: edge.Node2Policy.Disabled,			
		}
		directedEdges = append(directedEdges, inDir, outDir)
	}
	directedGraph.XMLName = g.XMLName
	directedGraph.Id = g.Id
	directedGraph.NumNodes = g.NumNodes
	directedGraph.NumEdges = len(directedEdges)
	directedGraph.NodeIds = g.NodeIds
	directedGraph.EdgeIds = g.EdgeIds
	directedGraph.ParseOrder = g.ParseOrder
	directedGraph.Nodes = g.Nodes
	directedGraph.Edges = directedEdges

	return
}
