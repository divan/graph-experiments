package layout

// SpringForce calculates spring compression/extension force
// according to Hooke's law. Implements Force interface.
type SpringForce struct {
	Stiffness float64
	Length    int64 // each spring tends to have this length
}

var defaultSpringForce = &SpringForce{
	Stiffness: 0.011,
	Length:    20,
}

// Apply calculates the spring force between two nodes. Satisfies Force interface.
func (s *SpringForce) Apply(from, to *Node) *ForceVector {
	actualLength := distance(from.Point, to.Point)
	if actualLength == 0 {
		actualLength = springLength
	}

	x := actualLength - springLength // deformation distance
	c := springStiffness * float64(from.Mass) * x / actualLength

	return &ForceVector{
		DX: c * float64(to.X-from.X),
		DY: c * float64(to.Y-from.Y),
		DZ: c * float64(to.Z-from.Z),
	}
}
