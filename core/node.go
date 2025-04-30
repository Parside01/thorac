package core

import (
	"context"
	"go.uber.org/zap"
	"math/rand"
	"sync"
	raftpb "thorac/shared/proto/raft"
	"thorac/shared/types"
	"thorac/storage"
	"thorac/transport"
	"time"
)

type RaftNode struct {
	state storage.StateStorage
	log   storage.LogStorage

	leaderID string
	ID       string
	peers    []string

	lastAppliedIndex  uint64
	lastCommitedIndex uint64

	nextLogIndex   map[string]uint64
	commitLogIndex map[string]uint64

	currentTerm uint64
	votedFor    string

	strategy  RoleStrategy
	transport transport.Transport

	role      types.Role
	eventChan chan NodeEvent

	wg      sync.WaitGroup
	mutex   sync.RWMutex
	context context.Context
	cancel  context.CancelFunc

	lastElection time.Time

	//electionTimer  *time.Timer
	//heartbeatTimer *time.Timer
}

func NewRaftNode() (*RaftNode, error) {
	// TODO: Все значения брать из конфига или глобальных переменных.
	stateStorage, err := storage.NewBoltStateStorage("raft_data/data.db", "state")
	if err != nil {
		return nil, err
	}

	logStorage, err := storage.NewBoltLogStorage("raft_data/data.db", "log")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	peersAddrs := make(map[string]string)

	node := &RaftNode{
		wg:      sync.WaitGroup{},
		state:   stateStorage,
		log:     logStorage,
		context: ctx,
		cancel:  cancel,
		mutex:   sync.RWMutex{},

		leaderID: "",
		ID:       "1",
		peers:    make([]string, 0),

		nextLogIndex:   make(map[string]uint64),
		commitLogIndex: make(map[string]uint64),

		// TODO: Это потом из состояния наверное брать надо и из конфига.
		lastAppliedIndex:  0,
		lastCommitedIndex: 0,
		eventChan:         make(chan NodeEvent, 100),
		transport:         transport.NewGRPCTransport("localhost:8080", peersAddrs),
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

func (node *RaftNode) ResetLeader() {
	node.leaderID = ""
	node.nextLogIndex = make(map[string]uint64)
	node.commitLogIndex = make(map[string]uint64)

	if node.heartbeatTimer == nil {
		return
	}

	if !node.heartbeatTimer.Stop() {
		select {
		case <-node.heartbeatTimer.C:
		default:
		}
	}
	node.heartbeatTimer = nil

	//if node.electionTimer == nil {
	//  return
	//}
	//if !node.electionTimer.Stop() {
	//  select {
	//  case <-node.electionTimer.C:
	//  default:
	//  }
	//}
}

func (node *RaftNode) ResetHeartbeatTimer() {
	duration := node.randomDuration(10*time.Millisecond, 50*time.Millisecond)
	if node.heartbeatTimer == nil {
		return
	}
	if !node.heartbeatTimer.Stop() {
		select {
		case <-node.heartbeatTimer.C:
		default:
		}
	}
	node.heartbeatTimer.Reset(duration)
}

func (node *RaftNode) StartElectionTimer() {
	duration := node.randomDuration(20*time.Millisecond, 100*time.Millisecond)

	node.mutex.RLock()
	startedTerm := node.currentTerm
	node.mutex.RUnlock()

	ticker := time.NewTimer(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C

		node.mutex.RLock()
		if node.strategy.Type() != types.Candidate && node.strategy.Type() != types.Follower {
			L.Warn("The status has been changed, bailing out", zap.String("Node ID", node.ID), zap.String("state", node.strategy.Type().ToString()))
			node.mutex.RUnlock()
			return
		}

		if startedTerm != node.currentTerm {
			L.Warn("The term has been changed, bailing out", zap.String("Node ID", node.ID), zap.Uint64("current term", node.currentTerm), zap.Uint64("expected term", startedTerm))
			node.mutex.RUnlock()
			return
		}

		if elapsed := time.Since(node.lastElection); elapsed >= duration {
			node.mutex.RUnlock()
			node.StartElection()
			return
		}
	}
}

// TODO: Это по идее стратегия кандидата.
func (node *RaftNode) StartElection() {
	node.role = types.Candidate
	node.currentTerm++

	startedTerm := node.currentTerm
	node.votedFor = node.ID

	L.Debug("Node become to Candidate", zap.String("Node ID", node.ID))

	votesReceived := 1

	wg := sync.WaitGroup{}
	for _, peerID := range node.peers {
		wg.Add(1)
		go func(peerID string) {
			defer wg.Done()

			req := &raftpb.RequestVoteRequest{
				CandidateId:  node.ID,
				Term:         node.currentTerm,
				LastLogIndex: node.lastCommitedIndex,
			}
			resp, err := node.transport.SendRequestVote(node.context, peerID, req)
			if err != nil {
				L.Warn("Failed to send request vote", zap.String("Node ID", node.ID), zap.String("Peer", peerID), zap.Error(err))
				return
			}

			node.mutex.Lock()
			defer node.mutex.Unlock()

			if node.strategy.Type() != types.Candidate {
				L.Warn("The status has been changed, bailing out", zap.String("Node ID", node.ID), zap.String("state", node.strategy.Type().ToString()))
				return
			}

			if resp.Term > startedTerm {
				L.Warn("Term out of date, bailing out", zap.String("Node ID", node.ID), zap.Uint64("Response Term", resp.Term), zap.Uint64("Started Term", startedTerm))
				node.TransitionToRole(types.Follower)
				return
			}

			// Нам не нужен голос от ноды, которая с неполным журналом.
			if resp.Term < node.currentTerm {
				return
			}

			if resp.IsGranted {
				votesReceived++
			}
		}(peerID)
	}
	wg.Wait()

	if votesReceived*2 > len(node.peers)+1 {
		L.Info("Won election", zap.String("Node ID", node.ID), zap.Int("Votes", votesReceived))
		node.TransitionToRole(types.Leader)
		return
	}

	node.TransitionToRole(types.Follower)
}

//func (node *RaftNode) ResetElectionTimer() {
//  duration := node.randomDuration(20*time.Millisecond, 100*time.Millisecond)
//  if node.electionTimer == nil {
//    return
//  }
//  if !node.electionTimer.Stop() {
//    select {
//    case <-node.electionTimer.C:
//    default:
//    }
//  }
//  node.electionTimer.Reset(duration)
//}

func (node *RaftNode) randomDuration(min, max time.Duration) time.Duration {
	rangeNanos := max.Nanoseconds() - min.Nanoseconds()
	if rangeNanos <= 0 {
		return min
	}
	return min + time.Duration(rand.Int63n(rangeNanos))
}

func (node *RaftNode) Shutdown() {
	node.cancel()
	node.wg.Wait()

	if node.log != nil {
		node.log.Close()
	}
	if node.log != nil {
		node.state.Close()
	}
}
