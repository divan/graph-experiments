package p2p

import (
	"fmt"
	"log"
	"time"
)

// LogEntry defines the reporting log entry for one
// p2p message sending.
type LogEntry struct {
	From int
	To   int
	Ts   time.Duration
}

// String implements Stringer interface for LogEntry.
func (l LogEntry) String() string {
	return fmt.Sprintf("%s: %d -> %d", l.Ts.String(), l.From, l.To)
}

// NewLogEntry creates new log entry.
func NewLogEntry(start time.Time, from, to int) LogEntry {
	return LogEntry{
		Ts:   time.Since(start) / time.Millisecond,
		From: from,
		To:   to,
	}
}

// LogEntries2PropagationLog converts raw slice of LogEntries to PropagationLog,
// aggregating by timestamps and converting nodes indices to link indices.
// We expect that timestamps already bucketed into Nms groups.
func (s *Simulator) LogEntries2PropagationLog(entries []*LogEntry) *PropagationLog {
	findLink := func(from, to int) int {
		for i := range s.links {
			if s.links[i].From == from && s.links[i].To == to {
				return i
			}
		}
		return -1
	}

	tss := make(map[time.Duration][]int)
	for _, entry := range entries {
		idx := findLink(entry.From, entry.To)
		if idx == -1 {
			log.Println("[EE] Wrong link", entry)
			continue
		}

		if _, ok := tss[entry.Ts]; !ok {
			tss[entry.Ts] = make([]int, 0)
		}

		values := tss[entry.Ts]
		values = append(values, idx)
		tss[entry.Ts] = values
	}

	var ret = &PropagationLog{
		Timestamps: make([]int, 0, len(tss)),
		Indices:    make([][]int, 0, len(tss)),
	}

	for ts, links := range tss {
		ret.Timestamps = append(ret.Timestamps, int(ts))
		ret.Indices = append(ret.Indices, links)
		fmt.Println("Adding", ts*time.Millisecond, int(ts), links)
	}

	return ret
}