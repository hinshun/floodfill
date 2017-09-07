package floodfill

import (
	"sync"
)

// Node is a node in a directed graph which are 'filled' as they are visited.
type Node interface {
	// Visit marks the node as visited, allowing the node to be lazily loaded
	// from an external source.
	Visit()
	// GetNeighbors retrieves the nodes that are directly connected with the node.
	GetNeighbors() []Node
}

type floodfiller struct {
	wg         sync.WaitGroup
	visitLock  sync.Mutex
	visitQueue []Node
	visited    map[Node]struct{}
}

// Floodfill determines the areas connected to a given list of nodes.
// Neighboring nodes are all visited in parallel, which is particularly useful
// if visiting and computing the node's neighbors is expensive or latency
// bound.
func Floodfill(nodes []Node) {
	f := &floodfiller{
		visited: make(map[Node]struct{}),
	}

	for _, node := range nodes {
		f.enqueue(node)
	}
	f.wg.Wait()
}

func (f *floodfiller) enqueue(node Node) {
	f.visitLock.Lock()
	defer f.visitLock.Unlock()

	// Add node to visit queue and fire off goroutine to visit.
	f.visitQueue = append(f.visitQueue, node)
	f.wg.Add(1)
	go f.visitNext()
}

func (f *floodfiller) dequeue() (Node, bool) {
	f.visitLock.Lock()
	defer f.visitLock.Unlock()

	// Dequeue the next node.
	node := f.visitQueue[0]
	f.visitQueue = f.visitQueue[1:]

	// Check whether node has been visited before and mark as visited.
	_, ok := f.visited[node]
	f.visited[node] = struct{}{}

	return node, ok
}

func (f *floodfiller) visitNext() {
	defer f.wg.Done()
	node, ok := f.dequeue()
	if ok {
		return
	}

	node.Visit()
	neighbors := node.GetNeighbors()
	for _, neighbor := range neighbors {
		f.enqueue(neighbor)
	}
}
