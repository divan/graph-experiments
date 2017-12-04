package layout

import "math"

func distance(from, to *Point) float64 {
	dx := float64(to.X - from.X)
	dy := float64(to.Y - from.Y)
	dz := float64(to.Z - from.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
