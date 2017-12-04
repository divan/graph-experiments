package layout

import (
	"math"
)

// integrate performs forces integration using Euler numerical
// integration method.
func (l *Layout3D) integrate(forces []*force) {
	const timeStep = float64(20) // FIXME: 20 what?
	for i := 0; i < len(l.nodes); i++ {
		body := l.nodes[i]
		coeff := timeStep / float64(body.Mass)

		vx := coeff * forces[i].dx
		vy := coeff * forces[i].dy
		vz := coeff * forces[i].dz
		v := math.Sqrt(vx*vx + vy*vy + vz*vz)

		l.nodes[i].VelocityX += vx
		l.nodes[i].VelocityY += vy
		l.nodes[i].VelocityZ += vz

		if v > 1 {
			vx = vx / v
			vy = vy / v
			vz = vz / v
			l.nodes[i].VelocityX = vx
			l.nodes[i].VelocityY = vy
			l.nodes[i].VelocityZ = vz
		}

		dx := timeStep * vx
		dy := timeStep * vy
		dz := timeStep * vz

		l.nodes[i].X += int32(dx)
		l.nodes[i].Y += int32(dy)
		l.nodes[i].Z += int32(dz)
	}
}
