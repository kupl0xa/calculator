package models

import "errors"

var ErrNotFound = errors.New("note not found")

type Task struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Operator string `json:"operator"`
	Result   int    `json:"result"`
	Status   string `json:"status"`
}
