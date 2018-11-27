package model

type HiveQuery struct {
	CompileTime    string `json:"compile_time"`
	Timestamp      string `json:"timestamp"`
	Queue          string `json:"queue"`
	QueryID        string `json:"query_id"`
	User           string `json:"user"`
	QueryStartTime string `json:"query_start_time"`
	Query          string `json:"query"`
	TimeTaken      string `json:"time_taken"`
}
