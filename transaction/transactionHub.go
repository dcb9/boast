package transaction

import (
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const MAX_TRANSACTIONS_LEN int = 8 * 1024

var s sync.Mutex

var AddChannel chan *Ts = make(chan *Ts)

type Hub struct {
	Transactions map[uuid.UUID]*Ts
	SortID       []uuid.UUID
}

func (h *Hub) Add(t Ts) error {
	if t.ID == uuid.Nil {
		return errors.New("Transcation id MUST BE set.")
	}

	if len(h.SortID) >= MAX_TRANSACTIONS_LEN {
		deleteID := h.SortID[0]
		h.SortID = h.SortID[1:]
		delete(h.Transactions, deleteID)
	}

	s.Lock()
	h.Transactions[t.ID] = &t
	s.Unlock()
	h.SortID = append(h.SortID, t.ID)

	AddChannel <- &t
	return nil
}

func (h *Hub) List() []*Ts {
	length := len(h.SortID)
	list := make([]*Ts, 0, len(h.SortID))
	for i := 0; i < length; i++ {
		transaction := h.Transactions[h.SortID[i]]
		list = append(list, transaction)
	}
	return list
}
