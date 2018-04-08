package graph

// NodeHasLinks implements fast check if given node has any links/
func (d *Data) NodeHasLinks(nodeID string) bool {
	if d.nodeHasLinks == nil {
		d.prepare()
	}

	return d.nodeHasLinks[nodeID]
}

// LinkExists returns true if there is a link between source and target node IDs.
func (d *Data) LinkExists(source, target string) bool {
	for _, link := range d.Links {
		if link.Source == source && link.Target == target ||
			link.Target == source && link.Source == target {
			return true
		}
	}
	return false
}
