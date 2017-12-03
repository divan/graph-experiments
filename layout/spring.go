package layout

import "math"

func springForce(from, to *Node) *force {
	dx := float64(to.X - from.X)
	dy := float64(to.Y - from.Y)
	dz := float64(to.Z - from.Z)
	r := math.Sqrt(dx*dx + dy*dy + dz*dz)

	if r == 0 {
		r = 10
	}

	const (
		length = 20
		coeff  = 0.0008
		weight = 1
	)

	d := r - length
	c := coeff * d / r * weight

	return &force{
		dx: c * dx,
		dy: c * dy,
		dz: c * dz,
	}
}
