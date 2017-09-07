package floodfill

import (
	"strings"
	"testing"
)

const (
	start = '@'
	floor = '.'
	fill  = 'x'
	wall  = '#'
)

type grid struct {
	tiles  [][]*tile
	starts []Node
}

func parseGrid(ascii string) *grid {
	rows := strings.Split(ascii, "\n")

	g := &grid{}
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

func (t *tile) Visit() {
	t.symbol = fill
}

func (t *tile) GetNeighbors() []Node {
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
	return neighbors
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
	} {
		g := parseGrid(testcase.initial)

		Floodfill(g.startingNodes())
		actual := g.toAscii()
		if testcase.expected != actual {
			t.Errorf(`
Expected: %s
Actual: %s`, testcase.expected, actual)
		}
	}
}