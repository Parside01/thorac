package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"go.uber.org/zap"
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrWaitNodeTimeout = errors.New("timed out waiting for node")
)

type Node interface {
	Join(id, address string) error
	Health() map[string]string
	Close() error
}

type node struct {
	raft   *raft.Raft
	store  raft.LogStore
	snap   raft.SnapshotStore
	trans  raft.Transport
	logger *zap.Logger
	config *Config
}

func NewNode(config *Config, logger *zap.Logger, fsm raft.FSM) (Node, error) {
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(config.Id)

	logger.Info("initializing raft node")

	store, err := raftboltdb.NewBoltStore(filepath.Join(config.BaseDataDir, "raft.db"))
	if err != nil {
		logger.Error("failed to create raft log store", zap.Error(err))
		return nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStore(filepath.Join(config.BaseDataDir, "snapshots"), 3, os.Stderr)
	if err != nil {
		logger.Error("failed to create snapshot store", zap.Error(err))
		return nil, err
	}

	advertiseAddr, err := net.ResolveTCPAddr("tcp", config.Advertise)
	if err != nil {
		logger.Error("invalid advertise address", zap.Error(err))
		return nil, err
	}

	transportConfig := &raft.NetworkTransportConfig{
		MaxPool: 3,
		Timeout: 1 * time.Second,
	}

	transport, err := raft.NewTCPTransportWithConfig(
		config.Address,
		advertiseAddr,
		transportConfig,
	)
	if err != nil {
		logger.Error("failed to create raft transport", zap.Error(err))
		return nil, err
	}

	n, err := raft.NewRaft(raftConfig, fsm, store, store, snapshotStore, transport)
	if err != nil {
		logger.Error("failed to create raft instance", zap.Error(err))
		return nil, err
	}

	newNode := &node{
		raft:   n,
		store:  store,
		snap:   snapshotStore,
		trans:  transport,
		logger: logger,
		config: config,
	}

	if err = newNode.JoinPeers(config.Peers); err != nil {
		logger.Error("failed to join peers", zap.String("node_id", config.Id), zap.Error(err))
		return nil, err
	}

	if !config.IsLeader {
		logger.Info("started follower node", zap.String("node_id", config.Id))
		return newNode, nil
	}

	logger.Info("attempting to bootstrap cluster", zap.String("node_id", config.Id))
	configFuture := n.GetConfiguration()
	if err = configFuture.Error(); err != nil {
		logger.Error("failed to get raft configuration", zap.Error(err))
		return nil, err
	}

	if len(configFuture.Configuration().Servers) == 0 {
		clusterConfig := raft.Configuration{
			Servers: []raft.Server{{
				ID:      raft.ServerID(config.Id),
				Address: raft.ServerAddress(config.Address),
			}},
		}
		if err = n.BootstrapCluster(clusterConfig).Error(); err != nil {
			logger.Error("failed to bootstrap raft cluster", zap.Error(err))
			return nil, err
		}
		logger.Info("bootstrap completed successfully")
	}

	return newNode, nil
}

func (n *node) JoinPeers(peers []Peer) error {
	for _, peer := range peers {
		if err := n.Join(peer.Id, peer.Address); err != nil {
			return err
		}
	}
	return nil
}

func (n *node) Join(id, address string) error {
	logger := n.logger.With(
		zap.String("peer_id", id),
		zap.String("peer_address", address),
	)

	logger.Debug("attempting to join peer")

	configFuture := n.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		logger.Error("failed to get configuration during join", zap.Error(err))
		return err
	}

	for _, server := range configFuture.Configuration().Servers {
		if server.ID == raft.ServerID(id) && server.Address == raft.ServerAddress(address) {
			logger.Info("peer already part of cluster")
			return nil
		}
		if server.ID == raft.ServerID(id) || server.Address == raft.ServerAddress(address) {
			logger.Warn("removing conflicting peer",
				zap.String("conflict_id", string(server.ID)),
				zap.String("conflict_address", string(server.Address)),
			)
			future := n.raft.RemoveServer(server.ID, 0, 0)
			if err := future.Error(); err != nil {
				logger.Error("failed to remove conflicting peer", zap.Error(err))
				return err
			}
		}
	}

	future := n.raft.AddVoter(raft.ServerID(id), raft.ServerAddress(address), 0, 0)
	if err := future.Error(); err != nil {
		logger.Error("failed to add peer as voter", zap.Error(err))
		return err
	}

	logger.Info("peer successfully joined")
	return nil
}

func (n *node) WaitForLeader(wait time.Duration) error {
	logger := n.logger.With(zap.Duration("timeout", wait))
	logger.Debug("waiting for leader election")

	timeout := time.After(wait)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			logger.Warn("leader election timed out")
			return ErrWaitNodeTimeout
		case <-ticker.C:
			if n.raft.State() == raft.Leader {
				logger.Info("node became leader")
				return nil
			}
		}
	}
}

func (n *node) Close() error {
	logger := n.logger
	logger.Info("shutting down")

	var errs []error

	if err := n.raft.Shutdown().Error(); err != nil {
		logger.Error("failed to shutdown", zap.Error(err))
		errs = append(errs, err)
	}
	if closer, ok := n.store.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			logger.Error("failed to close store", zap.Error(err))
			errs = append(errs, err)
		}
	}
	if closer, ok := n.trans.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			logger.Error("failed to close transport", zap.Error(err))
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("multiple shutdown errors: %+v", errs)
	}

	logger.Info("shutdown complete")
	return nil
}

func (n *node) Health() map[string]string {
	addr, id := n.raft.LeaderWithID()
	leader := ""
	if id != "" {
		leader = fmt.Sprintf("%s (%s)", id, addr)
	}
	return map[string]string{
		"state":  n.raft.State().String(),
		"leader": leader,
		"id":     n.config.Id,
	}
}
