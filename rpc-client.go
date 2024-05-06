package pcommon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn              *websocket.Conn
	requests          map[string]chan Response
	mu                sync.Mutex
	reconnect         bool
	connected         bool
	parserServerURL   string
	wg                sync.WaitGroup
	logging           bool
	reconnectInterval time.Duration
}

func NewRPCClient(url string, reconnectInterval time.Duration, logging bool) *Client {
	s := &Client{
		parserServerURL:   url,
		reconnect:         true,
		requests:          make(map[string]chan Response),
		logging:           logging,
		reconnectInterval: reconnectInterval,
	}
	return s
}

func (s *Client) Connect() {
	u, _ := url.Parse(s.parserServerURL)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		if s.logging {
			log.Println("Connection failed:", err)
		}
		if s.reconnect {
			time.AfterFunc(s.reconnectInterval, s.Connect)
		}
		return
	}
	s.conn = conn
	s.connected = true
	if s.logging {
		log.Println("Connection with parser is open")
	}
	go s.readMessages()
	time.Sleep(100 * time.Millisecond)
	go func() {
		for {
			if _, _, err := s.conn.NextReader(); err != nil {
				s.connected = false
				if s.reconnect {
					if s.logging {
						log.Println("Connection closed:", err)
					}
					time.AfterFunc(s.reconnectInterval, s.Connect)
				}
				break
			}
		}
	}()
}

func (s *Client) CheckConnectedError() error {
	if !s.connected {
		return fmt.Errorf("service is not connected")
	}
	return nil
}

func (s *Client) readMessages() {
	s.wg.Add(1)
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			s.wg.Done()
			if s.logging {
				log.Println("read:", err)
			}
			return
		}

		var resp Response
		if err := json.Unmarshal(message, &resp); err != nil {
			if s.logging {
				log.Println("Error unmarshaling response:", err)
			}
			continue
		}
		s.mu.Lock()
		if ch, ok := s.requests[resp.Id]; ok {
			ch <- resp
			close(ch)
			delete(s.requests, resp.Id)
		}
		s.mu.Unlock()
	}
}

func (s *Client) Request(method string, payload map[string]interface{}) (Response, error) {
	id := hashMethodAndPayload(method, payload)
	req := map[string]interface{}{
		"id":      id,
		"method":  method,
		"payload": payload,
	}
	reqData, _ := json.Marshal(req)
	s.mu.Lock()
	ch := make(chan Response)
	s.requests[id] = ch
	s.mu.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.conn.WriteMessage(websocket.TextMessage, reqData); err != nil {
			resp := Response{
				Id:    id,
				Data:  nil,
				Error: err.Error(),
			}
			ch <- resp
		}
	}()

	resp := <-ch
	return resp, nil
}

func hashMethodAndPayload(method string, payload map[string]interface{}) string {
	payloadData, _ := json.Marshal(payload)
	data := fmt.Sprintf("%s:%s", method, payloadData)

	//return a 15 length hash
	return fmt.Sprintf("%x", data)
}

func (s *Client) Stop() {
	s.reconnect = false
	if s.conn != nil && s.connected {
		message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
		s.conn.WriteMessage(websocket.CloseMessage, message)
		s.conn.Close()
		s.connected = false
	}
	s.wg.Wait()
}
