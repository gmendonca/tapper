package elasticsearch

import "time"

type PrestoQuery struct {
	TimeTaken      string    `json:"time_taken"`
	MemoryPool     string    `json:"memory_pool"`
	State          string    `json:"state"`
	QueryStartTime time.Time `json:"query_start_time"`
	Source         string    `json:"source"`
	Query          string    `json:"query"`
	QueryID        string    `json:"query_id"`
	Timestamp      string    `json:"timestamp"`
	User           string    `json:"user"`
	UserAgent      string    `json:"user_agent"`
}
