package ws

import (
	"log"
	"time"

	"github.com/dcb9/boast/transaction"
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

type ReceivedMessage struct {
	Action string `json:"Action"`
	ID     string `json:"ID"`
}

func (m *ReceivedMessage) Do() {
	switch m.Action {
	case "replay":
		if id, err := uuid.Parse(m.ID); err != nil {
			log.Println(err)
		} else {
			transaction.Replay(id)
		}
	}
}
