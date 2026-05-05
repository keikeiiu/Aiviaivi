package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ailivili/internal/metrics"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// OnMessage is a callback invoked when a client sends a valid message.
// The handler layer sets this to process danmaku (validate, persist, broadcast).
var OnMessage func(c *Client, msg Message)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	videoID  string
	userID   string
	username string
	send     chan []byte
}

func NewClient(hub *Hub, w http.ResponseWriter, r *http.Request, videoID string, userID string, username string) (*Client, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	c := &Client{
		hub:      hub,
		conn:     conn,
		videoID:  videoID,
		userID:   userID,
		username: username,
		send:     make(chan []byte, 64),
	}

	hub.Join(videoID, c)
	go c.writePump()
	go c.readPump()

	metrics.IncWS()
	log.Printf("ws: client joined video=%s user=%s", videoID, userID)
	return c, nil
}

func (c *Client) readPump() {
	defer func() {
		c.hub.Leave(c.videoID, c)
		c.conn.Close()
		metrics.DecWS()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("ws read error: %v", err)
			}
			break
		}

		msg, err := ParseClientMessage(data)
		if err != nil {
			c.sendMsg(MakeError("invalid message format"))
			continue
		}

		if msg.Type != "danmaku" {
			c.sendMsg(MakeError("unknown message type"))
			continue
		}

		msg.VideoID = c.videoID
		msg.UserID = c.userID
		msg.Username = c.username

		if OnMessage != nil {
			OnMessage(c, msg)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case data, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) sendMsg(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	select {
	case c.send <- data:
	default:
	}
}

func (c *Client) Send(msg Message) {
	c.sendMsg(msg)
}

func (c *Client) IsAuthenticated() bool {
	return c.userID != ""
}

func (c *Client) UserID() string {
	return c.userID
}

func (c *Client) VideoID() string {
	return c.videoID
}
