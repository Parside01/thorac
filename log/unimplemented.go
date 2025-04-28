package log

type unimplementedStorage struct{}

func (*unimplementedStorage) Append(entry *Entry) error {
	return nil
}

func (*unimplementedStorage) Get(index int) (*Entry, error) {
	return nil, nil
}

func (*unimplementedStorage) FirstIndex() (int, error) {
	return -1, nil
}

func (*unimplementedStorage) LastIndex() (int, error) {
	return -1, nil
}

func (*unimplementedStorage) TermByIndex(index int) (int, error) {
	return -1, nil
}

func (*unimplementedStorage) Truncate(index int) error {
	return nil
}
