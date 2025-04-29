package storage

import "thorac/core"

type unimplementedLogStorage struct{}

func (*unimplementedLogStorage) Append(entry *core.LogEntry) error {
	return nil
}

func (*unimplementedLogStorage) Get(index int) (*core.LogEntry, error) {
	return nil, nil
}

func (*unimplementedLogStorage) FirstIndex() (int, error) {
	return -1, nil
}

func (*unimplementedLogStorage) LastIndex() (int, error) {
	return -1, nil
}

func (*unimplementedLogStorage) TermByIndex(index int) (int, error) {
	return -1, nil
}

func (*unimplementedLogStorage) Truncate(index int) error {
	return nil
}

func (*unimplementedLogStorage) Close() error {
	return nil
}

type unimplementedStateStorage struct {
}

func (*unimplementedStateStorage) GetTerm() (int, error) {
	return -1, nil
}
func (*unimplementedStateStorage) SetTerm(int) error {
	return nil
}

func (*unimplementedStateStorage) GetVotedFor() (int, error) {
	return -1, nil
}

func (*unimplementedStateStorage) SetVotedFor(int) error {
	return nil
}

func (*unimplementedStateStorage) Close() error {
	return nil
}
