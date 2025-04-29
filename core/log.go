package core

type LogEntry struct {
	Term    int
	Command any
	Index   int
}

type LogStorage interface {
	Append(entry *LogEntry) error
	Get(index int) (*LogEntry, error)
	FirstIndex() (int, error)
	LastIndex() (int, error)
	TermByIndex(index int) (int, error)
	Truncate(index int) error
	Close() error
}
