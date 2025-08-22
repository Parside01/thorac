package main

import (
	"github.com/hashicorp/raft"
	"io"
)

type Fsm struct{}

func (f *Fsm) Apply(log *raft.Log) interface{} {
	return nil
}

func (f *Fsm) Snapshot() (raft.FSMSnapshot, error) {
	return &Snapshot{}, nil
}

func (f *Fsm) Restore(snapshot io.ReadCloser) error {
	return snapshot.Close()
}
