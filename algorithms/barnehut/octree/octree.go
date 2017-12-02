package octree

// Octree represents Octree data structure.
// See https://en.wikipedia.org/wiki/Octree for details.
type Octree struct {
	root octant
}

// Point represents 3D point with mass, that'd be used
// to calculate center of the mass of octants.
type Point struct {
	X, Y, Z int32
	Mass    int32
}

// octant represent a node in octree, which is an octant of a cube.
// See: http://en.wikipedia.org/wiki/Octant_(solid_geometry)
type octant interface {
	Center() *Point
	Insert(p *Point) octant
}

// node represents octant with children, "internal node". Satisifies octant.
type node struct {
	leafs      *[8]octant
	massCenter *Point
}

// Center returns center of the mass of the node. Implements octant interface.
func (n *node) Center() *Point {
	return n.massCenter
}

// make sure node satisfies octant interface at compile time.
var _ = octant(&node{})

// leaf represents octant without children, "external node". Satisfies octant.
type leaf struct {
	point *Point
}

// Center returns point of the leaf. Implements octant interface.
func (l *leaf) Center() *Point {
	return l.point
}

// make sure leaf satisfies octant interface at compile time.
var _ = octant(&leaf{})

// New inits new octree.
func New() *Octree {
	return &Octree{}
}

// newNode initializes a new node.
func newNode() *node {
	var leafs [8]octant
	for i := 0; i < 8; i++ {
		leafs[i] = newLeaf(nil)
	}
	return &node{
		leafs: &leafs,
	}
}

// newLeaf initializes a new leaf.
func newLeaf(p *Point) *leaf {
	return &leaf{
		point: p,
	}
}

// Insert adds new Point into the Octree data structure.
func (o *Octree) Insert(p *Point) {
	if o.root == nil {
		o.root = newLeaf(p)
		return
	}

	o.root = o.root.Insert(p)
}

// Insert inserts new Point into existing node and returns
// updated node. Implements octant interface.
func (n *node) Insert(p *Point) octant {
	idx := findOctantIdx(n, p)
	child := n.leafs[idx]
	child.Insert(p)
	n.updateMassCenter()
	return n
}

// Insert inserts new Point into existing leaf and returns updated
// node, which may be transformed into node. Implements octant interface.
func (l *leaf) Insert(p *Point) octant {
	if l == nil || l.Center() == nil {
		l = newLeaf(p)
		return l
	}

	//external node, and we have two points in one octant.
	//need to convert it to internal node and divide
	node := newNode()
	node.massCenter = l.Center()
	idx1 := findOctantIdx(node, l.Center())
	idx2 := findOctantIdx(node, p)
	node.leafs[idx1] = l
	node.leafs[idx2] = newLeaf(p)
	node.updateMassCenter()
	return node
}

// update center of the mass of the given node, calculating it from
// leaf centers of the mass.
func (n *node) updateMassCenter() {
	var (
		p          = &Point{}
		xm, ym, zm int64
	)

	for _, leaf := range n.leafs {
		if leaf == nil || leaf.Center() == nil {
			continue
		}
		c := leaf.Center()
		p.Mass += c.Mass
		xm += int64(c.X) * int64(c.Mass)
		ym += int64(c.Y) * int64(c.Mass)
		zm += int64(c.Z) * int64(c.Mass)
	}

	p.X = int32(xm / int64(p.Mass))
	p.Y = int32(ym / int64(p.Mass))
	p.Z = int32(zm / int64(p.Mass))

	n.massCenter = p
}

// findOctantIdx returns index of 8-length array with children of the
// given octant. It's in following order:
// 0 - Top, Front, Left
// 1 - Top, Front, Right
// 2 - Top, Back, Left
// 3 - Top, Back, Right
// 4 - Bottom, Front, Left
// 5 - Bottom, Front, Right
// 6 - Bottom, Back, Left
// 7 - Bottom, Back, Right
func findOctantIdx(o octant, p *Point) int {
	center := o.Center()

	var i int
	if p.X > center.X {
		i |= 1
	}

	if p.Y > center.Y {
		i |= 2
	}

	if p.Z > center.Z {
		i |= 4
	}
	return i
}
