package main

import (
	"fmt"
)

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.
	return map[int][]int{}
}

func main() {
	// You can insert optional local tests here if desired.
	graph := map[int][]int{
		0: {1, 2},
		1: {0, 3, 4},
		2: {0, 5},
		3: {1, 5},
		4: {1, 5},
		5: {2, 3, 6, 7},
		6: {5},
		7: {5},
		8: {5},
		9: {8},
	}
	_, path := standardBFS(graph, 0, 8)
	fmt.Println(path)
}

func standardBFS(graph map[int][]int, start int, goal int) (bool, map[int][]int) {
	queue := []int{start}
	path := make(map[int][]int)
	visited := make(map[int]bool)

	// init start node infomation
	visited[start] = true
	path[start] = []int{start}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node == goal {
			return true, path
		}

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true

				// * note: output : [xx, xx, xx, xx, <not have goal>]
				// path[neighbor] = append(path[node], node)

				// * note: output : [xx, xx, xx, xx, <have goal>]
				newPath := append(append([]int(nil), path[node]...), neighbor)
				path[neighbor] = newPath

				queue = append(queue, neighbor)
			}
		}
	}
	return false, path
}
