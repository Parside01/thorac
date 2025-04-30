package transport

import (
	"context"
	raftpb "thorac/shared/proto/raft"
)

type gRPCHandler struct {
	raftpb.UnimplementedRaftServer
	handler RPCHandler
}

func (h *gRPCHandler) RequestVote(ctx context.Context, req *raftpb.RequestVoteRequest) (*raftpb.RequestVoteResponse, error) {
	return h.handler.HandleRequestVote(req), nil
}

func (h *gRPCHandler) AppendEntries(ctx context.Context, req *raftpb.AppendEntriesRequest) (*raftpb.AppendEntriesResponse, error) {
	return h.handler.HandleAppendEntries(req), nil
}
