package main

import (
	"fmt"
	"sync"
)

type BFSResult struct {
	StartNode int
	Traversal []int
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	jobChan := make(chan int, len(queries))
	resultChan := make(chan BFSResult, len(queries))

	var wg sync.WaitGroup

	for i := range numWorkers {
		wg.Add(1)

		go func(workID int) {
			defer wg.Done()
			for startNode := range jobChan {
				bfsTraversal := BFS(graph, startNode)
				resultChan <- BFSResult{
					StartNode: startNode,
					Traversal: bfsTraversal,
				}
			}
		}(i)
	}

	for _, query := range queries {
		jobChan <- query
	}
	close(jobChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	finalResults := make(map[int][]int)
	for result := range resultChan {
		finalResults[result.StartNode] = result.Traversal
	}

	return finalResults
}

// BFS performs breadth-first search starting from the given node
func BFS(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	result := []int{}

	// start node initformation
	visited[start] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		if neighbors, exist := graph[node]; exist {
			for _, neighbor := range neighbors {
				if !visited[neighbor] {
					visited[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}
	}

	return result
}

func main() {
	// You can insert optional local tests here if desired.
	graph := map[int][]int{
		0: {1, 2},
		1: {0, 3, 4},
		2: {0, 5},
		3: {1, 5},
		4: {1, 5},
		5: {2, 3, 6, 7, 8, 9},
		6: {5},
		7: {5},
		8: {5},
		9: {8},
	}
	_, path := standardBFS(graph, 5, 0)
	fmt.Println(path)
	_, path = standardBFS(graph, 0, 5)
	// path := BFS(graph, 5)
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

		// if node == goal {
		// 	return true, path
		// }

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
