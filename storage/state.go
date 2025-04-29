package storage

type StateStorage interface {
	GetTerm() (uint64, error)
	SetTerm(uint64) error
	GetVotedFor() (string, error)
	SetVotedFor(string) error
	Close() error
}
