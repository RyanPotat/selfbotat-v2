package pubsub

import (
	"os"
	"os/signal"
	"sync"
	"time"

	"selfbotat-v2/bot/logger"

	"github.com/gorilla/websocket"
)

const (
	PingInterval = 200 * 1000
	MaxSize      = 50
)

type ReadyState int

const (
	Connecting ReadyState = 0
	Open       ReadyState = 1
	Closing    ReadyState = 2
	Closed     ReadyState = 3
)

var Pool = make(map[int]WSPool)

type WSPool struct {
	Connections map[int]*SingleCon
	mutex sync.Mutex
}

type SingleCon struct {
	websocket.Conn
	Topics []string
	ID int 
}

func (p *WSPool) AddConn(ws *SingleCon) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Connections[ws.ID] = ws
}

func (p *WSPool) RemoveConn(id int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.Connections, id)
}

func (p *WSPool) CreateClient(uri string, topics []string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ws, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		Log.Error.Println("Error connecting to PubSub: ", err)
	}
	defer ws.Close()

	id := len(p.Connections) + 1

	done := make(chan struct{})

  go CreateReader(ws, done)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			p.RemoveConn(id)
			return
		case t := <-ticker.C:
			err := ws.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				Log.Debug.Println("write:", err)
				return
			}
		case <-interrupt:
			Log.Debug.Println("interrupt")

			err := ws.WriteMessage(
				websocket.CloseMessage, 
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				Log.Debug.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func CreateReader(ws *websocket.Conn, done chan struct{}) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			Log.Debug.Println("read:", err)
			close(done)
			return
		}
		Log.Debug.Printf("recv: %s", message)
	}
}