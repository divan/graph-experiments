package naivep2p

import (
	"fmt"
	"sync"
	"time"

	"github.com/divan/graph-experiments/graph"
)

// Simulator is responsible for running propagation simulation.
type Simulator struct {
	data            *graph.Data
	delay           time.Duration
	links           []LinkIndex
	peers           map[int][]int
	nodesCh         []chan Message
	reportCh        chan LogEntry
	peersToSendTo   int // number of peers to propagate message
	wg              *sync.WaitGroup
	simulationStart time.Time
}

// Message represents the message propagated in the simulation.
type Message struct {
	Content string
	TTL     int
}

// NewSimulator initializes new simulator.
func NewSimulator(data *graph.Data, N int, delay time.Duration) *Simulator {
	nodeCount := len(data.Nodes)
	sim := &Simulator{
		data:          data,
		delay:         delay,
		links:         PrecalculateLinkIndexes(data),
		peers:         PrecalculatePeers(data),
		peersToSendTo: N,
		reportCh:      make(chan LogEntry),
		nodesCh:       make([]chan Message, nodeCount), // one channel per node
		wg:            new(sync.WaitGroup),
	}
	fmt.Println("[DD] Peers", sim.peers)
	fmt.Println("[DD] Links", sim.links)
	sim.wg.Add(nodeCount)
	for i := 0; i < nodeCount; i++ {
		ch := sim.startNode(i)
		sim.nodesCh[i] = ch // this channel will be used to talk to node by index
	}
	return sim
}

func (s *Simulator) Run(startNodeIdx, ttl int) []*LogEntry {
	message := Message{
		Content: "dummy",
		TTL:     ttl,
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
}

func (s *Simulator) startNode(i int) chan Message {
	ch := make(chan Message)
	go s.runNode(i, ch)
	return ch
}

// runNode does actual node processing part
func (s *Simulator) runNode(i int, ch chan Message) {
	defer s.wg.Done()
	t := time.NewTimer(4 * time.Second)

	cache := make(map[string]bool)
	for {
		select {
		case message := <-ch:
			if cache[message.Content] {
				continue
			}
			cache[message.Content] = true
			message.TTL--
			if message.TTL == 0 {
				return
			}
			s.propagateMessage(i, message)
		case <-t.C:
			return
		}
	}
}

// propagateMessage simulates message sending from node to its peers.
func (s *Simulator) propagateMessage(from int, message Message) {
	time.Sleep(s.delay)
	peers := s.peers[from]
	for i := range peers {
		go s.sendMessage(from, peers[i], message)
	}
}

// sendMessage simulates message sending for given from and to indexes.
func (s *Simulator) sendMessage(from, to int, message Message) {
	s.nodesCh[to] <- message
	s.reportCh <- NewLogEntry(s.simulationStart, from, to)
}
