package domain

type ApiResponse struct {
	Success     bool                   `json:"success"`
	Status      int                    `json:"statusCode"`
	Message     string                 `json:"message"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Errors      map[string]string      `json:"errors,omitempty"`
}