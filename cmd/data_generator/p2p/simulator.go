package p2p

import (
	"log"
	"sync"
	"time"

	"github.com/divan/graph-experiments/cmd/data_generator/net"
)

// Simulator is responsible for running propagation simulation.
type Simulator struct {
	data          *net.Data
	delay         time.Duration
	links         []LinkIndex
	peers         map[int][]int
	nodesCh       []chan Message
	reportCh      chan int
	peersToSendTo int // number of peers to propagate message
	wg            sync.WaitGroup
}

// Message represents the message propagated in the simulation.
type Message struct {
	Content []byte
	TTL     int
}

// NewSimulator initializes new simulator.
func NewSimulator(data *net.Data, N int, delay time.Duration) *Simulator {
	nodeCount := len(data.Nodes)
	sim := &Simulator{
		data:          data,
		delay:         delay,
		peers:         PrecalculatePeers(data),
		peersToSendTo: N,
		reportCh:      make(chan int),
		nodesCh:       make([]chan Message, nodeCount), // one channel per node
	}
	sim.wg.Add(N)
	for i := 0; i < nodeCount; i++ {
		ch := sim.startNode(i)
		sim.nodesCh[i] = ch // this channel will be used to talk to node by index
	}
	return sim
}

func (s *Simulator) Run(startNodeIdx int) {
	message := Message{
		Content: []byte("dummy"),
		TTL:     10,
	}
	s.propagateMessage(startNodeIdx, message)
	s.wg.Wait()
}

func (s *Simulator) startNode(i int) chan Message {
	ch := make(chan Message)
	go s.runNode(i, ch)
	return ch
}

// runNode does actual node processing part
func (s *Simulator) runNode(i int, ch chan Message) {
	defer s.wg.Done()
	for message := range ch {
		log.Printf("Node %d received message %s with TTL %d", i, message.Content, message.TTL)
		message.TTL--
		if message.TTL == 0 {
			break
		}
		time.Sleep(s.delay)
		s.propagateMessage(i, message)
	}
}

// propagateMessage simulates message sending from node to its peers.
func (s *Simulator) propagateMessage(from int, message Message) {
	peers := s.peers[from]
	for i := range peers {
		go s.sendMessage(from, peers[i], message)
	}
}

// sendMessage simulates message sending for given from and to indexes.
func (s *Simulator) sendMessage(from, to int, message Message) {
	s.nodesCh[to] <- message
	// report sending here to be added to the log
}
