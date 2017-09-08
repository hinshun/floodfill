package floodfill

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

const (
	start = '@'
	floor = '.'
	fill  = 'x'
	wall  = '#'
)

type grid struct {
	visitLock sync.Mutex
	visited   map[string]struct{}
	tiles     [][]*tile
	starts    []Node
}

func parseGrid(ascii string) *grid {
	rows := strings.Split(ascii, "\n")

	g := &grid{
		visited: make(map[string]struct{}),
	}
	offset := 0
	for y, row := range rows {
		if len(row) == 0 {
			offset++
			continue
		}
		var tileRow []*tile
		for x, symbol := range row {
			t := &tile{
				grid:   g,
				symbol: symbol,
				x:      x,
				y:      y - offset,
			}
			tileRow = append(tileRow, t)
			if symbol == start {
				g.starts = append(g.starts, t)
			}
		}
		g.tiles = append(g.tiles, tileRow)
	}

	return g
}

func (g *grid) startingNodes() []Node {
	return g.starts
}

func (g *grid) toAscii() string {
	var tiles []string
	tiles = append(tiles, "")
	for _, row := range g.tiles {
		var tileRow []byte
		for _, t := range row {
			tileRow = append(tileRow, byte(t.symbol))
		}
		tiles = append(tiles, string(tileRow))
	}
	return strings.Join(tiles, "\n")
}

type tile struct {
	grid   *grid
	symbol rune
	x      int
	y      int
}

func (t *tile) GetID() string {
	return fmt.Sprintf("(%d,%d)", t.x, t.y)
}

func (t *tile) Visit() error {
	t.grid.visitLock.Lock()
	defer t.grid.visitLock.Unlock()

	_, ok := t.grid.visited[t.GetID()]
	if ok {
		return fmt.Errorf("tile %s visited before", t.GetID())
	}
	t.grid.visited[t.GetID()] = struct{}{}

	t.symbol = fill
	return nil
}

func (t *tile) GetNeighbors() ([]Node, error) {
	var neighbors []Node
	if t.y-1 >= 0 && t.grid.tiles[t.y-1][t.x].symbol != wall {
		neighbors = append(neighbors, t.grid.tiles[t.y-1][t.x])
	}
	if t.y+1 < len(t.grid.tiles) && t.grid.tiles[t.y+1][t.x].symbol != wall {
		neighbors = append(neighbors, t.grid.tiles[t.y+1][t.x])
	}
	if t.x-1 >= 0 && t.grid.tiles[t.y][t.x-1].symbol != wall {
		neighbors = append(neighbors, t.grid.tiles[t.y][t.x-1])
	}
	if t.x+1 < len(t.grid.tiles[0]) && t.grid.tiles[t.y][t.x+1].symbol != wall {
		neighbors = append(neighbors, t.grid.tiles[t.y][t.x+1])
	}
	return neighbors, nil
}

func TestFloodfill(t *testing.T) {
	for _, testcase := range []struct {
		initial  string
		expected string
	}{
		{
			initial: `
@`,
			expected: `
x`,
		},
		{
			initial: `
@@`,
			expected: `
xx`,
		},
		{
			initial: `
.#
#@`,
			expected: `
.#
#x`,
		},
		{
			initial: `
...
.@.
...`,
			expected: `
xxx
xxx
xxx`,
		},
		{
			initial: `
.#.
.#@
.#.`,
			expected: `
.#x
.#x
.#x`,
		},
		{
			initial: `
.#.
@#@
.#.`,
			expected: `
x#x
x#x
x#x`,
		},
		{
			initial: `
#####
#...#
#.#.#
#.@.#
#####`,
			expected: `
#####
#xxx#
#x#x#
#xxx#
#####`,
		},
	} {
		g := parseGrid(testcase.initial)

		err := Floodfill(g.startingNodes(), 4)
		if err != nil {
			t.Errorf(`
Expected: %s
But error has occured: %s`, testcase.expected, err)
		}

		actual := g.toAscii()
		if testcase.expected != actual {
			t.Errorf(`
Expected: %s
Actual: %s`, testcase.expected, actual)
		}
	}
}
