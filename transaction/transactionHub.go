package transaction

import (
	"sync"

	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
)

const MAX_TRANSACTIONS_LEN int = 8 * 1024

var s sync.Mutex

var AddChannel chan *Tx = make(chan *Tx)

type Hub struct {
	Transactions map[uuid.UUID]*Tx
	SortID       []uuid.UUID
}

func NewTxHub() *Hub {
	hub := &Hub{
		Transactions: make(map[uuid.UUID]*Tx),
		SortID:       make([]uuid.UUID, 0, 32*1024),
	}
	hub.Init()
	return hub
}

var db *bolt.DB

func (h *Hub) Init() {
	var err error
	db, err = bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("transactions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("transactions"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var t Tx
			err = json.Unmarshal(v, &t)
			if err != nil {
				log.Println("json.Unmarshal err ", err)
				continue
			}

			h.Transactions[t.ID] = &t
			h.SortID = append(h.SortID, t.ID)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Hub) Add(t Tx) error {
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
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("transactions"))
		bytes, err := json.Marshal(t)
		if err != nil {
			return err
		}
		return b.Put(t.ID[:], bytes)
	})
	s.Unlock()
	if err != nil {
		return err
	}
	h.SortID = append(h.SortID, t.ID)

	AddChannel <- &t
	return nil
}

func (h *Hub) List() []*Tx {
	length := len(h.SortID)
	list := make([]*Tx, 0, len(h.SortID))
	for i := 0; i < length; i++ {
		transaction := h.Transactions[h.SortID[i]]
		list = append(list, transaction)
	}
	return list
}
