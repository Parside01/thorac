package types

type Role uint8

const (
	Follower Role = iota
	Candidate
	Leader
)
