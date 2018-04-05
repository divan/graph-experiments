package p2p

import (
	"fmt"
	"sync"
	"time"

	"github.com/divan/graph-experiments/cmd/data_generator/net"
)

// Simulator is responsible for running propagation simulation.
type Simulator struct {
	data            *net.Data
	delay           time.Duration
	links           []LinkIndex
	peers           map[int][]int
	nodesCh         []chan Message
	reportCh        chan LogEntry
	peersToSendTo   int // number of peers to propagate message
	wg              sync.WaitGroup
	simulationStart time.Time
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
		reportCh:      make(chan LogEntry),
		nodesCh:       make([]chan Message, nodeCount), // one channel per node
	}
	sim.wg.Add(N)
	for i := 0; i < nodeCount; i++ {
		ch := sim.startNode(i)
		sim.nodesCh[i] = ch // this channel will be used to talk to node by index
	}
	return sim
}

func (s *Simulator) Run(startNodeIdx int) []*LogEntry {
	message := Message{
		Content: []byte("dummy"),
		TTL:     10,
	}
	s.simulationStart = time.Now()
	s.propagateMessage(startNodeIdx, message)

	done := make(chan bool)
	go func() {
		s.wg.Wait()
		done <- true
	}()

	var ret []*LogEntry
	for {
		select {
		case val := <-s.reportCh:
			ret = append(ret, &val)
		case <-done:
			return ret
		}
	}
	return ret
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

type LogEntry struct {
	From int
	To   int
	Ts   time.Duration
}

func (l LogEntry) String() string {
	return fmt.Sprintf("%s: %d -> %d", l.Ts.String(), l.From, l.To)
}

// sendMessage simulates message sending for given from and to indexes.
func (s *Simulator) sendMessage(from, to int, message Message) {
	s.nodesCh[to] <- message
	s.reportCh <- LogEntry{
		Ts:   time.Since(s.simulationStart) / time.Millisecond,
		From: from,
		To:   to,
	}
}
