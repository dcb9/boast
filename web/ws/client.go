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

var tsHub = transaction.TsHub

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []*transaction.Ts
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
		m := &ReceivedMessage{}
		if err := json.Unmarshal(message, m); err != nil {
			log.Println(err)
		}
		m.Do()
	}
}

func (c *Client) sendTss(tss []*transaction.Ts) error {
	d := make([]Transaction, 0, len(tss))
	for _, ts := range tss {

		var rawReq string
		{
			rawHeaders := make([]string, 0, len(ts.Req.Header))
			for key, vals := range ts.Req.Header {
				rawHeaders = append(rawHeaders, key+": "+strings.Join(vals, ",")+"\r\n")
			}

			rawReq = fmt.Sprintf(
				"%s %s %s\r\n%s\r\n%s",
				ts.Req.Method, ts.Req.URL.String(), ts.Req.Proto,
				strings.Join(rawHeaders, ""), string(ts.Req.Body),
			)
		}

		var rawResp string
		{
			rawHeaders := make([]string, 0, len(ts.Resp.Header))
			for key, vals := range ts.Resp.Header {
				rawHeaders = append(rawHeaders, key+": "+strings.Join(vals, ",")+"\r\n")
			}

			contenType := ts.Resp.Header.Get("Content-Type")
			var body string
			if strings.Contains(contenType, "text") ||
				strings.Contains(contenType, "html") ||
				strings.Contains(contenType, "xml") ||
				strings.Contains(contenType, "json") ||
				strings.Contains(contenType, "javascript") {
				body = string(ts.Resp.Body)
			}

			rawResp = fmt.Sprintf(
				"%s %s\r\n%s\r\n%s",
				ts.Resp.Proto, ts.Resp.Status,
				strings.Join(rawHeaders, ""), body,
			)
		}

		body := ioutil.NopCloser(bytes.NewReader(ts.Req.Body))
		req, err := http.NewRequest(ts.Req.Method, ts.Req.URL.String(), body)
		req.Header = transaction.CopyHeader(ts.Req.Header)
		if err != nil {
			log.Println("http.NewRequest err", err)
			continue
		}

		curlCommand, _ := http2curl.GetCurlCommand(req)
		t := Transaction{
			ID: ts.ID,
			Request: Request{
				Method:      ts.Req.Method,
				Path:        ts.Req.URL.Path,
				RawText:     base64.StdEncoding.EncodeToString([]byte(rawReq)),
				CurlCommand: base64.StdEncoding.EncodeToString([]byte(curlCommand.String())),
			},
			Response: Response{
				Status:  ts.Resp.Status,
				RawText: base64.StdEncoding.EncodeToString([]byte(rawResp)),
			},
			ClientIP: "127.0.0.1",
			BeginAt:  ts.BeginAt,
			EndAt:    ts.EndAt,
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

	if err := c.sendTss(tsHub.List()); err != nil {
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

			if err := c.sendTss(tss); err != nil {
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
	client := &Client{hub: hub, conn: conn, send: make(chan []*transaction.Ts)}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}
