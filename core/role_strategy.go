package core

import (
	"context"
	"thorac/shared/types"
)

type NodeRoleStrategy interface {
	Run(ctx context.Context)
	Type() types.Role
	Enter(node *RaftNode)
	Leave(node *RaftNode)
}

type NodeEventType uint8

const (
	ElectionTimeout NodeEventType = iota
	HeartbeatTimeout
	ClientCommand
	RequestVote
	AppendLogEntry
)

type NodeEvent struct {
	Type NodeEventType
	Data interface{}
}

func NewNodeRoleStrategy(role types.Role) NodeRoleStrategy {
	switch role {
	case types.Candidate:
	case types.Follower:
	case types.Leader:
	}
	return nil
}
