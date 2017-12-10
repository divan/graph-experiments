package layout

import "fmt"

type Force struct {
	Name string  `json:"name"`
	DX   float64 `json:"dx"`
	DY   float64 `json:"dy"`
	DZ   float64 `json:"dz"`
}

func (f Force) String() string {
	return fmt.Sprintf("f(%.03f, %.03f, %.03f)", f.DX, f.DY, f.DZ)
}

// Add adds new force to f.
func (f *Force) Add(f1 *Force) *Force {
	f.DX += f1.DX
	f.DY += f1.DY
	f.DZ += f1.DZ
	return f
}

// Sub substracts new force from f.
func (f *Force) Sub(f1 *Force) *Force {
	f.DX -= f1.DX
	f.DY -= f1.DY
	f.DZ -= f1.DZ
	return f
}
