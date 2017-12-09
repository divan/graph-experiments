package layout

import "math"

const (
	springStiffness = 0.011
	springLength    = 20 // each spring tends to have this length

	gravityConst = -1.2 // coulumb's coeff, negative, thus nodes repel
)

// springForce calculates spring compression/extension force
// according to Hooke's law.
func springForce(from, to *Node) *force {
	actualLength := distance(from.Point, to.Point)
	if actualLength == 0 {
		actualLength = springLength
	}

	x := actualLength - springLength // deformation distance
	c := springStiffness * float64(from.Mass) * x / actualLength

	return &force{
		dx: c * float64(to.X-from.X),
		dy: c * float64(to.Y-from.Y),
		dz: c * float64(to.Z-from.Z),
	}
}

func gravity(from, to *Point) *force {
	xx := float64(to.X - from.X)
	yy := float64(to.Y - from.Y)
	zz := float64(to.Z - from.Z)
	// distance calculates distance between points.
	r := int32(math.Sqrt(float64(xx*xx) + float64(yy*yy) + float64(zz*zz)))
	if r == 0 {
		r = 20
	}

	v := gravityConst * float64(from.Mass*to.Mass) / float64(r*r*r)
	return &force{
		dx: (xx * v),
		dy: (yy * v),
		dz: (zz * v),
	}
}
