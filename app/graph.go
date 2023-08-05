package main

import (
)

// mutation is for suckers

type Graph[T comparable] struct {
	Nodes map[T]int // int here represents color
	Edges map[T]map[T]bool
}

func NewGraph[T comparable]() *Graph[T] {
	return &Graph[T]{
		Nodes: make(map[T]int),
		Edges: make(map[T]map[T]bool),
	}
}

func AddNode[T comparable](graph Graph[T], node T) *Graph[T] {
	graph.Nodes[node] = -1
	return &graph
}

func HasNode[T comparable](graph *Graph[T], node T) bool {
	_, present := graph.Nodes[node]
	return present
}

func AddEdge[T comparable](graph Graph[T], node1 T, node2 T) *Graph[T] {
	graph = *AddNode(graph, node1)
	graph = *AddNode(graph, node2)

	if node1 == node2 {
		return &graph
	}

	_, node1Present := graph.Edges[node1]
	_, node2Present := graph.Edges[node2]
	if !node1Present {
		graph.Edges[node1] = make(map[T]bool)
	}
	if !node2Present {
		graph.Edges[node2] = make(map[T]bool)
	}

	graph.Edges[node1][node2] = true
	graph.Edges[node2][node1] = true
	return &graph
}

func GetNodes[T comparable](graph Graph[T]) []T {
	nodes := make([]T, len(graph.Nodes))
	i := 0
	for key, _ := range(graph.Nodes) {
		nodes[i] = key
		i++
	}
	return nodes
}

func GetEdges[T comparable](graph *Graph[T], node T) []T {
	var edges []T
	for key, _ := range(graph.Edges[node]) {
		edges = append(edges, key)
	}
	return edges
}

// this is a greedy search and is unlikely to come up with the optimal coloring
func ColorGraph[T comparable](graph *Graph[T]) map[T]int {
	uncoloredNodes := make([]T, len(graph.Nodes))
	saturation := make(map[T][]int)
	for key, _ := range(graph.Nodes) {
		uncoloredNodes = append(uncoloredNodes, key)
		saturation[key] = make([]int, 0)
	}
	return ColorGraphHelper(graph, saturation, make(map[T]int), []int{0})
}

func ColorGraphHelper[T comparable](graph *Graph[T], saturation map[T][]int, colorings map[T]int, colors []int) map[T]int {
	if len(saturation) == 0 {
		return colorings
	}

	node := GetMaxSaturatedNode(saturation)
	adjacentColors := GetAdjacentColors(node, graph, colorings)
	nodeColor, colors := GetMinimalColor(adjacentColors, colors)

	colorings[node] = nodeColor
	delete(saturation, node)

	for _, adjacentNode := range GetEdges(graph, node) {
		_, present := saturation[adjacentNode]
		if present {
			saturation[adjacentNode] = append(saturation[adjacentNode], nodeColor)
		}
	}

	return ColorGraphHelper(graph, saturation, colorings, colors)
}

// second return value is in case a new color had to be added
// as stated above, this is not an optimal solution
func GetMinimalColor(adjacentColors []int, availableColors []int) (int, []int) {
	for available := range availableColors {
		colorAvailable := true
		for adjacent := range adjacentColors {
			if available == adjacent {
				colorAvailable = false
				break
			}
		}
		if colorAvailable {
			return available, availableColors
		}
	}
	newColor := availableColors[len(availableColors) - 1] + 1
	return newColor, append(availableColors, newColor)
}

func GetAdjacentColors[T comparable](node T, graph *Graph[T], colorings map[T]int) []int {
	var adjacentColors []int
	edges := GetEdges(graph, node)
	for _, edge := range(edges) {
		color, present := colorings[edge]
		if present {
			adjacentColors = append(adjacentColors, color)
		}
	}
	return adjacentColors
}

// the V generic is unnecessary but I mean why not when we're having fun
func GetMaxSaturatedNode[T comparable, V any](saturations map[T][]V) T {
	var maxSaturatedNode T
	maxSaturation := -1
	for key, val := range(saturations) {
		if len(val) > maxSaturation {
			maxSaturatedNode = key
		}
	}
	return maxSaturatedNode
}
