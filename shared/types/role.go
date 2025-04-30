package types

type Role uint8

const (
	Follower Role = iota
	Candidate
	Leader
)

var roleNames = map[Role]string{
	Follower:  "Follower",
	Candidate: "Candidate",
	Leader:    "Leader",
}

func (r Role) ToString() string {
	return roleNames[r]
}
