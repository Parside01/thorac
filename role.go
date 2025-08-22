package main

type Role uint8

const (
	Follower Role = iota
	Candidate
	Leader
	Dead
)

var roleNames = map[Role]string{
	Follower:  "Follower",
	Candidate: "Candidate",
	Leader:    "Leader",
	Dead:      "Dead",
}

func (r Role) ToString() string {
	return roleNames[r]
}
