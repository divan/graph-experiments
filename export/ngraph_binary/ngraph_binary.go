package ngraph_binary

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/divan/graph-experiments/graph"
	"github.com/divan/graph-experiments/layout"
)

// NgraphBinaryOutput stores graph data as binary files, compatible
// with ngraph library: positions.bin, links.bin, labels.json and meta.json
type NgraphBinaryOutput struct {
	dir string
}

func NewExporter(dir string) *NgraphBinaryOutput {
	fs, err := os.Stat(dir)
	if err != nil {
		log.Fatalf("Failed to prepare output dir: %v", err)
	}
	if !fs.IsDir() {
		log.Fatalf("'%s' is not a dir, aborting...", dir)
	}
	return &NgraphBinaryOutput{
		dir: dir,
	}
}

func (o *NgraphBinaryOutput) Save(l layout.Layout, data *graph.Data) error {
	err := o.WritePositionsBin(l)
	if err != nil {
		return err
	}

	err = o.WriteLinksBin(data)
	if err != nil {
		return err
	}

	return o.WriteLabels(data)
}

// WritePositionsBin writes position data into 'positions.bin' file in the
// following way: XYZXYZXYZ... where X, Y and Z are coordinates
// for each node in signed 32 bit integer Little Endian format.
func (o *NgraphBinaryOutput) WritePositionsBin(l layout.Layout) error {
	file := filepath.Join(o.dir, "positions.bin")
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	iw := NewInt32LEWriter(fd)

	nodes := l.Nodes()
	for i, _ := range nodes {
		iw.Write(nodes[i].X)
		iw.Write(nodes[i].Y)
		iw.Write(nodes[i].Z)
		if iw.err != nil {
			return err
		}
	}

	return nil
}

// WriteLinksBin writes links information into `links.bin` file in the
// following way: Sidx,L1idx,L2idx,S2idx,L1idx... where SNidx - is the
// start node index, and LNidx - is the other link end node index.
func (o *NgraphBinaryOutput) WriteLinksBin(data *graph.Data) error {
	file := filepath.Join(o.dir, "links.bin")
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	iw := NewInt32LEWriter(fd)
	for i, node := range data.Nodes {
		if !data.NodeHasLinks(node.ID) {
			continue
		}

		iw.Write(int32(-(i + 1)))
		for _, link := range data.Links {
			if link.Source == node.ID {
				iw.Write(int32(link.ToIdx + 1))
			}
		}
		if iw.err != nil {
			return err
		}
	}
	return nil
}

// WriteLabels writes node ids (labels) information into `labels.json` file
// as an array of strings.
func (o *NgraphBinaryOutput) WriteLabels(data *graph.Data) error {
	file := filepath.Join(o.dir, "labels.json")
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	var labels []string
	for i, _ := range data.Nodes {
		labels = append(labels, data.Nodes[i].ID)
	}
	return json.NewEncoder(fd).Encode(labels)
}

type Int32LEWriter struct {
	w   io.Writer
	err error
}

func NewInt32LEWriter(w io.Writer) *Int32LEWriter {
	return &Int32LEWriter{
		w: w,
	}
}

func (iw *Int32LEWriter) Write(number int32) {
	if iw.err != nil {
		return
	}

	err := binary.Write(iw.w, binary.LittleEndian, number)
	if err != nil {
		iw.err = err
	}

	return
}
