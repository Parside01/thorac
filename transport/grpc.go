package transport

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	raftpb "thorac/shared/proto/raft"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Transport interface {
	SendRequestVote(ctx context.Context, peerID string, req *raftpb.RequestVoteRequest) (*raftpb.RequestVoteResponse, error)
	SendAppendEntries(ctx context.Context, peerID string, req *raftpb.AppendEntriesRequest) (*raftpb.AppendEntriesResponse, error)

	RegisterHandler(handler RPCHandler)
	Start() error
	Shutdown() error
}

type RPCHandler interface {
	HandleRequestVote(req *raftpb.RequestVoteRequest) *raftpb.RequestVoteResponse
	HandleAppendEntries(req *raftpb.AppendEntriesRequest) *raftpb.AppendEntriesResponse
}

type RPCTransport struct {
	listenAddr string
	peerAddrs  map[string]string

	server *grpc.Server

	peerClients map[string]raftpb.RaftClient
	peerConns   map[string]*grpc.ClientConn

	rpcHandler RPCHandler

	isShutdown   bool
	shutdownChan chan struct{}
	wg           sync.WaitGroup
}

func NewGRPCTransport(listenAddr string, peerAddrs map[string]string) Transport {
	return &RPCTransport{
		listenAddr:   listenAddr,
		peerAddrs:    peerAddrs,
		peerClients:  make(map[string]raftpb.RaftClient),
		peerConns:    make(map[string]*grpc.ClientConn),
		shutdownChan: make(chan struct{}),
	}
}

func (t *RPCTransport) Start() error {
	if t.isShutdown {
		return fmt.Errorf("transport is already shutdown")
	}

	if t.rpcHandler == nil {
		return fmt.Errorf("gRPC handler is not registered")
	}

	lis, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %s", t.listenAddr, err.Error())
	}

	serverOpts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    1 * time.Minute,
			Timeout: 20 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	t.server = grpc.NewServer(serverOpts...)

	grpcHandler := &gRPCHandler{handler: t.rpcHandler}
	raftpb.RegisterRaftServer(t.server, grpcHandler)

	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		if err := t.server.Serve(lis); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	clientOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	for peerID, peerAddr := range t.peerAddrs {
		conn, err := grpc.Dial(peerAddr, clientOpts...)
		if err != nil {
			log.Printf("Warning: failed to dial peer %s at %s: %w", peerID, peerAddr, err.Error())
			continue
		}
		t.peerConns[peerID] = conn
		t.peerClients[peerID] = raftpb.NewRaftClient(conn)
		log.Printf("gRPC client created for peer %s at %s", peerID, peerAddr)
	}

	return nil
}

func (t *RPCTransport) RegisterHandler(handler RPCHandler) {
	t.rpcHandler = handler
	log.Printf("RPC handler registered for gRPC transport")
}

func (t *RPCTransport) SendRequestVote(ctx context.Context, peerID string, req *raftpb.RequestVoteRequest) (*raftpb.RequestVoteResponse, error) {
	if t.isShutdown {
		return nil, fmt.Errorf("transport is shut down")
	}

	client, ok := t.peerClients[peerID]
	if !ok {
		return nil, fmt.Errorf("client for peer %s not found", peerID)
	}

	resp, err := client.RequestVote(ctx, req)
	if err != nil {
		log.Printf("Error sending RequestVote to %s: %v", peerID, err)
		return nil, fmt.Errorf("send RequestVote to %s: %w", peerID, err)
	}

	return resp, nil
}

func (t *RPCTransport) SendAppendEntries(ctx context.Context, peerID string, req *raftpb.AppendEntriesRequest) (*raftpb.AppendEntriesResponse, error) {
	if t.isShutdown {
		return nil, fmt.Errorf("transport is shut down")
	}

	client, ok := t.peerClients[peerID]
	if !ok {
		return nil, fmt.Errorf("client for peer %s not found", peerID)
	}

	resp, err := client.AppendEntries(ctx, req)
	if err != nil {
		log.Printf("Error sending AppendEntries to %s: %v", peerID, err)
		return nil, fmt.Errorf("send AppendEntries to %s: %w", peerID, err)
	}

	return resp, nil
}

func (t *RPCTransport) Shutdown() error {
	if t.isShutdown {
		return fmt.Errorf("transport is already shut down")
	}
	t.isShutdown = true
	close(t.shutdownChan)

	if t.server != nil {
		t.server.GracefulStop()
		t.server = nil
	}

	for peerID, conn := range t.peerConns {
		if conn != nil {
			log.Printf("Closing connection to peer %s", peerID)
			conn.Close()
		}
	}
	t.peerClients = make(map[string]raftpb.RaftClient)
	t.peerConns = make(map[string]*grpc.ClientConn)

	t.wg.Wait()

	return nil
}
