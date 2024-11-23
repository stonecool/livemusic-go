package client

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/message"
	"log"
	"sync"
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
	conn    *websocket.Conn
	account account.IAccount
	done    chan struct{}
}

var (
	clients = make(map[account.IAccount]*Client)
	mu      sync.Mutex
)

func newClient(account account.IAccount, ctx *gin.Context) (*Client, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:    conn,
		account: account,
		done:    make(chan struct{}),
	}

	return client, nil
}

func (c *Client) readPump() {
	defer func() {
		close(c.done)
		c.conn.Close()
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		msg := &message.Message{}
		if err := json.Unmarshal(data, msg); err != nil {
			log.Printf("error parsing data: %v", err)
			continue
		}

		if err := c.handleWebSocketMessage(msg); err != nil {
			log.Printf("error handling websocket data: %v", err)
			continue
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
		case msg := <-c.account.GetMsgChan():
			if err := c.handleAccountMessage(msg); err != nil {
				log.Printf("handle account message error: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.done:
			return
		}
	}
}

func (c *Client) handleWebSocketMessage(msg *message.Message) error {
	switch msg.Cmd {
	case message.CrawlCmd_Login:
		asyncMessage := message.NewAsyncMessageWithMsg(msg, nil)

		c.account.GetMsgChan() <- asyncMessage
		result := <-asyncMessage.Result

		msg := &message.Message{
			Cmd:  message.CrawlCmd_Login,
			Data: []byte(fmt.Sprintf("%v", result)),
		}
		c.account.GetMsgChan() <- message.NewAsyncMessageWithMsg(msg, nil)
	}
	return nil
}

func (c *Client) handleAccountMessage(msg *message.AsyncMessage) error {
	return c.conn.WriteMessage(websocket.TextMessage, msg.Data)
}

func (c *Client) Close() {
	close(c.done)

	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.conn.Close()

	mu.Lock()
	defer mu.Unlock()
	delete(clients, c.account)
}

func HandleWebsocket(accountId int, ctx *gin.Context) error {
	acc, err := account.GetAccount(accountId)
	if err != nil {
		return err
	}

	client, err := newClient(acc, ctx)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	if oldClient, ok := clients[acc]; ok {
		oldClient.Close()
	}

	clients[acc] = client

	go client.readPump()
	go client.writePump()

	return nil
}
