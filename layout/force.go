package layout

import "fmt"

type ForceDebugInfo struct {
	Name string
	ForceVector
}

type ForceVector struct {
	DX float64 `json:"dx"`
	DY float64 `json:"dy"`
	DZ float64 `json:"dz"`
}

type Force interface {
	Name() string
	Apply(from, to *Point) *ForceVector
}

func (f ForceVector) String() string {
	return fmt.Sprintf("f(%.03f, %.03f, %.03f)", f.DX, f.DY, f.DZ)
}

// Add adds new force to f.
func (f *ForceVector) Add(f1 *ForceVector) *ForceVector {
	f.DX += f1.DX
	f.DY += f1.DY
	f.DZ += f1.DZ
	return f
}

// Sub substracts new force from f.
func (f *ForceVector) Sub(f1 *ForceVector) *ForceVector {
	f.DX -= f1.DX
	f.DY -= f1.DY
	f.DZ -= f1.DZ
	return f
}
