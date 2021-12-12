package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Node struct {
	name     string
	is_small bool
}

func (node Node) String() string {
	return node.name
}

func newNode(name string) Node {
	return Node{
		name:     name,
		is_small: name == strings.ToLower(name),
	}
}

type Graph struct {
	children map[Node][]Node
}

func parseGraph(input string) (Graph, error) {
	input = strings.TrimSpace(input)
	input = strings.Trim(input, "\n")
	lines := strings.Split(input, "\n")

	children := map[Node][]Node{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		chunks := strings.SplitN(line, "-", 2)
		if len(chunks) != 2 {
			return (Graph{}), fmt.Errorf("invalid node format")
		}

		node_from := newNode(chunks[0])
		node_to := newNode(chunks[1])

		children[node_from] = append(children[node_from], node_to)
		children[node_to] = append(children[node_to], node_from)
	}

	return Graph{children: children}, nil
}

type Queued struct {
	node  Node
	depth int
}

func (graph *Graph) forEachPath(do_not_repeat func(Node, []Node) bool, f func([]Node)) {
	path := make([]Node, 0)
	queue := []Queued{{node: newNode("start"), depth: 0}}

	for len(queue) > 0 {
		curr := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		// trim the path to make it match the depth of the current node
		depth := len(path)
		if depth > curr.depth {
			path = path[:curr.depth]
			depth = curr.depth
		}

		// append the current node to the path
		path = append(path, curr.node)
		depth += 1

		if curr.node.name == "end" {
			f(path)
		} else {
		children:
			for _, child := range graph.children[curr.node] {
				if do_not_repeat(child, path) {
					for _, node := range path {
						if child.name == node.name {
							continue children
						}
					}
				}

				queue = append(queue, Queued{node: child, depth: depth})
			}
		}
	}
}

func hasSmallDuplicates(path []Node) bool {
	for i, node := range path {
		if node.is_small {
			for j := i + 1; j < len(path); j += 1 {
				if node == path[j] {
					return true
				}
			}
		}
	}
	return false
}

func main() {
	bytes, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	input := string(bytes)
	graph, err := parseGraph(input)
	if err != nil {
		panic(err)
	}

	count_1 := 0
	graph.forEachPath(
		func(node Node, _ []Node) bool { return node.is_small },
		func(_ []Node) { count_1 += 1 },
	)
	fmt.Printf("Part 1: %d\n", count_1)

	count_2 := 0
	graph.forEachPath(
		func(node Node, path []Node) bool {
			return node.is_small && (node.name == "start" || node.name == "end" || hasSmallDuplicates(path))
		},
		func(_ []Node) { count_2 += 1 },
	)
	fmt.Printf("Part 2: %d\n", count_2)
}
