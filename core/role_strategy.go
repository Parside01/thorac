package core

import (
	"context"
	"thorac/shared/types"
)

type RoleStrategy interface {
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
	Type         NodeEventType
	ResponseChan chan interface{}
	Data         interface{}
}

func NewNodeRoleStrategy(role types.Role) RoleStrategy {
	switch role {
	case types.Candidate:
	case types.Follower:
	case types.Leader:
	}
	return nil
}
