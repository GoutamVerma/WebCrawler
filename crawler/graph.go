package graph

import (
	"fmt"
	"strings"
)

type Vertex struct {
	value string
	edge  *Vertex
}

type Graph struct {
	Adjacency map[string][]string
}

func NewGraph() Graph {
	return Graph{
		Adjacency: make(map[string][]string),
	}
}

func (g *Graph) AddVertex(vertex string) bool {
	if _, ok := g.Adjacency[vertex]; ok {
		return false
	}
	g.Adjacency[vertex] = []string{}
	return true
}

func (g *Graph) AddEdge(vertex, node string) bool {
	if _, ok := g.Adjacency[vertex]; !ok {
		return false
	}
	if ok := contains(g.Adjacency[vertex], node); ok {
		return false
	}

	if _, ok := g.Adjacency[node]; !ok {
		g.AddVertex(node)
	}

	g.Adjacency[vertex] = append(g.Adjacency[vertex], node)
	return true
}

func contains(edges []string, node string) bool {
	set := make(map[string]struct{}, len(edges))
	for _, n := range edges {
		set[n] = struct{}{}
	}
	_, ok := set[node]
	return ok
}

func (g Graph) Print() {
	for i, val := range g.Adjacency {
		fmt.Printf("key: %s value %s \n \n", i, strings.Join(val, " -> "))
	}
}

func (g Graph) CreatePath(firstNode, secondNode string) bool {
	visited := g.createVisited()
	var (
		path []string
		q    []string
	)
	q = append(q, firstNode)
	visited[firstNode] = true

	for len(q) > 0 {
		var currentNode string
		currentNode, q = q[0], q[1:]
		path = append(path, currentNode)
		edges := g.Adjacency[currentNode]
		if contains(edges, secondNode) {
			path = append(path, secondNode)
			fmt.Println(strings.Join(path, "->"))
			return true
		}

		for _, node := range g.Adjacency[currentNode] {
			if !visited[node] {
				visited[node] = true
				q = append(q, node)
			}
		}
	}
	fmt.Println("no link found")
	return false
}

func (g Graph) createVisited() map[string]bool {
	visited := make(map[string]bool, len(g.Adjacency))
	for key := range g.Adjacency {
		visited[key] = false
	}
	return visited
}

func (g Graph) PrintSitemap(baseUrl string) {
	fmt.Println(baseUrl)
	g.printSitemapHelper(baseUrl, 1)
}

func (g Graph) printSitemapHelper(currentUrl string, depth int) {
	edges := g.Adjacency[currentUrl]
	if len(edges) == 0 {
		return
	}

	for _, edge := range edges {
		indent := strings.Repeat("\t", depth)
		fmt.Printf("%s- %s\n", indent, edge)
		g.printSitemapHelper(edge, depth+1)
	}
}
