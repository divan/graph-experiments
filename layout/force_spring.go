package layout

// SpringForce calculates spring compression/extension force
// according to Hooke's law. Implements Force interface.
type SpringForce struct {
	Stiffness float64
	Length    float64 // each spring tends to have this length
}

var defaultSpringForce = NewSpringForce(0.011, 20.0)

// NewSpringForce creates and inits new spring force.
func NewSpringForce(stiffness, length float64) Force {
	return &SpringForce{
		Stiffness: stiffness,
		Length:    length,
	}
}

// Apply calculates the spring force between two nodes. Satisfies Force interface.
func (s *SpringForce) Apply(from, to *Point) *ForceVector {
	actualLength := distance(from, to)
	if actualLength == 0 {
		actualLength = s.Length
	}

	x := actualLength - s.Length // deformation distance
	c := s.Stiffness * float64(from.Mass) * x / actualLength

	return &ForceVector{
		DX: c * float64(to.X-from.X),
		DY: c * float64(to.Y-from.Y),
		DZ: c * float64(to.Z-from.Z),
	}
}

// Name returns name of the force. Satisifies Force interface.
func (s *SpringForce) Name() string {
	return "spring"
}
