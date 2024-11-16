package client

import (
	"encoding/json"
	"fmt"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/account"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	account     account.ICrawlAccount
	conn        *websocket.Conn
	accountChan chan *internal.Message // 用于接收来自 Account 的消息
	done        chan struct{}          // 用于关闭客户端
}

var (
	clients = make(map[account.ICrawlAccount]*Client)
	mu      sync.Mutex
)

func newClient(account account.ICrawlAccount, ctx *gin.Context) (*Client, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		account:     account,
		conn:        conn,
		accountChan: make(chan *internal.Message),
		done:        make(chan struct{}),
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
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		// 解析消息并处理
		msg := &internal.Message{}
		if err := json.Unmarshal(message, msg); err != nil {
			log.Printf("error parsing message: %v", err)
			continue
		}

		if err := c.handleWebSocketMessage(msg); err != nil {
			log.Printf("error handling websocket message: %v", err)
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
		case msg := <-c.accountChan:
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

func (c *Client) handleWebSocketMessage(msg *internal.Message) error {
	switch msg.Cmd {
	case internal.CrawlCmd_Login:
		// 创建任务消息
		asyncMessage := NewAsyncMessage(msg)

		// 发送任务并等待结果
		c.account.GetMsgChan() <- asyncMessage
		result := <-asyncMessage.Result

		// 发送结果回 WebSocket
		c.accountChan <- &internal.Message{
			Cmd:  internal.CrawlCmd_Login,
			Data: []byte(fmt.Sprintf("%v", result)),
		}
	}
	return nil
}

func (c *Client) handleAccountMessage(msg *internal.Message) error {
	// 将消息写入 WebSocket
	return c.conn.WriteMessage(websocket.TextMessage, msg.Data)
}

func (c *Client) Close() {
	// 关闭 done channel 来通知所有 goroutine 退出
	close(c.done)

	// 关闭 WebSocket 连接
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.conn.Close()

	// 关闭 accountChan
	close(c.accountChan)

	// 从全局 clients map 中移除
	mu.Lock()
	delete(clients, c.account)
	mu.Unlock()
}

func HandleWebsocket(accountId int, ctx *gin.Context) error {
	account, err := GetCrawl(accountId)
	if err != nil {
		return err
	}

	mu.Lock()
	client, err := newClient(account, ctx)
	if err != nil {
		mu.Unlock()
		return err
	}

	if oldClient, ok := clients[account]; ok {
		oldClient.Close()
	}

	clients[account] = client
	mu.Unlock()

	go client.readPump()
	go client.writePump()

	return nil
}
