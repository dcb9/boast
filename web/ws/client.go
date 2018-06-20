package ws

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dcb9/boast/transaction"
	"github.com/gorilla/websocket"
	"github.com/moul/http2curl"
	"io/ioutil"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []*transaction.Tx
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error %v", err)
			}
			break
		}
		log.Println("received message: ", string(message))
	}
}

func (c *Client) sendTxs(txs []*transaction.Tx) error {
	d := make([]Transaction, 0, len(txs))
	for _, tx := range txs {

		var rawReq string
		{
			rawHeaders := make([]string, 0, len(tx.Req.Header))
			for key, vals := range tx.Req.Header {
				rawHeaders = append(rawHeaders, key+": "+strings.Join(vals, ",")+"\r\n")
			}

			rawReq = fmt.Sprintf(
				"%s %s %s\r\n%s\r\n%s",
				tx.Req.Method, tx.Req.URL.String(), tx.Req.Proto,
				strings.Join(rawHeaders, ""), string(tx.Req.Body),
			)
		}

		var rawResp string
		{
			rawHeaders := make([]string, 0, len(tx.Resp.Header))
			for key, vals := range tx.Resp.Header {
				rawHeaders = append(rawHeaders, key+": "+strings.Join(vals, ",")+"\r\n")
			}

			contenType := tx.Resp.Header.Get("Content-Type")
			var body string
			if strings.Contains(contenType, "text") ||
				strings.Contains(contenType, "html") ||
				strings.Contains(contenType, "xml") ||
				strings.Contains(contenType, "json") ||
				strings.Contains(contenType, "javascript") {
				body = string(tx.Resp.Body)
			}

			rawResp = fmt.Sprintf(
				"%s %s\r\n%s\r\n%s",
				tx.Resp.Proto, tx.Resp.Status,
				strings.Join(rawHeaders, ""), body,
			)
		}

		body := ioutil.NopCloser(bytes.NewReader(tx.Req.Body))
		req, err := http.NewRequest(tx.Req.Method, tx.Req.URL.String(), body)
		req.Header = transaction.CopyHeader(tx.Req.Header)
		if err != nil {
			log.Println("http.NewRequest err", err)
			continue
		}

		curlCommand, _ := http2curl.GetCurlCommand(req)
		t := Transaction{
			ID: tx.ID,
			Request: Request{
				Method:      tx.Req.Method,
				Path:        tx.Req.URL.Path,
				RawText:     base64.StdEncoding.EncodeToString([]byte(rawReq)),
				CurlCommand: base64.StdEncoding.EncodeToString([]byte(curlCommand.String())),
			},
			Response: Response{
				Status:  tx.Resp.Status,
				RawText: base64.StdEncoding.EncodeToString([]byte(rawResp)),
			},
			ClientIP: "127.0.0.1",
			BeginAt:  tx.BeginAt,
			EndAt:    tx.EndAt,
		}
		d = append(d, t)
	}

	bytes, err := json.Marshal(d)
	if err != nil {
		return err
	}
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	w.Write(bytes)

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	if err := c.sendTxs(transaction.TxHub.List()); err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case tss, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.sendTxs(tss); err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func Serve(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []*transaction.Tx)}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}
