package types

type LogEntry struct {
	Term    uint64
	Command any
	Index   uint64
}
