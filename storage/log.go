package storage

import "thorac/shared/types"

type LogStorage interface {
	Append(entry *types.LogEntry) error
	Get(index uint64) (*types.LogEntry, error)
	FirstIndex() (uint64, error)
	LastIndex() (uint64, error)
	TermByIndex(index uint64) (uint64, error)
	Truncate(index uint64) error
	Close() error
}
