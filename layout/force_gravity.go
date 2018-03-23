package layout

import "math"

// GravityForce represents gravity force,
// calculated by Coloumb's law. Implements Force interface.
type GravityForce struct {
	Coeff float64
}

var defaultGravity = &GravityForce{
	Coeff: -100.2, // repelling
}

// Apply calculates the gravity force between two nodes. Satisfies Force interface.
func (g *GravityForce) Apply(from, to *Point) *ForceVector {
	xx := float64(to.X - from.X)
	yy := float64(to.Y - from.Y)
	zz := float64(to.Z - from.Z)
	// distance calculates distance between points.
	r := int32(math.Sqrt(float64(xx*xx) + float64(yy*yy) + float64(zz*zz)))
	if r == 0 {
		r = 10
	}

	v := g.Coeff * float64(from.Mass*to.Mass) / float64(r*r*r)
	return &ForceVector{
		DX: (xx * v),
		DY: (yy * v),
		DZ: (zz * v),
	}
}
