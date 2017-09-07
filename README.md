# floodfill

Parallel floodfill implementation for Go

[![Build Status](https://travis-ci.org/hinshun/floodfill.svg?branch=master)](https://travis-ci.org/hinshun/floodfill)
[![GoDoc](https://godoc.org/github.com/hinshun/floodfill?status.svg)](https://godoc.org/github.com/hinshun/floodfill)

The [Flood fill](https://en.wikipedia.org/wiki/Flood_fill) algorithm is used to find the connected components of a graph. For example, it is used to "bucket" fill areas of similarly colored areas of a paint program.

For graphs where edges are not known ahead of time, or have nodes that have to be retrieved over a network, flood fill greatly benefits from a parallelized implementation.

# Usage

## Implement `Node` interface

```go
type Tile struct {
  X int
  Y int

  // Node data
  ...
}

func (t *Tile) Visit() error {
  // Retrieve node data
  ...
}

func (t *Tile) GetNeighbors() ([]Node, error) {
  // Parse node data and return list of neighbors
  ...
}
```

## Call `Floodfill` function

```go
// We know the coordinates of the starting tiles but no metadata about the tile
// or about its neighbors.
tiles := []Node{
  &Tile{X: 12, Y: 5},
  &Tile{X: 6, Y: 9},
}

err := floodfill.Floodfill(tiles, 50)
if err != nil {
  // In cases where errors are intermittent, like API throttling, you can rerun
  // floodfill on errored nodes.
  floodfillErr := err.(floodfill.ErrFloodfill)
  var nodes []floodfill.Node
  for _, visit := range floodfillErr.Visits {
    nodes = append(nodes, visit.Node)
  }

  err = floodfill.Floodfill(nodes)
  if err != nil {
    return err
  }
}
```
