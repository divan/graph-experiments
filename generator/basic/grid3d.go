package basic

import (
	"fmt"
	"math"

	"github.com/divan/graph-experiments/graph"
)

// Grid3DGenerator implements generator for 3D grid graph.
type Grid3DGenerator struct {
	rows   int
	cols   int
	levels int
}

// NewGrid3DGenerator creates new grid graph generator for known number of rows, cols and levels.
func NewGrid3DGenerator(rows, cols, levels int) *Grid3DGenerator {
	return &Grid3DGenerator{
		rows:   rows,
		cols:   cols,
		levels: levels,
	}
}

// NewGrid3DGeneratorN creates new grid 3D graph generator for N nodes.
func NewGrid3DGeneratorN(n int) *Grid3DGenerator {
	rows, cols, levels := estimateRowsColsLevels(n)
	return &Grid3DGenerator{
		rows:   rows,
		cols:   cols,
		levels: levels,
	}
}

// estimateRowsCols tries to find multiplies for n closest to cube.
// TODO: make it efficient and correct :)
func estimateRowsColsLevels(n int) (int, int, int) {
	root := math.Round(math.Cbrt(float64(n)))
	if root < 2 {
		return 1, 1, 1
	}

	return int(root), int(root), int(root)
}

// Generate generates the data for graph. Implements Generator interface.
func (l *Grid3DGenerator) Generate() *graph.Data {
	data := graph.NewData()

	for k := 0; k < l.levels; k++ {
		for i := 0; i < l.rows; i++ {
			for j := 0; j < l.cols; j++ {
				level := k * l.rows * l.cols
				idx := i + j*l.rows + level
				node := &graph.Node{
					ID: l.idxToID(idx),
				}
				data.Nodes = append(data.Nodes, node)

				if i > 0 {
					link := &graph.Link{
						Source: l.idxToID(idx),
						Target: l.idxToID(i - 1 + j*l.rows + level),
					}
					data.Links = append(data.Links, link)
				}
				if j > 0 {
					link := &graph.Link{
						Source: l.idxToID(idx),
						Target: l.idxToID(i + (j-1)*l.rows + level),
					}
					data.Links = append(data.Links, link)
				}
				if k > 0 {
					link := &graph.Link{
						Source: l.idxToID(idx),
						Target: l.idxToID(i + j*l.rows + (k-1)*l.rows*l.cols),
					}
					data.Links = append(data.Links, link)
				}
			}
		}
	}

	return data
}

func (l *Grid3DGenerator) idxToID(i int) string {
	return fmt.Sprintf("Node %d", i)
}
