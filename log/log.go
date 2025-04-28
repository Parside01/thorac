package log

type Entry struct {
	Term    int
	Command any
	Index   int
}

type Storage interface {
	Append(entry *Entry) error
	Get(index int) (*Entry, error)
	FirstIndex() (int, error)
	LastIndex() (int, error)
	TermByIndex(index int) (int, error)
	Truncate(index int) error
}
