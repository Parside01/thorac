package main

type SagaStep struct {
	Name       string       `json:"name"`
	Execute    *HttpRequest `json:"execute"`
	Compensate *HttpRequest `json:"compensate"`
}

type Saga struct {
	Id    string     `json:"id"`
	Steps []SagaStep `json:"steps"`
}
