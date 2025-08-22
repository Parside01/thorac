package main

import "github.com/hashicorp/raft"

type Snapshot struct{}

func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	return sink.Close()
}

func (s *Snapshot) Release() {}
