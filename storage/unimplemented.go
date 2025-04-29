package storage

type unimplementedLogStorage struct{}

func (*unimplementedLogStorage) Append(entry *LogEntry) error {
	return nil
}

func (*unimplementedLogStorage) Get(index uint64) (*LogEntry, error) {
	return nil, nil
}

func (*unimplementedLogStorage) FirstIndex() (uint64, error) {
	return -1, nil
}

func (*unimplementedLogStorage) LastIndex() (uint64, error) {
	return -1, nil
}

func (*unimplementedLogStorage) TermByIndex(index uint64) (uint64, error) {
	return -1, nil
}

func (*unimplementedLogStorage) Truncate(index uint64) error {
	return nil
}

func (*unimplementedLogStorage) Close() error {
	return nil
}

type unimplementedStateStorage struct {
}

func (*unimplementedStateStorage) GetTerm() (uint64, error) {
	return -1, nil
}
func (*unimplementedStateStorage) SetTerm(uint64) error {
	return nil
}

func (*unimplementedStateStorage) GetVotedFor() (string, error) {
	return "", nil
}

func (*unimplementedStateStorage) SetVotedFor(string) error {
	return nil
}

func (*unimplementedStateStorage) Close() error {
	return nil
}
