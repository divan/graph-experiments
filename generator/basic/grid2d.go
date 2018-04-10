package basic

import (
	"math"

	"github.com/divan/graph-experiments/graph"
)

// Grid2DGenerator implements generator for 2D grid graph.
type Grid2DGenerator struct {
	rows int
	cols int
}

// NewGrid2DGenerator creates new grid graph generator for known number of rows and cols.
func NewGrid2DGenerator(rows, cols int) *Grid2DGenerator {
	return &Grid2DGenerator{
		rows: rows,
		cols: cols,
	}
}

// NewGrid2DGeneratorN creates new graph generator for N nodes.
func NewGrid2DGeneratorN(n int) *Grid2DGenerator {
	rows, cols := estimateRowsCols(n)
	return &Grid2DGenerator{
		rows: rows,
		cols: cols,
	}
}

// estimateRowsCols tries to find multiplies for n closest to square.
// TODO: make it efficient and correct :)
func estimateRowsCols(n int) (int, int) {
	root := math.Round(math.Sqrt(float64(n)))
	if root < 2 {
		return 1, 1
	}

	res := math.Mod(root, float64(n))
	if math.Round(res) < 2 {
		return int(root), 1
	}

	return int(root), int(res)
}

// Generate generates the data for graph. Implements Generator interface.
func (l *Grid2DGenerator) Generate() *graph.Data {
	data := graph.NewData()

	for i := 0; i < l.rows; i++ {
		for j := 0; j < l.cols; j++ {
			idx := i + j*l.rows
			addNode(data, idx)

			if i > 0 {
				addLink(data, idx, i-1+j*l.rows)
			}
			if j > 0 {
				addLink(data, idx, i+(j-1)*l.rows)
			}
		}
	}

	return data
}
