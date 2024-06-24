package crawl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	crawl ICrawl
	conn  *websocket.Conn
}

func NewClient(crawl ICrawl, ctx *gin.Context) (*Client, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		crawl: crawl,
		conn:  conn,
	}

	return client, nil
}

func (c *Client) Read() {
	defer func() {
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("error: %v", err)
	}

	c.conn.SetPongHandler(func(string) error {
		err = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Printf("error: %v", err)

			return err
		}
		return nil
	})

	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch msgType {
		case websocket.TextMessage:
			fmt.Println("Received Text Message:", string(msg))

		case websocket.BinaryMessage:
			fmt.Println("Received Binary Message:", msg)

		default:
			fmt.Println("Received Unknown Message Type")
		}
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		//case message, ok := <-c.send:
		//	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		//	if !ok {
		//		// The hub closed the channel.
		//		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		//		return
		//	}
		//
		//	w, err := c.conn.NextWriter(websocket.TextMessage)
		//	if err != nil {
		//		return
		//	}
		//	w.Write(message)
		//
		//	// Add queued chat messages to the current websocket message.
		//	n := len(c.send)
		//	for i := 0; i < n; i++ {
		//		w.Write(newline)
		//		w.Write(<-c.send)
		//	}
		//
		//	if err := w.Close(); err != nil {
		//		return
		//	}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
