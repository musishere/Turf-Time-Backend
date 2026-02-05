package types

type Job struct {
	JobType string                 `json:"job_type"`
	Payload map[string]interface{} `json:"payload"`
}
