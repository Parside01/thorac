package storage

type StateStorage interface {
	GetTerm() (int, error)
	SetTerm(int) error
	GetVotedFor() (int, error)
	SetVotedFor(int) error
	Close() error
}
