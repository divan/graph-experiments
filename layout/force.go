package layout

import "fmt"

type force struct {
	dx, dy, dz float64
}

func (f force) String() string {
	return fmt.Sprintf("f(%.03f, %.03f, %.03f)", f.dx, f.dy, f.dz)
}

// Add adds new force to f.
func (f *force) Add(f1 *force) *force {
	f.dx += f1.dx
	f.dy += f1.dy
	f.dz += f1.dz
	return f
}

// Sub substracts new force from f.
func (f *force) Sub(f1 *force) *force {
	f.dx -= f1.dx
	f.dy -= f1.dy
	f.dz -= f1.dz
	return f
}
