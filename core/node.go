package core

import (
	"context"
	"sync"
	"thorac/shared/types"
	"thorac/storage"
	"time"
)

type RaftNode struct {
	state storage.StateStorage
	log   storage.LogStorage

	lastAppliedIndex  uint64
	lastCommitedIndex uint64

	currentTerm uint64
	votedFor    string

	strategy  NodeRoleStrategy
	eventChan chan NodeEvent

	wg      sync.WaitGroup
	context context.Context

	electionTimer  *time.Timer
	heartbeatTimer *time.Timer
}

func NewRaftNode() (*RaftNode, error) {
	stateStorage, err := storage.NewBoltStateStorage("raft_data/data.db", "state")
	if err != nil {
		return nil, err
	}

	logStorage, err := storage.NewBoltLogStorage("raft_data/data.db", "log")
	if err != nil {
		return nil, err
	}

	node := &RaftNode{
		wg:      sync.WaitGroup{},
		state:   stateStorage,
		log:     logStorage,
		context: context.Background(),

		// TODO: Это потом из состояния наверное брать надо.
		lastAppliedIndex:  0,
		lastCommitedIndex: 0,
	}

	if err = node.loadState(); err != nil {
		node.Shutdown()
		return nil, err
	}

	node.strategy = NewNodeRoleStrategy(types.Follower)

	return node, nil
}

func (node *RaftNode) TransitionToRole(newRole types.Role) {
	if node.strategy != nil && node.strategy.Type() == newRole {
		return
	}

	if node.strategy != nil {
		node.strategy.Leave(node)
	}

	newStrategy := NewNodeRoleStrategy(newRole)
	node.strategy = newStrategy
	newStrategy.Enter(node)

	node.wg.Add(1)
	go newStrategy.Run(node.context)
}

func (node *RaftNode) loadState() error {
	term, err := node.state.GetTerm()
	if err != nil {
		return err
	}
	node.currentTerm = term

	votedFor, err := node.state.GetVotedFor()
	if err != nil {
		return err
	}
	node.votedFor = votedFor

	return nil
}

func (node *RaftNode) saveState() error {
	if err := node.state.SetTerm(node.currentTerm); err != nil {
		return err
	}
	if err := node.state.SetVotedFor(node.votedFor); err != nil {
		return err
	}
	return nil
}

func (node *RaftNode) Shutdown() {
	node.wg.Wait()
	node.log.Close()
	node.state.Close()
}
