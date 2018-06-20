package ws

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	Method      string `json:"Method"`
	Path        string `json:"Path"`
	RawText     string `json:"RawText"`
	CurlCommand string `json:"CurlCommand"`
}

type Response struct {
	Status  string `json:"Status"`
	RawText string `json:"RawText"`
}

type Transaction struct {
	ID       uuid.UUID `json:"ID"`
	Request  Request   `json:"Req"`
	Response Response  `json:"Resp"`
	ClientIP string    `json:"ClientIP"`
	BeginAt  time.Time `json:"BeginAt"`
	EndAt    time.Time `json:"EndAt"`
}
