package layout

import (
	"math"
)

// integrate performs forces integration using Euler numerical
// integration method.
//
// F = d(m * v) / dt
//  (mass is constant in our case)
// v = d{x,y,z}/dt
//
// dv = dt * F / m
//
// d{x,y,z} = v * dt
func (l *Layout3D) integrate() {
	const dt = float64(3)
	for i := 0; i < len(l.nodes); i++ {
		body := l.nodes[i]
		force := l.forceVectors[i]
		coeff := dt / float64(body.Mass)

		if force == nil {
			force = &ForceVector{}
		}

		body.Velocity.X += coeff * force.DX
		body.Velocity.Y += coeff * force.DY
		body.Velocity.Z += coeff * force.DZ
		dvx, dvy, dvz := body.Velocity.X, body.Velocity.Y, body.Velocity.Z
		v := math.Sqrt(dvx*dvx + dvy*dvy + dvz*dvz)

		if v > 1 {
			body.Velocity.X = dvx / v
			body.Velocity.Y = dvy / v
			body.Velocity.Z = dvz / v
		}

		dx := dt * body.Velocity.X
		dy := dt * body.Velocity.Y
		dz := dt * body.Velocity.Z

		l.nodes[i].X += int32(dx)
		l.nodes[i].Y += int32(dy)
		l.nodes[i].Z += int32(dz)
	}
}
