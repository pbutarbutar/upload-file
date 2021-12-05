package models

type ApiResponse struct {
	ProcessTime int64                  `json:"process_time"`
	Success     bool                   `json:"success"`
	Message     string                 `json:"message"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Errors      map[string]string      `json:"errors,omitempty"`
}
