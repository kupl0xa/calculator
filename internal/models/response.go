package models

type TaskResultResponse struct {
	ID     int `json:"id"`
	Result int `json:"result"`
}

type TaskStatusResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}
