/*
Package floodfill provides an implementation of a parallel flood fill algorithm
of lazily loaded nodes.
*/
package floodfill

import (
	"strings"
	"sync"
)

// Node is a node in a directed graph which are 'filled' as they are visited.
type Node interface {
	// Visit marks the node as visited, allowing the node to be lazily loaded
	// from an external source.
	Visit() error

	// GetNeighbors retrieves the nodes that are directly connected with the node.
	GetNeighbors() ([]Node, error)
}

// ErrFloodfill is returned if any of the visited nodes returned an error
// upon visited or getting its neighbors.
type ErrFloodfill struct {
	Visits []ErrVisit
}

// Error returns the joined error strings of its visited nodes that errored
// during floodfill.
func (e ErrFloodfill) Error() string {
	var errs []string
	for _, visit := range e.Visits {
		errs = append(errs, visit.Error())
	}
	return strings.Join(errs, ", ")
}

// ErrVisit is a wrapper of the node visited and the error returned from
// either visiting the node or getting its neighbors.
type ErrVisit struct {
	Node Node
	Err  error
}

// Error returns the error of the visited node.
func (e ErrVisit) Error() string {
	return e.Err.Error()
}

type floodfiller struct {
	wg         sync.WaitGroup
	permitCh   chan struct{}
	errCh      chan ErrVisit
	visitLock  sync.Mutex
	visitQueue []Node
	visited    map[Node]struct{}
}

// Floodfill determines the areas connected to a given list of nodes.
// Neighboring nodes are all visited in parallel, which is particularly useful
// if visiting and computing the node's neighbors is expensive or latency
// bound. The number of goroutines spawned by floodfill can be limited by
// the given parallelism limit.
func Floodfill(nodes []Node, parallelism int) error {
	f := &floodfiller{
		errCh:    make(chan ErrVisit),
		permitCh: make(chan struct{}, parallelism),
		visited:  make(map[Node]struct{}),
	}
	defer close(f.permitCh)

	for i := 0; i < parallelism; i++ {
		f.permitCh <- struct{}{}
	}

	for _, node := range nodes {
		f.enqueue(node)
	}

	var errs []ErrVisit
	go func() {
		for err := range f.errCh {
			errs = append(errs, err)
		}
	}()

	f.wg.Wait()
	close(f.errCh)

	if len(errs) > 0 {
		return ErrFloodfill{errs}
	}

	return nil
}

func (f *floodfiller) enqueue(node Node) {
	f.visitLock.Lock()
	defer f.visitLock.Unlock()

	// Add node to visit queue and fire off goroutine to visit.
	f.visitQueue = append(f.visitQueue, node)
	f.wg.Add(1)
	go func() {
		err := f.visitNext()
		if err != nil {
			f.errCh <- ErrVisit{
				Node: node,
				Err:  err,
			}
		}
	}()
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

func (f *floodfiller) visitNext() error {
	defer f.wg.Done()

	node, ok := f.dequeue()
	if ok {
		return nil
	}

	// Wait for a parallelism permit to perform work.
	<-f.permitCh
	defer func() {
		f.permitCh <- struct{}{}
	}()

	err := node.Visit()
	if err != nil {
		return err
	}

	neighbors, err := node.GetNeighbors()
	if err != nil {
		return err
	}

	for _, neighbor := range neighbors {
		f.enqueue(neighbor)
	}
	return nil
}
